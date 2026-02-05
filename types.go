package ctranslate2ffi

import "github.com/jupiterrider/ffi"

type Ct2storageview uintptr

type Ct2whisper uintptr

type Ct2translator uintptr

type Ct2generator uintptr

type Ct2modelconfig struct {
	Device      Ct2device
	ComputeType Ct2computetype
	DeviceIndex int32
	NumThreads  uint64
	NumReplicas uint64
}

var FFITypeCt2modelconfig = ffi.NewType(
	&ffi.TypePointer,
	&ffi.TypePointer,
	&ffi.TypeSint32,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
)

type Ct2stringarray struct {
	Strings string
	Count   uint64
}

var FFITypeCt2stringarray = ffi.NewType(
	&ffi.TypePointer,
	&ffi.TypeUint64,
)

type Ct2floatarray struct {
	Values uintptr
	Count  uint64
}

var FFITypeCt2floatarray = ffi.NewType(
	&ffi.TypePointer,
	&ffi.TypeUint64,
)

type Ct2whisperoptions struct {
	BeamSize                 uint64
	Patience                 float32
	LengthPenalty            float32
	RepetitionPenalty        float32
	NoRepeatNgramSize        uint64
	MaxLength                uint64
	SamplingTopk             uint64
	SamplingTemperature      float32
	NumHypotheses            uint64
	ReturnScores             bool
	ReturnNoSpeechProb       bool
	MaxInitialTimestampIndex uint64
	SuppressBlank            bool
}

var FFITypeCt2whisperoptions = ffi.NewType(
	&ffi.TypeUint64,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeFloat,
	&ffi.TypeUint64,
	&ffi.TypeUint8,
	&ffi.TypeUint8,
	&ffi.TypeUint64,
	&ffi.TypeUint8,
)

type Ct2whisperresult struct {
	Sequences    string
	NumSequences uint64
	Scores       uintptr
	NumScores    uint64
	NoSpeechProb float32
}

var FFITypeCt2whisperresult = ffi.NewType(
	&ffi.TypePointer,
	&ffi.TypeUint64,
	&ffi.TypePointer,
	&ffi.TypeUint64,
	&ffi.TypeFloat,
)

type Ct2translationoptions struct {
	BeamSize            uint64
	Patience            float32
	LengthPenalty       float32
	CoveragePenalty     float32
	RepetitionPenalty   float32
	NoRepeatNgramSize   uint64
	DisableUnk          bool
	MaxInputLength      uint64
	MaxDecodingLength   uint64
	MinDecodingLength   uint64
	SamplingTopk        uint64
	SamplingTopp        float32
	SamplingTemperature float32
	UseVmap             bool
	NumHypotheses       uint64
	ReturnScores        bool
	ReturnAttention     bool
	ReplaceUnknowns     bool
}

var FFITypeCt2translationoptions = ffi.NewType(
	&ffi.TypeUint64,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeUint64,
	&ffi.TypeUint8,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeUint8,
	&ffi.TypeUint64,
	&ffi.TypeUint8,
	&ffi.TypeUint8,
	&ffi.TypeUint8,
)

type Ct2translationresult struct {
	Hypotheses        string
	HypothesesLengths uintptr
	NumHypotheses     uint64
	Scores            uintptr
	NumScores         uint64
}

var FFITypeCt2translationresult = ffi.NewType(
	&ffi.TypePointer,
	&ffi.TypePointer,
	&ffi.TypeUint64,
	&ffi.TypePointer,
	&ffi.TypeUint64,
)

type Ct2generationoptions struct {
	BeamSize              uint64
	Patience              float32
	LengthPenalty         float32
	RepetitionPenalty     float32
	NoRepeatNgramSize     uint64
	DisableUnk            bool
	MaxLength             uint64
	MinLength             uint64
	SamplingTopk          uint64
	SamplingTopp          float32
	SamplingTemperature   float32
	NumHypotheses         uint64
	ReturnScores          bool
	IncludePromptInResult bool
}

var FFITypeCt2generationoptions = ffi.NewType(
	&ffi.TypeUint64,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeUint64,
	&ffi.TypeUint8,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeUint64,
	&ffi.TypeFloat,
	&ffi.TypeFloat,
	&ffi.TypeUint64,
	&ffi.TypeUint8,
	&ffi.TypeUint8,
)

type Ct2generationresult struct {
	Sequences       string
	SequenceLengths uintptr
	NumSequences    uint64
	Scores          uintptr
	NumScores       uint64
}

var FFITypeCt2generationresult = ffi.NewType(
	&ffi.TypePointer,
	&ffi.TypePointer,
	&ffi.TypeUint64,
	&ffi.TypePointer,
	&ffi.TypeUint64,
)

type Ct2device int32

const (
	Ct2DeviceCpu  Ct2device = 0
	Ct2DeviceCuda Ct2device = 1
)

type Ct2computetype int32

const (
	Ct2ComputeDefault      Ct2computetype = 0
	Ct2ComputeAuto         Ct2computetype = 1
	Ct2ComputeFloat32      Ct2computetype = 2
	Ct2ComputeInt8         Ct2computetype = 3
	Ct2ComputeInt8Float32  Ct2computetype = 4
	Ct2ComputeInt8Float16  Ct2computetype = 5
	Ct2ComputeInt8Bfloat16 Ct2computetype = 6
	Ct2ComputeInt16        Ct2computetype = 7
	Ct2ComputeFloat16      Ct2computetype = 8
	Ct2ComputeBfloat16     Ct2computetype = 9
)
