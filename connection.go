package ar

import (
	"database/sql"
	"errors"
	"fmt"
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
	models         map[string]*Mapper
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

	if db, err := sql.Open(conn.Dialect.SqlName(), connector); err == nil {
		db.SetMaxIdleConns(100)
		conn.DB = db
	} else {
		return nil, err
	}

	conn.dbName = dbName
	conn.stmtMap = make(map[string]*sql.Stmt)

	return conn, nil
}

func (c *Connection) IndexExists(table, index string) bool {
	return c.Dialect.IndexExists(c.DB, dbName, table, index)
}

func (c *Connection) ColumnsForTable(table interface{}) []string {
	columnMap := c.Dialect.ColumnsInTable(c.DB, c.dbName, table)
	columns := make([]string, len(columnMap))
	i := 0
	for k, _ := range columnMap {
		columns[i] = k
		i++
	}

	return columns
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
