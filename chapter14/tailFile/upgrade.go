package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"syscall"
)

func upgrade(filename string) error {
	file, _ := os.Open(filename)
	fd, _ := syscall.InotifyInit()
	_, _ = syscall.InotifyAddWatch(fd, filename, syscall.IN_MODIFY)
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
		if err = waitForChange(fd); err != nil {
			return err
		}
	}
}

func waitForChange(fd int) error {
	for {
		var buf [syscall.SizeofInotifyEvent]byte
		_, err := syscall.Read(fd, buf[:])
		if err != nil {
			return err
		}
		r := bytes.NewReader(buf[:])
		var ev = syscall.InotifyEvent{}
		_ = binary.Read(r, binary.LittleEndian, &ev)
		if ev.Mask&syscall.IN_MODIFY == syscall.IN_MODIFY {
			return nil
		}
	}
}
