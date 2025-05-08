package main

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	todos := Todos{}
	task := "new task"

	todos.Add(task)

	task_from_todos := todos[0]

	if task_from_todos.Task != task {
		t.Error("Task was not created.")
	}
}

func TestComplete(t *testing.T) {
	todos := Todos{
		item{
			Task:        "new task",
			Done:        false,
			CreatedAt:   time.Now(),
			CompletedAt: time.Time{},
		},
	}
	// Correct index starts from 1, not from 0
	t.Run("incorrect index", func(t *testing.T) {
		incorrect_index := 0
		got := todos.Complete(incorrect_index)
		want := "invalid index"

		if got.Error() != want {
			t.Error("Task has not been completed.")
		}
	})
	t.Run("correct index", func(t *testing.T) {
		correct_index := 1
		todos.Complete(correct_index)

		new_task := todos[0]

		current_time_str := time.Now().String()[:11]

		if new_task.Done != true {
			t.Error("Task status didn't change.")
		}

		if new_task.CompletedAt.String()[:11] != current_time_str {
			t.Error("Task completedAt time didn't change.")
		}
	})
}

func TestDelete(t *testing.T) {
	todos := Todos{
		item{
			Task:        "new task",
			Done:        false,
			CreatedAt:   time.Now(),
			CompletedAt: time.Time{},
		},
	}
	// Correct index starts from 1, not from 0
	t.Run("incorrect index", func(t *testing.T) {
		incorrect_index := 0
		got := todos.Delete(incorrect_index)
		want := "invalid index"

		if got.Error() != want {
			t.Error("Task has not been completed.")
		}
	})
	t.Run("correct index", func(t *testing.T) {
		incorrect_index := 1
		todos.Delete(incorrect_index)

		if len(todos) > 0 {
			t.Error("Task was not deleted.")
		}

	})

}
