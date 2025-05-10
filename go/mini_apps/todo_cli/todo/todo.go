package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {

	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	ls := *t

	// index has to start from 1
	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	*t = slices.Delete(ls, index-1, index)

	return nil
}

func (t *Todos) Load(filename string) error {

	// read file and store data in bytes
	data, err := os.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return errors.New("could not read file")
	}

	// Parse bytes and store data into Todos
	err = json.Unmarshal(data, t)
	if err != nil {
		return errors.New("could not read bytes")
	}

	return nil
}

func (t *Todos) Store(filename string) error {

	// convert into bytes
	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}

	// write bytes into file
	err = os.WriteFile(filename, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Print() {
	if len(*t) == 0 {
		fmt.Println("No tasks...")
	} else {
		for i, item := range *t {
			i++ // to start count from 1 instead 0
			fmt.Printf("%d - %s - %v - %s\n", i, item.Task, item.Done,
				item.CreatedAt.Format("2006-01-02"))
		}
	}
}
