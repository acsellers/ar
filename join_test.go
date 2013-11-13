package db

import (
	. "github.com/acsellers/assert"
	"testing"
)

func TestAliased(t *testing.T) {
	Within(t, func(test *Test) {
		sm := &sourceMapping{
			&structOptions{
				Name:     "User",
				FullName: "gf/models:User",
			},
			&ColumnInfo{},
		}

		sm2 := &sourceMapping{
			&structOptions{
				Name:     "Owner",
				FullName: "gf/models:User",
			},
			&ColumnInfo{},
		}

		test.IsTrue(sm2.Aliased())
		test.IsFalse(sm.Aliased())

	})
}
