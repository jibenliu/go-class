package main

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"os"
)

func follow2(filename string) error {
	file, _ := os.Open(filename)
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	_ = watcher.Add(filename)

	r := bufio.NewReader(file)
	for {
		by, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		fmt.Print(string(by))
		if err != io.EOF {
			continue
		}
		if err = waitForChange2(watcher); err != nil {
			return err
		}
	}
}

func waitForChange2(w *fsnotify.Watcher) error {
	for {
		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				return nil
			}
		case err := <-w.Errors:
			return err
		}
	}
}
