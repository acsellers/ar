package ar

import (
	"testing"
	"time"
)

func TestSqlTypeForOrDialect(t *testing.T) {
	Within(t, func(test *Test) {
		d := newOracle()
		test.AreEqual("NUMBER", d.SqlType(uint32(2), 0))
		test.AreEqual("NUMBER", d.SqlType(int64(1), 0))
		test.AreEqual("NUMBER(16,2)", d.SqlType(1.8, 0))
		test.AreEqual("CLOB", d.SqlType([]byte("asdf"), 0))
		test.AreEqual("VARCHAR2(255)", d.SqlType("a", 255))
		test.AreEqual("VARCHAR2(128)", d.SqlType("b", 128))
		test.AreEqual("DATE", d.SqlType(time.Now(), 0))
	})
}
