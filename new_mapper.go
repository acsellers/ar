package db

import (
	"reflect"
)

// CreateMapper returns a Mapper instance for the mapee struct you passed
func (c *Connection) CreateMapper(name string, mapee interface{}) (Mapper, error) {
	// create and save the source (primary mapper interfacee)
	ms := c.newSource(name, mapee)
	c.sources[name] = ms

	c.createRelations(ms)

	return ms, nil
}

// This is the same as CreateMapper, but will panic on an error instead
// of returning it
func (c *Connection) MustCreateMapper(name string, v interface{}) Mapper {
	m, e := c.CreateMapper(name, v)
	if e != nil {
		panic(e)
	}
	return m
}

func (c *Connection) createRelations(s *source) {
	if _, found := mappedStructs[s.FullName]; found {
		mappedStructs[s.FullName] = &source{multiMapped: true}
	} else {
		mappedStructs[s.FullName] = s
	}
	c.mappedStructs[s.FullName] = s

	if dependents, ok := c.mappableStructs[s.FullName]; ok {
		for _, d := range dependents {
			d.refreshRelated(s.FullName)
		}
	}
	delete(c.mappableStructs, s.FullName)

	s.loadRelated()
}

// this function is to make testing with multiple
// RDBMS's simpler
func (c *Connection) m(name string) Mapper {
	if s, ok := c.sources[name]; ok {
		return s
	}
	return nil
}

// Initialize a MapperPlus instance, need more documentation
// and tests on this
func (c *Connection) InitMapperPlus(name string, v interface{}, Options ...map[string]map[string]interface{}) {
	rv := reflect.ValueOf(v).Elem()
	fv := rv.Field(0)
	mp := new(mapperPlus)
	vmp := reflect.ValueOf(mp)
	if fv.Type().Kind() == reflect.Ptr {
		fv.Set(vmp)
	}
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
