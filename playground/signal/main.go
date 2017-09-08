package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// ReloadImage for signal handler
func ReloadImage() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		s := <-c
		fmt.Println("Got signal:", s)
	}
}

func main() {

	fmt.Println(runtime.GOOS)

	//	go ReloadImage()

	//	time.Sleep(time.Hour)
}
