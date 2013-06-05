package ar

import (
	"database/sql"
	"errors"
	"testing"
	"time"
)

var testDbName = "qbs_test"
var testDbUser = "qbs_test"

type addColumn struct {
	Prim   int64  `qbs:"pk"`
	First  string `qbs:"size:64,notnull"`
	Last   string `qbs:"size:128,default:'defaultValue'"`
	Amount int
}

func (table *addColumn) Indexes(indexes *Indexes) {
	indexes.AddUnique("first", "last")
}

func doTestTransaction(t *testing.T) {
	Within(t, func(test *Test) {
		type txModel struct {
			Id int64
			A  string
		}
		table := txModel{
			A: "A",
		}
		err := WithMigration(func(mg *Migration) error {
			mg.dropTableIfExists(&table)
			mg.CreateTableIfNotExists(&table)
			return nil
		})
		test.IsNil(err)

		WithQbs(func(q *Qbs) error {
			q.Begin()
			test.IsNotNil(q.tx)
			_, err := q.Save(&table)
			test.IsNil(err)
			err = q.Rollback()
			test.IsNil(err)
			out := new(txModel)
			err = q.Find(out)
			test.AreEqual(sql.ErrNoRows, err)
			q.Begin()
			table.Id = 0
			_, err = q.Save(&table)
			test.IsNil(err)
			err = q.Commit()
			test.IsNil(err)
			out.Id = table.Id
			err = q.Find(out)
			test.IsNil(err)
			test.AreEqual("A", out.A)
			return nil
		})
	})

}

func doTestSaveAndDelete(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		x := time.Now()
		test.AreEqual(0, x.Sub(x.UTC()))
		now := time.Now()
		type saveModel struct {
			Id      int64
			A       string
			B       int
			Updated time.Time
			Created time.Time
		}
		model1 := saveModel{
			A: "banana",
			B: 5,
		}
		model2 := saveModel{
			A: "orange",
			B: 4,
		}

		mg.dropTableIfExists(&model1)
		mg.CreateTableIfNotExists(&model1)

		affected, err := q.Save(&model1)
		test.IsNil(err)
		test.AreEqual(1, affected)
		test.IsTrue(model1.Created.Sub(now) > 0)
		test.IsTrue(model1.Updated.Sub(now) > 0)
		// make sure created/updated values match the db
		var model1r []*saveModel
		err = q.WhereEqual("id", model1.Id).FindAll(&model1r)

		test.IsNil(err)
		test.AreEqual(1, len(model1r))
		test.AreEqual(model1.Created.Unix(), model1r[0].Created.Unix())
		test.AreEqual(model1.Updated.Unix(), model1r[0].Updated.Unix())

		oldCreate := model1.Created
		oldUpdate := model1.Updated
		model1.A = "grape"
		model1.B = 9

		time.Sleep(time.Second * 1) // sleep for 1 sec
		affected, err = q.Save(&model1)
		test.IsNil(err)
		test.AreEqual(1, affected)
		test.IsTrue(model1.Created.Equal(oldCreate))
		test.IsTrue(model1.Updated.Sub(oldUpdate) > 0)

		// make sure created/updated values match the db
		var model1r2 []*saveModel
		err = q.Where("id = ?", model1.Id).FindAll(&model1r2)
		test.IsNil(err)
		test.AreEqual(1, len(model1r2))
		test.IsTrue(model1r2[0].Updated.Sub(model1r2[0].Created) >= 1)
		test.AreEqual(model1.Created.Unix(), model1r2[0].Created.Unix())
		test.AreEqual(model1.Updated.Unix(), model1r2[0].Updated.Unix())

		affected, err = q.Save(&model2)
		test.IsNil(err)
		test.AreEqual(1, affected)

		affected, err = q.Delete(&model2)
		test.IsNil(err)
		test.AreEqual(1, affected)
	})
}

