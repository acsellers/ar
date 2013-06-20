package ar

import (
	"database/sql"
	"reflect"
)

var registeredDialects map[string]Dialect

func init() {
	registeredDialects = make(map[string]Dialect)

	registeredDialects["mysql"] = newMysql()
}

// If you have an external dialect, use this function to
// load it in so that consumers can create Connections using
// it. You can override builtin dialects by naming your
// dialect the same as a builtin dialect.
func RegisterDialect(name string, dialect Dialect) {
	registeredDialects[name] = dialect
}

type Dialect interface {
	ColumnsInTable(db *sql.DB, dbName string, tableName interface{}) map[string]*columnInfo
	CompatibleSqlTypes(f reflect.Type) []string
	FormatQuery(query string) string
	Query(queryable *Queryable) (string, []interface{})
	Update(queryable *Queryable, values map[string]interface{}) (string, []interface{})

	// Quote will quote identifiers in a SQL statement.
	Quote(s string) string

	// Being Replaced
	//InsertSql(queryable *Queryable) (sql string, args []interface{})
	//DeleteSql(queryable *Queryable) (string, []interface{})
	//UpdateSql(queryable *Queryable) (string, []interface{})
}
