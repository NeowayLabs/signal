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
