package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func producer(queue chan string, waitGroup *sync.WaitGroup, quit <-chan bool) {
	defer close(queue)
	for i := 0; i < 10; i++ {
		select {
		case <-quit:
			waitGroup.Done()
			return
		default:
			fmt.Println("Виробництво предмета: ", (i + 1))
			queue <- fmt.Sprintf("Предмет: %d", (i + 1))
		}
	}
	waitGroup.Done()
}

func consumer(queue chan string, waitGroup *sync.WaitGroup) {
	for val := range queue {
		time.Sleep(time.Second)
		fmt.Println("Споживання ", val)
	}
	waitGroup.Done()
}

func handleSigInt(sigInt chan os.Signal, queue chan string, quit chan<- bool) {
	_ = <-sigInt
	fmt.Println("Споживайте очікувані продукти та закінчуйте процес")
	quit <- true
}

func main() {
	fmt.Println("Початок роботи. Виробництво продукції та очікування споживання")
	queue := make(chan string)
	sigInt, quit := make(chan os.Signal), make(chan bool)
	signal.Notify(sigInt, syscall.SIGINT, syscall.SIGTERM)

	go handleSigInt(sigInt, queue, quit)

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go producer(queue, &waitGroup, quit)
	go consumer(queue, &waitGroup)

	waitGroup.Wait()
	fmt.Println("Кінець роботи")
}
