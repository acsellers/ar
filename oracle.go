package ar

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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

func (d oracleDialect) SqlType(f interface{}, size int) string {
	switch f.(type) {
	case time.Time:
		return "DATE"
	/*
		        case bool:
				return "boolean"
	*/
	case int, int8, int16, int32, uint, uint8, uint16, uint32, int64, uint64:
		if size > 0 {
			return fmt.Sprintf("NUMBER(%d)", size)
		}
		return "NUMBER"
	case float32, float64:
		if size > 0 {
			return fmt.Sprintf("NUMBER(%d,%d)", size/10, size%10)
		}
		return "NUMBER(16,2)"
	case []byte, string:
		if size > 0 && size < 4000 {
			return fmt.Sprintf("VARCHAR2(%d)", size)
		}
		return "CLOB"
	}
	panic("invalid sql type")
}

func (d oracleDialect) Insert(q *Qbs) (int64, error) {
	sql, args := d.dialect.InsertSql(q.criteria)
	row := q.QueryRow(sql, args...)
	value := q.criteria.model.pk.value
	var err error
	var id int64
	if _, ok := value.(int64); ok {
		err = row.Scan(&id)
	} else if _, ok := value.(string); ok {
		var str string
		err = row.Scan(&str)
	}
	return id, err
}

func (d oracleDialect) InsertSql(criteria *criteria) (string, []interface{}) {
	sql, values := d.base.InsertSql(criteria)
	sql += " RETURNING " + d.dialect.Quote(criteria.model.pk.name)
	return sql, values
}

func (d oracleDialect) IndexExists(db *sql.DB, dbName, tableName, indexName string) bool {
	var row *sql.Row
	var name string
	query := "SELECT INDEX_NAME FROM USER_INDEXES "
	query += "WHERE TABLE_NAME = ? AND INDEX_NAME = ?"
	query = d.SubstituteMarkers(query)
	row = db.QueryRow(query, tableName, indexName)
	row.Scan(&name)
	return name != ""
}

func (d oracleDialect) SubstituteMarkers(query string) string {
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

func (d oracleDialect) ColumnsInTable(db *sql.DB, dbName string, table interface{}) map[string]bool {
	tn := tableName(table)
	columns := make(map[string]bool)
	query := "SELECT COLUMN_NAME FROM USER_TAB_COLUMNS WHERE TABLE_NAME = ?"
	query = d.SubstituteMarkers(query)
	rows, err := db.Query(query, tn)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		column := ""
		err := rows.Scan(&column)
		if err == nil {
			columns[column] = true
		}
	}
	return columns
}

func (d oracleDialect) PrimaryKeySql(isString bool, size int) string {
	if isString {
		return fmt.Sprintf("VARCHAR2(%d) PRIMARY KEY NOT NULL", size)
	}
	if size == 0 {
		size = 16
	}
	return fmt.Sprintf("NUMBER(%d) PRIMARY KEY NOT NULL", size)
}

func (d oracleDialect) CreateTableSql(model *model, ifNotExists bool) string {
	baseSql := d.base.CreateTableSql(model, false)
	if _, isString := model.pk.value.(string); isString {
		return baseSql
	}
	table_pk := model.table + "_" + model.pk.name
	sequence := " CREATE SEQUENCE " + table_pk + "_seq" +
		" MINVALUE 1 NOMAXVALUE START WITH 1 INCREMENT BY 1 NOCACHE CYCLE"
	trigger := " CREATE TRIGGER " + table_pk + "_triger BEFORE INSERT ON " + table_pk +
		" FOR EACH ROW WHEN (new.id is null)" +
		" begin" +
		" select " + table_pk + "_seq.nextval into: new.id from dual " +
		" end "
	return baseSql + ";" + sequence + ";" + trigger
}

func (d oracleDialect) CatchMigrationError(err error) bool {
	errString := err.Error()
	return strings.Contains(errString, "ORA-00955") || strings.Contains(errString, "ORA-00942")
}

func (d oracleDialect) DropTableSql(table string) string {
	a := []string{"DROP TABLE"}
	a = append(a, d.dialect.Quote(table))
	return strings.Join(a, " ")
}
