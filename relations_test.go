package db

import (
	. "github.com/acsellers/assert"
	"testing"
)

func TestPostRelations(t *testing.T) {
	Within(t, func(test *Test) {
		for _, conn := range availableTestConns() {
			test.AreEqual(1, len(conn.m("Post").(*source).relations))
			test.IsNotNil(conn.m("Post").(*source).relations[0].ForeignKey)
			test.AreEqual(1, len(conn.m("User").(*source).relations))
			test.IsNotNil(conn.m("User").(*source).relations[0].ForeignKey)
		}
	})
}
