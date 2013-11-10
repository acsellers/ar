package db

import (
	"github.com/acsellers/inflections"
	"strings"
)

/*
The RailsConfig tries to replicate Rails practices for
naming.

Table names are pluralized lowercase versions of struct
names and field names are lowercased versions of field names.
Foreign key names are the lowercased name of field with '_id'
appended.

Ids are identified with the field name "Id" and the database
column "id", while the timestamp fields are "CreatedAt" and
"UpdatedAt".
*/
func NewRailsConfig() *Config {
	c := new(Config)
	c.StructToTable = pluralizeStruct
	c.FieldToColumn = func(s, f string) string {
		return inflections.Underscore(f)
	}
	c.ForeignKeyName = func(fn, sn string) string {
		return strings.ToLower(fn) + "_id"
	}
	c.IdName = func(s string) (string, string) {
		return "Id", "id"
	}
	c.CreatedColumn = "CreatedAt"
	c.UpdatedColumn = "UpdatedAt"

	return c
}

func pluralizeStruct(s string) string {
	return inflections.Pluralize(strings.ToLower(s))
}

func singularizeTable(s string) string {
	return strings.ToTitle(s)
}
