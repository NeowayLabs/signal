package wave_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/NeowayLabs/signal/encoding/wave"
)

func TestEncoder(t *testing.T) {
	f, err := os.Open("testdata/audios/sint16le.wav")
	assertNoError(t, err)

	defer f.Close()

	expectedData := make([]int16, 0)
	d := wave.NewDecoder(f)
	expectedHdr, err := d.DecodeInt16(&expectedData)
	assertNoError(t, err)

	enc := wave.NewEncoder(
		wave.NewPCM(
			int(expectedHdr.RiffChunkFmt.NumChannels),
			int(expectedHdr.RiffChunkFmt.SampleRate),
			int(expectedHdr.RiffChunkFmt.BitsPerSample)),
	)

	// encode a new audio from expected data
	audioBytes, err := enc.EncodeInt16(expectedData)
	assertNoError(t, err)

	d = wave.NewDecoder(bytes.NewBuffer(audioBytes))
	gotData := make([]int16, 0)
	_, err = d.DecodeInt16(&gotData)
	assertNoError(t, err)

	if !reflect.DeepEqual(gotData, expectedData) {
		t.Fatalf("generated wave differs")
	}

	// TODO(i4k): arrange a clean WAVE (header + data) to make the
	// comparisons.
	// if !reflect.DeepEqual(gotHdr, expectedHdr) {
	// 	t.Fatalf("headers differs: %v != %v", gotHdr, expectedHdr)
	// }
}
