package ar

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

// Each Connection can have multiple loggers, and each logger
// will be logged to based on the logType that is passed. One
// of the following log types are the possibilities.
const (
	LOG_ALL = iota
	LOG_QUERY
	LOG_ERROR
)

var TransactionLimitError = errors.New("Transaction Limit Reached")

type Connection struct {
	DB             *sql.DB
	Dialect        Dialect
	dbName         string
	stmtMap        map[string]*sql.Stmt
	combinedLogs   []Logger
	errorLogs      []Logger
	queryLogs      []Logger
	txBlock        bool
	txCount, txMax int
	txMutex        sync.Mutex
	txWait         []chan int
	sources        map[string]*source
	Config         *Config
}

type Tx struct {
	Tx        *sql.Tx
	Conn      *Connection
	txStmtMap map[string]*sql.Stmt
}

func NewConnection(dialectName, dbName, connector string) (*Connection, error) {
	conn := new(Connection)
	if dialect, found := registeredDialects[dialectName]; found {
		conn.Dialect = dialect
	} else {
		return nil, fmt.Errorf("Could not locate dialect '%s'", dialectName)
	}

	if db, err := sql.Open(dialectName, connector); err == nil {
		db.SetMaxIdleConns(100)
		conn.DB = db
	} else {
		return nil, err
	}

	conn.dbName = dbName
	conn.stmtMap = make(map[string]*sql.Stmt)

	return conn, nil
}

func (c *Connection) StartTransaction() (*Tx, error) {
	c.txCount++
	if c.txMax > 0 {
		c.txMutex.Lock()
		if c.txCount >= c.txMax {
			if c.txBlock {
				waiter := make(chan int)
				c.txWait = append(c.txWait, waiter)
				c.txMutex.Unlock()
				<-waiter
			} else {
				c.txMutex.Unlock()
				c.txCount--
				return nil, TransactionLimitError
			}
		}
	}
	sqlTx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{Tx: sqlTx, Conn: c, txStmtMap: make(map[string]*sql.Stmt)}, nil
}

func (tx *Tx) Commit() error {
	err := tx.Tx.Commit()
	tx.release()
	return err
}

func (tx *Tx) release() {
	tx.Conn.txCount--
	if tx.Conn.txMax > 0 {
		tx.Conn.txMutex.Lock()
		if len(tx.Conn.txWait) > 0 {
			tx.Conn.txWait[0] <- 1
			tx.Conn.txWait = tx.Conn.txWait[1:]
		}
		tx.Conn.txMutex.Unlock()
	}
}

func (tx *Tx) Rollback() error {
	err := tx.Tx.Rollback()
	tx.release()
	return err
}

// ar will set the pool size for a connection to 100, if you need
// a different pool size, you can do it with this function
func (c *Connection) ChangePoolSize(size int) {
	c.DB.SetMaxIdleConns(size)
}

// A Connection can be closed, which essentially means that the
// *sql.DB connection is closed, though it should not be counted
// on that the Close operation will not clear out other structures
// or will clear out other structures.
func (c *Connection) Close() error {
	return c.DB.Close()
}

// A Logger is a struct that has a subset of a log.Logger, you
// can use a log.Logger for it, but you can substitute a different
// struct of your own imagining if you wish.
func (c *Connection) SetLogger(logger Logger, logType int) {
	switch logType {
	case LOG_ALL:
		c.combinedLogs = append(c.combinedLogs, logger)
	case LOG_QUERY:
		c.queryLogs = append(c.queryLogs, logger)
	case LOG_ERROR:
		c.errorLogs = append(c.errorLogs, logger)
	}
}

// Sets the limit of concurrent transactions, blocking determines
// whether Connection.StartTransaction will block on waiting
// or return a TransactionLimitError. blocking is set to false
// by default for a new Connection. Setting max to 0 will not
// disable Transactions. To disable Transactions, it would be
// suggested that you set the max to 1, call StartTransaction, then
// call Commit on the sql.Tx manually. So long as Commit or
// Rollback is not called on the ar.Tx object, all future
// StartTransaction calls will not succeed.
func (c *Connection) SetTransactionLimit(max int, blocking bool) {
	c.txBlock = blocking
	currentTx := c.txMax - c.txCount
	if currentTx < 0 {
		currentTx = -currentTx
	}

	c.txMax = max
	c.txCount = max - currentTx
}

func (c *Connection) newSource(name string, ptr interface{}, Options []map[string]map[string]interface{}) *source {
	structType := c.getType(ptr)

	s := new(source)
	s.conn = c
	s.Name = name
	s.SqlName = c.Config.StructToTable(name)
	s.config = c.Config
	s.Fields = c.createMappingsFromType(structType)
	c.createSqlMappings(s)
	c.propagateOptions(s, Options)
	s.structName = structType.Name()

	return s
}

func (c *Connection) getType(ptr interface{}) reflect.Type {
	currentType := reflect.TypeOf(ptr)
	for currentType.Kind() == reflect.Ptr {
		currentType = currentType.Elem()
	}

	return currentType
}

func (c *Connection) createMappingsFromType(structType reflect.Type) []*sourceMapping {
	output := make([]*sourceMapping, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		mapping := new(sourceMapping)
		options := new(structOptions)
		mapping.structOptions = options
		field := structType.Field(i)

		options.Name = field.Name
		options.Index = i
		options.Kind = field.Type.Kind()
		options.ColumnHint = field.Tag.Get("colName")
		options.Options = c.parseFieldOptions(field.Tag)

		output = append(output, mapping)
	}

	return output
}

func (c *Connection) propagateOptions(s *source, Options []map[string]map[string]interface{}) {
	for _, optionSet := range Options {
		if allOptions, ok := optionSet["all"]; ok {
			for _, field := range s.Fields {
				for key, value := range allOptions {
					field.structOptions.Options[key] = value
				}
			}
		}
		for column, colOptions := range optionSet {
			for _, field := range s.Fields {
				if field.structOptions.Name == column {
					for key, value := range colOptions {
						field.structOptions.Options[key] = value
					}
				}
			}
		}
	}
}

func (c *Connection) parseFieldOptions(tag reflect.StructTag) map[string]interface{} {
	options := make(map[string]interface{})
	optionString := string(tag)

	for optionString != "" {
		// following code is adapted from the golang reflect package
		i := 0
		for i < len(optionString) && optionString[i] == ' ' {
			i++
		}
		optionString = optionString[i:]
		if optionString == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(optionString) && optionString[i] != ' ' && optionString[i] != ':' && optionString[i] != '"' {
			i++
		}
		if i+1 >= len(optionString) || optionString[i] != ':' || optionString[i+1] != '"' {
			break
		}
		name := string(optionString[:i])
		optionString = optionString[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(optionString) && optionString[i] != '"' {
			if optionString[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(optionString) {
			break
		}
		qvalue := string(optionString[:i+1])
		optionString = optionString[i+1:]

		if name[:2] == "ar" {
			options[name[2:]], _ = strconv.Unquote(qvalue)
		} else {
			options[name[2:]], _ = strconv.Unquote(qvalue)
		}
	}

	return options
}

func (c *Connection) createSqlMappings(s *source) {
	for _, column := range c.Dialect.ColumnsInTable(c, c.dbName, s.SqlName) {
		if column.Number+1 > s.ColNum {
			s.ColNum = column.Number + 1
		}

		for _, field := range s.Fields {
			if c.Config.FieldToColumn(field.structOptions.Name) == column.Name {
				field.columnInfo = column
				break
			}
		}
	}
}
