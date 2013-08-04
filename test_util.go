package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func stringMatch(items []string, wanted string) bool {
	for _, item := range items {
		if item == wanted {
			return true
		}
	}

	return false
}

func setupDefaultConn() *Connection {
	db, err := sql.Open("mysql", "root:toor@/ar_test?charset=utf8")
	if err != nil {
		panic(err)
	}
	conn := new(Connection)
	conn.DB = db
	conn.Dialect = newMysql()
	conn.Config = NewSimpleConfig()
	conn.sources = make(map[string]*source)
	conn.CreateMapper("Post", &post{})

	return conn
}

type post struct {
	ID        int
	Title     string
	Permalink string
	Body      string
	Views     int
}
