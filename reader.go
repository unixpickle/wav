package wav

import (
	"encoding/binary"
	"io"
)

type Sample float64

type Reader struct {
	header    Header
	reader    io.Reader
	remaining uint32
}

// NewReader wraps an io.Reader with a Reader.
func NewReader(r io.Reader) (*Reader, error) {
	h, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}
	return &Reader{h, r, h.Data.Size}, nil
}

// Header returns the header that a Reader read during NewReader().
func (r *Reader) Header() Header {
	return r.header
}

// Read returns a single sample for each channel.
// The samples are signed values ranging between -1.0 and 1.0.
// If the end of stream is reached, ErrDone will be returned.
func (r *Reader) Read() ([]Sample, error) {
	if r.remaining == 0 {
		return nil, ErrDone
	}
	h := r.header
	r.remaining -= uint32(h.Format.BlockSize())

	// Decode the list of samples
	res := make([]Sample, h.Format.NumChannels)
	if h.Format.BitsPerSample == 8 {
		raw := make([]uint8, len(res))
		if err := binary.Read(r.reader, binary.LittleEndian, raw); err != nil {
			return nil, err
		}
		for i, x := range raw {
			res[i] = (Sample(x) - 0x80) / 0x80
		}
	} else if h.Format.BitsPerSample == 16 {
		raw := make([]int16, len(res))
		if err := binary.Read(r.reader, binary.LittleEndian, raw); err != nil {
			return nil, err
		}
		for i, x := range raw {
			res[i] = Sample(x) / 0x8000
		}
	} else {
		return nil, ErrSampleSize
	}
	return res, nil
}

// Remaining returns the number of samples left to read.
func (r *Reader) Remaining() int {
	return int(r.remaining / uint32(r.header.Format.BlockAlign))
}
