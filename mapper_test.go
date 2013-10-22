package db

import (
	. "github.com/acsellers/assert"
	"testing"
)

func TestPostMapper(t *testing.T) {
	Within(t, func(test *Test) {
		test.Section("Setup")
		conn := setupDefaultConn()
		Posts := conn.m("Post")
		test.IsNotNil(Posts)
		test.AreEqual(5, len(Posts.(*source).Fields))

		test.Section("Finding All Posts")
		var posts []post
		err := Posts.RetrieveAll(&posts)
		test.NoError(err)
		test.AreEqual(len(posts), 2)
		test.AreEqual(posts[0].Title, "First Post")
		test.AreEqual(posts[0].Views, 1)
		test.AreEqual(posts[1].Title, "Second Post")

		test.Section("Finding First Post")
		var singlePost post
		test.NoError(Posts.Find(1, &singlePost))
		test.AreEqual(singlePost.Title, "First Post")

		test.NoError(Posts.Find(2, &singlePost))
		test.AreEqual(singlePost.Title, "Second Post")
		test.AreEqual(singlePost.Views, 1)

		Posts.EqualTo("id", 2).UpdateAttribute("views", 2)
		test.NoError(Posts.Find(2, &singlePost))
		test.AreEqual(singlePost.Views, 2)
		Posts.EqualTo("id", 2).UpdateAttributes(map[string]interface{}{
			"views":     1,
			"permalink": "invalid",
		})

		test.NoError(Posts.Find(2, &singlePost))
		test.AreEqual(singlePost.Views, 1)
		test.AreEqual(singlePost.Permalink, "invalid")
		Posts.EqualTo("id", 2).UpdateAttribute("permalink", "second_post")
	})
}
