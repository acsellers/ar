package ar

import (
	"database/sql"
)

type Connection struct {
	DB        *sql.DB
	Dialect   Dialect
	dbName    string
	txStmtMap map[string]*sql.Stmt
}

func NewConnection(dialectName, dbName, connector string) *Connection {
	return nil
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
