package main

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	time.Sleep(100 * time.Millisecond)
}
