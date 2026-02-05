package ctranslate2ffi

import (
	"errors"
	"unsafe"
)

// Device represents the compute device.
type Device int

const (
	DeviceCPU  Device = 0
	DeviceCUDA Device = 1
)

// ComputeType represents the compute precision.
type ComputeType int

const (
	ComputeDefault     ComputeType = 0
	ComputeAuto        ComputeType = 1
	ComputeFloat32     ComputeType = 2
	ComputeInt8        ComputeType = 3
	ComputeInt8Float32 ComputeType = 4
	ComputeInt8Float16 ComputeType = 5
	ComputeInt16       ComputeType = 7
	ComputeFloat16     ComputeType = 8
	ComputeBFloat16    ComputeType = 9
)

// ModelConfig holds configuration for loading models.
type ModelConfig struct {
	Device      Device
	ComputeType ComputeType
	DeviceIndex int
	NumThreads  int
	NumReplicas int
}

// DefaultModelConfig returns sensible default configuration.
func DefaultModelConfig() ModelConfig {
	return ModelConfig{
		Device:      DeviceCPU,
		ComputeType: ComputeDefault,
		DeviceIndex: 0,
		NumThreads:  0, // auto
		NumReplicas: 1,
	}
}

func (c ModelConfig) toC() Ct2modelconfig {
	return Ct2modelconfig{
		Device:      Ct2device(c.Device),
		ComputeType: Ct2computetype(c.ComputeType),
		DeviceIndex: int32(c.DeviceIndex),
		NumThreads:  uint64(c.NumThreads),
		NumReplicas: uint64(c.NumReplicas),
	}
}

// WhisperOptions holds generation options for Whisper.
type WhisperOptions struct {
	BeamSize                 int
	Patience                 float32
	LengthPenalty            float32
	RepetitionPenalty        float32
	NoRepeatNgramSize        int
	MaxLength                int
	SamplingTopK             int
	SamplingTemperature      float32
	NumHypotheses            int
	ReturnScores             bool
	ReturnNoSpeechProb       bool
	MaxInitialTimestampIndex int
	SuppressBlank            bool
}

// DefaultWhisperOptions returns sensible default options.
func DefaultWhisperOptions() WhisperOptions {
	return WhisperOptions{
		BeamSize:                 5,
		Patience:                 1.0,
		LengthPenalty:            1.0,
		RepetitionPenalty:        1.0,
		NoRepeatNgramSize:        0,
		MaxLength:                448,
		SamplingTopK:             1,
		SamplingTemperature:      1.0,
		NumHypotheses:            1,
		ReturnScores:             false,
		ReturnNoSpeechProb:       false,
		MaxInitialTimestampIndex: 50,
		SuppressBlank:            true,
	}
}

func (o WhisperOptions) toC() Ct2whisperoptions {
	return Ct2whisperoptions{
		BeamSize:                 uint64(o.BeamSize),
		Patience:                 o.Patience,
		LengthPenalty:            o.LengthPenalty,
		RepetitionPenalty:        o.RepetitionPenalty,
		NoRepeatNgramSize:        uint64(o.NoRepeatNgramSize),
		MaxLength:                uint64(o.MaxLength),
		SamplingTopk:             uint64(o.SamplingTopK),
		SamplingTemperature:      o.SamplingTemperature,
		NumHypotheses:            uint64(o.NumHypotheses),
		ReturnScores:             o.ReturnScores,
		ReturnNoSpeechProb:       o.ReturnNoSpeechProb,
		MaxInitialTimestampIndex: uint64(o.MaxInitialTimestampIndex),
		SuppressBlank:            o.SuppressBlank,
	}
}

// WhisperResult holds the result of Whisper transcription.
type WhisperResult struct {
	Sequences    []string
	Scores       []float32
	NoSpeechProb float32
}

// Whisper wraps a CTranslate2 Whisper model.
type Whisper struct {
	handle Ct2whisper
}

// NewWhisper loads a Whisper model from the given path.
func NewWhisper(modelPath string, config ModelConfig) (*Whisper, error) {
	handle := Ct2WhisperCreate(modelPath, config.toC())
	if handle == 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "failed to load whisper model"
		}
		return nil, errors.New(errMsg)
	}
	return &Whisper{handle: handle}, nil
}

