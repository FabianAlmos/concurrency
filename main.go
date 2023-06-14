package main

import (
	worker "concuLec/workers"
	"fmt"
)

func main() {
	queue := worker.NewAsyncQueue()

	l := []string{"1", "2", "3"}
	queue.Start()

	for i := 0; i < 2; i++ {
		queue.Append(func() {
			func(k string) {
				fmt.Println("line", k)
			}(l[i])
		})
	}

	queue.Shutdown()
}
