package ar

import (
	"testing"
	"time"
)

func TestSqlTypeForOrDialect(t *testing.T) {
	Within(t, func(test *Test) {
		d := newOracle()
		test.AreEqual("NUMBER", d.sqlType(uint32(2), 0))
		test.AreEqual("NUMBER", d.sqlType(int64(1), 0))
		test.AreEqual("NUMBER(16,2)", d.sqlType(1.8, 0))
		test.AreEqual("CLOB", d.sqlType([]byte("asdf"), 0))
		test.AreEqual("VARCHAR2(255)", d.sqlType("a", 255))
		test.AreEqual("VARCHAR2(128)", d.sqlType("b", 128))
		test.AreEqual("DATE", d.sqlType(time.Now(), 0))
	})
}
