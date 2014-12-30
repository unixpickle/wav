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

// Valid returns true only if the ID and format match the expected values for
// WAVE files.
func (h FileHeader) Valid() bool {
	return h.ID == 0x46464952 && h.Format == 0x45564157
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

// Returns true only if the ID, size, and audio format match those of a PCM WAV
// audio file.
func (f FormatHeader) Valid() bool {
	return f.ID == 0x20746d66 && f.Size == 0x10 && f.AudioFormat == 1
}

// Header is the canonical header for all WAV files
type Header struct {
	File   FileHeader
	Format FormatHeader
	Data   ChunkHeader
}

// NewHeader creates a header with some reasonable defaults.
func NewHeader() *Header {
	var result Header
	result.File.ID = 0x46464952
	result.File.Format = 0x45564157
	result.Format.ID = 0x20746d66
	result.Format.Size = 0x10
	result.Format.AudioFormat = 1
	result.Data.ID = 0x61746164
	return &result
}

// ReadHeader reads a header from a reader.
// This does basic verification to make sure the header is valid.
func ReadHeader(r io.Reader) (*Header, error) {
	var h Header

	// Attempt to read the header
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	// Skip over arbitrary chunks that are not the data chunk
	for h.File.Valid() && h.Format.Valid() && h.Data.ID != 0x61746164 {
		unused := int(h.Data.Size)
		_, err := io.ReadFull(r, make([]byte, unused))
		if err != nil {
			return nil, err
		}
		err = binary.Read(r, binary.LittleEndian, &h.Data)
		if err != nil {
			return nil, err
		}
	}

	// Make sure the header is valid
	if !h.Valid() {
		return nil, ErrInvalid
	}

	// Make sure we support the bitrate
	sSize := h.Format.BitsPerSample
	if sSize != 8 && sSize != 16 {
		return nil, ErrSampleSize
	}

	return &h, nil
}

// Duration returns the duration of the WAV file.
func (h *Header) Duration() time.Duration {
	samples := h.Data.Size / uint32(h.Format.BlockSize())
	seconds := float64(samples) / float64(h.Format.SampleRate)
	return time.Duration(seconds * float64(time.Second))
}

// Valid returns true only if the header is for a valid WAV PCM audio file.
func (h *Header) Valid() bool {
	return h.File.Valid() && h.Format.Valid() && h.Data.ID == 0x61746164
}

// Write writes the header to a writer.
func (h *Header) Write(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, h)
}

type fmtContent struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}