func doTestSaveAgain(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		b := new(basic)
		mg.dropTableIfExists(b)
		mg.CreateTableIfNotExists(b)
		b.Name = "a"
		b.State = 2
		affected, err := q.Save(b)
		test.IsNil(err)
		test.AreEqual(1, affected)
		affected, err = q.Save(b)
		test.IsNil(err)
		if _, ok := q.Dialect.(*mysql); ok {
			test.AreEqual(0, affected)
		} else {
			test.AreEqual(1, affected)
		}
	})
}

func doTestForeignKey(t *testing.T) {
	Within(t, func(test *Test) {
		type User struct {
			Id   int64
			Name string
		}
		type Post struct {
			Id       int64
			Title    string
			AuthorId int64
			Author   *User
		}
		aUser := &User{
			Name: "john",
		}
		aPost := &Post{
			Title: "A Title",
		}
		WithMigration(func(mg *Migration) error {
			mg.dropTableIfExists(aPost)
			mg.dropTableIfExists(aUser)
			mg.CreateTableIfNotExists(aUser)
			mg.CreateTableIfNotExists(aPost)
			return nil
		})
		WithQbs(func(q *Qbs) error {
			affected, err := q.Save(aUser)
			test.IsNil(err)
			aPost.AuthorId = int64(aUser.Id)
			affected, err = q.Save(aPost)
			test.AreEqual(1, affected)
			pst := new(Post)
			pst.Id = aPost.Id
			err = q.Find(pst)
			test.IsNil(err)
			test.AreEqual(aPost.Id, pst.Id)
			test.AreEqual("john", pst.Author.Name)

			pst.Author = nil
			err = q.OmitFields("Author").Find(pst)
			test.IsNil(err)
			test.IsNil(pst.Author)

			err = q.OmitJoin().Find(pst)
			test.IsNil(err)
			test.IsNil(pst.Author)

			var psts []*Post
			err = q.FindAll(&psts)
			test.IsNil(err)
			test.AreEqual(1, len(psts))
			test.AreEqual("john", psts[0].Author.Name)
			return nil
		})
	})
}

func doTestFind(t *testing.T) {
	Within(t, func(test *Test) {
		now := time.Now()
		type types struct {
			Id    int64
			Str   string
			Intgr int64
			Flt   float64
			Bytes []byte
			Time  time.Time
		}
		modelData := &types{
			Str:   "string!",
			Intgr: -1,
			Flt:   3.8,
			Bytes: []byte("bytes!"),
			Time:  now,
		}
		WithMigration(func(mg *Migration) error {
			mg.dropTableIfExists(modelData)
			mg.CreateTableIfNotExists(modelData)
			return nil
		})
		WithQbs(func(q *Qbs) error {
			out := new(types)
			condition := NewCondition("str = ?", "string!").And("intgr = ?", -1)
			err := q.Condition(condition).Find(out)
			test.AreEqual(sql.ErrNoRows, err)

			affected, err := q.Save(modelData)
			test.IsNil(err)
			test.AreEqual(1, affected)
			out.Id = modelData.Id
			err = q.Condition(condition).Find(out)
			test.IsNil(err)
			test.AreEqual(1, out.Id)
			test.AreEqual("string!", out.Str)
			test.AreEqual(-1, out.Intgr)
			test.AreEqual(3.8, out.Flt)
			test.AreEqual([]byte("bytes!"), out.Bytes)
			diff := now.Sub(out.Time)
			test.IsTrue(diff < time.Second && diff > -time.Second)

			modelData.Id = 5
			modelData.Str = "New row"
			_, err = q.Save(modelData)
			test.IsNil(err)

			out = new(types)
			condition = NewCondition("str = ?", "New row").And("flt = ?", 3.8)
			err = q.Condition(condition).Find(out)
			test.IsNil(err)
			test.AreEqual(5, out.Id)

			out = new(types)
			out.Id = 100
			err = q.Find(out)
			test.IsNotNil(err)

			allOut := []*types{}
			err = q.WhereEqual("intgr", -1).FindAll(&allOut)
			test.IsNil(err)
			test.AreEqual(2, len(allOut))
			return nil
		})
	})
}

