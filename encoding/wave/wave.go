package wave

type (
	// Header of Wave
	Header struct {
		RiffHeader   RiffHeader
		RiffChunkFmt RiffChunkFmt

		DataBlockSize uint32 // size of sample data (PCM data)
	}

	// RiffHeader is the header of RIFF
	RiffHeader struct {
		Ident     [4]byte // RIFF
		ChunkSize uint32
		FileType  [4]byte // WAVE
	}

	RiffChunkFmt struct {
		LengthOfHeader uint32
		AudioFormat    uint16
		NumChannels    uint16
		SampleRate     uint32
		BytesPerSec    uint32
		BytesPerBloc   uint16
		BitsPerSample  uint16
	}
)

const (
	FormatPCM        = 0x0001
	FormatIEEEFloat  = 0x0003
	FormatALAW       = 0x0006
	FormatMULAW      = 0x0007
	FormatExtensible = 0xFFFE
)

// NewPCM creates a new PCM wave header
func NewPCM(nchannels, samplerate, bits int) Header {
	return Header{
		RiffHeader: waveRiff(),
		RiffChunkFmt: RiffChunkFmt{
			LengthOfHeader: 16,
			AudioFormat:    FormatPCM,
			NumChannels:    uint16(nchannels),
			SampleRate:     uint32(samplerate),
			BytesPerSec:    uint32(bits/8) * uint32(samplerate),
			BytesPerBloc:   uint16(bits / 8),
			BitsPerSample:  uint16(bits),
		},
	}
}

// NewIEEEFloat creates a new WAVE storing IEEE float data
func NewIEEEFloat(nchannels, samplerate, bits int) Header {
	return Header{
		RiffHeader: waveRiff(),
		RiffChunkFmt: RiffChunkFmt{
			LengthOfHeader: 16,
			AudioFormat:    FormatIEEEFloat,
			NumChannels:    uint16(nchannels),
			SampleRate:     uint32(samplerate),
			BytesPerSec:    uint32(bits/8) * uint32(samplerate),
			BytesPerBloc:   uint16(bits / 8),
			BitsPerSample:  uint16(bits),
		},
	}
}

func waveRiff() RiffHeader {
	return RiffHeader{
		Ident:    [4]byte{'R', 'I', 'F', 'F'},
		FileType: [4]byte{'W', 'A', 'V', 'E'},
	}
}
