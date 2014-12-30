package wav

import (
	"encoding/binary"
	"io"
)

type Sample float64

type Reader interface {
	// Header returns the reader's WAV header.
	Header() *Header

	// Read reads as many as len(out) samples.
	// The samples are signed values ranging between -1.0 and 1.0.
	// Channels are packed side-by-side, so for a stereo track it'd be LRLRLR...
	// If the end of stream is reached, ErrDone will be returned.
	Read(out []Sample) (int, error)

	// Remaining returns the number of samples left to read.
	Remaining() int
}

// NewReader wraps an io.Reader with a Reader.
func NewReader(r io.Reader) (Reader, error) {
	h, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}
	remaining := int(h.Data.Size / uint32(h.Format.BitsPerSample/8))
	if h.Format.BitsPerSample == 8 {
		return &pcm8Reader{reader{h, r, remaining}}, nil
	} else if h.Format.BitsPerSample == 16 {
		return &pcm16Reader{reader{h, r, remaining}}, nil
	}
	return nil, ErrSampleSize
}

type pcm8Reader struct {
	reader
}

func (r pcm8Reader) Read(out []Sample) (int, error) {
	if r.remaining == 0 {
		return 0, ErrDone
	}

	toRead := len(out)
	if toRead > r.remaining {
		toRead = r.remaining
	}

	// Decode the list of raw samples
	raw := make([]uint8, toRead)
	if err := binary.Read(r.input, binary.LittleEndian, raw); err != nil {
		return 0, err
	}
	for i, x := range raw {
		out[i] = (Sample(x) - 0x80) / 0x80
	}

	// Return the amount read and a possible ErrDone error.
	r.remaining -= toRead
	if r.remaining == 0 {
		return toRead, ErrDone
	}
	return toRead, nil
}

type pcm16Reader struct {
	reader
}

func (r pcm16Reader) Read(out []Sample) (int, error) {
	if r.remaining == 0 {
		return 0, ErrDone
	}

	toRead := len(out)
	if toRead > r.remaining {
		toRead = r.remaining
	}

	// Decode the list of raw samples
	raw := make([]int16, toRead)
	if err := binary.Read(r.input, binary.LittleEndian, raw); err != nil {
		return 0, err
	}
	for i, x := range raw {
		out[i] = Sample(x) / 0x8000
	}

	// Return the amount read and a possible ErrDone error.
	r.remaining -= toRead
	if r.remaining == 0 {
		return toRead, ErrDone
	}
	return toRead, nil
}

type reader struct {
	header    *Header
	input     io.Reader
	remaining int
}

func (r *reader) Header() *Header {
	return r.header
}

func (r *reader) Remaining() int {
	return r.remaining
}
