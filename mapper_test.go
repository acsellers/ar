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
		var posts []post
		err := Posts.RetrieveAll(&posts)
		test.IsNil(err)
		test.AreEqual(len(posts), 2)
		test.AreEqual(posts[0].Title, "First Post")
		test.AreEqual(posts[0].Views, 1)
		test.AreEqual(posts[1].Title, "Second Post")
	})
}
