package main

import (
	"errors"
)

type Post struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type InMemoryDB []Post

func (i InMemoryDB) GetAll() []Post {
	return i
}

func (i InMemoryDB) GetById(id int) (Post, error) {
	post := Post{}

	for _, p := range i {
		if p.Id == id {
			post.Id = id
			post.Author = p.Author
			post.Title = p.Title
			post.Body = p.Body
		}
	}

	if post.Author == "" || post.Title == "" || post.Body == "" {
		return post, errors.New("Post not found")
	}

	return post, nil
}
