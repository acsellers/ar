package ar

import (
	"database/sql"
	"fmt"
)

type Connection struct {
	DB      *sql.DB
	Dialect Dialect
	dbName  string
	stmtMap map[string]*sql.Stmt
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
	sqlTx, err := c.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{Tx: sqlTx, Conn: c, txStmtMap: make(map[string]*sql.Stmt)}, nil
}
