package wav

import "errors"

var (
	ErrDone          = errors.New("Done reading audio data.")
	ErrSampleSize    = errors.New("Unsupported sample rate.")
	ErrChunkID       = errors.New("Read unexpected chunk ID.")
	ErrUnknownFormat = errors.New("Unsupported or invalid audio format.")
)
