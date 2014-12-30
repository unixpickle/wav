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
		return errors.New("Usage: wav-crop <file.wav> <start> <end> " +
			"<output.wav>")
	}
	s, err := wav.ReadSound(os.Args[1])
	if err != nil {
		return err
	}
	start, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		return err
	}
	end, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		return err
	}
	wav.Crop(s, time.Duration(start*float64(time.Second)),
		time.Duration(end*float64(time.Second)))
	return wav.WriteFile(s, os.Args[4])
}
