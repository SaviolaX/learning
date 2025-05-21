package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var storage = InMemoryDB{}

func deletePost(storage *InMemoryDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postId := c.Param("id")

		idInt, err := strconv.Atoi(postId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID param is not an integer"})
			return
		}

		err = storage.Delete(idInt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Post %d deleted", idInt),
		})
	}
}

func addPost(storage *InMemoryDB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var post Post
		c.BindJSON(&post)

		storage.Add(post)

		c.JSON(http.StatusCreated, gin.H{
			"message": "A new post added",
		})
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
		postId := c.Param("id")

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
