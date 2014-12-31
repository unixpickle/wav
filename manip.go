package wav

import "time"

// Append appends sounds to a sound.
// This method adds and removes channels and modifies sample rates as needed.
func Append(dest Sound, sounds ...Sound) {
	for _, source := range sounds {
		if dest.SampleRate() == source.SampleRate() &&
			dest.Channels() == source.Channels() {
			// This is the simple, fast case
			dest.SetSamples(append(dest.Samples(), source.Samples()...))
			continue
		}
		// Generic conversion algorithm.
		ratio := float64(dest.SampleRate()) / float64(source.SampleRate())
		sourceBlocks := len(source.Samples()) / source.Channels()
		destBlocks := int(float64(sourceBlocks) * ratio)
		mutualChannels := source.Channels()
		if dest.Channels() < mutualChannels {
			mutualChannels = dest.Channels()
		}
		for i := 0; i < destBlocks; i++ {
			sourceIdx := source.Channels() * int(float64(i)/ratio)
			newSamples := source.Samples()[sourceIdx : sourceIdx+mutualChannels]
			dest.SetSamples(append(dest.Samples(), newSamples...))
			// Duplicate the first source channel for the remaining channels in
			// the destination.
			for i := mutualChannels; i < dest.Channels(); i++ {
				dest.SetSamples(append(dest.Samples(), newSamples[0]))
			}
		}
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
// The voume is 0% at start and 100% at end.
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
	for i := startIdx; i < endIdx; i += s.Channels() {
		value := Sample(i-startIdx) / Sample(endIdx-startIdx)
		if !upwards {
			value = 1.0 - value
		}
		for j := 0; j < s.Channels(); j++ {
			s.Samples()[i+j] *= value
		}
	}
}

// Overlay overlays a sound over another sound at a certain offset.
// The overlaying sound will be converted to match the destination's sample
// rate and channel count.
func Overlay(s, o Sound, delay time.Duration) {
	// Convert the overlay if needed.
	if o.Channels() != s.Channels() || o.SampleRate() != s.SampleRate() {
		dest := NewPCM16Sound(s.Channels(), s.SampleRate())
		Append(dest, o)
		Overlay(s, dest, delay)
		return
	}

	// Figure out the length of the new sound.
	start := unclippedSampleIndex(s, delay)
	sSize := len(s.Samples())
	oSize := len(o.Samples())
	totalSize := sSize
	if start+oSize > totalSize {
		totalSize = start + oSize
	}

	// Perform the actual overlay
	for i := 0; i < totalSize; i++ {
		if i >= sSize {
			s.SetSamples(append(s.Samples(), 0))
		}
		if i >= start && i < start+oSize {
			sample := o.Samples()[i-start]
			s.Samples()[i] = clamp(s.Samples()[i] + sample)
		}
	}
}

// Volume scales all the samples in a Sound.
func Volume(s Sound, scale float64) {
	sScale := Sample(scale)
	for i, sample := range s.Samples() {
		s.Samples()[i] = clamp(sample * sScale)
	}
}

func clamp(s Sample) Sample {
	if s < -1.0 {
		return -1
	} else if s > 1.0 {
		return 1
	}
	return s
}

func sampleIndex(s Sound, t time.Duration) int {
	index := unclippedSampleIndex(s, t)
	if index > len(s.Samples()) {
		return len(s.Samples())
	}
	return index
}

func unclippedSampleIndex(s Sound, t time.Duration) int {
	secs := float64(t) / float64(time.Second)
	index := int(secs * float64(s.SampleRate())) * s.Channels()
	if index < 0 {
		return 0
	}
	return index
}
