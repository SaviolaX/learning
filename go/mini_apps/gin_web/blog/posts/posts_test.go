package posts

import (
	"testing"

	"github.com/quii/go-graceful-shutdown/assert"
)

func TestGetAll(t *testing.T) {
	t.Run("not empty list of posts", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})
		testDB.Posts = append(testDB.Posts, Post{
			Id:     2,
			Author: "testUser1",
			Title:  "testTitle2",
			Body:   "testBody2",
		})

		posts := testDB.GetAll()

		assert.Equal(t, len(posts), 2)

	})
	t.Run("empty list of posts", func(t *testing.T) {
		testDB := InMemoryDB{}

		posts := testDB.GetAll()

		assert.Equal(t, len(posts), 0)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("get post with correct ID", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})

		correctID := 1
		want := Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		}

		postById, _ := testDB.GetById(correctID)

		assert.Equal(t, postById, want)
	})
	t.Run("incorrect post ID", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})

		incorrectID := 0000
		want := "Post not found"

		_, err := testDB.GetById(incorrectID)

		assert.Equal(t, err.Error(), want)
	})
}

func TestAdd(t *testing.T) {

	testDB := InMemoryDB{}

	post := Post{
		Id:     1,
		Author: "testUser1",
		Title:  "testTitle1",
		Body:   "testBody1",
	}

	wantPostsLen := 1

	testDB.Add(post)

	assert.Equal(t, len(testDB.Posts), wantPostsLen)
}

func TestDelete(t *testing.T) {
	t.Run("post deleted", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})

		correctID := 1
		wantPostsLen := 0

		err := testDB.Delete(correctID)

		assert.Equal(t, err, nil)
		assert.Equal(t, len(testDB.Posts), wantPostsLen)
	})
	t.Run("post not found", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})

		incorrectID := 0000
		want := "Post not found"

		err := testDB.Delete(incorrectID)

		assert.Equal(t, err.Error(), want)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("post updated", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})

		updatedPost := Post{
			Id:     1,
			Author: "testUser1",
			Title:  "updatedTitle1",
			Body:   "updatedBody1",
		}

		correctID := 1
		wantPostTitle := updatedPost.Title
		wantPostBody := updatedPost.Body

		err := testDB.Update(correctID, updatedPost)

		assert.Equal(t, err, nil)

		post, _ := testDB.GetById(correctID)
		assert.Equal(t, post.Title, wantPostTitle)
		assert.Equal(t, post.Body, wantPostBody)
	})
	t.Run("post not found", func(t *testing.T) {
		testDB := InMemoryDB{}
		testDB.Posts = append(testDB.Posts, Post{
			Id:     1,
			Author: "testUser1",
			Title:  "testTitle1",
			Body:   "testBody1",
		})

		updatedPost := Post{
			Id:     1,
			Author: "testUser1",
			Title:  "updatedTitle1",
			Body:   "updatedBody1",
		}

		incorrectID := 0000
		want := "Post not found"

		err := testDB.Update(incorrectID, updatedPost)

		assert.Equal(t, err.Error(), want)
	})

}
