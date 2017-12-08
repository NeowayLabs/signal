package wave_test

import (
	"bytes"
	"fmt"

	"github.com/NeowayLabs/signal/encoding/wave"
)

func ExampleDecodeHeader() {
	samples := []int16{
		-255, 0, 255, 0, -255, 0, 255,
	}
	enc := wave.NewEncoder(wave.NewPCM(1, 8000, 16))
	audio, err := enc.EncodeInt16(samples)
	if err != nil {
		panic(err)
	}
	hdr, err := wave.DecodeHeader(bytes.NewReader(audio))
	if err != nil {
		panic(err)
	}

	fmt.Printf("SampleRate: %dhz\n", hdr.SampleRate)
	fmt.Printf("Number of channels: %d\n", hdr.NumChannels)
	fmt.Printf("Bytes/block: %d\n", hdr.BytesPerBloc)

	// Output: SampleRate: 8000hz
	// Number of channels: 1
	// Bytes/block: 2
}
