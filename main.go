package main

import (
	"fmt"
	"time"
)

func func1(channel chan string) {
	channel <- "boom"
}

func func2(channel chan string) {
	channel <- "What is this nice thing"
}

func main() {
	channel := make(chan string)
	channel2 := make(chan string)
	go func2(channel)
	go func1(channel2)

	for i := 0; i < 10; i++ {
		select {
		case v := <-channel:
			fmt.Println("This is message from the channel", v)
		case v := <-channel:
			fmt.Println("This is message from the channel2", v)
		default:
			fmt.Println("Waiting for message")
		}

		time.Sleep(time.Millisecond)
	}
}
