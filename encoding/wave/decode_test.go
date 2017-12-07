package wave_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/NeowayLabs/signal/encoding/wave"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func fixHdr(hdr *wave.Header) {
	hdr.RiffHeader.Ident[0] = 'R'
	hdr.RiffHeader.Ident[1] = 'I'
	hdr.RiffHeader.Ident[2] = 'F'
	hdr.RiffHeader.Ident[3] = 'F'
	hdr.RiffHeader.FileType[0] = 'W'
	hdr.RiffHeader.FileType[1] = 'A'
	hdr.RiffHeader.FileType[2] = 'V'
	hdr.RiffHeader.FileType[3] = 'E'
}

func testParseWAV(t *testing.T, filename string) {
	r, err := os.Open(filename)
	assertNoError(t, err)

	d := wave.NewDecoder(r)
	data := make([]int16, 0)
	hdr, hdrerr := d.Decode(&data)

	ext := filepath.Ext(filename)
	noext := strings.TrimSuffix(filename, ext)
	expectedHdrFile := noext + ".hdr.expected"
	errFile := noext + ".err"

	expectedErr, err := ioutil.ReadFile(errFile)
	if err == nil {
		if hdrerr == nil {
			t.Fatalf("Expected error: %s but ran successfully...", string(expectedErr))
		}
		if hdrerr.Error() != string(expectedErr) {
			t.Fatalf("Error differs: '%s' != '%s'", hdrerr, string(expectedErr))
		}
		return
	} else if hdrerr != nil {
		t.Fatalf("Error: %s", hdrerr)
	}

	expectedHdrContent, err := ioutil.ReadFile(expectedHdrFile)
	if err == nil {
		var expected wave.Header
		err = json.Unmarshal(expectedHdrContent, &expected)
		assertNoError(t, err)

		fixHdr(&expected)

		if !reflect.DeepEqual(hdr, expected) {
			t.Fatalf("WAV header differs:\n\n%#v\n\n!=\n\n%#v\n", hdr, expected)
		}
		return
	}

	t.Fatalf("no error file nor expected file found for input: %s", filename)
}

func TestParseWAVHeaders(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	assertNoError(t, err)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "wav") {
			fname := file.Name()
			t.Run(fmt.Sprintf("header-%s", fname), func(t *testing.T) {
				testParseWAV(t, filepath.Join("testdata", fname))
			})
		}
	}
}

func TestSignedInt16LittleEndianSamples(t *testing.T) {
	f, err := os.Open("testdata/audios/sint16le.wav")
	assertNoError(t, err)

	d := wave.NewDecoder(f)
	samples := make([]int16, 0)
	_, err = d.Decode(&samples)

	assertNoError(t, err)

	gotbuf := &bytes.Buffer{}
	err = binary.Write(gotbuf, binary.LittleEndian, samples)
	assertNoError(t, err)

	expected, err := ioutil.ReadFile("testdata/audios/sint16le.raw")
	assertNoError(t, err)

	assertBytesEqual(t, expected, gotbuf.Bytes())
}

func TestFloat32LittleEndianSamples(t *testing.T) {
	f, err := os.Open("testdata/audios/float32le.wav")
	d := wave.NewDecoder(f)

	samples := make([]float32, 0)
	_, err = d.Decode(&samples)
	assertNoError(t, err)

	gotbuf := &bytes.Buffer{}
	err = binary.Write(gotbuf, binary.LittleEndian, samples)
	assertNoError(t, err)

	expected, err := ioutil.ReadFile("testdata/audios/float32le.raw")
	assertNoError(t, err)

	assertBytesEqual(t, expected, gotbuf.Bytes())
}

func assertBytesEqual(t *testing.T, expected []byte, got []byte) {
	if len(expected) != len(got) {
		t.Fatalf("expected len[%d] != got len[%d]", len(expected), len(got))
	}

	for i, expectedByte := range expected {
		gotByte := got[i]
		if expectedByte != gotByte {
			t.Fatalf("got wrong byte at index[%d] expected[%d] got[%d]", i, expectedByte, gotByte)
		}
	}
}
