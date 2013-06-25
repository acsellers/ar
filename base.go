package ar

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type base struct {
	dialect Dialect
}

func (d base) FormatQuery(query string) string {
	return query
}

func (d base) Quote(s string) string {
	segs := strings.Split(s, ".")
	buf := new(bytes.Buffer)
	buf.WriteByte('`')
	buf.WriteString(segs[0])
	for i := 1; i < len(segs); i++ {
		buf.WriteString("`.`")
		buf.WriteString(segs[i])
	}
	buf.WriteByte('`')
	return buf.String()
}

func (d base) Query(queryable *Queryable) (string, []interface{}) {
	output := "SELECT " + queryable.selectorSql() + " FROM " + queryable.source.SqlName
	output += queryable.joinSql()
	conditions, values := queryable.conditionSql()
	if conditions != "" {
		output += " WHERE " + conditions
	}
	output += queryable.endingSql()

	return output, values
}

func (d base) Update(queryable *Queryable, values map[string]interface{}) (string, []interface{}) {
	output := "UPDATE " + queryable.source.SqlName + " SET "
	columns := make([]string, 0, len(values))
	args := make([]interface{}, 0, len(values))
	for c, v := range values {
		columns = append(columns, c)
		args = append(args, v)
	}
	output += strings.Join(columns, "= ?, ") + " = ?"
	conditions, sqlArgs := queryable.conditionSql()
	output += " WHERE " + conditions

	return output, append(args, sqlArgs...)
}

/*

func (d base) InsertSql(queryable *Queryable) (string, []interface{}) {
	columns, values := criteria.model.columnsAndValues(false)
	quotedColumns := make([]string, 0, len(columns))
	markers := make([]string, 0, len(columns))
	for _, c := range columns {
		quotedColumns = append(quotedColumns, d.dialect.Quote(c))
		markers = append(markers, "?")
	}
	sql := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		d.dialect.Quote(criteria.model.table),
		strings.Join(quotedColumns, ", "),
		strings.Join(markers, ", "),
	)
	return sql, values
}


func (d base) UpdateSql(queryable *Queryable) (string, []interface{}) {
	columns, values := criteria.model.columnsAndValues(true)
	pairs := make([]string, 0, len(columns))
	for _, column := range columns {
		pairs = append(pairs, fmt.Sprintf("%v = ?", d.dialect.Quote(column)))
	}
	conditionSql, args := queryable.conditionSql()
	sql := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v",
		d.dialect.Quote(criteria.model.table),
		strings.Join(pairs, ", "),
		conditionSql,
	)
	values = append(values, args...)
	return sql, values
}
func (d base) DeleteSql(queryable *Queryable) (string, []interface{}) {
	conditionSql, args := queryable.conditionSql()
	sql := "DELETE FROM " + d.dialect.Quote(criteria.model.table) + " WHERE " + conditionSql
	return sql, args
}
*/
func (d base) ColumnsInTable(db *sql.DB, dbName string, table string) map[string]*columnInfo {
	columns := make(map[string]*columnInfo)
	query := "SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	query = d.FormatQuery(query)
	rows, err := db.Query(query, dbName, table)
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

func (d base) printArg(v interface{}) string {
	switch t := v.(type) {
	case Formula:
		return string(t)
	case string:
		return "'" + t + "'"
	default:
		return fmt.Sprint(v)
	}
}
