package wav

import "errors"

var (
	ErrDone          = errors.New("Done reading audio data.")
	ErrSampleSize    = errors.New("Unsupported sample size.")
	ErrInvalid       = errors.New("The input data was invalid.")
)
