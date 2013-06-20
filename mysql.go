package ar

import (
	"reflect"
)

type mysqlDialect struct {
	base
}

func newMysql() Dialect {
	d := new(mysqlDialect)
	d.base.dialect = d
	return d
}

func (d mysqlDialect) CompatibleSqlTypes(f reflect.Type) []string {
	switch f.Kind() {
	case reflect.Struct:
		if f.String() == "time.Time" {
			return []string{"timestamp"}
		}
	case reflect.Bool:
		return []string{"boolean"}
	case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
		return []string{"int", "integer", "bigint"}
	case reflect.Int64, reflect.Uint64:
		return []string{"bigint"}
	case reflect.Int8, reflect.Uint8:
		return []string{"tinyint", "smallint", "mediumint", "int", "integer", "bigint"}
	case reflect.Int16, reflect.Uint16:
		return []string{"mediumint", "int", "integer", "bigint"}
	case reflect.Float32:
		return []string{"double", "float"}
	case reflect.Float64:
		return []string{"double"}
	case reflect.Slice:
		if f.String() == "[]uint8" { //[]byte
			return []string{"varbinary", "longblob"}
		}
	case reflect.String:
		return []string{"varchar", "text", "longtext"}
	}
	panic("invalid sql type")
}
