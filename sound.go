package wav

import "os"

type Sound struct {
	sampleRate int
	samples    [][]int32
}

// ReadSound reads a sound from a file.
func ReadSound(path string) (*Sound, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := NewReader(f)
	if err != nil {
		return nil, err
	}
	samples := make([][]int32, r.Remaining())
	for i := 0; i < len(samples); i++ {
		samples[i], err = r.Read()
		if err != nil {
			return nil, err
		}
	}
	sr := r.Header().FormatHeader.SampleRate
	return &Sound{int(sr), samples}, nil
}

// SampleRate returns the number of samples per second in the sound.
// This number is not affected by the number of channels (i.e. 10 samples
// per second with two channels means that 20 samples total occur each
// second, 10 for each channel).
func (s *Sound) SampleRate() int {
	return s.sampleRate
}

// Samples returns the sample data for the sound.
// Each element in the outer array is an array of channel samples.
// Each sample is an integer which ranges from -0x80000000 to 0x7fffffff.
func (s *Sound) Samples() [][]int32 {
	return s.samples
}
