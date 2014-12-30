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

// Convert converts a sound to meet another sound's parameters.
// Both the channel count and the sample rate will be converted.
// The converted source samples will be appended to the destination.
func Convert(dest, source Sound) {
	ratio := float64(dest.SampleRate()) / float64(source.SampleRate())
	numSamples := int(ratio * float64(len(source.Samples())))
	for i := 0; i < numSamples; i++ {
		sourceSample := int(float64(i) / ratio)
		newRes := diffChannel(source.Samples()[sourceSample], dest.Channels())
		dest.SetSamples(append(dest.Samples(), newRes))
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

// Overlay overlays a sound over another sound at a certain offset.
// If the sounds are a different number of channels, only some channels will
// be overlayed.
func Overlay(s, o Sound, delay time.Duration) {
	start := sampleIndex(s, delay)
	sSize := len(s.Samples())
	oSize := len(o.Samples())
	totalSize := sSize
	if start+oSize > totalSize {
		totalSize = start + oSize
	}
	for i := 0; i < totalSize; i++ {
		if i >= sSize {
			zeroes := make([]Sample, s.Channels())
			s.SetSamples(append(s.Samples(), zeroes))
		}
		if i >= start && i < start+oSize {
			add := diffChannel(o.Samples()[i-start], s.Channels())
			for j, sample := range add {
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

func diffChannel(packet []Sample, channels int) []Sample {
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

func diffChannelAppend(s1, s2 Sound) {
	for _, x := range s2.Samples() {
		fixed := diffChannel(x, s1.Channels())
		s1.SetSamples(append(s1.Samples(), fixed))
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