func doTestCreateTable(t *testing.T, mg *Migration) {
	Within(t, func(test *Test) {
		defer mg.Close()
		{
			type AddColumn struct {
				Prim int64 `qbs:"pk"`
			}
			table := &AddColumn{}
			mg.dropTableIfExists(table)
			mg.CreateTableIfNotExists(table)
			columns := mg.dialect.columnsInTable(mg, table)
			test.AreEqual(1, len(columns))
			test.IsTrue(columns["prim"])
		}
		table := &addColumn{}
		mg.CreateTableIfNotExists(table)
		test.IsTrue(mg.dialect.indexExists(mg, "add_column", "add_column_first_last"))
		columns := mg.dialect.columnsInTable(mg, table)
		test.AreEqual(4, len(columns))
	})
}

type basic struct {
	Id    int64
	Name  string `qbs:"size:64"`
	State int64
}

func doTestUpdate(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		mg.dropTableIfExists(&basic{})
		mg.CreateTableIfNotExists(&basic{})
		_, err := q.Save(&basic{Name: "a", State: 1})
		_, err = q.Save(&basic{Name: "b", State: 1})
		_, err = q.Save(&basic{Name: "c", State: 0})
		test.IsNil(err)
		{
			// define a temporary struct in a block to update partial columns of a table
			// as the type is in a block, so it will not conflict with other types with the same name in the same method
			type basic struct {
				Name string
			}
			affected, err := q.WhereEqual("state", 1).Update(&basic{Name: "d"})
			test.IsNil(err)
			test.AreEqual(2, affected)

			var datas []*basic
			q.WhereEqual("state", 1).FindAll(&datas)
			test.AreEqual(2, len(datas))
			test.AreEqual("d", datas[0].Name)
			test.AreEqual("d", datas[1].Name)
		}

		// if choose basic table type to update, all zero value in the struct will be updated too.
		// this may be cause problems, so define a temporary struct to update table is the recommended way.
		affected, err := q.Where("state = ?", 1).Update(&basic{Name: "e"})
		test.IsNil(err)
		test.AreEqual(2, affected)
		var datas []*basic
		q.WhereEqual("state", 1).FindAll(&datas)
		test.AreEqual(0, len(datas))
	})
}

type validatorTable struct {
	Id   int64
	Name string
}

func (v *validatorTable) Validate(q *Qbs) error {
	if q.ContainsValue(v, "name", v.Name) {
		return errors.New("name already taken")
	}
	return nil
}

func doTestValidation(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		valid := new(validatorTable)
		mg.dropTableIfExists(valid)
		mg.CreateTableIfNotExists(valid)
		valid.Name = "ok"
		q.Save(valid)
		valid.Id = 0
		_, err := q.Save(valid)
		test.IsNotNil(err)
		test.AreEqual("name already taken", err.Error())
	})
}

func doTestBoolType(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		type BoolType struct {
			Id     int64
			Active bool
		}
		bt := new(BoolType)
		mg.dropTableIfExists(bt)
		mg.CreateTableIfNotExists(bt)
		bt.Active = true
		q.Save(bt)
		bt.Active = false
		q.WhereEqual("active", true).Find(bt)
		test.IsTrue(bt.Active)
	})
}

func doTestStringPk(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		type StringPk struct {
			Tag   string `qbs:"pk,size:16"`
			Count int32
		}
		spk := new(StringPk)
		spk.Tag = "health"
		spk.Count = 10
		mg.dropTableIfExists(spk)
		mg.CreateTableIfNotExists(spk)
		affected, _ := q.Save(spk)
		test.AreEqual(1, affected)
		spk.Count = 0
		q.Find(spk)
		test.AreEqual(10, spk.Count)
	})
}

func doTestCount(t *testing.T) {
	Within(t, func(test *Test) {
		setupBasicDb()
		WithQbs(func(q *Qbs) error {
			basic := new(basic)
			basic.Name = "name"
			basic.State = 1
			q.Save(basic)
			for i := 0; i < 5; i++ {
				basic.Id = 0
				basic.State = 2
				q.Save(basic)
			}
			count1 := q.Count("basic")
			test.AreEqual(6, count1)
			count2 := q.WhereEqual("state", 2).Count(basic)
			test.AreEqual(5, count2)
			return nil
		})
	})
}

