package ar

import (
	"testing"
)

func TestPostMapper(t *testing.T) {
	Within(t, func(test *Test) {
		conn := setupDefaultConn()
		Posts := conn.M("Post")
		test.IsNotNil(Posts)
		test.AreEqual(5, len(Posts.source.Fields))
	})
}
