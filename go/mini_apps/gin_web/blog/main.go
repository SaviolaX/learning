package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var storage = InMemoryDB{
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

func addPost(storage *InMemoryDB) gin.HandlerFunc {
	return func(c *gin.Context) {

		idInt, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID param is not an integer"})
			return
		}

		newPost := Post{
			Id:     idInt,
			Author: c.Param("author"),
			Title:  c.Param("title"),
			Body:   c.Param("body"),
		}

		updatedStorage := append(*storage, newPost)
		fmt.Println(updatedStorage)
	}
}

func getPosts(storage *InMemoryDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts := storage.GetAll()
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})
	}
}

func getPost(storage *InMemoryDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postId := c.Param("postId")

		idInt, err := strconv.Atoi(postId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID param is not an integer"})
			return
		}

		post, err := storage.GetById(idInt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func setupRouter() *gin.Engine {
	return gin.Default()
}

func main() {
	router := setupRouter()

	router.Run()
}
