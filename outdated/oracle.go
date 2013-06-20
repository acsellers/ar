package ar

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type oracleDialect struct {
	base
}

func newOracle() Dialect {
	d := &oracleDialect{}
	d.base.dialect = d
	return d
}

func (d oracleDialect) Quote(s string) string {
	sep := "."
	a := []string{}
	c := strings.Split(s, sep)
	for _, v := range c {
		a = append(a, fmt.Sprintf(`"%s"`, v))
	}
	return strings.Join(a, sep)
}

func (d oracleDialect) CompatibleSqlTypes(f reflect.Type) []string {
	return []string{d.sqlType(f, 250)}
}
func (d oracleDialect) sqlType(f reflect.Type, size int) string {
	switch f.Kind() {
	case reflect.Struct:
		if f.String() == "time.Time" {
			return "DATE"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		if size > 0 {
			return fmt.Sprintf("NUMBER(%d)", size)
		}
		return "NUMBER"
	case reflect.Float32, reflect.Float64:
		if size > 0 {
			return fmt.Sprintf("NUMBER(%d,%d)", size/10, size%10)
		}
		return "NUMBER(16,2)"
	case reflect.String:
		if size > 0 && size < 4000 {
			return fmt.Sprintf("VARCHAR2(%d)", size)
		}

	case reflect.Slice:
		if f.String() == "[]uint8" { //[]byte
			if size > 0 && size < 4000 {
				return fmt.Sprintf("VARCHAR2(%d)", size)
			}
			return "CLOB"
		}
	}
	panic("invalid sql type")
}

/*
func (d oracleDialect) InsertSql(criteria *criteria) (string, []interface{}) {
	sql, values := d.base.InsertSql(criteria)
	sql += " RETURNING " + d.dialect.Quote(criteria.model.pk.name)
	return sql, values
}
*/
func (d oracleDialect) FormatQuery(query string) string {
	position := 1
	chunks := make([]string, 0, len(query)*2)
	for _, v := range query {
		if v == '?' {
			chunks = append(chunks, fmt.Sprintf("$%d", position))
			position++
		} else {
			chunks = append(chunks, string(v))
		}
	}
	return strings.Join(chunks, "")
}

func (d oracleDialect) ColumnsInTable(db *sql.DB, dbName string, table interface{}) map[string]*columnInfo {
	tn := tableName(table)
	columns := make(map[string]*columnInfo)
	query := "SELECT COLUMN_NAME FROM USER_TAB_COLUMNS WHERE TABLE_NAME = ?"
	query = d.FormatQuery(query)
	rows, err := db.Query(query, tn)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		column := ""
		err := rows.Scan(&column)
		if err == nil {
			columns[column] = new(columnInfo)
		}
	}
	return columns
}
