package ar

import (
	"database/sql"
	"reflect"
)

type source struct {
	ID      *sourceMapping
	Name    string
	SqlName string
	ColNum  int
	config  *Config
	conn    *Connection
	Fields  []*sourceMapping

	structName, tableName string
}

type sourceMapping struct {
	*structOptions
	*columnInfo
}

func (sm *sourceMapping) Column() string {
	return sm.SqlTable + "." + sm.SqlColumn
}

type planner struct {
	scanners []interfaceable
}

func (s *source) mapPlan(v reflector) *planner {
	p := &planner{make([]interfaceable, s.ColNum)}
	for i, _ := range p.scanners {
		p.scanners[i] = new(nullScanner)
	}
	for _, col := range s.Fields {
		p.scanners[col.columnInfo.Number] = &reflectScanner{parent: v, index: col.columnInfo.Number}
	}

	return p
}

func (p *planner) Items() []interface{} {
	output := make([]interface{}, len(p.scanners))
	for i, _ := range output {
		output[i] = p.scanners[i].iface()
	}

	return output
}

type reflectScanner struct {
	parent reflector
	index  int
}

type interfaceable interface {
	iface() interface{}
}

type reflector struct {
	item reflect.Value
}

func (rf *reflectScanner) iface() interface{} {
	return rf.parent.item.Elem().Field(rf.index).Addr().Interface()
}

type nullScanner struct {
}

func (n *nullScanner) iface() interface{} {
	return interface{}(n)
}

func (*nullScanner) Scan(v interface{}) error {
	return nil
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
	SqlType   string
	Length    int
	Nullable  bool
	Number    int
}

func (s *source) runQuery(query string, values []interface{}) (*sql.Rows, error) {
	return s.conn.DB.Query(query, values...)
}

func (s *source) runExec(query string, values []interface{}) (sql.Result, error) {
	return s.conn.DB.Exec(query, values...)
}
