package ar

import (
	"reflect"
	"testing"
	"time"
)

func TestMysqlCompatibleSqlTypes(t *testing.T) {
	Within(t, func(test *Test) {
		d := newMysql()
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf(true)),
				"boolean",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf(uint32(2))),
				"int",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf(int64(1))),
				"bigint",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf(1.8)),
				"double",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf([]byte("asdf"))),
				"longblob",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf("astring")),
				"longtext",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf("a")),
				"text",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf("b")),
				"varchar",
			),
		)
		test.IsTrue(
			stringMatch(
				d.CompatibleSqlTypes(reflect.TypeOf(time.Now())),
				"timestamp",
			),
		)
	})
}
