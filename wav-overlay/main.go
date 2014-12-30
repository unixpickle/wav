package main

import (
	"errors"
	"fmt"
	"github.com/unixpickle/wav"
	"os"
	"strconv"
	"time"
)

func main() {
	if err := ErrMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ErrMain() error {
	if len(os.Args) != 5 {
		return errors.New("Usage: wav-overlay <input.wav> <overlay.wav> " +
			"<start> <output.wav>")
	}
	start, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		return err
	}
	s1, err := wav.ReadSound(os.Args[1])
	if err != nil {
		return err
	}
	s2, err := wav.ReadSound(os.Args[2])
	if err != nil {
		return err
	}
	wav.Volume(s1, 0.5)
	wav.Volume(s2, 0.5)
	wav.Overlay(s1, s2, time.Duration(start*float64(time.Second)))
	return wav.WriteFile(s1, os.Args[4])
}
