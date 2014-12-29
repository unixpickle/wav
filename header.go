package wav

import (
	"encoding/binary"
	"io"
)

// ChunkHeader is the generic 64-bit header in WAV files.
type ChunkHeader struct {
	ChunkID   uint32
	ChunkSize uint32
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

// Header is the canonical header for all WAV files
type Header struct {
	FileHeader   FileHeader
	FormatHeader FormatHeader
	DataHeader   ChunkHeader
}

func ReadHeader(r io.Reader) (Header, error) {
	var h Header
	err := binary.Read(r, binary.LittleEndian, &h)
	return h, err
}

func (h Header) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, h)
}
