package wav

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkReader(b *testing.B) {
	sound := NewPCM16Sound(2, 22050)
	sound.SetSamples(make([]Sample, 22050*2*10))
	var buf bytes.Buffer
	if err := sound.Write(&buf); err != nil {
		b.Fatal(err)
	}

	reader := bytes.NewReader(buf.Bytes())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		ReadSound(reader)
	}
}
