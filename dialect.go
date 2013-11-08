package db

import (
	"reflect"
)

var registeredDialects map[string]Dialect

func init() {
	registeredDialects = make(map[string]Dialect)

	registeredDialects["mysql"] = newMysql()
	registeredDialects["sqlite3"] = newSqlite()
	registeredDialects["postgres"] = newPostgres()
}

// If you have an external dialect, use this function to
// load it in so that consumers can create Connections using
// it. You can override builtin dialects by naming your
// dialect the same as a builtin dialect.
func RegisterDialect(name string, dialect Dialect) {
	registeredDialects[name] = dialect
}

type Dialect interface {
	ColumnsInTable(conn *Connection, dbName string, tableName string) map[string]*ColumnInfo
	CompatibleSqlTypes(f reflect.Type) []string
	FormatQuery(query string) string
	Query(scope Scope) (string, []interface{})
	Create(mapper Mapper, values map[string]interface{}) (string, []interface{})
	CreateExec() bool
	Update(scope Scope, values map[string]interface{}) (string, []interface{})
	Delete(scope Scope) (string, []interface{})

	// Quote will quote identifiers in a SQL statement.
	Quote(s string) string

	// Being Replaced
	//InsertSql(queryable *Queryable) (sql string, args []interface{})
	//DeleteSql(queryable *Queryable) (string, []interface{})
	//UpdateSql(queryable *Queryable) (string, []interface{})
}