// Close releases the model resources.
func (w *Whisper) Close() {
	if w.handle != 0 {
		Ct2WhisperFree(w.handle)
		w.handle = 0
	}
}

// IsMultilingual returns whether the model supports multiple languages.
func (w *Whisper) IsMultilingual() bool {
	return Ct2WhisperIsMultilingual(w.handle)
}

// NumMels returns the number of mel filterbanks expected by the model.
func (w *Whisper) NumMels() int {
	return int(Ct2WhisperNMels(w.handle))
}

// NumLanguages returns the number of supported languages.
func (w *Whisper) NumLanguages() int {
	return int(Ct2WhisperNumLanguages(w.handle))
}

// Generate transcribes audio features.
func (w *Whisper) Generate(features *StorageView, prompts []string, opts WhisperOptions) (*WhisperResult, error) {
	if w.handle == 0 {
		return nil, errors.New("whisper model is closed")
	}

	// Prepare prompts as C strings
	var promptsPtr uintptr
	numPrompts := uint64(len(prompts))
	if len(prompts) > 0 {
		cStrings := make([]*byte, len(prompts))
		for i, p := range prompts {
			b := append([]byte(p), 0)
			cStrings[i] = &b[0]
		}
		promptsPtr = uintptr(unsafe.Pointer(&cStrings[0]))
	}

	var result Ct2whisperresult
	ret := Ct2WhisperGenerate(w.handle, features.handle, promptsPtr, numPrompts, opts.toC(), &result)
	if ret != 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "whisper generation failed"
		}
		return nil, errors.New(errMsg)
	}

	// Convert result - note: the C code joins tokens, so we get one string per sequence
	wr := &WhisperResult{
		NoSpeechProb: result.NoSpeechProb,
	}

	// TODO: Properly extract sequences and scores from the raw pointers
	// This requires understanding the memory layout

	Ct2WhisperResultFree(&result)
	return wr, nil
}

// StorageView wraps a CTranslate2 storage view (tensor).
type StorageView struct {
	handle Ct2storageview
}

// NewStorageViewFloat creates a storage view from float data.
func NewStorageViewFloat(data []float32, shape []int64, device Device) (*StorageView, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	dataPtr := uintptr(unsafe.Pointer(&data[0]))
	shapePtr := uintptr(unsafe.Pointer(&shape[0]))

	handle := Ct2StorageCreateFloat(dataPtr, shapePtr, uint64(len(shape)), Ct2device(device))
	if handle == 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "failed to create storage view"
		}
		return nil, errors.New(errMsg)
	}

	return &StorageView{handle: handle}, nil
}

// Close releases the storage view.
func (s *StorageView) Close() {
	if s.handle != 0 {
		Ct2StorageFree(s.handle)
		s.handle = 0
	}
}

// Size returns the total number of elements.
func (s *StorageView) Size() int64 {
	return Ct2StorageSize(s.handle)
}

// Shape returns the shape of the storage view.
func (s *StorageView) Shape() ([]int64, error) {
	shape := make([]int64, 8) // max 8 dimensions
	var ndims uint64

	ret := Ct2StorageGetShape(s.handle, uintptr(unsafe.Pointer(&shape[0])), uintptr(unsafe.Pointer(&ndims)))
	if ret != 0 {
		return nil, errors.New(Ct2GetLastError())
	}

	return shape[:ndims], nil
}

// ToFloat copies the data to a float slice.
func (s *StorageView) ToFloat() ([]float32, error) {
	size := s.Size()
	if size <= 0 {
		return nil, errors.New("empty storage view")
	}

	data := make([]float32, size)
	ret := Ct2StorageToFloat(s.handle, uintptr(unsafe.Pointer(&data[0])))
	if ret != 0 {
		return nil, errors.New(Ct2GetLastError())
	}

	return data, nil
}

// Version returns the library version.
func Version() string {
	return Ct2Version()
}

// CUDAAvailable returns whether CUDA is available.
func CUDAAvailable() bool {
	return Ct2CudaAvailable()
}

// CUDADeviceCount returns the number of CUDA devices.
func CUDADeviceCount() int {
	return int(Ct2CudaDeviceCount())
}

// GetLastError returns the last error message.
func GetLastError() string {
	return Ct2GetLastError()
}

// ClearError clears the last error.
func ClearError() {
	Ct2ClearError()
}
