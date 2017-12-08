package wave

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

type (
	// Decoder is a WAVE decoder
	Decoder struct {
		input     io.Reader        // input stream
		byteOrder binary.ByteOrder // decoder's byte order for data samples
	}
)

// NewDecoder creates a new WAVE decoder using Little-Endian as default
// byte order of data samples (use d.BigEndian() to opt for big-endian).
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		input:     r,
		byteOrder: binary.LittleEndian,
	}
}

// LittleEndian configures to decode data chunk using little-endian.
func (d *Decoder) LittleEndian() {
	d.byteOrder = binary.LittleEndian
}

// BigEndian configures to decode data chunk using big-endian.
func (d *Decoder) BigEndian() {
	d.byteOrder = binary.BigEndian
}

// DecodeInt16 decodes the WAV buffer, returning the wave header and
// filling data with the audio samples.
// In case the data chunk is corrupted or there's some other error
// parsing the samples, the parsed header is returned to inspection also
// (useful to check corrupted WAV files).
func (d *Decoder) DecodeInt16(data *[]int16) (hdr Header, err error) {
	hdr, err = d.DecodeHeader()
	if err != nil {
		return Header{}, err
	}

	var bytesRead uint32
	for bytesRead < hdr.DataBlockSize {
		var sample int16
		err := binary.Read(d.input, d.byteOrder, &sample)
		if err != nil {
			return hdr, fmt.Errorf("decoding int16: bytes read[%d], total[%d]: %s",
				bytesRead,
				hdr.DataBlockSize,
				err)
		}
		*data = append(*data, sample)
		bytesRead += uint32(unsafe.Sizeof(sample))
	}

	return hdr, nil
}

// DecodeFloat32 decodes the WAV buffer, returning the wave header and
// filling data with the audio samples.
// In case the data chunk is corrupted or there's some other error
// parsing the samples, the parsed header is returned to inspection also
// (useful to check corrupted WAV files).
func (d *Decoder) DecodeFloat32(data *[]float32) (hdr Header, err error) {
	hdr, err = d.DecodeHeader()
	if err != nil {
		return Header{}, err
	}

	const maxval float32 = 1.0
	const minval float32 = -1.0

	var bytesRead uint32
	for bytesRead < hdr.DataBlockSize {
		var sample float32
		err := binary.Read(d.input, d.byteOrder, &sample)
		if err != nil {
			return hdr, fmt.Errorf("decoding float32: bytes read[%d], total[%d]: %s",
				bytesRead,
				hdr.DataBlockSize,
				err)
		}

		if sample < minval || sample > maxval {
			return hdr, fmt.Errorf(
				"sample[%f] is outside the valid value range for a PCM float",
				sample,
			)
		}

		*data = append(*data, sample)
		bytesRead += uint32(unsafe.Sizeof(sample))
	}

	return hdr, nil
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

// DecodeHeader decodes just the header of the WAV.
func (d *Decoder) DecodeHeader() (Header, error) {
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
