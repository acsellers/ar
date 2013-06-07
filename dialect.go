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

type Dialect interface {

	//Substitute "?" marker if database use other symbol as marker
	SubstituteMarkers(query string) string

	// Quote will quote identifiers in a SQL statement.
	Quote(s string) string

	sqlType(f interface{}, size int) string

	parseBool(value reflect.Value) bool

	setModelValue(value reflect.Value, field reflect.Value) error

	querySql(criteria *criteria) (sql string, args []interface{})

	insert(q *Qbs) (int64, error)

	insertSql(criteria *criteria) (sql string, args []interface{})

	update(q *Qbs) (int64, error)

	updateSql(criteria *criteria) (string, []interface{})

	delete(q *Qbs) (int64, error)

	deleteSql(criteria *criteria) (string, []interface{})

	createTableSql(model *model, ifNotExists bool) string

	dropTableSql(table string) string

	addColumnSql(table, column string, typ interface{}, size int) string

	createIndexSql(name, table string, unique bool, columns ...string) string

	indexExists(db *sql.DB, dbName, tableName, indexName string) bool

	columnsInTable(db *sql.DB, dbName string, tableName interface{}) map[string]bool

	primaryKeySql(isString bool, size int) string

	catchMigrationError(err error) bool
}
