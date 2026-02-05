// Example: Transcribe an MP3 file using CTranslate2 Whisper bindings
//
// This example demonstrates how to:
// 1. Load an MP3 audio file
// 2. Convert it to mel spectrogram features
// 3. Use the CTranslate2 Whisper model to transcribe
// 4. Output the transcription to stdout
//
// Requirements:
// - CTranslate2 library with C API built and installed
// - A Whisper model converted to CTranslate2 format
// - The audio file (tts-sample.mp3)
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ardanlabs/ctranslate2ffi"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func main() {
	// Command-line flags
	libPath := flag.String("lib", "/usr/local/lib", "Path to CTranslate2 library directory")
	modelPath := flag.String("model", "", "Path to Whisper CTranslate2 model directory")
	audioFile := flag.String("audio", "tts-sample.mp3", "Audio file to transcribe")
	language := flag.String("lang", "en", "Language code (e.g., en, es, fr)")
	flag.Parse()

	if *modelPath == "" {
		log.Fatal("Please provide the path to a Whisper model with -model flag")
	}

	// Load the CTranslate2 library
	if err := ctranslate2ffi.Load(*libPath); err != nil {
		log.Fatalf("Failed to load CTranslate2 library: %v", err)
	}

	fmt.Printf("CTranslate2 version: %s\n", ctranslate2ffi.Version())
	fmt.Printf("CUDA available: %v\n", ctranslate2ffi.CUDAAvailable())

	// Convert MP3 to WAV if needed
	wavFile := *audioFile
	if strings.HasSuffix(strings.ToLower(*audioFile), ".mp3") {
		var err error
		wavFile, err = convertMP3ToWAV(*audioFile)
		if err != nil {
			log.Fatalf("Failed to convert MP3 to WAV: %v", err)
		}
		defer os.Remove(wavFile) // Clean up temp file
	}

	// Load audio samples
	samples, sampleRate, err := loadWAV(wavFile)
	if err != nil {
		log.Fatalf("Failed to load audio: %v", err)
	}
	fmt.Printf("Loaded %d samples at %d Hz (%.2f seconds)\n",
		len(samples), sampleRate, float64(len(samples))/float64(sampleRate))

	// Resample to 16kHz if needed (Whisper expects 16kHz)
	if sampleRate != 16000 {
		samples = resample(samples, sampleRate, 16000)
		sampleRate = 16000
		fmt.Printf("Resampled to %d samples at %d Hz\n", len(samples), sampleRate)
	}

	// Compute mel spectrogram features
	// Whisper expects 80 mel bands, 3000 time frames for 30 seconds of audio
	melFeatures, err := computeMelSpectrogram(samples, sampleRate, 80)
	if err != nil {
		log.Fatalf("Failed to compute mel spectrogram: %v", err)
	}
	fmt.Printf("Computed mel spectrogram: %d mel bands x %d time frames\n",
		len(melFeatures), len(melFeatures[0]))

	// Flatten mel spectrogram for CTranslate2
	// Shape: [1, n_mels, time_frames]
	nMels := len(melFeatures)
	timeFrames := len(melFeatures[0])
	flatMel := make([]float32, nMels*timeFrames)
	for i := 0; i < nMels; i++ {
		for j := 0; j < timeFrames; j++ {
			flatMel[i*timeFrames+j] = float32(melFeatures[i][j])
		}
	}

	// Load Whisper model
	fmt.Println("Loading Whisper model...")
	config := ctranslate2ffi.DefaultModelConfig()
	config.Device = ctranslate2ffi.DeviceCPU
	config.ComputeType = ctranslate2ffi.ComputeFloat32

	whisper, err := ctranslate2ffi.NewWhisper(*modelPath, config)
	if err != nil {
		log.Fatalf("Failed to load Whisper model: %v", err)
	}
	defer whisper.Close()

	fmt.Printf("Model loaded - Multilingual: %v, Mels: %d, Languages: %d\n",
		whisper.IsMultilingual(), whisper.NumMels(), whisper.NumLanguages())

	// Create storage view for mel features
	shape := []int64{1, int64(nMels), int64(timeFrames)}
	features, err := ctranslate2ffi.NewStorageViewFloat(flatMel, shape, ctranslate2ffi.DeviceCPU)
	if err != nil {
		log.Fatalf("Failed to create storage view: %v", err)
	}
	defer features.Close()

	// Prepare prompts (Whisper special tokens)
	prompts := []string{
		"<|startoftranscript|>",
		"<|" + *language + "|>",
		"<|transcribe|>",
		"<|notimestamps|>",
	}

	// Transcribe
	fmt.Println("Transcribing...")
	opts := ctranslate2ffi.DefaultWhisperOptions()
	opts.BeamSize = 5
	opts.ReturnScores = true

	result, err := whisper.Generate(features, prompts, opts)
	if err != nil {
		log.Fatalf("Transcription failed: %v", err)
	}

	// Output transcription
	fmt.Println("\n=== Transcription ===")
	for _, seq := range result.Sequences {
		fmt.Println(seq)
	}
	if len(result.Scores) > 0 {
		fmt.Printf("\nScore: %.4f\n", result.Scores[0])
	}
}

