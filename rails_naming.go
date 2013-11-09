package db

import (
	"github.com/acsellers/inflections"
	"strings"
)

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
