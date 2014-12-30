package wav

import "time"

// Append appends sounds to a sound.
// If the various sounds have different numbers of channels, channels will be
// added or removed as needed.
// This works best if all the sounds have the same sample rate. Otherwise,
// the result will have some faster and some slower parts.
func Append(s1 Sound, sounds ...Sound) {
	for _, s := range sounds {
		if s.Channels() != s1.Channels() {
			diffChannelAppend(s1, s)
			continue
		}
		s1.SetSamples(append(s1.Samples(), s.Samples()...))
	}
}

// Crop isolates a time segment in a sound.
func Crop(s Sound, start, end time.Duration) {
	// Cannot crop an empty sound.
	if len(s.Samples()) == 0 {
		return
	}

	// Figure out sample indexes.
	startIdx := sampleIndex(s, start)
	endIdx := sampleIndex(s, end)

	// Clamp indexes
	if endIdx < startIdx {
		startIdx, endIdx = endIdx, startIdx
	}

	// Perform crop
	s.SetSamples(s.Samples()[startIdx:endIdx])
}

// Gradient creates a linear fade-in gradient for an audio file.
// The gradient will start at 0% volume at start and 100% volume at end.
func Gradient(s Sound, start, end time.Duration) {
	if len(s.Samples()) == 0 {
		return
	}
	startIdx := sampleIndex(s, start)
	endIdx := sampleIndex(s, end)
	upwards := (startIdx < endIdx)
	if !upwards {
		startIdx, endIdx = endIdx, startIdx
	}
	for i := startIdx; i < endIdx; i++ {
		value := float64(i-startIdx) / float64(endIdx-startIdx)
		if !upwards {
			value = 1.0 - value
		}
		for j, sample := range s.Samples()[i] {
			s.Samples()[i][j] = sample * Sample(value)
		}
	}
}

func diffChannelAppend(s1, s2 Sound) {
	if s2.Channels() > s1.Channels() {
		for _, x := range s2.Samples() {
			cut := x[0:s1.Channels()]
			s1.SetSamples(append(s1.Samples(), cut))
		}
	} else {
		for _, x := range s2.Samples() {
			bigger := make([]Sample, s1.Channels())
			copy(bigger, x)
			for i := len(x); i < len(bigger); i++ {
				bigger[i] = bigger[0]
			}
			s1.SetSamples(append(s1.Samples(), bigger))
		}
	}
}

func sampleIndex(s Sound, t time.Duration) int {
	secs := float64(t) / float64(time.Second)
	index := int(secs * float64(s.SampleRate()))
	if index < 0 {
		return 0
	} else if index > len(s.Samples()) {
		return len(s.Samples())
	}
	return index
}
