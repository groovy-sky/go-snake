package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	ch := make(chan string)
	go func(ch chan string) {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b = make([]byte, 1)
		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)
	keyTimer := time.Tick(1 * time.Second)
	for {
		select {
		case <-keyTimer:
			fmt.Println("Hello")
		case stdin, _ := <-ch:
			fmt.Println(stdin)
		}

	}
}
