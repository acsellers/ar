package ar

import (
	"testing"
)

func TestInCondition(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup In Condition")
		ic := &inCondition{"test_tbl.test_col", []interface{}{1, 2, 3, 4}}
		test.AreEqual(len(ic.vals), 4, "Vals not equal")
		test.AreEqual(ic.Values(), []interface{}{[]interface{}{1, 2, 3, 4}})

		test.Section("Test In Condition SQL")
		test.F.AreEqual(
			ic.Fragment(),
			"test_tbl.test_col IN (?)",
			"Fragment not correctly generated for IN condition, expected %s, got %s",
			"test_tbl.test_col IN (?)",
			ic.Fragment(),
		)

		test.Section("Test In Condition Log")
		test.F.AreEqual(
			ic.String(),
			"test_tbl.test_col IN (1,2,3,4)",
			"String not correctly generated for IN condition, expected %s, got %s",
			"test_tbl.test_col IN (1,2,3,4)",
			ic.String(),
		)
	})
}

func TestBetweenCondition(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup Between Condition")
		bc := &betweenCondition{"test_tbl.test_col", 1, 5}
		test.AreEqual(bc.Values(), []interface{}{1, 5})

		test.Section("Test Between Condition SQL")
		test.F.AreEqual(
			bc.Fragment(),
			"test_tbl.test_col BETWEEN ? AND ?",
			"Fragment not correctly generated for BETWEEN condition, expected %s, got %s",
			"test_tbl.test_col BETWEEN ? AND ?",
			bc.Fragment(),
		)

		test.Section("Test Between Condition Log")
		test.F.AreEqual(
			bc.String(),
			"test_tbl.test_col BETWEEN 1 AND 5",
			"String not correctly generate for BETWEEN condition, expected %s, got %s",
			"test_tbl.test_col BETWEEN 1 AND 5",
			bc.String(),
		)
	})
}

func TestEqualCondition(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup Equal Condition")
		ec := &equalCondition{"test_tbl.test_col", 1}
		test.AreEqual(ec.Values(), []interface{}{1})

		test.Section("Test Equal Condition SQL")
		test.F.AreEqual(
			ec.Fragment(),
			"test_tbl.test_col = ?",
			"Fragment not correctly generated for EQUAL condition, expected %s, got %s",
			"test_tbl.test_col = ?",
			ec.Fragment(),
		)

		test.Section("Test Equal Condition Log")
		test.F.AreEqual(
			ec.String(),
			"test_tbl.test_col = 1",
			"String not correctly generate for EQUAL condition, expected %s, got %s",
			"test_tbl.test_col = 1",
			ec.String(),
		)
	})
}

func TestOrCondition(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup OR Condition")
		ec1 := &equalCondition{"test_tbl.test_col", 1}
		ec2 := &equalCondition{"test_tbl.test_col", 2}
		oc := &orCondition{[]condition{ec1, ec2}}
		test.F.AreEqual(
			oc.Values(),
			[]interface{}{1, 2},
			"Values not correctly generated for OR meta-condition, expected %v, got %v",
			[]interface{}{1, 2},
			oc.Values(),
		)

		test.Section("Test OR SQL")
		test.AreEqual(
			oc.Fragment(),
			"(test_tbl.test_col = ? OR test_tbl.test_col = ?)",
		)

		test.Section("Test OR Log")
		test.AreEqual(
			oc.String(),
			"(test_tbl.test_col = 1 OR test_tbl.test_col = 2)",
		)
	})
}

func TestAndCondition(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup AND Condition")
		ec1 := &equalCondition{"test_tbl.test_col", 1}
		ec2 := &equalCondition{"test_tbl.test_col", 2}
		ac := &andCondition{[]condition{ec1, ec2}}
		test.F.AreEqual(
			ac.Values(),
			[]interface{}{1, 2},
			"Values not correctly generated for AND meta-condition, expected %v, got %v",
			[]interface{}{1, 2},
			ac.Values(),
		)

		test.Section("Test AND SQL")
		test.AreEqual(
			ac.Fragment(),
			"(test_tbl.test_col = ? AND test_tbl.test_col = ?)",
		)

		test.Section("Test AND Log")
		test.AreEqual(
			ac.String(),
			"(test_tbl.test_col = 1 AND test_tbl.test_col = 2)",
		)
	})
}
func TestWhereCondition(t *testing.T) {
	Within(t, func(test *Test) {
	})
}