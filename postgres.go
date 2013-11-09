package db

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

type postgresDialect struct {
	Base
}

func newPostgres() Dialect {
	d := new(postgresDialect)
	d.Base.Dialect = d
	return d
}

func (d postgresDialect) CompatibleSqlTypes(f reflect.Type) []string {
	switch f.Kind() {
	case reflect.Struct:
		if f.String() == "time.Time" {
			return []string{"timestamp", "time", "timetz", "timestamptz", "date"}
		}
	case reflect.Bool:
		return []string{"boolean", "bool"}
	case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
		return []string{"int", "integer", "bigint", "int", "int4"}
	case reflect.Int64, reflect.Uint64:
		return []string{"bigint"}
	case reflect.Int8, reflect.Uint8:
		return []string{"int2", "smallint"}
	case reflect.Int16, reflect.Uint16:
		return []string{"int", "integer"}
	case reflect.Float32, reflect.Float64:
		return []string{"double precision", "float8"}
	case reflect.Slice:
		if f.String() == "[]uint8" { //[]byte
			return []string{"bytea"}
		}
	case reflect.String:
		return []string{"character varying", "text"}
	}
	return []string{}
}

func (d postgresDialect) ColumnsInTable(conn *Connection, dbName string, table string) map[string]*ColumnInfo {
	query := `SELECT column_name, data_type, is_nullable, COALESCE(character_maximum_length, -1), ordinal_position
FROM information_schema.columns WHERE table_catalog = $1 AND table_name = $2`
	output := make(map[string]*ColumnInfo)
	rows, err := conn.DB.Query(query, dbName, table)
	if err != nil {
		panic(err)
		return nil
	}
	defer rows.Close()

	var name, sqlType string
	var nullable string
	var number, length int
	for rows.Next() {
		ci := new(ColumnInfo)

		err = rows.Scan(&name, &sqlType, &nullable, &length, &number)
		if err == nil {
			ci.Name = name
			ci.SqlTable = table
			ci.SqlColumn = name
			if nullable == "YES" {
				ci.Nullable = true
			} else {
				ci.Nullable = false
			}
			ci.SqlType = sqlType
			ci.Length = length
			ci.Number = number - 1
			output[name] = ci
		} else {
			fmt.Println(err)
		}
	}
	return output
}

func (d postgresDialect) FormatQuery(query string) string {
	parts := strings.Split(query, "?")
	var newQuery []string
	for i, part := range parts[:len(parts)-1] {
		newQuery = append(newQuery, fmt.Sprintf("%s$%d", part, i+1))
	}
	newQuery = append(newQuery, parts[len(parts)-1])

	return strings.Join(newQuery, "")
}

func (d postgresDialect) CreateExec() bool {
	return false
}

func (d postgresDialect) Create(mapper Mapper, values map[string]interface{}) (string, []interface{}) {
	output := "INSERT INTO " + mapper.TableName() + " ("
	sqlVals := make([]interface{}, len(values))
	current := 0
	var holders, cols []string
	for col, val := range values {
		sqlVals[current] = val
		cols = append(cols, col)
		holders = append(holders, "?")
		current++
	}
	output += strings.Join(cols, ",") + ") VALUES (" + strings.Join(holders, ",") + ")"
	output += " RETURNING " + mapper.PrimaryKeyColumn()

	return d.Dialect.FormatQuery(output), sqlVals
}
