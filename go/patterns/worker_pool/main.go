package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

const workers = 3

func putTasksInChanel() chan int {
	lst := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	taskChan := make(chan int)
	go func() {
		for _, i := range lst {
			taskChan <- i
		}
		close(taskChan)
	}()

	return taskChan

}

func worker(ctx context.Context, id int, taskChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case x, ok := <-taskChan:
			if !ok {
				return
			}
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Worker %d finished task: %d\n", id, x)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	taskChan := putTasksInChanel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	var wg sync.WaitGroup

	for i := range workers {
		wg.Add(1)
		go func(id int) {
			worker(ctx, id, taskChan, &wg)
		}(i)
	}
	go func() {
		<-quit
		cancel()
	}()
	wg.Wait()
	fmt.Println("all workers finished")
}