func doTestQueryMap(t *testing.T, mg *Migration, q *Qbs) {
	Within(t, func(test *Test) {
		defer closeMigrationAndQbs(mg, q)
		type types struct {
			Id      int64
			Name    string `qbs:"size:64"`
			Created time.Time
		}
		tp := new(types)
		mg.dropTableIfExists(tp)
		mg.CreateTableIfNotExists(tp)
		result, err := q.QueryMap("SELECT * FROM types")
		test.IsNil(result)
		test.AreEqual(sql.ErrNoRows, err)
		for i := 0; i < 3; i++ {
			tp.Id = 0
			tp.Name = "abc"
			q.Save(tp)
		}
		result, err = q.QueryMap("SELECT * FROM types")
		test.IsNotNil(result)
		test.AreEqual(1, result["id"])
		test.AreEqual("abc", result["name"])
		/*
			if _, sql3 := q.Dialect.(*sqlite3); !sql3 {
				_, ok := result["created"].(time.Time)
				test.IsTrue(ok)
			} else {
				_, ok := result["created"].(string)
				test.IsTrue(ok)
			}
		*/
		results, err := q.QueryMapSlice("SELECT * FROM types")
		test.AreEqual(3, len(results))
	})
}

func doTestBulkInsert(t *testing.T) {
	Within(t, func(test *Test) {
		setupBasicDb()
		WithQbs(func(q *Qbs) error {
			var bulk []*basic
			for i := 0; i < 10; i++ {
				b := new(basic)
				b.Name = "basic"
				b.State = int64(i)
				bulk = append(bulk, b)
			}
			err := q.BulkInsert(bulk)
			test.IsNil(err)
			for i := 0; i < 10; i++ {
				test.AreEqual(i+1, bulk[i].Id)
			}
			return nil
		})
	})
}

func doTestQueryStruct(t *testing.T) {
	Within(t, func(test *Test) {
		setupBasicDb()
		WithQbs(func(q *Qbs) error {
			b := new(basic)
			b.Name = "abc"
			b.State = 2
			q.Save(b)
			b = new(basic)
			err := q.QueryStruct(b, "SELECT * FROM basic")
			test.IsNil(err)
			test.AreEqual(1, b.Id)
			test.AreEqual("abc", b.Name)
			test.AreEqual(2, b.State)
			var slice []*basic
			q.QueryStruct(&slice, "SELECT * FROM basic")
			test.AreEqual(1, len(slice))
			test.AreEqual("abc", slice[0].Name)
			return nil
		})
	})
}

func doTestConnectionLimit(t *testing.T) {
	Within(t, func(test *Test) {
		SetConnectionLimit(2, false)
		q0, _ := GetQbs()
		GetQbs()
		GetQbs()
		_, err := GetQbs()
		test.AreEqual(ConnectionLimitError, err)
		q0.Close()
		q4, _ := GetQbs()
		test.IsNotNil(q4)
		SetConnectionLimit(0, true)
		a := 0
		go func() {
			a = 1
			q4.Close()
		}()
		GetQbs()
		test.AreEqual(1, a)
		SetConnectionLimit(-1, false)
		test.IsNil(connectionLimit)
	})
}

func doTestIterate(t *testing.T) {
	Within(t, func(test *Test) {
		setupBasicDb()
		q, _ := GetQbs()
		for i := 0; i < 4; i++ {
			b := new(basic)
			b.State = int64(i)
			q.Save(b)
		}
		var stateSum int64
		b := new(basic)
		err := q.Iterate(b, func() error {
			if b.State == 3 {
				return errors.New("A error")
			}
			stateSum += b.State
			return nil
		})
		test.AreEqual("A error", err.Error())
		test.AreEqual(3, stateSum)
	})
}

func setupBasicDb() {
	WithMigration(func(mg *Migration) error {
		b := new(basic)
		mg.dropTableIfExists(b)
		mg.CreateTableIfNotExists(b)
		return nil
	})
}

func closeMigrationAndQbs(mg *Migration, q *Qbs) {
	mg.Close()
	q.Close()
}

func noConvert(s string) string {
	return s
}
