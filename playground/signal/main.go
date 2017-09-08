package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ReloadImage() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		s := <-c
		fmt.Println("Got signal:", s)
	}
}

func main() {

	go ReloadImage()

	time.Sleep(time.Hour)
}
