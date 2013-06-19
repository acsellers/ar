
type MigrateDialect interface {
	// Migration Related
	PrimaryKeySql(isString bool, size int) string
	IndexExists(db *sql.DB, dbName, tableName, indexName string) bool
	CreateTableSql(model *model, ifNotExists bool) string
	DropTableSql(table string) string
	AddColumnSql(table, column string, typ interface{}, size int) string
	CreateIndexSql(name, table string, unique bool, columns ...string) string
}

type base struct {
	mg MigrateDialect
}

func (d base) CreateTableSql(model *model, ifNotExists bool) string {
	a := []string{"CREATE TABLE "}
	if ifNotExists {
		a = append(a, "IF NOT EXISTS ")
	}
	a = append(a, d.dialect.Quote(model.table), " ( ")
	for i, field := range model.fields {
		b := []string{
			d.dialect.Quote(field.name),
		}
		if field.pk {
			_, ok := field.value.(string)
			b = append(b, d.dialect.PrimaryKeySql(ok, field.size()))
		} else {
			b = append(b, d.dialect.CompatibleSQLTypes(field.value, field.size())[0])
			if field.notNull() {
				b = append(b, "NOT NULL")
			}
			if x := field.dfault(); x != "" {
				b = append(b, "DEFAULT "+x)
			}
		}
		a = append(a, strings.Join(b, " "))
		if i < len(model.fields)-1 {
			a = append(a, ", ")
		}
	}
	for _, v := range model.refs {
		if v.foreignKey {
			a = append(a, ", FOREIGN KEY (", d.dialect.Quote(v.refKey), ") REFERENCES ")
			a = append(a, d.dialect.Quote(v.model.table), " (", d.dialect.Quote(v.model.pk.name), ") ON DELETE CASCADE")
		}
	}
	a = append(a, " )")
	return strings.Join(a, "")
}

func (d base) DropTableSql(table string) string {
	a := []string{"DROP TABLE IF EXISTS"}
	a = append(a, d.dialect.Quote(table))
	return strings.Join(a, " ")
}

func (d base) AddColumnSql(table, column string, typ interface{}, size int) string {
	return fmt.Sprintf(
		"ALTER TABLE %v ADD COLUMN %v %v",
		d.dialect.Quote(table),
		d.dialect.Quote(column),
		d.dialect.SqlType(typ, size),
	)
}

func (d base) CreateIndexSql(name, table string, unique bool, columns ...string) string {
	a := []string{"CREATE"}
	if unique {
		a = append(a, "UNIQUE")
	}
	quotedColumns := make([]string, 0, len(columns))
	for _, c := range columns {
		quotedColumns = append(quotedColumns, d.dialect.Quote(c))
	}
	a = append(a, fmt.Sprintf(
		"INDEX %v ON %v (%v)",
		d.dialect.Quote(name),
		d.dialect.Quote(table),
		strings.Join(quotedColumns, ", "),
	))
	return strings.Join(a, " ")
}

type oracleDialect struct {
	base
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

func (d oracleDialect) DropTableSql(table string) string {
	a := []string{"DROP TABLE"}
	a = append(a, d.dialect.Quote(table))
	return strings.Join(a, " ")
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



