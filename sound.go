package wav

import (
	"encoding/binary"
	"io"
	"os"
	"time"
)

// Sound represents and abstract list of samples which can be encoded to a
// file.
type Sound interface {
	Channels() int
	Clone() Sound
	Duration() time.Duration
	SampleRate() int
	Samples() []Sample
	SetSamples([]Sample)
	Write(io.Writer) error
}

// WriteFile saves a sound to a file.
func WriteFile(s Sound, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.Write(f)
}

// NewPCM8Sound creates a new empty Sound with given parameters.
func NewPCM8Sound(channels int, sampleRate int) Sound {
	res := wavSound8{wavSound{NewHeader(), []Sample{}}}
	res.header.Format.BitsPerSample = 8
	res.header.Format.BlockAlign = uint16(channels)
	res.header.Format.ByteRate = uint32(sampleRate * channels)
	res.header.Format.SampleRate = uint32(sampleRate)
	res.header.Format.NumChannels = uint16(channels)
	return &res
}

// NewPCM16Sound creates a new empty Sound with given parameters.
func NewPCM16Sound(channels int, sampleRate int) Sound {
	res := wavSound16{wavSound{NewHeader(), []Sample{}}}
	res.header.Format.BitsPerSample = 16
	res.header.Format.BlockAlign = uint16(channels * 2)
	res.header.Format.ByteRate = uint32(sampleRate * channels * 2)
	res.header.Format.SampleRate = uint32(sampleRate)
	res.header.Format.NumChannels = uint16(channels)
	return &res
}

// ReadSound reads a sound from an io.Reader.
func ReadSound(f io.Reader) (Sound, error) {
	r, err := NewReader(f)
	if err != nil {
		return nil, err
	}
	if r.Header().Format.BitsPerSample != 8 &&
		r.Header().Format.BitsPerSample != 16 {
		return nil, ErrSampleSize
	}
	samples := make([]Sample, r.Remaining())
	_, err = r.Read(samples)
	if err != ErrDone && err != nil {
		return nil, err
	}
	if r.Header().Format.BitsPerSample == 8 {
		return &wavSound8{wavSound{r.Header(), samples}}, nil
	} else {
		return &wavSound16{wavSound{r.Header(), samples}}, nil
	}
}

func ReadSoundFile(path string) (Sound, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadSound(f)
}

type wavSound struct {
	header  *Header
	samples []Sample
}

func (s *wavSound) Channels() int {
	return int(s.header.Format.NumChannels)
}

func (s *wavSound) Clone() wavSound {
	newSamples := make([]Sample, len(s.samples))
	copy(newSamples, s.samples)
	return wavSound{s.header, newSamples}
}

func (s *wavSound) Duration() time.Duration {
	return s.Header().Duration()
}

func (s *wavSound) Header() *Header {
	h := s.header
	h.Data.Size = uint32(s.header.Format.BitsPerSample/8) *
		uint32(len(s.Samples()))
	h.File.Size = 36 + s.header.Data.Size
	return h
}

func (s *wavSound) SampleRate() int {
	return int(s.header.Format.SampleRate)
}

func (s *wavSound) Samples() []Sample {
	return s.samples
}

func (s *wavSound) SetSamples(ss []Sample) {
	s.samples = ss
}

type wavSound8 struct {
	wavSound
}

func (s *wavSound8) Clone() Sound {
	return &wavSound8{s.wavSound.Clone()}
}

func (s *wavSound8) Write(w io.Writer) error {
	// Write the header
	if err := binary.Write(w, binary.LittleEndian, s.Header()); err != nil {
		return err
	}
	data := make([]byte, len(s.Samples()))
	// Write the actual data
	for i, sample := range s.Samples() {
		data[i] = byte(sample*0x80 + 0x80)
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

type wavSound16 struct {
	wavSound
}

func (s *wavSound16) Clone() Sound {
	return &wavSound16{s.wavSound.Clone()}
}

func (s *wavSound16) Write(w io.Writer) error {
	// Write the header
	if err := binary.Write(w, binary.LittleEndian, s.Header()); err != nil {
		return err
	}
	data := make([]int16, len(s.Samples()))
	for i, sample := range s.Samples() {
		data[i] = int16(sample * 0x8000)
	}
	if err := binary.Write(w, binary.LittleEndian, data); err != nil {
		return err
	}
	return nil
}
