package db

import (
	. "github.com/acsellers/assert"
	"testing"
)

func TestPostMapper(t *testing.T) {
	Within(t, func(test *Test) {
		for _, conn := range availableTestConns() {
			test.Section("Setup")
			Posts := conn.m("Post")
			test.IsNotNil(Posts)
			test.AreEqual(5, len(Posts.(*source).Fields))

			test.Section("Finding All Posts")
			var posts []post
			err := Posts.RetrieveAll(&posts)
			test.NoError(err)
			test.AreEqual(len(posts), 2)
			if len(posts) == 2 {
				test.AreEqual(posts[0].Title, "First Post")
				test.AreEqual(posts[0].Views, 1)
				test.AreEqual(posts[1].Title, "Second Post")
			}

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
		}
	})
}

func TestEqualTo(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			test.IsNotNil(Posts)

			var posts []post
			Posts.EqualTo("permalink", "first_post").RetrieveAll(&posts)
			test.AreEqual(1, len(posts))
			test.AreEqual("First Post", posts[0].Title)

			var singlepost post
			Posts.EqualTo("permalink", "second_post").Retrieve(&singlepost)
			test.AreEqual(2, singlepost.Id)
			test.AreEqual("Second Post", singlepost.Title)
		}
	})
}

func TestCond(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			test.IsNotNil(Posts)

			var posts []post
			Posts.Cond("id", GTE, 1).RetrieveAll(&posts)
			test.AreEqual(2, len(posts))
			if len(posts) == 2 {
				test.AreEqual("First Post", posts[0].Title)
			}

			Posts.Cond("id", GT, 1).RetrieveAll(&posts)
			test.AreEqual(1, len(posts))
			if len(posts) == 1 {
				test.AreEqual("Second Post", posts[0].Title)
			}
		}
	})
}

func TestBetween(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			test.IsNotNil(Posts)

			var posts []post
			Posts.Between("id", 0, 2).RetrieveAll(&posts)
			test.AreEqual(2, len(posts))
			if len(posts) == 2 {
				test.AreEqual("First Post", posts[0].Title)
			}

			Posts.Between("id", 2, 10).RetrieveAll(&posts)
			test.AreEqual(1, len(posts))
			if len(posts) == 1 {
				test.AreEqual("Second Post", posts[0].Title)
			}

		}
	})
}
