package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type input struct {
	filename, tag string
}

var wait sync.WaitGroup

var pollTime = 100 * time.Millisecond

func main() {
	inputs := make([]input, 0, len(os.Args)-1)

	tagLength := 0

	for _, arg := range os.Args[1:] {
		var filename, tag string

		split := strings.SplitN(arg, ":", 2)

		filename = split[0]
		if len(split) > 1 {
			tag = split[1]
		} else {
			tag = filename
		}

		inputs = append(inputs, input{filename, tag})

		if len(tag) > tagLength {
			tagLength = len(tag)
		}
	}

	if len(inputs) == 0 {
		fmt.Fprintln(os.Stderr, "please specify at least one input file")
		os.Exit(2)
	}

	for _, input := range inputs {
		// Pad the specified tag so that all the output lines line up.
		tag := input.tag
		for i := 0; i < tagLength-len(input.tag); i += 1 {
			tag += " "
		}
		tag += " | "

		wait.Add(1)
		go func(filename, tag string) {
			follow(filename, tag)
			wait.Done()
		}(input.filename, tag)
	}

	go func() {
		for line := range buffer {
			fmt.Println(line)
		}
	}()

	wait.Wait()
	close(buffer)
}

func follow(filename, tag string) {
	var lastRead time.Time
	var offset int64

	for {
		f, err := os.Open(filename)
		if err != nil {
			// Maybe the file doesn't exist yet; wait a little bit and try again.
			time.Sleep(pollTime)
			continue
		}

		fi, err := f.Stat()
		if err != nil {
			time.Sleep(pollTime)
			continue
		}

		if fi.ModTime().After(lastRead) {
			// If the file has gotten smaller all of a sudden, assume it's been
			// rotated, and start from the beginning.
			if fi.Size() < offset {
				offset = 0
			}

			f.Seek(offset, 0)

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				buffer <- fmt.Sprintf("%s%s\n", tag, scanner.Text())
			}

			offset, _ = f.Seek(0, 1)
			lastRead = fi.ModTime()
		} else {
			time.Sleep(pollTime)
		}
	}
}

var buffer = make(chan string, 1024)
