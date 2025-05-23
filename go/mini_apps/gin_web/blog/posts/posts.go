package posts

import (
	"errors"
	"slices"
)

type Post struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type InMemoryDB struct {
	Posts []Post
}

func (i *InMemoryDB) Update(id int, updatedPost Post) error {
	for idx, post := range i.Posts {
		if post.Id == id {
			i.Posts[idx] = updatedPost
			return nil
		}
	}
	return errors.New("Post not found")
}

func (i *InMemoryDB) Delete(id int) error {
	for idx, post := range i.Posts {
		if post.Id == id {
			i.Posts = slices.Delete(i.Posts, idx, idx+1)
			return nil
		}
	}
	return errors.New("Post not found")
}

func (i *InMemoryDB) Add(post Post) {
	i.Posts = append(i.Posts, post)
}

func (i *InMemoryDB) GetAll() []Post {
	return i.Posts
}

func (i *InMemoryDB) GetById(id int) (Post, error) {
	post := Post{}

	for _, p := range i.Posts {
		if p.Id == id {
			post = p
		}
	}

	if post.Author == "" || post.Title == "" || post.Body == "" {
		return post, errors.New("Post not found")
	}

	return post, nil
}
