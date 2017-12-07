package main

import (
	"flag"
	"io/ioutil"
	"math"

	"github.com/NeowayLabs/signal/encoding/wave"
)

var (
	filename   string
	sampleRate uint
	volume     uint
	frequency  float64
)

func init() {
	flag.StringVar(&filename, "output", "out.wav", "output file")
	flag.UintVar(&sampleRate, "samplerate", 8000,
		"The sample rate. Eg.: 8000, 44100, 48000, 96000, etc")
	flag.UintVar(&volume, "volume", 28000, "volume range from 0 to 32767")
	flag.Float64Var(&frequency, "frequency", 440.0, "sine wave frequency")
}

func main() {
	flag.Parse()
	var freqRadPerSample = frequency * 2 * math.Pi / float64(sampleRate)
	var phase float64
	enc := wave.NewEncoder(wave.NewPCM(1, int(sampleRate), 16))

	var data []int16

	// TODO(i4k) add option to duration/numberOfSamples
	for i := 0; i < 10000; i++ {
		phase += freqRadPerSample
		sample := float64(volume) * math.Sin(phase)
		data = append(data, int16(sample))
	}
	audioBytes, err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, audioBytes, 0664)
	if err != nil {
		panic(err)
	}
}
