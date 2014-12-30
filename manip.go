package wav

import "time"

// Append appends sounds to a sound.
// This method adds and removes channels and modifies sample rates as needed.
func Append(s Sound, sounds ...Sound) {
	for _, x := range sounds {
		if s.SampleRate() == x.SampleRate() && s.Channels() == x.Channels() {
			// This is the simple, fast case
			s.SetSamples(append(s.Samples(), x.Samples()...))
			continue
		}
		// Generic conversion algorithm.
		ratio := float64(s.SampleRate()) / float64(x.SampleRate())
		numSamples := int(ratio * float64(len(x.Samples())))
		for i := 0; i < numSamples; i++ {
			source := int(float64(i) / ratio)
			newRes := convPacket(x.Samples()[source], s.Channels())
			s.SetSamples(append(s.Samples(), newRes))
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
	start := sampleIndex(s, delay)
	sSize := len(s.Samples())
	oSize := len(o.Samples())
	totalSize := sSize
	if start+oSize > totalSize {
		totalSize = start + oSize
	}

	// Perform the actual overlay
	for i := 0; i < totalSize; i++ {
		if i >= sSize {
			zeroes := make([]Sample, s.Channels())
			s.SetSamples(append(s.Samples(), zeroes))
		}
		if i >= start && i < start+oSize {
			for j, sample := range o.Samples()[i-start] {
				s.Samples()[i][j] = clamp(s.Samples()[i][j] + sample)
			}
		}
	}
}

// Volume scales all the samples in a Sound.
func Volume(s Sound, scale float64) {
	for _, x := range s.Samples() {
		for i, sample := range x {
			x[i] = clamp(sample * Sample(scale))
		}
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

func convPacket(packet []Sample, channels int) []Sample {
	if len(packet) == channels {
		return packet
	} else if len(packet) > channels {
		return packet[0:channels]
	}
	bigger := make([]Sample, channels)
	copy(bigger, packet)
	for i := len(packet); i < channels; i++ {
		bigger[i] = bigger[0]
	}
	return bigger
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
