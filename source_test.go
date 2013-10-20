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
				&sourceMapping{&structOptions{}, &columnInfo{SqlColumn: "id"}},
				&sourceMapping{&structOptions{}, &columnInfo{SqlColumn: "name"}},
				&sourceMapping{&structOptions{}, &columnInfo{SqlColumn: "email"}},
				&sourceMapping{&structOptions{}, &columnInfo{SqlColumn: "password"}},
			},
		}

		test.AreEqual(
			[]string{"users.id", "users.name", "users.email", "users.password"},
			s.selectColumns(),
		)
	})
}
