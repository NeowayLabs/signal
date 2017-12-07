package wave

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

const (
	waveHdrSize = 44 // Riff header + FmtChunk
)

// Encoder of WAVE audio format
type Encoder struct {
	hdr Header
}

// NewEncoder creates a new encoder for header hdr.
// The header doesn't need to be fully initialized, you
// could use NewPCM() or NewIEEEFloat() to setup a
// header with only the required fields.
func NewEncoder(hdr Header) *Encoder {
	return &Encoder{
		hdr: hdr,
	}
}

// Encode the audio using data samples. Note that type of
// data slice is validated against the format of WAVE header.
func (e *Encoder) Encode(data interface{}) ([]byte, error) {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice:
		// get the underlying type of slice elements
		switch typ := v.Type().Elem().Kind(); typ {
		case reflect.Int16, reflect.Float32:
			return e.encode(v, typ)
		default:
			return nil, fmt.Errorf("unknown slice type: %v",
				v.Type().Elem().Kind())
		}
	}
	return nil, fmt.Errorf("impossible to encode %s", v.Type())
}

func (e *Encoder) encode(v reflect.Value, typ reflect.Kind) ([]byte, error) {
	data := v.Interface()
	if typ == reflect.Int16 {
		return e.encodePCM(data.([]int16))
	}

	return e.encodeIEEEFloat(data.([]float32))
}

func (e *Encoder) writeHeader(buf *bytes.Buffer, datasz uint32) error {
	lewrite := func(d interface{}) error {
		return binary.Write(buf, binary.LittleEndian, d)
	}

	err := lewrite(e.hdr.RiffHeader.Ident)
	if err != nil {
		return err
	}

	err = lewrite(uint32(waveHdrSize + 8 + datasz)) // 8 = ChID+CkSize, 2 = sizeof(int16)
	if err != nil {
		return err
	}

	err = lewrite(e.hdr.RiffHeader.FileType)
	if err != nil {
		return err
	}

	err = lewrite([4]byte{'f', 'm', 't', ' '})
	if err != nil {
		return err
	}

	return lewrite(e.hdr.RiffChunkFmt)
}

func (e *Encoder) encodePCM(data []int16) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, waveHdrSize))

	lewrite := func(d interface{}) error {
		return binary.Write(buf, binary.LittleEndian, d)
	}

	err := e.writeHeader(buf, uint32(len(data)*2))
	if err != nil {
		return nil, err
	}

	err = lewrite([4]byte{'d', 'a', 't', 'a'})
	if err != nil {
		return nil, err
	}

	// chunk size
	err = lewrite(uint32(2 * len(data)))
	if err != nil {
		return nil, err
	}

	// write data
	for _, d := range data {
		err := lewrite(d)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (e *Encoder) encodeIEEEFloat(data []float32) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, waveHdrSize))
	datasz := uint32(4 * len(data))

	lewrite := func(d interface{}) error {
		return binary.Write(buf, binary.LittleEndian, d)
	}

	err := e.writeHeader(buf, datasz)
	if err != nil {
		return nil, err
	}

	err = lewrite([4]byte{'d', 'a', 't', 'a'})
	if err != nil {
		return nil, err
	}

	// chunk size
	err = lewrite(datasz)
	if err != nil {
		return nil, err
	}

	// write data
	for _, d := range data {
		err := lewrite(d)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
