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
	sampleCount := len(s.Samples())
	if sampleCount == 0 {
		return
	}

	// Figure out sample indexes.
	startSecs := float64(start) / float64(time.Second)
	endSecs := float64(end) / float64(time.Second)
	startIdx := int(startSecs * float64(s.SampleRate()))
	endIdx := int(endSecs * float64(s.SampleRate()))

	// Clamp indexes
	if endIdx < startIdx {
		startIdx, endIdx = endIdx, startIdx
	}
	if startIdx < 0 {
		startIdx = 0
	} else if startIdx >= sampleCount {
		startIdx = sampleCount - 1
	}
	if endIdx < 0 {
		endIdx = 0
	} else if endIdx > sampleCount {
		endIdx = sampleCount
	}

	// Perform crop
	s.SetSamples(s.Samples()[startIdx:endIdx])
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
