package ar

import (
	"reflect"
)

type source struct {
	ID     *sourceMapping
	Name   string
	config *Config
	Fields []*sourceMapping

	structName, tableName string
}

type sourceMapping struct {
	*structOptions
	*columnInfo
}

func (sm *sourceMapping) Column() string {
	return sm.SqlTable + "." + sm.SqlColumn
}

type structOptions struct {
	Name       string
	Index      int
	Kind       reflect.Kind
	ColumnHint string
	Options    map[string]interface{}
}

type columnInfo struct {
	Name      string
	SqlTable  string
	SqlColumn string
	SqlType   int
	Length    int
	Nullable  bool
	Key       int
	Default   interface{}
	Extra     map[string]interface{}
}
