package db

import (
	. "github.com/acsellers/assert"
	"testing"
)

func TestMixin(t *testing.T) {
	Within(t, func(test *Test) {
		for _, c := range availableTestConns() {
			Posts := c.m("Post")
			newpost := post{
				Title:     "New Post",
				Permalink: "new_post",
				Body:      "LOOK AT THIS POAST",
			}
			test.IsNil(newpost.Mixin)

			test.NoError(newpost.InitWithConn(c, &newpost))
			test.IsNotNil(newpost.Mixin)

			newpost.Mixin = nil
			Posts.Initialize(&newpost)
			test.IsNotNil(newpost.Mixin)

			c, e := Posts.Count()
			test.NoError(e)
			test.AreEqual(2, c)

			test.NoError(newpost.Save())
			c, e = Posts.Count()
			test.NoError(e)
			test.AreEqual(3, c)

			test.NoError(newpost.Delete())

			c, e = Posts.Count()
			test.NoError(e)
			test.AreEqual(2, c)

			var firstPost, secondPost post
			test.NoError(Posts.Find(1, &firstPost))
			test.NoError(Posts.Find(2, &secondPost))
			test.IsNotNil(firstPost.Mixin)

			var posts []post
			test.NoError(Posts.RetrieveAll(&posts))
			for _, p := range posts {
				test.IsNotNil(p.Mixin, "RetrieveAll: Mixin should not be nil")
			}

			c, e = Posts.EqualTo("title", "Temp Name").Count()
			test.AreEqual(0, c)
			test.NoError(e)
			test.NoError(firstPost.UpdateAttribute("title", "Temp Name"))
			c, e = Posts.EqualTo("title", "Temp Name").Count()
			test.AreEqual(1, c)
			test.NoError(e)

			test.NoError(firstPost.UpdateAttributes(Attributes{"title": "First Post"}))
			c, e = Posts.EqualTo("title", "Temp Name").Count()
			test.AreEqual(0, c)
			test.NoError(e)
			c, e = Posts.EqualTo("title", "First Post").Count()
			test.AreEqual(1, c)
			test.NoError(e)

			test.IsFalse(firstPost.IsNull("user_id"))
			test.IsTrue(secondPost.IsNull("user_id"))
		}
	})
}
