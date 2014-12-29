package main

import (
	"errors"
	"fmt"
	"github.com/unixpickle/wav"
	"os"
)

func main() {
	if err := ErrMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ErrMain() error {
	if len(os.Args) != 2 {
		return errors.New("Usage: wav-duration <file.wav>")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	header, err := wav.ReadHeader(f)
	f.Close()
	if err != nil {
		return err
	}
	fmt.Println("header", header)
	fmt.Println("duration is", header.Duration())
	return nil
}
