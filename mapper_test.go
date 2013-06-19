package ar

import (
	"strings"
	"testing"
)

var testConfig = &Config{
	StructToTable: strings.ToLower,
	TableToStruct: strings.Title,
	FieldToColumn: func(s string) string {
		return s
	},
	ColumnToField: func(s string) string {
		return s
	},
	IdName: "ID",
}

func TestMapperBasic(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup")
		c := &Connection{
			Dialect: newMysql(),
			dbName:  "testdb",
			sources: make(map[string]*source),
			Config:  testConfig,
		}

		userMapper, err := c.CreateMapper("User", &user{})
		test.IsNil(err)
		test.IsNotNil(userMapper)

	})
}

type user struct {
	ID   int64
	Name string
}
