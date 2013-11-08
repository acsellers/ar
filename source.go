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
	FullName   string
	Index      int
	Kind       reflect.Kind
	ColumnHint string
	Options    map[string]interface{}
}

type ColumnInfo struct {
	Name      string
	SqlTable  string
	SqlColumn string
	SqlType   string
	Length    int
	Nullable  bool
	Number    int
}

func (s *source) runQuery(query string, values []interface{}) (*sql.Rows, error) {
	return s.conn.DB.Query(query, values...)
}

func (s *source) runQueryRow(query string, values []interface{}) *sql.Row {
	return s.conn.DB.QueryRow(query, values...)
}

func (s *source) runExec(query string, values []interface{}) (sql.Result, error) {
	return s.conn.DB.Exec(query, values...)
}
