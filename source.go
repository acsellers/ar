package db

import (
	"database/sql"
	"reflect"
)

type source struct {
	ID          *sourceMapping
	Name        string
	FullName    string
	SqlName     string
	ColNum      int
	hasMixin    bool
	mixinField  int
	multiMapped bool
	config      *Config
	conn        *Connection
	Fields      []*sourceMapping

	structName, tableName string
}

type sourceMapping struct {
	*structOptions
	*ColumnInfo
}

func (sm *sourceMapping) Column() string {
	return sm.SqlTable + "." + sm.SqlColumn
}

type structOptions struct {
	Name       string
	Mapped     bool
	Relation   *source
	FullName   string
	Index      int
	Kind       reflect.Kind
	ColumnHint string
	Options    map[string]interface{}
	ForeignKey *ColumnInfo
}

// ColumnInfo is the data returned by a ColumnsInTable function which is
// implemented by the database Dialect's.
type ColumnInfo struct {
	// Name will be set to the struct field name, you can set it
	// to the column name if you wish without harming anything
	Name string
	// The table that this column belongs to
	SqlTable string
	// The name for the database column
	SqlColumn string
	// The type, this should correlate to a type given by the Dialect
	// function CompatibleSqlTypes
	SqlType string
	// The Length of the field, you should set this to -1 for fields that
	// have no effective limit
	Length int
	// Whether this field could return a NULL value, this is safe to mark
	// as true if in doubt. It activates nil protection for mapping
	Nullable bool
	// The index of the column within the table, it is an optional field
	Number int
}

func (s *source) runQuery(query string, values []interface{}) (*sql.Rows, error) {
	return s.conn.Query(query, values...)
}

func (s *source) runQueryRow(query string, values []interface{}) *sql.Row {
	return s.conn.QueryRow(query, values...)
}

func (s *source) runExec(query string, values []interface{}) (sql.Result, error) {
	return s.conn.Exec(query, values...)
}