// convertMP3ToWAV converts an MP3 file to WAV using ffmpeg
func convertMP3ToWAV(mp3Path string) (string, error) {
	wavPath := filepath.Join(os.TempDir(), "whisper_temp.wav")

	cmd := exec.Command("ffmpeg", "-y", "-i", mp3Path,
		"-ar", "16000", // Resample to 16kHz
		"-ac", "1", // Mono
		"-f", "wav",
		wavPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %w\n%s", err, output)
	}

	return wavPath, nil
}

// loadWAV loads a WAV file and returns samples as float32 normalized to [-1, 1]
func loadWAV(path string) ([]float32, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, 0, fmt.Errorf("invalid WAV file")
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	// Convert to float32
	samples := make([]float32, len(buf.Data))
	bitDepth := decoder.BitDepth
	maxVal := float32(int(1)<<(bitDepth-1) - 1)

	for i, s := range buf.Data {
		samples[i] = float32(s) / maxVal
	}

	return samples, int(decoder.SampleRate), nil
}

// resample resamples audio from srcRate to dstRate using linear interpolation
func resample(samples []float32, srcRate, dstRate int) []float32 {
	if srcRate == dstRate {
		return samples
	}

	ratio := float64(srcRate) / float64(dstRate)
	newLen := int(float64(len(samples)) / ratio)
	result := make([]float32, newLen)

	for i := 0; i < newLen; i++ {
		srcIdx := float64(i) * ratio
		srcIdxInt := int(srcIdx)
		frac := float32(srcIdx - float64(srcIdxInt))

		if srcIdxInt+1 < len(samples) {
			result[i] = samples[srcIdxInt]*(1-frac) + samples[srcIdxInt+1]*frac
		} else if srcIdxInt < len(samples) {
			result[i] = samples[srcIdxInt]
		}
	}

	return result
}

