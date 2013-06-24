package ar

import (
	"testing"
	"time"
)

func TestParseTags(t *testing.T) {
	Within(t, func(test *Test) {
		m := parseTags(`fk:User`)
		_, ok := m["fk"]
		test.IsTrue(ok)
		m = parseTags(`notnull,default:'banana'`)
		_, ok = m["notnull"]
		test.IsTrue(ok)
		x, _ := m["default"]
		test.AreEqual("'banana'", x)
	})
}

func TestFieldOmit(t *testing.T) {

	Within(t, func(test *Test) {

		type Schema struct {
			A string `qbs:"-"`
			B string
			C string
		}
		m := structPtrToModel(&Schema{}, true, []string{"C"})
		test.AreEqual(1, len(m.fields))
	})
}

func TestInterfaceToModelWithReference(t *testing.T) {
	Within(t, func(test *Test) {
		type parent struct {
			Id    int64
			Name  string
			Value string
		}
		type table struct {
			ColPrimary int64 `qbs:"pk"`
			FatherId   int64 `qbs:"fk:Father"`
			Father     *parent
		}
		table1 := &table{
			6, 3, &parent{3, "Mrs. A", "infinite"},
		}
		m := structPtrToModel(table1, true, nil)
		ref, ok := m.refs["Father"]
		test.IsTrue(ok)
		f := ref.model.fields[1]
		x, ok := f.value.(string)
		test.IsTrue(ok)
		test.AreEqual("Mrs. A", x)
	})
}

type indexedTable struct {
	ColPrimary int64  `qbs:"pk"`
	ColNotNull string `qbs:"notnull,default:'banana'"`
	ColVarChar string `qbs:"size:64"`
	ColTime    time.Time
}

func (table *indexedTable) Indexes(indexes *Indexes) {
	indexes.Add("col_primary", "col_time")
	indexes.AddUnique("col_var_char", "col_time")
}

func TestInterfaceToModel(t *testing.T) {
	Within(t, func(test *Test) {
		now := time.Now()
		table1 := &indexedTable{
			ColPrimary: 6,
			ColVarChar: "orange",
			ColTime:    now,
		}
		m := structPtrToModel(table1, true, nil)
		test.AreEqual("col_primary", m.pk.name)
		test.AreEqual(4, len(m.fields))
		test.AreEqual(2, len(m.indexes))
		test.AreEqual("col_primary_col_time", m.indexes[0].name)
		test.IsTrue(!m.indexes[0].unique)
		test.AreEqual("col_var_char_col_time", m.indexes[1].name)
		test.IsTrue(m.indexes[1].unique)

		f := m.fields[0]
		test.AreEqual(6, f.value)
		test.IsTrue(f.pk)

		f = m.fields[1]
		test.AreEqual("'banana'", f.dfault())

		f = m.fields[2]
		str, _ := f.value.(string)
		test.AreEqual("orange", str)
		test.AreEqual(64, f.size())

		f = m.fields[3]
		tm, _ := f.value.(time.Time)
		test.AreEqual(now, tm)
	})
}

func TestInterfaceToSubModel(t *testing.T) {
	Within(t, func(test *Test) {
		type User struct {
			Id   int64
			Name string
		}
		type Post struct {
			Id       int64
			AuthorId int64 `qbs:"fk:Author"`
			Author   *User
			Content  string
		}
		pst := new(Post)
		model := structPtrToModel(pst, true, nil)
		test.AreEqual(1, len(model.refs))
	})
}

func TestColumnsAndValues(t *testing.T) {
	Within(t, func(test *Test) {
		type User struct {
			Id   int64
			Name string
		}
		user := new(User)
		model := structPtrToModel(user, true, nil)
		columns, values := model.columnsAndValues(false)
		test.AreEqual(1, len(columns))
		test.AreEqual(1, len(values))
	})
}
