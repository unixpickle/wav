package wav

import (
	"encoding/binary"
	"io"
	"os"
)

type Sound struct {
	header  Header
	samples [][]Sample
}

// NewPCM8Sound creates a new empty Sound with given parameters.
func NewPCM8Sound(channels int, sampleRate int) *Sound {
	res := Sound{NewHeader(), [][]Sample{}}
	res.header.Format.BitsPerSample = 8
	res.header.Format.BlockAlign = uint16(channels)
	res.header.Format.ByteRate = uint32(sampleRate * channels)
	res.header.Format.SampleRate = uint32(sampleRate)
	return &res
}

// NewPCM16Sound creates a new empty Sound with given parameters.
func NewPCM16Sound(channels int, sampleRate int) *Sound {
	res := Sound{NewHeader(), [][]Sample{}}
	res.header.Format.BitsPerSample = 16
	res.header.Format.BlockAlign = uint16(channels * 2)
	res.header.Format.ByteRate = uint32(sampleRate * channels * 2)
	res.header.Format.SampleRate = uint32(sampleRate)
	return &res
}

// ReadSound reads a sound from a file.
func ReadSound(path string) (*Sound, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := NewReader(f)
	if err != nil {
		return nil, err
	}
	samples := make([][]Sample, r.Remaining())
	for i := 0; i < len(samples); i++ {
		samples[i], err = r.Read()
		if err != nil {
			return nil, err
		}
	}
	return &Sound{r.Header(), samples}, nil
}

// Header returns the header for the sound.
func (s *Sound) Header() Header {
	h := s.header
	h.Data.Size = uint32(s.header.Format.BlockSize()) * uint32(len(s.samples))
	h.File.Size = 36 + s.header.Data.Size
	return h
}

// SampleRate returns the number of samples per second per channel.
func (s *Sound) SampleRate() int {
	return int(s.header.Format.SampleRate)
}

// Samples returns the sample data for the sound.
// Each element in the outer array is an array of channel samples.
func (s *Sound) Samples() [][]Sample {
	return s.samples
}

// Write writes a WAV file (including its header) to an io.Writer.
func (s *Sound) Write(w io.Writer) error {
	// Write the header
	if err := binary.Write(w, binary.LittleEndian, s.Header()); err != nil {
		return err
	}
	// Write the actual data
	if s.header.Format.BitsPerSample == 8 {
		for _, block := range s.samples {
			for _, sample := range block {
				data := []byte{byte(sample*0x80 + 0x80)}
				if _, err := w.Write(data); err != nil {
					return err
				}
			}
		}
	} else {
		for _, block := range s.samples {
			for _, sample := range block {
				num := uint16(sample * 0x8000)
				data := []byte{byte(num & 0xff), byte((num >> 8) & 0xff)}
				if _, err := w.Write(data); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
