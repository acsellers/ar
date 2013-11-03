package db

import (
	"fmt"
	"reflect"
)

type Mixin struct {
	instance interface{}
	model    *source
}

func (m *Mixin) Init(instance interface{}) error {
	tn := fullNameFor(getType(instance))
	s := mappedStructs[tn]
	if s == nil {
		return fmt.Errorf("Could not locate a mapper for this struct")
	}
	if s.multiMapped {
		return fmt.Errorf("Must use InitWithConn for initializing this struct")
	}

	s.Initialize(instance)
	return nil
}

func (m *Mixin) InitWithConn(conn *Connection, instance interface{}) error {
	tn := fullNameFor(getType(instance))
	s := conn.mappedStructs[tn]
	if s == nil {
		return fmt.Errorf("Could not locate a mapper for this struct")
	}

	s.Initialize(instance)

	return nil
}

func (m *Mixin) Save() error {
	return m.model.SaveAll(m.instance)
}

func (m *Mixin) selfScope() Scope {
	id := m.model.extractID(reflect.ValueOf(m.instance).Elem())
	return m.model.EqualTo(m.model.ID.SqlColumn, id)
}
func (m *Mixin) Delete() error {
	return m.selfScope().Delete()
}

func (m *Mixin) UpdateAttribute(attr string, value interface{}) error {
	return m.selfScope().UpdateAttribute(attr, value)
}

func (m *Mixin) UpdateAttributes(values map[string]interface{}) error {
	return m.selfScope().UpdateAttributes(values)
}
