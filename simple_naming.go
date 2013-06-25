package ar

import (
	"strings"
)

func NewSimpleConfig() *Config {
	c := new(Config)
	c.StructToTable = strings.ToLower
	c.TableToStruct = strings.ToTitle
	c.FieldToColumn = strings.ToLower
	c.ColumnToField = strings.ToTitle
	c.IdName = "ID"
	c.CreatedColumn = "Creation"
	c.UpdatedColumn = "Modified"

	return c
}
