package db

import (
	"database/sql"
	"reflect"
	"strings"
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
	relations   []*sourceMapping
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

func (sm *sourceMapping) Aliased() bool {
	return !strings.HasSuffix(sm.FullName, ":"+sm.structOptions.Name)
}
func (sm *sourceMapping) MappedColumn() bool {
	return sm.structOptions != nil && sm.ColumnInfo != nil
}
func (sm *sourceMapping) NamedBy(s string) bool {
	switch {
	case sm.MappedColumn():
		return sm.structOptions.Name == s || sm.ColumnInfo.SqlColumn == s
	case sm.structOptions != nil:
		return sm.structOptions.Name == s
	case sm.ColumnInfo != nil:
		return sm.ColumnInfo.SqlColumn == s
	default:
		return false
	}
}

type structOptions struct {
	Name       string
	Mapped     bool
	FullName   string
	Index      int
	Kind       reflect.Kind
	ColumnHint string
	Options    map[string]interface{}
	Relation   *source
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

func (s *source) loadRelated() {
	var reload bool
	for _, f := range s.Fields {
		if !f.Mapped && f.ColumnInfo == nil {
			if rs, ok := s.conn.mappedStructs[f.FullName]; ok {
				f.Relation = rs
				s.relations = append(s.relations, f)
				reload = true
			} else {
				s.conn.mappableStructs[f.FullName] = append(s.conn.mappableStructs[f.FullName], s)
			}
		}
	}
	if reload {
		s.locateForeignKeys()
	}
}
func (s *source) locateForeignKeys() {
	for _, f := range s.relations {
		if f.ForeignKey == nil {
			// we either have a has_many or a habtm
			if f.Kind == reflect.Slice {
				// we're going to search through the fields and try to find a matching
				// field for our struct, we'll need to beef this up so we have 1 path
				// for unaliased structs and 1 for aliased structs (aliased would be like
				// Author User vs User User), if we find a matching field, we'll get the
				// foreign key name for that field, then re-iterate through our relation's
				// fields to find the one that matches the foreign key name
				for _, rf := range f.Relation.Fields {
					if rf.ColumnInfo == nil && rf.FullName == s.FullName {
						kn := s.conn.Config.ForeignKeyName(rf.structOptions.Name, rf.FullName)
						for _, pfk := range f.Relation.Fields {
							if pfk.ColumnInfo != nil && pfk.ColumnInfo.SqlColumn == kn {
								f.ForeignKey = pfk.ColumnInfo
								continue
							}
						}
					}
				}
			} else {
				// find belongs to relation fields
				kn := s.conn.Config.ForeignKeyName(f.structOptions.Name, f.FullName)
				for _, pfk := range s.Fields {
					if pfk.ColumnInfo != nil && pfk.ColumnInfo.SqlColumn == kn {
						f.ForeignKey = pfk.ColumnInfo
						continue
					}
				}
				// find has one relations
				for _, pfk := range f.Relation.Fields {
					if pfk.ColumnInfo != nil && pfk.ColumnInfo.SqlColumn == kn {
						f.ForeignKey = pfk.ColumnInfo
						continue
					}
				}

			}
		}
	}
}
func (s *source) refreshRelated(sn string) {
	ns := s.conn.mappedStructs[sn]
	for _, f := range s.Fields {
		if f.FullName == sn {
			f.Relation = ns
			s.relations = append(s.relations, f)
		}
	}
	s.locateForeignKeys()
}
