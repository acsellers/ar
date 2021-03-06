package db

import (
	"reflect"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var typeRegex = regexp.MustCompile("^([a-zA-Z0-9]+)(\\([0-9]+\\))?(.*)")

type mysqlDialect struct {
	Base
}

func newMysql() Dialect {
	d := new(mysqlDialect)
	d.Base.Dialect = d
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

func (d mysqlDialect) ColumnsInTable(conn *Connection, dbName string, table string) map[string]*ColumnInfo {
	query := "SHOW COLUMNS FROM " + table
	if dbName != "" {
		query = "SHOW COLUMNS FROM " + table + " IN " + dbName
	}

	output := make(map[string]*ColumnInfo)
	rows, err := conn.DB.Query(query)
	if err != nil {
		defer rows.Close()
		panic(err)
		return nil
	}

	var name, sqlType, key, extra string
	var def string
	var notnull bool
	var number int
	for rows.Next() {
		ci := new(ColumnInfo)

		rows.Scan(&name, &sqlType, &notnull, &key, &def, &extra)
		ci.Name = name
		ci.SqlTable = table
		ci.SqlColumn = name
		ci.Nullable = !notnull
		ci.SqlType = d.sqlTypeFrom(sqlType)
		ci.Length = d.sqlLengthFrom(sqlType)
		ci.Number = number
		output[name] = ci
		number++
	}

	return output
}

func (d mysqlDialect) sqlTypeFrom(st string) string {
	if typeRegex.MatchString(st) {
		return typeRegex.FindStringSubmatch(st)[1]
	}
	return st
}

func (d mysqlDialect) sqlLengthFrom(st string) int {
	if typeRegex.MatchString(st) {
		matches := typeRegex.FindStringSubmatch(st)
		if len(matches) > 3 && len(matches[2]) > 0 {
			strLength := len(matches[2])
			lenStr := matches[2][1 : strLength-1]
			i, err := strconv.ParseInt(lenStr, 10, 32)
			if err != nil {
				return -1
			}
			return int(i)
		}
	}

	return 0
}
