package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	GET_post_path = "/posts/:postId"
	GET_all_posts = "/posts"
	AddPostPath   = "/posts/add"
)

var testStorage = InMemoryDB{
	Post{
		Id:     1,
		Author: "testUser",
		Title:  "testPost1",
		Body:   "testBody",
	},
	Post{
		Id:     2,
		Author: "testUser",
		Title:  "testPost2",
		Body:   "testBody",
	},
}

func TestAddPost(t *testing.T) {
	router := setupRouter()
	router.POST(AddPostPath, addPost(&testStorage))

	newPost := Post{
		Id:     3,
		Author: "testUser",
		Title:  "testPost3",
		Body:   "testBody",
	}

	newPostJson, _ := json.Marshal(newPost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts/add", strings.NewReader(string(newPostJson)))
	router.ServeHTTP(w, req)

	//	expectedPosts := 3

	assert.Equal(t, w.Code, 201)
}

// test for GET /posts/:id
func TestPost(t *testing.T) {
	router := setupRouter()
	router.GET(GET_post_path, getPost(&testStorage))

	testPost := Post{
		Id:     1,
		Author: "testUser",
		Title:  "testPost1",
		Body:   "testBody",
	}

	postJson, _ := json.Marshal(testPost)

	t.Run("correct request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(postJson))

	})
	t.Run("incorrect id param", func(t *testing.T) {
		incorrectParam := "incorrect"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts/"+incorrectParam, nil)
		router.ServeHTTP(w, req)

		expectErrorResp := map[string]string{"error": "ID param is not an integer"}
		want, _ := json.Marshal(expectErrorResp)

		assert.Equal(t, w.Code, 404)
		assert.Equal(t, w.Body.String(), string(want))
	})
	t.Run("post not found", func(t *testing.T) {
		incorrectId := "00000"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts/"+incorrectId, nil)
		router.ServeHTTP(w, req)

		expectErrorResp := map[string]string{"error": "Post not found"}
		want, _ := json.Marshal(expectErrorResp)

		assert.Equal(t, w.Code, 404)
		assert.Equal(t, w.Body.String(), string(want))

	})
}

// test for GET /posts
func TestPosts(t *testing.T) {

	t.Run("return all posts", func(t *testing.T) {
		posts := map[string][]Post{
			"posts": {
				{
					Id:     1,
					Author: "testUser",
					Title:  "testPost1",
					Body:   "testBody",
				},
				{
					Id:     2,
					Author: "testUser",
					Title:  "testPost2",
					Body:   "testBody",
				},
			},
		}

		router := setupRouter()
		router.GET(GET_all_posts, getPosts(&testStorage))

		postsJson, _ := json.Marshal(posts)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(postsJson))
	})
	t.Run("no posts", func(t *testing.T) {
		emptyStorage := InMemoryDB{}

		router := setupRouter()
		router.GET(GET_all_posts, getPosts(&emptyStorage))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts", nil)
		router.ServeHTTP(w, req)

		want := map[string][]Post{"posts": {}}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
}
