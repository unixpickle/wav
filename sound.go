package wav

import (
	"encoding/binary"
	"io"
	"os"
)

// Sound holds a list of PCM samples and a WAV header.
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
	res.header.Format.NumChannels = uint16(channels)
	return &res
}

// NewPCM16Sound creates a new empty Sound with given parameters.
func NewPCM16Sound(channels int, sampleRate int) *Sound {
	res := Sound{NewHeader(), [][]Sample{}}
	res.header.Format.BitsPerSample = 16
	res.header.Format.BlockAlign = uint16(channels * 2)
	res.header.Format.ByteRate = uint32(sampleRate * channels * 2)
	res.header.Format.SampleRate = uint32(sampleRate)
	res.header.Format.NumChannels = uint16(channels)
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
// The header's size data will be modified to fit the sound's sample data.
func (s *Sound) Header() Header {
	h := s.header
	h.Data.Size = uint32(s.header.Format.BlockSize()) *
		uint32(len(s.Samples()))
	h.File.Size = 36 + s.header.Data.Size
	return h
}

// NumChannels returns the number of channels in a sound.
func (s *Sound) NumChannels() int {
	return int(s.header.Format.NumChannels)
}

// SampleRate returns the number of samples per second per channel.
func (s *Sound) SampleRate() int {
	return int(s.header.Format.SampleRate)
}

// Samples returns an array of arrays.
// The inner arrays contain a single sample per channel.
func (s *Sound) Samples() [][]Sample {
	return s.samples
}

// SetSamples sets the sample data for the sound
func (s *Sound) SetSamples(ss [][]Sample) {
	s.samples = ss
}

// Write writes a WAV file (including its header) to an io.Writer.
func (s *Sound) Write(w io.Writer) error {
	// Write the header
	if err := binary.Write(w, binary.LittleEndian, s.Header()); err != nil {
		return err
	}
	// Write the actual data
	if s.header.Format.BitsPerSample == 8 {
		for _, block := range s.Samples() {
			for _, sample := range block {
				data := []byte{byte(sample*0x80 + 0x80)}
				if _, err := w.Write(data); err != nil {
					return err
				}
			}
		}
	} else if s.header.Format.BitsPerSample == 16 {
		for _, block := range s.Samples() {
			for _, sample := range block {
				num := uint16(sample * 0x8000)
				data := []byte{byte(num & 0xff), byte((num >> 8) & 0xff)}
				if _, err := w.Write(data); err != nil {
					return err
				}
			}
		}
	} else {
		return ErrSampleSize
	}
	return nil
}
