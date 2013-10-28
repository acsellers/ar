package db

import (
	"errors"
	"reflect"
)

func (s *source) Identity() Scope {
	return &queryable{source: s}
}

func (s *source) Where(fragment string, args ...interface{}) Scope {
	return s.Identity().Where(fragment, args...)
}

func (s *source) Cond(column string, condition int, val interface{}) Scope {
	return s.Identity().Cond(column, condition, val)
}

func (m *source) EqualTo(column string, val interface{}) Scope {
	return m.Identity().EqualTo(column, val)
}

func (m *source) Between(column string, lower, upper interface{}) Scope {
	return m.Identity().Between(column, lower, upper)
}

func (m *source) Limit(limit int) Scope {
	return m.Identity().Limit(limit)
}

func (m *source) Offset(offset int) Scope {
	return m.Identity().Offset(offset)
}

func (m *source) OrderBy(column, direction string) Scope {
	return m.Identity().OrderBy(column, direction)
}

func (m *source) Order(ordering string) Scope {
	return m.Identity().Order(ordering)
}

func (m *source) Reorder(ordering string) Scope {
	return m.Identity().Reorder(ordering)
}

// Find looks for the record with primary key equal to val
func (m *source) Find(id interface{}, val interface{}) error {
	return m.Identity().Find(id, val)
}

func (m *source) Retrieve(val interface{}) error {
	return m.Identity().Retrieve(val)
}

func (m *source) RetrieveAll(val interface{}) error {
	return m.Identity().RetrieveAll(val)
}

func (m *source) SaveAll(val interface{}) error {
	vv := reflect.ValueOf(val)
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		vv = vv.Elem()
	}
	vk := vv.Type().Kind()
	if !(vk == reflect.Slice || vk == reflect.Struct) {
		return errors.New("Was not passed mappable values")
	}
	if vk == reflect.Slice {
		return m.saveSlice(vv)
	} else {
		return m.saveItem(vv)
	}
}

func (m *source) saveSlice(v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		vi := v.Index(i)
		if vi.Type().Kind() == reflect.Ptr {
			vi = vi.Elem()
		}
		err := m.saveItem(vi)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *source) saveItem(v reflect.Value) error {
	ident := m.extractID(v)
	if ident == 0 {
		return m.createItem(v)
	}
	values := m.extractColumnValues(v)
	return m.EqualTo(m.ID.SqlColumn, ident).UpdateAttributes(values)
}

func (m *source) createItem(v reflect.Value) error {
	values := m.extractColumnValues(v)
	query, vals := m.conn.Dialect.Create(m, values)
	result, err := m.runExec(query, vals)
	if err != nil {
		return err
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.setID(v, newId)

	return nil
}

func (m *source) extractID(v reflect.Value) int64 {
	fieldIndex := m.ID.Index
	return v.Field(fieldIndex).Int()
}

func (m *source) setID(v reflect.Value, id int64) {
	fieldIndex := m.ID.Index
	v.Field(fieldIndex).SetInt(id)
}

func (m *source) extractColumnValues(v reflect.Value) map[string]interface{} {
	output := make(map[string]interface{})
	for _, field := range m.Fields {
		if field.columnInfo != nil && field.structOptions != nil {
			value := v.Field(field.structOptions.Index)
			output[field.columnInfo.Name] = value.Interface()
		}
	}
	return output
}

func (s *source) TableName() string {
	return s.SqlName
}
