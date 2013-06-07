package ar

import (
	"testing"
	"time"
)

func TestMysqlSqlType(t *testing.T) {
	Within(t, func(test *Test) {
		d := newMysql()
		test.AreEqual("boolean", d.sqlType(true, 0))
		var indirect interface{} = true
		test.AreEqual("boolean", d.sqlType(indirect, 0))
		test.AreEqual("int", d.sqlType(uint32(2), 0))
		test.AreEqual("bigint", d.sqlType(int64(1), 0))
		test.AreEqual("double", d.sqlType(1.8, 0))
		test.AreEqual("longblob", d.sqlType([]byte("asdf"), 0))
		test.AreEqual("longtext", d.sqlType("astring", 0))
		test.AreEqual("longtext", d.sqlType("a", 65536))
		test.AreEqual("varchar(128)", d.sqlType("b", 128))
		test.AreEqual("timestamp", d.sqlType(time.Now(), 0))
	})
}
