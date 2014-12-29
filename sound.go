package wav

type Sound struct {
	SampleRate int32
	Samples    [][]int
}

func ReadSound(path string) *Sound {
	return nil
}
