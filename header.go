package wav

import (
	"encoding/binary"
	"io"
	"time"
)

// ChunkHeader is the generic 64-bit header in WAV files.
type ChunkHeader struct {
	ID   uint32
	Size uint32
}

// FileHeader is the "RIFF" chunk
type FileHeader struct {
	ChunkHeader
	Format uint32
}

// FormatHeader is the "fmt" sub-chunk
type FormatHeader struct {
	ChunkHeader
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

// BlockSize returns the number of bytes per sample-channel.
func (f FormatHeader) BlockSize() uint16 {
	return (f.BitsPerSample / 8) * f.NumChannels
}

// Header is the canonical header for all WAV files
type Header struct {
	File   FileHeader
	Format FormatHeader
	Data   ChunkHeader
}

// NewHeader creates a header with some reasonable defaults.
func NewHeader() Header {
	var result Header
	result.File.ID = 0x46464952
	result.File.Format = 0x45564157
	result.Format.ID = 0x20746d66
	result.Format.Size = 0x10
	result.Format.AudioFormat = 1
	result.Data.ID = 0x61746164
	return result
}

// ReadHeader reads a header from a reader.
// This does basic verification to make sure the header is valid.
func ReadHeader(r io.Reader) (Header, error) {
	var h Header
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return h, err
	}
	if h.File.ID != 0x46464952 || h.File.Format != 0x45564157 ||
		h.Format.ID != 0x20746d66 || h.Data.ID != 0x61746164 {
		return h, ErrChunkID
	}
	sSize := h.Format.BitsPerSample
	if sSize != 8 && sSize != 16 {
		return h, ErrSampleSize
	}
	if h.Format.AudioFormat != 1 {
		return h, ErrUnknownFormat
	}
	return h, nil
}

// Duration returns the duration of the WAV file.
func (h Header) Duration() time.Duration {
	samples := h.Data.Size / uint32(h.Format.BlockSize())
	seconds := float64(samples) / float64(h.Format.SampleRate)
	return time.Duration(seconds * float64(time.Second))
}

// Write writes the header to a writer.
func (h Header) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, h)
}
