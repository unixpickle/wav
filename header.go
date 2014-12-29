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

// ReadHeader reads a header from a reader.
// This does not validate any part of the header.
func ReadHeader(r io.Reader) (Header, error) {
	var h Header
	err := binary.Read(r, binary.LittleEndian, &h)
	return h, err
}

// Write writes the header to a writer.
func (h Header) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, h)
}

// Duration returns the duration of the WAV file.
func (h Header) Duration() time.Duration {
	samples := h.Data.Size / uint32(h.Format.BlockSize())
	seconds := float64(samples) / float64(h.Format.SampleRate)
	return time.Duration(seconds * float64(time.Second))
}
