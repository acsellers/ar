package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/dchest/cache"
)

// Each Connection can have multiple loggers, and each logger
// will be logged to based on the logType that is passed. One
// of the following log types are the possibilities.
const (
	LOG_ALL = iota
	LOG_QUERY
	LOG_ERROR
)

var mappedStructs = make(map[string]*source)

type Connection struct {
	// This is the main Database connection used by the Connection
	DB *sql.DB
	// The Dialect that the Connection uses to format queries
	Dialect         Dialect
	mappedStructs   map[string]*source
	mappableStructs map[string][]*source
	dbName          string
	stmtMap         map[string]*sql.Stmt
	combinedLogs    []Logger
	errorLogs       []Logger
	queryLogs       []Logger
	sources         map[string]*source
	// The Config for mapping structs to database tables and records
	Config *Config
	// The QueryCache for saving and reusing Queries. To disable, simply set this to nil
	// By default, it will store up to 4096 distinct queries, you can use the
	// CacheSize(n int) to change the query storage number
	QueryCache *cache.Cache
}

/*
NewConnection creates a connection to a specific database server with
a specific database. Connections are used to create Mappers, which are
then used to retrieve records. Note that NewConnection will set the
MaxIdleConns to 100 for the database connection.
*/
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
	conn.mappedStructs = make(map[string]*source)
	conn.mappableStructs = make(map[string][]*source)
	conn.sources = make(map[string]*source)
	conn.QueryCache = cache.New(cache.Config{
		MaxItems:        4096,
		TrackAccessTime: true,
		RemoveHandler: func(p cache.Item) {
			if q, ok := p.Value.(*sql.Stmt); ok {
				q.Close()
			}
		},
	})

	return conn, nil
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

func getType(ptr interface{}) reflect.Type {
	currentType := reflect.TypeOf(ptr)
	for currentType.Kind() == reflect.Ptr {
		currentType = currentType.Elem()
	}
	return currentType
}
func fullNameFor(t reflect.Type) string {
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() + ":" + t.Name()
}
func (c *Connection) newSource(name string, ptr interface{}) *source {
	structType := getType(ptr)

	s := new(source)
	s.conn = c
	s.FullName = fullNameFor(structType)
	s.Name = name
	s.SqlName = c.Config.StructToTable(name)
	s.config = c.Config
	s.Fields = c.createMappingsFromType(structType)
	mn := fullNameFor(reflect.TypeOf(Mixin{}))
	idName, _ := c.Config.IdName(name)
	for _, field := range s.Fields {
		if field.structOptions.Name == idName {
			s.ID = field
		}
		if field.FullName != "" {
			if field.FullName == mn {
				s.hasMixin = true
				s.mixinField = field.Index
				field.Mapped = true
			}
		}
	}
	c.createSqlMappings(s)
	//c.propagateOptions(s, Options)
	s.structName = structType.Name()

	return s
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
		if options.Kind == reflect.Ptr || options.Kind == reflect.Struct || options.Kind == reflect.Slice {
			rt := field.Type
			for rt.Kind() == reflect.Ptr {
				rt = rt.Elem()
			}
			options.FullName = fullNameFor(rt)
		}
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
			if c.Config.FieldToColumn(s.Name, field.structOptions.Name) == column.Name {
				field.ColumnInfo = column
				break
			}
		}
	}
}

// This is almost the same as Connection.DB.Query(query, args), but will
// prepare the statement and cache it in the connection's query cache.
// The resemblence to the database/sql DB interface is intentional.
//
// rows, e := pgConn.Query(bigQuery, values...)
func (c *Connection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, e := c.getQuery(query)
	if e == nil {
		return stmt.Query(args...)
	} else {
		return c.DB.Query(query, args...)
	}
}

// This is almost the same as Connection.DB.QueryRow(query, args), but will
// prepare the statement and cache it in the connection's query cache.
// The resemblence to the database/sql DB interface is intentional.
//
// row := pgConn.QueryRow(complicated, values...)
func (c *Connection) QueryRow(query string, args ...interface{}) *sql.Row {
	stmt, e := c.getQuery(query)
	if e == nil {
		return stmt.QueryRow(args...)
	} else {
		return c.DB.QueryRow(query, args...)
	}
}

// This is almost the same as Connection.DB.Exec(query, args), but will
// prepare the statement and cache it in the connection's query cache.
// The resemblence to the database/sql DB interface is intentional.
//
// result, e := pgConn.Exec(createThings, values...)
func (c *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	stmt, e := c.getQuery(query)
	if e == nil {
		return stmt.Exec(args...)
	} else {
		return c.DB.Exec(query, args...)
	}
}

func (c *Connection) CacheSize(n int) {
	c.QueryCache.Reconfigure(cache.Config{
		MaxItems:        n,
		TrackAccessTime: true,
		RemoveHandler: func(p cache.Item) {
			if q, ok := p.Value.(*sql.Stmt); ok {
				q.Close()
			}
		},
	})
}

func (c *Connection) getQuery(query string) (*sql.Stmt, error) {
	if c.QueryCache == nil {
		return nil, fmt.Errorf("QueryCache not enabled")
	}

	i, ok := c.QueryCache.Get(query)
	if ok {
		if q, ok := i.(*sql.Stmt); ok {
			return q, nil
		}
	}
	q, e := c.DB.Prepare(query)
	if e != nil {
		return nil, e
	}
	c.QueryCache.Set(query, q, 0)

	return q, nil
}
