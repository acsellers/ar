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

// New dialects must fufill this interface, sqlite, mysql and postgres
// have internal dialects that implement this interface
type Dialect interface {
	// Given a database name and table name, output the information for each
	// column listed in the table
	ColumnsInTable(conn *Connection, dbName string, tableName string) map[string]*ColumnInfo
	// This function isn't used at the moment, but when db gets migration
	// support, it will be used
	CompatibleSqlTypes(f reflect.Type) []string
	// Format query takes in a query that uses ?'s as placeholders for
	// query parameters, and formats the query as needed for the database
	FormatQuery(query string) string
	// Query creates a simple select query
	Query(scope Scope) (string, []interface{})
	// Create takes a mapper and a map of the sql column names with the
	// values set on a struct instance corresponding to the mapped
	// fields corresponding to those columns and returns the query and
	// any necessary parameters
	Create(mapper Mapper, values map[string]interface{}) (string, []interface{})
	// If the database doesn't repsond to LastInsertId when using Exec,
	// you can add a RETURNING predicate in Create and return false here
	// so struct instances can get their primary key values
	CreateExec() bool
	// Write an update statement similar to the create statement and
	// return the parameterized query form
	Update(scope Scope, values map[string]interface{}) (string, []interface{})
	// Formulate a delete from statement from a scope and return
	// the sql and parameters
	Delete(scope Scope) (string, []interface{})
	// The GROUP BY sql syntax has some differences from database to database
	// Mysql allows you to specify a single column that determines the table
	// to group by, while sqlite, postgres, others require every non-aggregated
	// field to be present. It would be onerous to ask developers to write
	// every field for database systems that require it, so we'll do it for
	// them, but we need to know whether we need to do it for this database
	// system
	ExpandGroupBy() bool
}
