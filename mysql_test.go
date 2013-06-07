package ar

import (
	"testing"
	"time"
)

func TestMysqlSqlType(t *testing.T) {
	Within(t, func(test *Test) {
		d := newMysql()
		test.AreEqual("boolean", d.SqlType(true, 0))
		var indirect interface{} = true
		test.AreEqual("boolean", d.SqlType(indirect, 0))
		test.AreEqual("int", d.SqlType(uint32(2), 0))
		test.AreEqual("bigint", d.SqlType(int64(1), 0))
		test.AreEqual("double", d.SqlType(1.8, 0))
		test.AreEqual("longblob", d.SqlType([]byte("asdf"), 0))
		test.AreEqual("longtext", d.SqlType("astring", 0))
		test.AreEqual("longtext", d.SqlType("a", 65536))
		test.AreEqual("varchar(128)", d.SqlType("b", 128))
		test.AreEqual("timestamp", d.SqlType(time.Now(), 0))
	})
}
