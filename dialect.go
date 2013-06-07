package ar

import (
	"database/sql"
	"reflect"
)

var registeredDialects map[string]Dialect

func init() {
	registeredDialects = make(map[string]Dialect)

	registeredDialects["mysql"] = newMysql()
	registeredDialects["oracle"] = newOracle()
}

// If you have an external dialect, use this function to
// load it in so that consumers can create Connections using
// it. You can override builtin dialects by naming your
// dialect the same as a builtin dialect.
func RegisterDialect(name string, dialect Dialect) {
	registeredDialects[name] = dialect
}

type Dialect interface {

	//Substitute "?" marker if database use other symbol as marker
	SubstituteMarkers(query string) string

	// Quote will quote identifiers in a SQL statement.
	Quote(s string) string

	SqlType(f interface{}, size int) string

	ParseBool(value reflect.Value) bool

	SetModelValue(value reflect.Value, field reflect.Value) error

	QuerySql(criteria *criteria) (sql string, args []interface{})

	Insert(q *Qbs) (int64, error)

	InsertSql(criteria *criteria) (sql string, args []interface{})

	Update(q *Qbs) (int64, error)

	UpdateSql(criteria *criteria) (string, []interface{})

	Delete(q *Qbs) (int64, error)

	DeleteSql(criteria *criteria) (string, []interface{})

	CreateTableSql(model *model, ifNotExists bool) string

	DropTableSql(table string) string

	AddColumnSql(table, column string, typ interface{}, size int) string

	CreateIndexSql(name, table string, unique bool, columns ...string) string

	IndexExists(db *sql.DB, dbName, tableName, indexName string) bool

	ColumnsInTable(db *sql.DB, dbName string, tableName interface{}) map[string]bool

	PrimaryKeySql(isString bool, size int) string

	CatchMigrationError(err error) bool
}
