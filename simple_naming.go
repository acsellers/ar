package db

import (
	"strings"
)

func NewSimpleConfig() *Config {
	c := new(Config)
	c.StructToTable = strings.ToLower
	c.FieldToColumn = strings.ToLower
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
