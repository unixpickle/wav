package wav

import (
	"encoding/binary"
	"io"
)

type Reader struct {
	h         Header
	r         io.Reader
	remaining int32
}

// NewReader creates a new stream for reading.
func NewReader(r io.Reader) (*Reader, error) {
	h, err := ReadHeader(r)
	if err != nil {
		return nil, err
	}
	if h.FileHeader.ChunkID != 0x46464952 {
		return nil, ErrChunkID
	} else if h.FormatHeader.ChunkID != 0x20746d66 {
		return nil, ErrChunkID
	}
	sSize := h.FormatHeader.BitsPerSample
	if sSize != 8 && sSize != 16 {
		return nil, ErrSampleSize
	}
	if h.FormatHeader.AudioFormat != 1 {
		return nil, ErrUnknownFormat
	}
	return &Reader{h, r, h.DataHeader.ChunkSize}, nil
}

// Header returns the header for a stream.
func (r *Reader) Header() Header {
	return r.h
}

// Read returns a single sample for each channel.
// The samples are signed 16-bit values but are represented as ints so pecision
// can be added later on.
// If the end of stream is reached, ErrDone will be returned.
func (r *Reader) Read() ([]int, error) {
	if r.remaining == 0 {
		return nil, ErrDone
	}
	r.remaining -= int32(r.h.FormatHeader.BlockAlign)

	// Decode the list of samples
	if r.h.FormatHeader.BitsPerSample == 8 {
		res := make([]uint8, r.h.FormatHeader.NumChannels)
		if err := binary.Read(r.r, binary.LittleEndian, res); err != nil {
			return nil, err
		}
		realRes := make([]int, len(res))
		for i, x := range res {
			realRes[i] = (int(x) - 0x80) * 0x100
		}
		return realRes, nil
	} else if r.h.FormatHeader.BitsPerSample == 16 {
		res := make([]int16, r.h.FormatHeader.NumChannels)
		if err := binary.Read(r.r, binary.LittleEndian, res); err != nil {
			return nil, err
		}
		realRes := make([]int, len(res))
		for i, x := range res {
			realRes[i] = int(x)
		}
		return realRes, nil
	}

	return nil, ErrSampleSize
}
