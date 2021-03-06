package db

import (
	"reflect"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

type sqliteDialect struct {
	Base
}

func newSqlite() Dialect {
	d := new(sqliteDialect)
	d.Base.Dialect = d
	return d
}

func (d sqliteDialect) CompatibleSqlTypes(f reflect.Type) []string {
	switch f.Kind() {
	case reflect.Struct:
		if f.String() == "time.Time" {
			return []string{"Integer"}
		}
	case reflect.Bool, reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return []string{"Integer"}
	case reflect.Float32, reflect.Float64:
		return []string{"Float"}
	case reflect.Slice:
		if f.String() == "[]uint8" { //[]byte
			return []string{"Blob"}
		}
	case reflect.String:
		return []string{"Text"}
	}
	return []string{}
}

func (d sqliteDialect) ColumnsInTable(conn *Connection, dbName string, table string) map[string]*ColumnInfo {
	query := "PRAGMA table_info(" + table + ")"

	output := make(map[string]*ColumnInfo)
	rows, err := conn.DB.Query(query)
	if err != nil {
		defer rows.Close()
		panic(err)
		return nil
	}

	var name, sqlType string
	var extra1, extra2 string
	var notnull bool
	var number int
	for rows.Next() {
		ci := new(ColumnInfo)

		rows.Scan(&number, &name, &sqlType, &notnull, &extra1, &extra2)
		ci.Number = number
		ci.Name = name
		ci.SqlTable = table
		ci.SqlColumn = name
		ci.Nullable = !notnull
		ci.SqlType = sqlType
		ci.Length = -1
		output[name] = ci
	}

	return output
}

func (d sqliteDialect) Query(scope Scope) (string, []interface{}) {
	out, args := d.Base.Query(scope)
	return out, args
}
