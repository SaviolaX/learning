package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second / 3)
	requests := 10

	for x := range requests {
		<-ticker.C
		fmt.Printf("request %d: %s\n", x+1, time.Now().Format("15:04:05.000"))
	}
	ticker.Stop()
}
