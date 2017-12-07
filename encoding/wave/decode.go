package wave

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type (
	// Decoder is a WAVE decoder
	Decoder struct {
		input io.Reader
	}
)

// NewDecoder creates a new WAVE decoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		input: r,
	}
}

// Decode the WAVE into data array and returns the WAVE header
func (d *Decoder) Decode(data interface{}) (Header, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return Header{}, fmt.Errorf("Decode expects a non-nil pointer")
	}
	deref := v.Elem()
	switch deref.Kind() {
	case reflect.Slice:
		// get the underlying type of slice elements
		switch typ := deref.Type().Elem().Kind(); typ {
		case reflect.Int16, reflect.Float32:
			return d.decode(v, typ)
		default:
			return Header{}, fmt.Errorf("unknown slice type: %v",
				deref.Type().Elem().Kind())
		}
	}

	return Header{}, fmt.Errorf("impossible to decode wave into %s. "+
		"Expect *[]int16, *[]float32", v.Kind())
}

func (d *Decoder) decode(v reflect.Value, typ reflect.Kind) (Header, error) {
	hdr, err := d.parseHeader()
	if err != nil {
		return Header{}, err
	}

	switch typ {
	case reflect.Int16:
		if hdr.RiffChunkFmt.AudioFormat != FormatPCM {
			return hdr, fmt.Errorf("[]int16 requires a PCM audio")
		}
		err = d.decodeInt16(v, hdr.DataBlockSize)
	case reflect.Float32:
		if hdr.RiffChunkFmt.AudioFormat != FormatIEEEFloat {
			return hdr, fmt.Errorf("[]float32 requires a IEEEFloat format")
		}
		err = d.decodeFloat32(v, hdr.DataBlockSize)
	}
	return hdr, err
}

func (d *Decoder) decodeInt16(v reflect.Value, datasz uint32) error {
	dataptr := v.Interface().(*[]int16)

	const typesize = 2
	var bytesRead uint32
	for bytesRead < datasz {
		var buf [typesize]byte
		n, err := d.input.Read(buf[:])
		if err != nil {
			return err
		}
		if n != typesize {
			return fmt.Errorf("corrupted audio")
		}

		sample := int16(binary.LittleEndian.Uint16(buf[:]))
		*dataptr = append(*dataptr, sample)
		bytesRead += typesize
	}

	return nil
}

func (d *Decoder) decodeFloat32(v reflect.Value, datasz uint32) error {
	const maxval float32 = 1.0
	const minval float32 = -1.0
	const typesz = 4

	dataptr := v.Interface().(*[]float32)
	var bytesRead uint32
	for bytesRead < datasz {
		var sample float32
		err := binary.Read(d.input, binary.LittleEndian, &sample)
		if err != nil {
			return fmt.Errorf("decoding floats: bytes read[%d], total[%d]: %s",
				bytesRead,
				datasz,
				err)
		}

		if sample < minval || sample > maxval {
			return fmt.Errorf(
				"sample[%f] is outside the valid value range for a PCM float",
				sample,
			)
		}

		*dataptr = append(*dataptr, sample)
		bytesRead += typesz
	}

	return nil
}

func (d *Decoder) parseRIFFHdr() (RiffHeader, error) {
	var hdr RiffHeader
	err := binary.Read(d.input, binary.LittleEndian, &hdr)
	if err != nil {
		return RiffHeader{}, fmt.Errorf("parsing riff: %s", err)
	}
	if string(hdr.Ident[:]) != "RIFF" {
		return RiffHeader{}, fmt.Errorf("Invalid RIFF ident: %s", string(hdr.Ident[:]))
	}
	return hdr, nil
}

func (d *Decoder) parseHeader() (Header, error) {
	riffhdr, err := d.parseRIFFHdr()
	if err != nil {
		return Header{}, err
	}

	// FMT chunk
	var chunk [4]byte
	var chunkFmt RiffChunkFmt

	err = binary.Read(d.input, binary.LittleEndian, &chunk)
	if err != nil {
		return Header{}, fmt.Errorf("parsing fmt chunk: %s", err)
	}

	if string(chunk[:]) != "fmt " {
		return Header{}, fmt.Errorf("Unexpected chunk type: %s", string(chunk[:]))
	}

	err = binary.Read(d.input, binary.LittleEndian, &chunkFmt)
	if err != nil {
		return Header{}, fmt.Errorf("parsing fmt chunk: %s", err)
	}

	if !isValidWavFormat(chunkFmt.AudioFormat) {
		return Header{}, fmt.Errorf("Isn't an audio format: format[%d]", chunkFmt.AudioFormat)
	}

	if chunkFmt.LengthOfHeader != 16 {
		var extraparams uint16
		// Get extra params size
		if err = binary.Read(d.input, binary.LittleEndian, &extraparams); err != nil {
			return Header{}, fmt.Errorf("error getting extra fmt params: %s", err)
		}
		// TODO: Skip for now
		_, err = d.input.Read(make([]byte, int(extraparams)))
		if err != nil {
			return Header{}, fmt.Errorf("error skipping extra params: %s", err)
		}
	}

	var chunkSize uint32
	for string(chunk[:]) != "data" {
		// Read chunkID
		err = binary.Read(d.input, binary.BigEndian, &chunk)
		if err != nil {
			return Header{}, fmt.Errorf("Expected data chunkid: %s", err)
		}

		err = binary.Read(d.input, binary.LittleEndian, &chunkSize)
		if err != nil {
			return Header{}, fmt.Errorf("Expected data chunkSize: %s", err)
		}

		// ignores LIST chunkIDs (unused for now)
		if string(chunk[:]) != "data" {
			_, err = d.input.Read(make([]byte, int(chunkSize)))
			if err != nil {
				return Header{}, fmt.Errorf("ignoring LIST chunks: %s", err)
			}
		}
	}

	return Header{
		RiffHeader:    riffhdr,
		RiffChunkFmt:  chunkFmt,
		DataBlockSize: uint32(chunkSize),
	}, nil
}

func isValidWavFormat(fmt uint16) bool {
	for _, valid := range []uint16{
		FormatMULAW,
		FormatALAW,
		FormatIEEEFloat,
		FormatPCM,
	} {
		if fmt == valid {
			return true
		}
	}

	return false
}
