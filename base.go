package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

/*
Base dialect implements default implementations of a Dialect
Each of the builtin dialects have a base embedded into them,
which simplifies their implementation to database specific
functions. An example implementation using Base would look like.

  type OracleDialect struct {
    Base
  }
  // Dialects need to be circular in structure, mainly for the
  // FormatQuery to be able to be called from Base into your
  // Dialect
  func NewOracle() Dialect {
    d := &OracleDialect{}
    d.Base.Dialect = d
    return d
  }

  func (d OracleDialect) CompatibleSqlTypes(f reflect.Type) []string {
    ...
  }

  func ...

*/
type Base struct {
	Dialect Dialect
}

func (d Base) FormatQuery(query string) string {
	return query
}

func (d Base) CreateExec() bool {
	return true
}

func (d Base) Quote(s string) string {
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

func (d Base) Query(scope Scope) (string, []interface{}) {
	output := "SELECT " + scope.SelectorSql() + " FROM " + scope.TableName()
	output += scope.JoinSql()
	conditions, values := scope.ConditionSql()
	if conditions != "" {
		output += " WHERE " + conditions
	}
	output += scope.EndingSql()

	return d.Dialect.FormatQuery(output), values
}

func (d Base) Create(mapper Mapper, values map[string]interface{}) (string, []interface{}) {
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

	return d.Dialect.FormatQuery(output), sqlVals
}

func (d Base) Update(scope Scope, values map[string]interface{}) (string, []interface{}) {
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

	return d.Dialect.FormatQuery(output), append(args, sqlArgs...)
}

func (d Base) Delete(scope Scope) (string, []interface{}) {
	output := "DELETE FROM " + scope.TableName()
	conditions, sqlArgs := scope.ConditionSql()
	output += " WHERE " + conditions

	return d.Dialect.FormatQuery(output), sqlArgs
}

func (d Base) ColumnsInTable(db *sql.DB, dbName string, table string) map[string]*ColumnInfo {
	columns := make(map[string]*ColumnInfo)
	query := "SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	query = d.Dialect.FormatQuery(query)
	rows, err := db.Query(query, dbName, table)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		column := ""
		err := rows.Scan(&column)
		if err == nil {
			columns[column] = new(ColumnInfo)
		}
	}
	return columns
}

func (d Base) printArg(v interface{}) string {
	switch t := v.(type) {
	case string:
		return "'" + t + "'"
	default:
		return fmt.Sprint(v)
	}
}

func (d Base) ExpandGroupBy() bool {
	return true
}
