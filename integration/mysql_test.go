package integration

import (
	"fmt"
	"github.com/acsellers/ar"
	_ "github.com/coocood/mysql"
	"testing"
	"time"
)

const (
	mysqlDrvformat = "%v@/%v?charset=utf8&parseTime=true&loc=Local"
	mysqlDriver    = "mysql"
)

var mysqlSyntax = dialectSyntax{
	ar.NewMysql(),
	"CREATE TABLE IF NOT EXISTS `without_pk` ( `first` longtext, `last` longtext, `amount` int )",
	"CREATE TABLE `with_pk` ( `primary` bigint PRIMARY KEY AUTO_INCREMENT, `first` longtext, `last` longtext, `amount` int )",
	"INSERT INTO `sql_gen_model` (`prim`, `first`, `last`, `amount`) VALUES (?, ?, ?, ?)",
	"UPDATE `sql_gen_model` SET `first` = ?, `last` = ?, `amount` = ? WHERE `prim` = ?",
	"DELETE FROM `sql_gen_model` WHERE `prim` = ?",
	"SELECT `post`.`id`, `post`.`author_id`, `post`.`content`, `author`.`id` AS author___id, `author`.`name` AS author___name FROM `post` LEFT JOIN `user` AS `author` ON `post`.`author_id` = `author`.`id`",
	"SELECT `name`, `grade`, `score` FROM `student` WHERE (grade IN (?, ?, ?)) AND ((score <= ?) OR (score >= ?)) ORDER BY `name`, `grade` DESC LIMIT ? OFFSET ?",
	"DROP TABLE IF EXISTS `drop_table`",
	"ALTER TABLE `a` ADD COLUMN `c` varchar(100)",
	"CREATE UNIQUE INDEX `iname` ON `itable` (`a`, `b`, `c`)",
	"CREATE INDEX `iname2` ON `itable2` (`d`, `e`)",
}

func setupMysqlDb() (*Migration, *Qbs) {
	registerMysqlTest()
	mg, _ := ar.GetMigration()
	q, _ := GetQbs()
	return mg, q
}

func registerMysqlTest() {
	ar.Register(mysqlDriver, fmt.Sprintf(mysqlDrvformat, "root:toor", testDbName), testDbName, NewMysql())
}

func TestMysqlTransaction(t *testing.T) {
	registerMysqlTest()
	doTestTransaction(t)
}

func TestMysqlSaveAndDelete(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestSaveAndDelete(t, mg, q)
}

func TestMysqlSaveAgain(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestSaveAgain(t, mg, q)
}

func TestMysqlForeignKey(t *testing.T) {
	registerMysqlTest()
	doTestForeignKey(t)
}

func TestMysqlFind(t *testing.T) {
	registerMysqlTest()
	doTestFind(t)
}

func TestMysqlCreateTable(t *testing.T) {
	mg, _ := setupMysqlDb()
	doTestCreateTable(t, mg)
}

func TestMysqlUpdate(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestUpdate(t, mg, q)
}

func TestMysqlValidation(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestValidation(t, mg, q)
}

func TestMysqlBoolType(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestBoolType(t, mg, q)
}

func TestMysqlStringPk(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestStringPk(t, mg, q)
}

func TestMysqlCount(t *testing.T) {
	registerMysqlTest()
	doTestCount(t)
}

func TestMysqlQueryMap(t *testing.T) {
	mg, q := setupMysqlDb()
	doTestQueryMap(t, mg, q)
}

func TestMysqlBulkInsert(t *testing.T) {
	registerMysqlTest()
	doTestBulkInsert(t)
}

func TestMysqlQueryStruct(t *testing.T) {
	registerMysqlTest()
	doTestQueryStruct(t)
}

func TestMysqlCustomNameConvertion(t *testing.T) {
	registerMysqlTest()
	ColumnNameToFieldName = noConvert
	FieldNameToColumnName = noConvert
	TableNameToStructName = noConvert
	StructNameToTableName = noConvert
	doTestForeignKey(t)
	ColumnNameToFieldName = snakeToUpperCamel
	FieldNameToColumnName = toSnake
	TableNameToStructName = snakeToUpperCamel
	StructNameToTableName = toSnake
}

func TestMysqlConnectionLimit(t *testing.T) {
	registerMysqlTest()
	doTestConnectionLimit(t)
}

func TestMysqlIterate(t *testing.T) {
	registerMysqlTest()
	doTestIterate(t)
}

func TestMysqlAddColumnSQL(t *testing.T) {
	doTestAddColumSQL(t, mysqlSyntax)
}

func TestMysqlCreateTableSQL(t *testing.T) {
	doTestCreateTableSQL(t, mysqlSyntax)
}

func TestMysqlCreateIndexSQL(t *testing.T) {
	doTestCreateIndexSQL(t, mysqlSyntax)
}

func TestMysqlInsertSQL(t *testing.T) {
	doTestInsertSQL(t, mysqlSyntax)
}

func TestMysqlUpdateSQL(t *testing.T) {
	doTestUpdateSQL(t, mysqlSyntax)
}

func TestMysqlDeleteSQL(t *testing.T) {
	doTestDeleteSQL(t, mysqlSyntax)
}

func TestMysqlSelectionSQL(t *testing.T) {
	doTestSelectionSQL(t, mysqlSyntax)
}

func TestMysqlQuerySQL(t *testing.T) {
	doTestQuerySQL(t, mysqlSyntax)
}
func TestMysqlDropTableSQL(t *testing.T) {
	doTestDropTableSQL(t, mysqlSyntax)
}

func BenchmarkMysqlFind(b *testing.B) {
	registerMysqlTest()
	doBenchmarkFind(b)
}
func BenchmarkMysqlDbQuery(b *testing.B) {
	registerMysqlTest()
	doBenchmarkDbQuery(b)
}

func BenchmarkMysqlStmtQuery(b *testing.B) {
	registerMysqlTest()
	doBenchmarkStmtQuery(b)
}

func BenchmarkMysqlTransaction(b *testing.B) {
	registerMysqlTest()
	doBenchmarkTransaction(b)
}
