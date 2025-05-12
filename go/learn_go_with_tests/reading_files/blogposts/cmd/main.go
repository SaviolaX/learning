package main

import (
	"log"
	"os"

	"github.com/SaviolaX/blogposts"
)

func main() {
	posts, err := blogposts.NewPostsFromFS(os.DirFS("posts"))
	if err != nil {
		log.Fatal(err)
	}
	for i, post := range posts {
		log.Println(i, "-", post)
	}
}
