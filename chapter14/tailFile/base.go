package main

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func follow(file io.Reader) error {
	r := bufio.NewReader(file)
	for {
		by, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		fmt.Print(string(by))
		if err == io.EOF {
			time.Sleep(time.Second)
		}
	}
}
