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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return []string{"int"}
	case reflect.Int64, reflect.Uint64:
		return []string{"bigint"}
	case reflect.Float32, reflect.Float64:
		return []string{"double"}
	case reflect.Slice:
		if f.String() == "[]uint8" { //[]byte
			return []string{"varbinary", "longblob"}
		}
	case reflect.String:
		return []string{"varchar", "longtext"}
	}
	panic("invalid sql type")
}
