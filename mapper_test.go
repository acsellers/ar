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
			test.AreEqual(7, len(Posts.(*source).Fields))
			test.AreEqual(true, Posts.(*source).hasMixin)

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

func TestWhere(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			test.IsNotNil(Posts)

			var posts []post
			Posts.Where("id BETWEEN ? AND ?", 0, 2).RetrieveAll(&posts)
			test.AreEqual(2, len(posts))
			if len(posts) == 2 {
				test.AreEqual("First Post", posts[0].Title)
			}

			Posts.Where("title LIKE ?", "%Post%").RetrieveAll(&posts)
			test.AreEqual(2, len(posts))
			if len(posts) == 2 {
				test.AreEqual("First Post", posts[0].Title)
			}

			Posts.Where("id = :id:", map[string]interface{}{"id": 1}).RetrieveAll(&posts)
			test.AreEqual(1, len(posts))
			if len(posts) == 1 {
				test.AreEqual("First Post", posts[0].Title)
			}
		}
	})
}

func TestCount(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			c, e := Posts.Count()
			test.NoError(e)
			test.AreEqual(2, c)

			c, e = Posts.EqualTo("id", 1).Count()
			test.NoError(e)
			test.AreEqual(1, c)

			c, e = Posts.EqualTo("title", "banana").Count()
			test.NoError(e)
			test.AreEqual(0, c)

		}
	})
}
func TestSave(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			newpost := post{
				Title:     "New Post",
				Permalink: "new_post",
				Body:      "LOOK AT THIS POAST",
			}

			test.NoError(Posts.SaveAll(&newpost))
			test.AreEqual(3, newpost.Id)

			var posts []post
			test.NoError(Posts.RetrieveAll(&posts))
			test.AreEqual(3, len(posts))
			test.AreEqual(1, posts[0].Id)
			test.AreEqual(2, posts[1].Id)
			test.AreEqual(3, posts[2].Id)

			newpost.Title = "Super Post"
			test.NoError(Posts.SaveAll(&newpost))
			posts = []post{}
			test.NoError(Posts.EqualTo("title", "Super Post").RetrieveAll(&posts))
			test.AreEqual(1, len(posts))
			test.AreEqual(3, posts[0].Id)

			test.NoError(Posts.EqualTo("id", 3).Delete())
			test.NoError(Posts.RetrieveAll(&posts))
			test.AreEqual(2, len(posts))
		}
	})
}

func TestPluck(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			var titles []string
			test.NoError(Posts.OrderBy("id", "ASC").Pluck("title", &titles))
			test.AreEqual(2, len(titles))
			test.AreEqual(titles[0], "First Post")
			test.AreEqual(titles[1], "Second Post")
		}
	})
}
