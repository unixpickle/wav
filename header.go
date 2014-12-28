package wav

import (
	"encoding/binary"
	"io"
)

// ChunkHeader is the generic 64-bit header in WAV files.
type ChunkHeader struct {
	ChunkID   int32
	ChunkSize int32
}

// FileHeader is the "RIFF" chunk
type FileHeader struct {
	ChunkHeader
	Format int32
}

// FormatHeader is the "fmt" sub-chunk
type FormatHeader struct {
	ChunkHeader
	AudioFormat   int16
	NumChannels   int16
	SampleRate    int32
	ByteRate      int32
	BlockAlign    int16
	BitsPerSample int16
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
