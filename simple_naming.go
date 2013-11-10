package db

import (
	"strings"
)

/*
Simple config is the simplest Config implementation
I could think of. Table and Column names are lowercase
versions of struct and field names, foreign keys are
lowercased field names with 'id' appended. Primary key
names are Id and id, with timestamps of Creation and
Modified.
*/
func NewSimpleConfig() *Config {
	c := new(Config)
	c.StructToTable = strings.ToLower
	c.FieldToColumn = func(s, f string) string {
		return strings.ToLower(f)
	}
	c.ForeignKeyName = func(fn, sn string) string {
		return strings.ToLower(fn) + "id"
	}

	c.IdName = func(s string) (string, string) {
		return "Id", "id"
	}
	c.CreatedColumn = "Creation"
	c.UpdatedColumn = "Modified"

	return c
}
