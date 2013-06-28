package ar

import (
	"errors"
	"reflect"
)

type Mapper struct {
	source *source
}

// Identity will create a queryable for a specific mapper,
// in this way, the different Queryable methods are simply
// implemented in the manner mapper.Identity().Method().
func (m *Mapper) Identity() *Queryable {
	return &Queryable{source: m.source}
}

func (m *Mapper) Where(fragment string, args ...interface{}) *Queryable {
	return m.Identity().Where(fragment, args...)
}

func (m *Mapper) EqualTo(column string, val interface{}) *Queryable {
	return m.Identity().EqualTo(column, val)
}

func (m *Mapper) Between(column string, lower, upper interface{}) *Queryable {
	return m.Identity().Between(column, lower, upper)
}

func (m *Mapper) In(column string, vals []interface{}) *Queryable {
	return m.Identity().In(column, vals)
}

func (m *Mapper) Limit(limit int) *Queryable {
	return m.Identity().Limit(limit)
}

func (m *Mapper) Offset(offset int) *Queryable {
	return m.Identity().Offset(offset)
}

func (m *Mapper) OrderBy(column, direction string) *Queryable {
	return m.Identity().OrderBy(column, direction)
}

func (m *Mapper) Order(ordering string) *Queryable {
	return m.Identity().Order(ordering)
}

func (m *Mapper) Reorder(ordering string) *Queryable {
	return m.Identity().Reorder(ordering)
}

// Find looks for the record with primary key equal to val
func (m *Mapper) Find(val interface{}) *Queryable {
	return m.Identity().Find(val)
}

func (m *Mapper) RetrieveAll(val interface{}) error {
	return m.Identity().RetrieveAll(val)
}

func (m *Mapper) Save(val interface{}) error {
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

func (m *Mapper) saveSlice(v reflect.Value) error {
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

func (m *Mapper) saveItem(v reflect.Value) error {
	ident := m.extractID(v)
	if ident == 0 {
		return m.createItem(v)
	}
	values := m.extractColumnValues(v)
	return m.Find(ident).UpdateAttributes(values)
}

func (m *Mapper) createItem(v reflect.Value) error {
	values := m.extractColumnValues(v)
	query, vals := m.source.conn.Dialect.Create(m, values)
	result, err := m.source.runExec(query, vals)
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

func (m *Mapper) extractID(v reflect.Value) int64 {
	fieldIndex := m.source.ID.Index
	return v.Field(fieldIndex).Int()
}

func (m *Mapper) setID(v reflect.Value, id int64) {
	fieldIndex := m.source.ID.Index
	v.Field(fieldIndex).SetInt(id)
}

func (m *Mapper) extractColumnValues(v reflect.Value) map[string]interface{} {
	output := make(map[string]interface{})
	for _, field := range m.source.Fields {
		if field.columnInfo != nil && field.structOptions != nil {
			value := v.Field(field.structOptions.Index)
			output[field.columnInfo.Name] = value.Interface()
		}
	}
	return output
}
