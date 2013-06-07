package ar

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type mysqlDialect struct {
	base
}

func newMysql() Dialect {
	d := new(mysqlDialect)
	d.base.dialect = d
	return d
}

func (mysqlDialect) ParseBool(value reflect.Value) bool {
	return value.Int() != 0
}

func (d mysqlDialect) SqlType(f interface{}, size int) string {
	switch f.(type) {
	case time.Time:
		return "timestamp"
	case bool:
		return "boolean"
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return "int"
	case int64, uint64:
		return "bigint"
	case float32, float64:
		return "double"
	case []byte:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("varbinary(%d)", size)
		}
		return "longblob"
	case string:
		if size > 0 && size < 65532 {
			return fmt.Sprintf("varchar(%d)", size)
		}
		return "longtext"
	}
	panic("invalid sql type")
}

func (d mysqlDialect) indexExists(db *sql.DB, dbName, tableName, indexName string) bool {
	var row *sql.Row
	var name string
	row = db.QueryRow("SELECT INDEX_NAME FROM INFORMATION_SCHEMA.STATISTICS "+
		"WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?", dbName, tableName, indexName)
	row.Scan(&name)
	return name != ""
}

func (d mysqlDialect) primaryKeySql(isString bool, size int) string {
	if isString {
		return fmt.Sprintf("varchar(%d) PRIMARY KEY", size)
	}
	return "bigint PRIMARY KEY AUTO_INCREMENT"
}
