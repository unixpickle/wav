package main

import (
	"errors"
	"fmt"
	"github.com/unixpickle/wav"
	"math"
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
	if len(os.Args) != 3 {
		return errors.New("Usage: wav-tone <frequency> <file.wav>")
	}
	freq, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}
	sampleRate := 44100
	sound := wav.NewPCM8Sound(1, sampleRate)
	for i := 0; i < sampleRate * 1; i++ {
		time := float64(i) / float64(sampleRate)
		value := wav.Sample(math.Sin(time * math.Pi * 2 * float64(freq)))
		sound.Samples = append(sound.Samples, []wav.Sample{value})
	}
	f, err := os.Create(os.Args[2])
	if err != nil {
		return err
	}
	defer f.Close()
	return sound.Write(f)
}
