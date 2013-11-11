package db

import (
	. "github.com/acsellers/assert"
	"testing"
)

func TestPostRelations(t *testing.T) {
	Within(t, func(test *Test) {
		for _, conn := range availableTestConns() {
			Posts := conn.m("Post")
			Users := conn.m("User")
			test.AreEqual(1, len(Posts.(*source).relations))
			test.IsNotNil(Posts.(*source).relations[0].ForeignKey)
			test.AreEqual(1, len(Users.(*source).relations))
			test.IsNotNil(Users.(*source).relations[0].ForeignKey)
			c, e := Users.Count()
			test.NoError(e)
			test.AreEqual(1, c)
		}
	})
}
