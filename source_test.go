package db

import (
	"github.com/acsellers/assert"
	"testing"
)

func TestSourceSelectColumns(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		s := &source{
			SqlName: "users",
			Fields: []*sourceMapping{
				&sourceMapping{&structOptions{}, &ColumnInfo{SqlColumn: "id"}},
				&sourceMapping{&structOptions{}, &ColumnInfo{SqlColumn: "name"}},
				&sourceMapping{&structOptions{}, &ColumnInfo{SqlColumn: "email"}},
				&sourceMapping{&structOptions{}, &ColumnInfo{SqlColumn: "password"}},
			},
		}

		test.AreEqual(
			[]string{"users.id", "users.name", "users.email", "users.password"},
			s.selectColumns(),
		)
	})
}
