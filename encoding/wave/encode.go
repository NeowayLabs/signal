package wave

import (
	"bytes"
	"encoding/binary"
)

const (
	waveHdrSize = 44 // Riff header + FmtChunk
)

// Encoder of WAVE audio format
type Encoder struct {
	hdr       Header           // hdr of output WAV
	byteOrder binary.ByteOrder // encoder's byte order for data samples
}

// NewEncoder creates a new encoder for header hdr.
// The header doesn't need to be fully initialized, you
// could use NewPCM() or NewIEEEFloat() to setup a
// header with only the required fields.
func NewEncoder(hdr Header) *Encoder {
	return &Encoder{
		hdr:       hdr,
		byteOrder: binary.LittleEndian,
	}
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

func (e *Encoder) EncodeInt16(data []int16) ([]byte, error) {
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

func (e *Encoder) EncodeFloat32(data []float32) ([]byte, error) {
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
