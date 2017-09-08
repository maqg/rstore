package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	s.Signal()
	fmt.Println("Got signal:", s)
}
