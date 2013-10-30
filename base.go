package db

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

func (d base) CreateExec() bool {
	return true
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

func (d base) Query(scope Scope) (string, []interface{}) {
	output := "SELECT " + scope.SelectorSql() + " FROM " + scope.TableName()
	output += scope.JoinSql()
	conditions, values := scope.ConditionSql()
	if conditions != "" {
		output += " WHERE " + conditions
	}
	output += scope.EndingSql()

	return d.dialect.FormatQuery(output), values
}

func (d base) Create(mapper Mapper, values map[string]interface{}) (string, []interface{}) {
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

	return d.dialect.FormatQuery(output), sqlVals
}

func (d base) Update(scope Scope, values map[string]interface{}) (string, []interface{}) {
	output := "UPDATE " + scope.TableName() + " SET "
	columns := make([]string, 0, len(values))
	args := make([]interface{}, 0, len(values))
	for c, v := range values {
		columns = append(columns, c)
		args = append(args, v)
	}
	output += strings.Join(columns, "= ?, ") + " = ?"
	conditions, sqlArgs := scope.ConditionSql()
	output += " WHERE " + conditions

	return d.dialect.FormatQuery(output), append(args, sqlArgs...)
}

func (d base) Delete(scope Scope) (string, []interface{}) {
	output := "DELETE FROM " + scope.TableName()
	conditions, sqlArgs := scope.ConditionSql()
	output += " WHERE " + conditions

	return d.dialect.FormatQuery(output), sqlArgs
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
	query = d.dialect.FormatQuery(query)
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
	case string:
		return "'" + t + "'"
	default:
		return fmt.Sprint(v)
	}
}
