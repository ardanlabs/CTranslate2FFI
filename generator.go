package ctranslate2ffi

import (
	"errors"
	"unsafe"
)

// GenerationOptions holds options for text generation.
type GenerationOptions struct {
	BeamSize              int
	Patience              float32
	LengthPenalty         float32
	RepetitionPenalty     float32
	NoRepeatNgramSize     int
	DisableUnk            bool
	MaxLength             int
	MinLength             int
	SamplingTopK          int
	SamplingTopP          float32
	SamplingTemp          float32
	NumHypotheses         int
	ReturnScores          bool
	IncludePromptInResult bool
}

// DefaultGenerationOptions returns sensible defaults.
func DefaultGenerationOptions() GenerationOptions {
	return GenerationOptions{
		BeamSize:              1,
		Patience:              1.0,
		LengthPenalty:         1.0,
		RepetitionPenalty:     1.0,
		NoRepeatNgramSize:     0,
		DisableUnk:            false,
		MaxLength:             512,
		MinLength:             0,
		SamplingTopK:          1,
		SamplingTopP:          1.0,
		SamplingTemp:          1.0,
		NumHypotheses:         1,
		ReturnScores:          false,
		IncludePromptInResult: true,
	}
}

func (o GenerationOptions) toC() Ct2generationoptions {
	return Ct2generationoptions{
		BeamSize:              uint64(o.BeamSize),
		Patience:              o.Patience,
		LengthPenalty:         o.LengthPenalty,
		RepetitionPenalty:     o.RepetitionPenalty,
		NoRepeatNgramSize:     uint64(o.NoRepeatNgramSize),
		DisableUnk:            o.DisableUnk,
		MaxLength:             uint64(o.MaxLength),
		MinLength:             uint64(o.MinLength),
		SamplingTopk:          uint64(o.SamplingTopK),
		SamplingTopp:          o.SamplingTopP,
		SamplingTemperature:   o.SamplingTemp,
		NumHypotheses:         uint64(o.NumHypotheses),
		ReturnScores:          o.ReturnScores,
		IncludePromptInResult: o.IncludePromptInResult,
	}
}

// GenerationResult holds the result of text generation.
type GenerationResult struct {
	Sequences [][]string
	Scores    []float32
}

// Generator wraps a CTranslate2 generator (language model).
type Generator struct {
	handle Ct2generator
}

// NewGenerator loads a generator model from the given path.
func NewGenerator(modelPath string, config ModelConfig) (*Generator, error) {
	handle := Ct2GeneratorCreate(modelPath, config.toC())
	if handle == 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "failed to load generator model"
		}
		return nil, errors.New(errMsg)
	}
	return &Generator{handle: handle}, nil
}

// Close releases the model resources.
func (g *Generator) Close() {
	if g.handle != 0 {
		Ct2GeneratorFree(g.handle)
		g.handle = 0
	}
}

// Generate generates text continuation from a prompt.
func (g *Generator) Generate(prompt []string, opts GenerationOptions) (*GenerationResult, error) {
	if g.handle == 0 {
		return nil, errors.New("generator is closed")
	}

	// Convert prompt to C strings
	cStrings := make([]*byte, len(prompt))
	for i, tok := range prompt {
		b := append([]byte(tok), 0)
		cStrings[i] = &b[0]
	}
	promptPtr := uintptr(unsafe.Pointer(&cStrings[0]))

	var result Ct2generationresult
	ret := Ct2GeneratorGenerate(g.handle, promptPtr, uint64(len(prompt)), opts.toC(), &result)
	if ret != 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "generation failed"
		}
		return nil, errors.New(errMsg)
	}

	gr := &GenerationResult{}
	// TODO: Properly extract sequences from the raw pointers

	Ct2GenerationResultFree(&result)
	return gr, nil
}
