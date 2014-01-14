package db

import (
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

// The Base FormatQuery will not do anything, just returns the input string
// This works for databases that parameterize with ?, but postgres and similar
// database will need to implement a transformation for this function
func (d Base) FormatQuery(query string) string {
	return query
}

// The Base CreateExec return value is true, so INSERT statements will
// be run by sql.Exec calls.
func (d Base) CreateExec() bool {
	return true
}

// Create a basic SELECT query using ScopeInformation functions
func (d Base) Query(scope Scope) (string, []interface{}) {
	output := "SELECT " + scope.SelectorSql() + " FROM " + scope.TableName()
	output += scope.JoinsSql()
	conditions, values := scope.ConditionSql()
	if conditions != "" {
		output += " WHERE " + conditions
	}
	ending, endValues := scope.EndingSql()
	if len(endValues) > 0 {
		values = append(values, endValues...)
	}
	output += ending

	return d.Dialect.FormatQuery(output), values
}

// The Base Create function uses the syntax of INSERT INTO `table` (col...) VALUES (...)
// if this syntax will not work or you need a RETURNING predicate to get the id of
// the inserted records, you should override this
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

// The Base Update sql is of the form UPDATE table SET col = ? WHERE ...
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
	if conditions != "" {
		output += " WHERE " + conditions
	}

	return d.Dialect.FormatQuery(output), append(args, sqlArgs...)
}

// The Base Delete sql takes the form of DELETE FROM `table` WHERE ...
func (d Base) Delete(scope Scope) (string, []interface{}) {
	output := "DELETE FROM " + scope.TableName()
	conditions, sqlArgs := scope.ConditionSql()
	if conditions != "" {
		output += " WHERE " + conditions
	}

	return d.Dialect.FormatQuery(output), sqlArgs
}

// The Base ColumnsInTable will attempt to use information_schema to
// retrieve column names, it will not try to guess types for columns
// It is in your best interest to implement this per database
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
			columns[column] = &ColumnInfo{
				Name:      column,
				SqlTable:  table,
				SqlColumn: column,
				Nullable:  true,
			}
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

// The Base ExpandGroupBy will return true
func (d Base) ExpandGroupBy() bool {
	return true
}
