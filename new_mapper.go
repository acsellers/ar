package db

import (
	"reflect"
	"strconv"
)

func (c *Connection) CreateMapper(name string, mapee interface{}, Options ...map[string]map[string]interface{}) (Mapper, error) {
	msource := c.newSource(name, mapee, Options)
	c.sources[name] = msource
	if _, found := mappedStructs[msource.FullName]; found {
		mappedStructs[msource.FullName] = &source{multiMapped: true}
	} else {
		mappedStructs[msource.FullName] = msource
	}
	c.mappedStructs[msource.FullName] = msource

	return msource, nil
}

func (c *Connection) m(name string) Mapper {
	if s, ok := c.sources[name]; ok {
		return s
	}
	return nil
}

func (c *Connection) CreateMapperPlus(name string, v interface{}, Options ...map[string]map[string]interface{}) {
	rv := reflect.ValueOf(v).Elem()
	fv := rv.Field(0)
	mp := new(mapperPlus)
	vmp := reflect.ValueOf(mp)
	if fv.Type().Kind() == reflect.Ptr {
		fv.Set(vmp)
	}
}

func (c *Connection) MustCreateMapper(name string, v interface{}) Mapper {
	m, e := c.CreateMapper(name, v)
	if e != nil {
		panic(e)
	}
	return m
}

func (c *Connection) createMapperForPtr(ptr interface{}) (string, *Mapper) {
	return "model", new(Mapper)
}

func (c *Connection) tableNameFor(ptr interface{}) string {
	if t, ok := ptr.(string); ok {
		return t
	}

	t := reflect.TypeOf(ptr).Elem()
	if t.Kind() != reflect.Struct {
		return ""
	}

	return c.Config.StructToTable(t.Name())
}
