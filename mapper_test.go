package ar

import (
	"testing"
)

func TestPostMapper(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup")
		conn := setupDefaultConn()
		Posts := conn.M("Post")
		test.IsNotNil(Posts)
		test.AreEqual(5, len(Posts.source.Fields))

		test.Section("Finding All Posts")
		var posts []post
		err := Posts.RetrieveAll(&posts)
		test.IsNil(err)
		test.AreEqual(len(posts), 2)
		test.AreEqual(posts[0].Title, "First Post")
		test.AreEqual(posts[0].Views, 1)
		test.AreEqual(posts[1].Title, "Second Post")

		test.Section("Finding First Post")
		var singlePost post
		test.IsNil(Posts.Find(1).Retrieve(&singlePost))
		test.AreEqual(singlePost.Title, "First Post")

		test.IsNil(Posts.Find(2).Retrieve(&singlePost))
		test.AreEqual(singlePost.Title, "Second Post")

	})
}
