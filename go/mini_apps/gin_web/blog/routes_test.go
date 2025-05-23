package main

import (
	"blog/posts"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/quii/go-graceful-shutdown/assert"
)

const (
	GET_post_path  = "/posts/:id"
	GET_all_posts  = "/posts"
	AddPostPath    = "/posts/add"
	DeletePostPath = "/posts/delete/:id"
	UpdatePostPath = "/posts/update/:id"
)

// test for POST /posts/add
func TestAddPost(t *testing.T) {
	testStorage := posts.InMemoryDB{}

	router := setupRouter()
	router.POST(AddPostPath, addPost(&testStorage))

	newPost := posts.Post{
		Id:     1,
		Author: "testUser",
		Title:  "testPost3",
		Body:   "testBody",
	}

	newPostJson, _ := json.Marshal(newPost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts/add", strings.NewReader(string(newPostJson)))
	router.ServeHTTP(w, req)

	expectedResponse := map[string]string{"message": "A new post added"}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, w.Body.String(), string(expectedResponseJson))
}

// test for GET /posts/:id
func TestPost(t *testing.T) {
	testPost := posts.Post{
		Id:     1,
		Author: "testUser",
		Title:  "testPost1",
		Body:   "testBody",
	}
	testStorage := posts.InMemoryDB{}
	testStorage.Posts = append(testStorage.Posts, testPost)

	router := setupRouter()
	router.GET(GET_post_path, getPost(&testStorage))

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
		testStorage := posts.InMemoryDB{}
		testStorage.Posts = append(testStorage.Posts, posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "testPost1",
			Body:   "testBody",
		})
		testStorage.Posts = append(testStorage.Posts, posts.Post{
			Id:     2,
			Author: "testUser",
			Title:  "testPost2",
			Body:   "testBody",
		})

		router := setupRouter()
		router.GET(GET_all_posts, getPosts(&testStorage))

		posts := map[string][]posts.Post{
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
		responsePostsJson, _ := json.Marshal(posts)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(responsePostsJson))
	})
	t.Run("no posts", func(t *testing.T) {
		emptyStorage := posts.InMemoryDB{}

		router := setupRouter()
		router.GET(GET_all_posts, getPosts(&emptyStorage))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/posts", nil)
		router.ServeHTTP(w, req)

		want := map[string][]posts.Post{"posts": nil}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
}

func TestDeletePost(t *testing.T) {
	testStorage := posts.InMemoryDB{}

	router := setupRouter()
	router.DELETE(DeletePostPath, deletePost(&testStorage))

	t.Run("post deleted", func(t *testing.T) {
		testStorage.Posts = append(testStorage.Posts, posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "testTitle",
			Body:   "testBody",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/posts/delete/1", nil)
		router.ServeHTTP(w, req)

		want := map[string]string{"message": "Post 1 deleted"}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(wantJSON))

		storagePostsLen := len(testStorage.Posts)
		wantStorageLen := 0
		assert.Equal(t, storagePostsLen, wantStorageLen)
	})
	t.Run("post not found", func(t *testing.T) {
		testStorage.Posts = append(testStorage.Posts, posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "testTitle",
			Body:   "testBody",
		})

		incorrectId := "0000"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/posts/delete/"+incorrectId, nil)
		router.ServeHTTP(w, req)

		want := map[string]string{"error": "Post not found"}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 404)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
	t.Run("id param is not an integer", func(t *testing.T) {
		incorrectId := "notIntegerID"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/posts/delete/"+incorrectId, nil)
		router.ServeHTTP(w, req)

		want := map[string]string{"error": "ID param is not an integer"}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 404)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
}

func TestUpdatePost(t *testing.T) {
	testStorage := posts.InMemoryDB{}

	router := setupRouter()
	router.PUT(UpdatePostPath, updatePost(&testStorage))

	t.Run("post updated", func(t *testing.T) {
		testStorage.Posts = append(testStorage.Posts, posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "testTitle",
			Body:   "testBody",
		})

		changedPost := posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "updatedTitle",
			Body:   "updatedBody",
		}
		changedPostJson, _ := json.Marshal(changedPost)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/posts/update/1", strings.NewReader(string(changedPostJson)))
		router.ServeHTTP(w, req)

		want := map[string]string{"message": "Post 1 updated"}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 200)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
	t.Run("post not found", func(t *testing.T) {
		testStorage.Posts = append(testStorage.Posts, posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "testTitle",
			Body:   "testBody",
		})

		changedPost := posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "updatedTitle",
			Body:   "updatedBody",
		}
		changedPostJson, _ := json.Marshal(changedPost)

		incorrectId := "0000"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/posts/update/"+incorrectId, strings.NewReader(string(changedPostJson)))
		router.ServeHTTP(w, req)

		want := map[string]string{"error": "Post not found"}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 404)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
	t.Run("incorrect id", func(t *testing.T) {
		changedPost := posts.Post{
			Id:     1,
			Author: "testUser",
			Title:  "updatedTitle",
			Body:   "updatedBody",
		}
		changedPostJson, _ := json.Marshal(changedPost)

		incorrectId := "incorrectId"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/posts/update/"+incorrectId, strings.NewReader(string(changedPostJson)))
		router.ServeHTTP(w, req)

		want := map[string]string{"error": "ID param is not an integer"}
		wantJSON, _ := json.Marshal(want)

		assert.Equal(t, w.Code, 404)
		assert.Equal(t, w.Body.String(), string(wantJSON))
	})
}
