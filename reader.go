package wav

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	h         Header
	r         io.Reader
	remaining uint32
}

// NewReader wraps an io.Reader with a Reader.
func NewReader(r io.Reader) (*Reader, error) {
	h, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}
	if h.File.ID != 0x46464952 {
		return nil, ErrChunkID
	} else if h.Format.ID != 0x20746d66 {
		return nil, ErrChunkID
	}
	sSize := h.Format.BitsPerSample
	if sSize != 8 && sSize != 16 {
		return nil, ErrSampleSize
	}
	if h.Format.AudioFormat != 1 {
		return nil, ErrUnknownFormat
	}
	return &Reader{h, r, h.Data.Size}, nil
}

// Header returns the header for a stream.
func (r *Reader) Header() Header {
	return r.h
}

// Read returns a single sample for each channel.
// The samples are signed 32-bit values.
// If the end of stream is reached, ErrDone will be returned.
func (r *Reader) Read() ([]int32, error) {
	if r.remaining == 0 {
		return nil, ErrDone
	}
	r.remaining -= uint32(r.h.Format.BlockSize())

	// Decode the list of samples
	if r.h.Format.BitsPerSample == 8 {
		res := make([]uint8, r.h.Format.NumChannels)
		if err := binary.Read(r.r, binary.LittleEndian, res); err != nil {
			return nil, err
		}
		realRes := make([]int32, len(res))
		for i, x := range res {
			realRes[i] = (int32(x) - 0x80) * 0x1000000
		}
		return realRes, nil
	} else if r.h.Format.BitsPerSample == 16 {
		res := make([]int16, r.h.Format.NumChannels)
		if err := binary.Read(r.r, binary.LittleEndian, res); err != nil {
			return nil, err
		}
		realRes := make([]int32, len(res))
		for i, x := range res {
			realRes[i] = int32(x) * 0x10000
		}
		return realRes, nil
	}

	return nil, ErrSampleSize
}

// Remaining returns the number of samples left to read.
func (r *Reader) Remaining() int {
	return int(r.remaining / uint32(r.h.Format.BlockAlign))
}
