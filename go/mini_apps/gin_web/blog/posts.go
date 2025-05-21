package main

import (
	"errors"
	"fmt"
	"slices"
)

type Post struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type InMemoryDB struct {
	posts []Post
}

func (i *InMemoryDB) Delete(id int) error {
	for i, post := range storage.posts {
		if post.Id == id {
			storage.posts = slices.Delete(storage.posts, i, i+1)
			fmt.Println(storage.posts)
			return nil
		}
	}
	return errors.New("Post not found")
}

func (i *InMemoryDB) Add(post Post) {
	i.posts = append(i.posts, post)
}

func (i *InMemoryDB) GetAll() []Post {
	return i.posts
}

func (i *InMemoryDB) GetById(id int) (Post, error) {
	post := Post{}

	for _, p := range i.posts {
		if p.Id == id {
			post = p
		}
	}

	if post.Author == "" || post.Title == "" || post.Body == "" {
		return post, errors.New("Post not found")
	}

	return post, nil
}
