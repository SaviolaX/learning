package main

import (
	"fmt"
)

func main() {
	lst := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	chn := generate(lst)

	res := square(chn)

	for x := range res {
		fmt.Println(x)
	}
}

func generate(lst []int) <-chan int {
	chn := make(chan int)

	go func() {
		for _, i := range lst {
			chn <- i
		}
		close(chn)
	}()

	return chn
}

func square(in <-chan int) <-chan int {
	res := make(chan int)

	go func() {
		for x := range in {
			res <- x * x
		}
		close(res)
	}()

	return res
}
