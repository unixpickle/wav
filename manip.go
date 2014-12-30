package wav

import "time"

func Crop(s Sound, start time.Duration, end time.Duration) {
	// Cannot crop an empty sound.
	sampleCount := len(s.Samples())
	if sampleCount == 0 {
		return
	}

	// Figure out sample indexes.
	startSecs := float64(start) / float64(time.Second)
	endSecs := float64(end) / float64(time.Second)
	h := s.Header()
	startIdx := int(startSecs * float64(h.Format.SampleRate))
	endIdx := int(endSecs * float64(h.Format.SampleRate))

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
