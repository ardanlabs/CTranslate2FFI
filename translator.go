package ctranslate2ffi

import (
	"errors"
	"unsafe"
)

// TranslationOptions holds options for translation.
type TranslationOptions struct {
	BeamSize          int
	Patience          float32
	LengthPenalty     float32
	CoveragePenalty   float32
	RepetitionPenalty float32
	NoRepeatNgramSize int
	DisableUnk        bool
	MaxInputLength    int
	MaxDecodingLength int
	MinDecodingLength int
	SamplingTopK      int
	SamplingTopP      float32
	SamplingTemp      float32
	UseVMap           bool
	NumHypotheses     int
	ReturnScores      bool
	ReturnAttention   bool
	ReplaceUnknowns   bool
}

// DefaultTranslationOptions returns sensible defaults.
func DefaultTranslationOptions() TranslationOptions {
	return TranslationOptions{
		BeamSize:          2,
		Patience:          1.0,
		LengthPenalty:     1.0,
		CoveragePenalty:   0.0,
		RepetitionPenalty: 1.0,
		NoRepeatNgramSize: 0,
		DisableUnk:        false,
		MaxInputLength:    1024,
		MaxDecodingLength: 256,
		MinDecodingLength: 1,
		SamplingTopK:      1,
		SamplingTopP:      1.0,
		SamplingTemp:      1.0,
		UseVMap:           false,
		NumHypotheses:     1,
		ReturnScores:      false,
		ReturnAttention:   false,
		ReplaceUnknowns:   false,
	}
}

func (o TranslationOptions) toC() Ct2translationoptions {
	return Ct2translationoptions{
		BeamSize:            uint64(o.BeamSize),
		Patience:            o.Patience,
		LengthPenalty:       o.LengthPenalty,
		CoveragePenalty:     o.CoveragePenalty,
		RepetitionPenalty:   o.RepetitionPenalty,
		NoRepeatNgramSize:   uint64(o.NoRepeatNgramSize),
		DisableUnk:          o.DisableUnk,
		MaxInputLength:      uint64(o.MaxInputLength),
		MaxDecodingLength:   uint64(o.MaxDecodingLength),
		MinDecodingLength:   uint64(o.MinDecodingLength),
		SamplingTopk:        uint64(o.SamplingTopK),
		SamplingTopp:        o.SamplingTopP,
		SamplingTemperature: o.SamplingTemp,
		UseVmap:             o.UseVMap,
		NumHypotheses:       uint64(o.NumHypotheses),
		ReturnScores:        o.ReturnScores,
		ReturnAttention:     o.ReturnAttention,
		ReplaceUnknowns:     o.ReplaceUnknowns,
	}
}

// TranslationResult holds the result of translation.
type TranslationResult struct {
	Hypotheses [][]string
	Scores     []float32
}

// Translator wraps a CTranslate2 translator model.
type Translator struct {
	handle Ct2translator
}

// NewTranslator loads a translator model from the given path.
func NewTranslator(modelPath string, config ModelConfig) (*Translator, error) {
	handle := Ct2TranslatorCreate(modelPath, config.toC())
	if handle == 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "failed to load translator model"
		}
		return nil, errors.New(errMsg)
	}
	return &Translator{handle: handle}, nil
}

// Close releases the model resources.
func (t *Translator) Close() {
	if t.handle != 0 {
		Ct2TranslatorFree(t.handle)
		t.handle = 0
	}
}

// Translate translates a single sequence of tokens.
func (t *Translator) Translate(tokens []string, opts TranslationOptions) (*TranslationResult, error) {
	if t.handle == 0 {
		return nil, errors.New("translator is closed")
	}

	// Convert tokens to C strings
	cStrings := make([]*byte, len(tokens))
	for i, tok := range tokens {
		b := append([]byte(tok), 0)
		cStrings[i] = &b[0]
	}
	tokensPtr := uintptr(unsafe.Pointer(&cStrings[0]))

	var result Ct2translationresult
	ret := Ct2TranslatorTranslate(t.handle, tokensPtr, uint64(len(tokens)), opts.toC(), &result)
	if ret != 0 {
		errMsg := Ct2GetLastError()
		if errMsg == "" {
			errMsg = "translation failed"
		}
		return nil, errors.New(errMsg)
	}

	tr := &TranslationResult{}
	// TODO: Properly extract hypotheses from the raw pointers

	Ct2TranslationResultFree(&result)
	return tr, nil
}
