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
	if len(os.Args) < 3 {
		return errors.New("Usage: wav-join <input.wav> [<another.wav> ...] " +
			"<output.wav>")
	}
	s, err := wav.ReadSoundFile(os.Args[1])
	if err != nil {
		return err
	}
	for _, f := range os.Args[2 : len(os.Args)-1] {
		nextS, err := wav.ReadSoundFile(f)
		if err != nil {
			return err
		}
		wav.Append(s, nextS)
	}
	return wav.WriteFile(s, os.Args[len(os.Args)-1])
}
