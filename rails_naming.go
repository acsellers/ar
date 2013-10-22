package db

import (
	"github.com/acsellers/inflections"
	"strings"
)

func NewRailsConfig() *Config {
	c := new(Config)
	c.StructToTable = pluralizeStruct
	c.TableToStruct = singularizeTable
	c.FieldToColumn = inflections.Underscore
	c.ColumnToField = inflections.Camelize
	c.IdName = "Id"
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
