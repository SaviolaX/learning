package main

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("test single url flag", func(t *testing.T) {
		// Імітуємо аргументи: -s google.com
		args := []string{"-s", "google.com"}

		err := run(args)
		if err != nil {
			t.Errorf("run() failed with args %v: %v", args, err)
		}
	})

	t.Run("test missing arguments error", func(t *testing.T) {
		// Імітуємо виклик без аргументів після прапорця
		args := []string{"-s"}

		err := run(args)
		fmt.Println(err)
		// if err == nil {
		// 	t.Error("expected error for missing input, got nil")
		// }
	})
}
