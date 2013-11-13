package db

import (
	"errors"
	"fmt"
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

func (m *source) Count() (int64, error) {
	return m.Identity().Count()
}

func (m *source) Delete() error {
	return m.Identity().Delete()
}

func (m *source) UpdateAttribute(column string, val interface{}) error {
	return m.Identity().UpdateAttribute(column, val)
}

func (m *source) UpdateAttributes(values Attributes) error {
	return m.Identity().UpdateAttributes(values)
}

func (m *source) UpdateSql(sql string, vals ...interface{}) error {
	return m.Identity().UpdateSql(sql, vals...)
}

func (m *source) Pluck(column, vals interface{}) error {
	return m.Identity().Pluck(column, vals)
}

func (m *source) Initialize(vals ...interface{}) error {
	for _, val := range vals {
		rt := reflect.TypeOf(val)
		switch rt.Kind() {
		case reflect.Ptr:
			if rt.Elem().Kind() == reflect.Struct {
				return m.initializeSingle(val)
			}
			if rt.Elem().Kind() == reflect.Slice {
				rv := reflect.ValueOf(val).Elem()
				var err error
				for i := 0; i < rv.Len(); i++ {
					err = m.initializeItem(rv.Index(i))
					if err != nil {
						return err
					}
				}

			}
		case reflect.Array, reflect.Slice:
			rv := reflect.ValueOf(val)
			var err error
			for i := 0; i < rv.Len(); i++ {
				err = m.initializeItem(rv)
				if err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("Could not recognize instance to be initialized")
		}
	}
	return nil
}

func (m *source) initializeSingle(val interface{}) error {
	vfn := fullNameFor(getType(val))
	if vfn != m.FullName {
		return fmt.Errorf("Mapper for %s cannot initialize %s", m.FullName, vfn)
	}

	rv := reflect.ValueOf(val)
	if rv.Elem().Type().Kind() != reflect.Struct {
		return fmt.Errorf("Could not map data of kind %s", rv.Kind().String())
	}

	if !m.hasMixin {
		return fmt.Errorf("Struct %s does not have a mixin to initialize", m.Name)
	}

	mx := new(Mixin)
	mx.model = m
	mx.instance = val
	if rv.Elem().Field(m.mixinField).IsNil() {
		rv.Elem().Field(m.mixinField).Set(reflect.ValueOf(mx))
	}

	return nil
}

func (m *source) initializeItem(val reflect.Value) error {
	for val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	vfn := fullNameFor(val.Type())
	if vfn != m.FullName {
		return fmt.Errorf("Slice Mapper for %s cannot initialize %s", m.FullName, vfn)
	}

	if val.Type().Kind() != reflect.Struct {
		return fmt.Errorf("Could not map data of kind %s", val.Kind().String())
	}

	if !m.hasMixin {
		return fmt.Errorf("Struct %s does not have a mixin to initialize", m.Name)
	}

	mx := new(Mixin)
	mx.model = m
	mx.instance = val
	val.Field(m.mixinField).Set(reflect.ValueOf(mx))

	return nil
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
	for c, _ := range values {
		if c == m.ID.SqlColumn {
			delete(values, c)
		}
	}
	query, vals := m.conn.Dialect.Create(m, values)
	if m.conn.Dialect.CreateExec() {
		result, err := m.runExec(query, vals)
		if err != nil {
			return err
		}
		newId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		m.setIntID(v, newId)
	} else {
		result := m.runQueryRow(query, vals)
		if m.ID.Kind == reflect.Int {
			var newId int64
			err := result.Scan(&newId)
			if err != nil {
				return err
			}
			m.setIntID(v, newId)
		} else {
			var newId interface{}
			err := result.Scan(&newId)
			if err != nil {
				return err
			}
			m.setID(v, newId)

		}
	}

	return nil
}

func (m *source) extractID(v reflect.Value) interface{} {
	fieldIndex := m.ID.Index
	return v.Field(fieldIndex).Interface()
}

func (m *source) setID(v reflect.Value, id interface{}) {
	fieldIndex := m.ID.Index
	v.Field(fieldIndex).Set(reflect.ValueOf(id))
}
func (m *source) setIntID(v reflect.Value, id int64) {
	fieldIndex := m.ID.Index
	v.Field(fieldIndex).SetInt(id)
}
func (m *source) extractColumnValues(v reflect.Value) map[string]interface{} {
	output := make(map[string]interface{})
	for _, field := range m.Fields {
		if field.ColumnInfo != nil && field.structOptions != nil {
			value := v.Field(field.structOptions.Index)
			output[field.ColumnInfo.Name] = value.Interface()
		}
	}
	if m.hasMixin {
		mxf := v.Field(m.mixinField)
		if !mxf.IsNil() {
			if mi, ok := mxf.Interface().(*Mixin); ok {
				for name, _ := range output {
					if mi.IsNull(name) {
						output[name] = nil
					}
				}
			}
		}
	}
	return output
}

func (s *source) TableName() string {
	return s.SqlName
}

func (s *source) PrimaryKeyColumn() string {
	return s.ID.SqlColumn
}

func (s *source) LeftJoin(joins ...interface{}) Scope {
	return s.Identity().LeftJoin(joins...)
}
func (s *source) InnerJoin(joins ...interface{}) Scope {
	return s.Identity().InnerJoin(joins...)
}
func (s *source) FullJoin(joins ...interface{}) Scope {
	return s.Identity().FullJoin(joins...)
}
func (s *source) RightJoin(joins ...interface{}) Scope {
	return s.Identity().RightJoin(joins...)
}
func (s *source) JoinSql(sql string, args ...interface{}) Scope {
	return s.Identity().JoinSql(sql, args...)
}

func (s *source) LeftInclude(include ...interface{}) Scope {
	return s.Identity().LeftInclude(include...)
}
func (s *source) InnerInclude(include ...interface{}) Scope {
	return s.Identity().InnerInclude(include...)
}
func (s *source) FullInclude(include interface{}, nullRecords interface{}) Scope {
	return s.Identity().FullInclude(include, nullRecords)
}
func (s *source) RightInclude(include interface{}, nullRecords interface{}) Scope {
	return s.Identity().RightInclude(include, nullRecords)
}
func (s *source) IncludeSql(il IncludeList, query string, args ...interface{}) Scope {
	return s.Identity().IncludeSql(il, query, args...)
}