// computeMelSpectrogram computes a log mel spectrogram from audio samples
// This is a simplified implementation - for production, use a proper DSP library
func computeMelSpectrogram(samples []float32, sampleRate, nMels int) ([][]float64, error) {
	// Whisper parameters
	nFFT := 400       // 25ms window at 16kHz
	hopLength := 160  // 10ms hop at 16kHz
	maxFrames := 3000 // 30 seconds

	// Calculate number of frames
	nFrames := (len(samples) - nFFT) / hopLength
	if nFrames > maxFrames {
		nFrames = maxFrames
	}
	if nFrames < 1 {
		nFrames = 1
	}

	// Create mel filterbank
	melFilters := createMelFilterbank(nMels, nFFT, sampleRate)

	// Initialize mel spectrogram
	melSpec := make([][]float64, nMels)
	for i := range melSpec {
		melSpec[i] = make([]float64, nFrames)
	}

	// Hann window
	window := make([]float64, nFFT)
	for i := 0; i < nFFT; i++ {
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(nFFT-1)))
	}

	// Process each frame
	for frame := 0; frame < nFrames; frame++ {
		start := frame * hopLength
		if start+nFFT > len(samples) {
			break
		}

		// Apply window and compute FFT magnitude
		magnitudes := make([]float64, nFFT/2+1)
		for i := 0; i < nFFT/2+1; i++ {
			var re, im float64
			for j := 0; j < nFFT; j++ {
				if start+j < len(samples) {
					x := float64(samples[start+j]) * window[j]
					angle := -2 * math.Pi * float64(i) * float64(j) / float64(nFFT)
					re += x * math.Cos(angle)
					im += x * math.Sin(angle)
				}
			}
			magnitudes[i] = math.Sqrt(re*re + im*im)
		}

		// Apply mel filterbank
		for m := 0; m < nMels; m++ {
			var energy float64
			for k := 0; k < len(magnitudes); k++ {
				energy += magnitudes[k] * magnitudes[k] * melFilters[m][k]
			}
			// Log mel spectrogram
			if energy > 1e-10 {
				melSpec[m][frame] = math.Log(energy)
			} else {
				melSpec[m][frame] = math.Log(1e-10)
			}
		}
	}

	// Normalize (Whisper-style)
	maxVal := -1e10
	for m := 0; m < nMels; m++ {
		for f := 0; f < nFrames; f++ {
			if melSpec[m][f] > maxVal {
				maxVal = melSpec[m][f]
			}
		}
	}
	for m := 0; m < nMels; m++ {
		for f := 0; f < nFrames; f++ {
			melSpec[m][f] = (melSpec[m][f] - maxVal) / 4.0
			if melSpec[m][f] < -1.0 {
				melSpec[m][f] = -1.0
			}
		}
	}

	return melSpec, nil
}

// createMelFilterbank creates a mel filterbank matrix
func createMelFilterbank(nMels, nFFT, sampleRate int) [][]float64 {
	fMin := 0.0
	fMax := float64(sampleRate) / 2.0

	// Convert to mel scale
	melMin := 2595.0 * math.Log10(1.0+fMin/700.0)
	melMax := 2595.0 * math.Log10(1.0+fMax/700.0)

	// Create mel points
	melPoints := make([]float64, nMels+2)
	for i := 0; i <= nMels+1; i++ {
		melPoints[i] = melMin + float64(i)*(melMax-melMin)/float64(nMels+1)
	}

	// Convert back to Hz
	hzPoints := make([]float64, nMels+2)
	for i := range hzPoints {
		hzPoints[i] = 700.0 * (math.Pow(10, melPoints[i]/2595.0) - 1.0)
	}

	// Convert to FFT bin indices
	binPoints := make([]int, nMels+2)
	for i := range binPoints {
		binPoints[i] = int(math.Floor((float64(nFFT) + 1) * hzPoints[i] / float64(sampleRate)))
	}

	// Create filterbank
	filterbank := make([][]float64, nMels)
	for m := 0; m < nMels; m++ {
		filterbank[m] = make([]float64, nFFT/2+1)
		for k := binPoints[m]; k < binPoints[m+1]; k++ {
			if k < nFFT/2+1 {
				filterbank[m][k] = float64(k-binPoints[m]) / float64(binPoints[m+1]-binPoints[m])
			}
		}
		for k := binPoints[m+1]; k < binPoints[m+2]; k++ {
			if k < nFFT/2+1 {
				filterbank[m][k] = float64(binPoints[m+2]-k) / float64(binPoints[m+2]-binPoints[m+1])
			}
		}
	}

	return filterbank
}

// Helper to format duration
func formatDuration(seconds float64) string {
	mins := int(seconds) / 60
	secs := seconds - float64(mins*60)
	return strconv.Itoa(mins) + "m" + strconv.FormatFloat(secs, 'f', 1, 64) + "s"
}

// Placeholder for audio buffer interface
var _ = audio.IntBuffer{}
