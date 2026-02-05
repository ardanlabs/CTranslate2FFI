# CTranslate2 Go Bindings

Go bindings for [CTranslate2](https://github.com/OpenNMT/CTranslate2) using libffi (no CGo required).

## Features

- **Whisper**: Speech recognition and transcription
- **Translator**: Sequence-to-sequence translation
- **Generator**: Language model text generation

## Requirements

### Runtime

- The CTranslate2 shared library (`libctranslate2.so`, `libctranslate2.dylib`, or `ctranslate2.dll`)
- libffi runtime:
  - **Linux**: `apt install libffi8` or `dnf install libffi`
  - **macOS**: Bundled with the system
  - **Windows**: Bundled with the ffi package

## Installation

```bash
go get github.com/ardanlabs/ctranslate2ffi
```

## Usage

### Load the Library

Before using any functions, load the CTranslate2 shared library:

```go
package main

import "github.com/ardanlabs/ctranslate2ffi"

func main() {
    // Provide the directory containing libctranslate2.so/dylib/dll
    if err := ctranslate2ffi.Load("/usr/local/lib"); err != nil {
        panic(err)
    }
    
    // Check version
    fmt.Println("CTranslate2 version:", ctranslate2.Version())
}
```

### Whisper (Speech Recognition)

```go
// Load a Whisper model
config := ctranslate2.DefaultModelConfig()
config.Device = ctranslate2.DeviceCPU
config.ComputeType = ctranslate2.ComputeFloat32

whisper, err := ctranslate2.NewWhisper("/path/to/whisper-model", config)
if err != nil {
    panic(err)
}
defer whisper.Close()

// Check model properties
fmt.Println("Multilingual:", whisper.IsMultilingual())
fmt.Println("Num mels:", whisper.NumMels())

// Create audio features (mel spectrogram)
// Shape should be [batch_size, n_mels, time_frames]
features, err := ctranslate2.NewStorageViewFloat(melData, []int64{1, 80, 3000}, ctranslate2.DeviceCPU)
if err != nil {
    panic(err)
}
defer features.Close()

// Transcribe
opts := ctranslate2.DefaultWhisperOptions()
result, err := whisper.Generate(features, []string{"<|startoftranscript|>", "<|en|>"}, opts)
if err != nil {
    panic(err)
}
```

### Translator

```go
// Load a translation model
translator, err := ctranslate2.NewTranslator("/path/to/model", config)
if err != nil {
    panic(err)
}
defer translator.Close()

// Translate tokenized input
opts := ctranslate2.DefaultTranslationOptions()
result, err := translator.Translate([]string{"Hello", "world"}, opts)
if err != nil {
    panic(err)
}
```

### Generator (Language Model)

```go
// Load a language model
generator, err := ctranslate2.NewGenerator("/path/to/model", config)
if err != nil {
    panic(err)
}
defer generator.Close()

// Generate text
opts := ctranslate2.DefaultGenerationOptions()
opts.MaxLength = 100
result, err := generator.Generate([]string{"Once", "upon", "a", "time"}, opts)
if err != nil {
    panic(err)
}
```

## Building CTranslate2 with C API

The C API wrapper must be compiled with CTranslate2:

1. Copy `ctranslate2_c.h` to `CTranslate2/include/`
2. Copy `ctranslate2_c.cc` to `CTranslate2/src/`
3. Add the source file to CMakeLists.txt
4. Build CTranslate2:

```bash
cd CTranslate2
mkdir build && cd build
cmake -DBUILD_SHARED_LIBS=ON ..
make -j$(nproc)
```

## API Reference

### Types

- `Device` - CPU or CUDA
- `ComputeType` - Computation precision (float32, int8, float16, etc.)
- `ModelConfig` - Model loading configuration
- `WhisperOptions` - Whisper generation options
- `TranslationOptions` - Translation options
- `GenerationOptions` - Text generation options

### Functions

- `Load(path string)` - Load the CTranslate2 shared library
- `Version()` - Get library version
- `CUDAAvailable()` - Check if CUDA is available
- `CUDADeviceCount()` - Get number of CUDA devices

### Error Handling

- `GetLastError()` - Get the last error message
- `ClearError()` - Clear the error state

## License

Apache 2.0
