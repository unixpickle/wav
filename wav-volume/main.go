package main

import (
	"errors"
	"fmt"
	"github.com/unixpickle/wav"
	"os"
	"strconv"
)

func main() {
	if err := ErrMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ErrMain() error {
	if len(os.Args) != 4 {
		return errors.New("Usage: wav-volume <file.wav> <volume> <output.wav>")
	}
	s, err := wav.ReadSoundFile(os.Args[1])
	if err != nil {
		return err
	}
	scale, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		return err
	}
	wav.Volume(s, scale)
	return wav.WriteFile(s, os.Args[3])
}
