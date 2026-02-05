package ctranslate2ffi

import (
	"fmt"
	"unsafe"

	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/unix"
)

var _ = unix.BytePtrFromString

var (
	ct2GetLastErrorFunc              ffi.Fun
	ct2ClearErrorFunc                ffi.Fun
	ct2ModelConfigDefaultFunc        ffi.Fun
	ct2StringsFreeFunc               ffi.Fun
	ct2FloatsFreeFunc                ffi.Fun
	ct2StorageCreateFloatFunc        ffi.Fun
	ct2StorageGetShapeFunc           ffi.Fun
	ct2StorageSizeFunc               ffi.Fun
	ct2StorageToFloatFunc            ffi.Fun
	ct2StorageFreeFunc               ffi.Fun
	ct2WhisperOptionsDefaultFunc     ffi.Fun
	ct2WhisperResultFreeFunc         ffi.Fun
	ct2WhisperCreateFunc             ffi.Fun
	ct2WhisperIsMultilingualFunc     ffi.Fun
	ct2WhisperNMelsFunc              ffi.Fun
	ct2WhisperNumLanguagesFunc       ffi.Fun
	ct2WhisperGenerateFunc           ffi.Fun
	ct2WhisperDetectLanguageFunc     ffi.Fun
	ct2WhisperEncodeFunc             ffi.Fun
	ct2WhisperFreeFunc               ffi.Fun
	ct2TranslationOptionsDefaultFunc ffi.Fun
	ct2TranslationResultFreeFunc     ffi.Fun
	ct2TranslatorCreateFunc          ffi.Fun
	ct2TranslatorTranslateBatchFunc  ffi.Fun
	ct2TranslatorTranslateFunc       ffi.Fun
	ct2TranslatorFreeFunc            ffi.Fun
	ct2GenerationOptionsDefaultFunc  ffi.Fun
	ct2GenerationResultFreeFunc      ffi.Fun
	ct2GeneratorCreateFunc           ffi.Fun
	ct2GeneratorGenerateFunc         ffi.Fun
	ct2GeneratorGenerateBatchFunc    ffi.Fun
	ct2GeneratorFreeFunc             ffi.Fun
	ct2VersionFunc                   ffi.Fun
	ct2CudaAvailableFunc             ffi.Fun
	ct2CudaDeviceCountFunc           ffi.Fun
)

func loadFuncs() error {
	var err error

	if ct2GetLastErrorFunc, err = lib.Prep("ct2_get_last_error", &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_get_last_error: %w", err)
	}

	if ct2ClearErrorFunc, err = lib.Prep("ct2_clear_error", &ffi.TypeVoid); err != nil {
		return fmt.Errorf("ct2_clear_error: %w", err)
	}

	if ct2ModelConfigDefaultFunc, err = lib.Prep("ct2_model_config_default", &FFITypeCt2modelconfig); err != nil {
		return fmt.Errorf("ct2_model_config_default: %w", err)
	}

	if ct2StringsFreeFunc, err = lib.Prep("ct2_strings_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_strings_free: %w", err)
	}

	if ct2FloatsFreeFunc, err = lib.Prep("ct2_floats_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_floats_free: %w", err)
	}

	if ct2StorageCreateFloatFunc, err = lib.Prep("ct2_storage_create_float", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_storage_create_float: %w", err)
	}

	if ct2StorageGetShapeFunc, err = lib.Prep("ct2_storage_get_shape", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_storage_get_shape: %w", err)
	}

	if ct2StorageSizeFunc, err = lib.Prep("ct2_storage_size", &ffi.TypeSint64, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_storage_size: %w", err)
	}

	if ct2StorageToFloatFunc, err = lib.Prep("ct2_storage_to_float", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_storage_to_float: %w", err)
	}

	if ct2StorageFreeFunc, err = lib.Prep("ct2_storage_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_storage_free: %w", err)
	}

	if ct2WhisperOptionsDefaultFunc, err = lib.Prep("ct2_whisper_options_default", &FFITypeCt2whisperoptions); err != nil {
		return fmt.Errorf("ct2_whisper_options_default: %w", err)
	}

	if ct2WhisperResultFreeFunc, err = lib.Prep("ct2_whisper_result_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_result_free: %w", err)
	}

	if ct2WhisperCreateFunc, err = lib.Prep("ct2_whisper_create", &ffi.TypePointer, &ffi.TypePointer, &FFITypeCt2modelconfig); err != nil {
		return fmt.Errorf("ct2_whisper_create: %w", err)
	}

	if ct2WhisperIsMultilingualFunc, err = lib.Prep("ct2_whisper_is_multilingual", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_is_multilingual: %w", err)
	}

	if ct2WhisperNMelsFunc, err = lib.Prep("ct2_whisper_n_mels", &ffi.TypeUint64, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_n_mels: %w", err)
	}

	if ct2WhisperNumLanguagesFunc, err = lib.Prep("ct2_whisper_num_languages", &ffi.TypeUint64, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_num_languages: %w", err)
	}

	if ct2WhisperGenerateFunc, err = lib.Prep("ct2_whisper_generate", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64, &FFITypeCt2whisperoptions, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_generate: %w", err)
	}

	if ct2WhisperDetectLanguageFunc, err = lib.Prep("ct2_whisper_detect_language", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_detect_language: %w", err)
	}

	if ct2WhisperEncodeFunc, err = lib.Prep("ct2_whisper_encode", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint8); err != nil {
		return fmt.Errorf("ct2_whisper_encode: %w", err)
	}

	if ct2WhisperFreeFunc, err = lib.Prep("ct2_whisper_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_whisper_free: %w", err)
	}

	if ct2TranslationOptionsDefaultFunc, err = lib.Prep("ct2_translation_options_default", &FFITypeCt2translationoptions); err != nil {
		return fmt.Errorf("ct2_translation_options_default: %w", err)
	}

	if ct2TranslationResultFreeFunc, err = lib.Prep("ct2_translation_result_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_translation_result_free: %w", err)
	}

	if ct2TranslatorCreateFunc, err = lib.Prep("ct2_translator_create", &ffi.TypePointer, &ffi.TypePointer, &FFITypeCt2modelconfig); err != nil {
		return fmt.Errorf("ct2_translator_create: %w", err)
	}

	if ct2TranslatorTranslateBatchFunc, err = lib.Prep("ct2_translator_translate_batch", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64, &FFITypeCt2translationoptions, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_translator_translate_batch: %w", err)
	}

	if ct2TranslatorTranslateFunc, err = lib.Prep("ct2_translator_translate", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64, &FFITypeCt2translationoptions, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_translator_translate: %w", err)
	}

	if ct2TranslatorFreeFunc, err = lib.Prep("ct2_translator_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_translator_free: %w", err)
	}

	if ct2GenerationOptionsDefaultFunc, err = lib.Prep("ct2_generation_options_default", &FFITypeCt2generationoptions); err != nil {
		return fmt.Errorf("ct2_generation_options_default: %w", err)
	}

	if ct2GenerationResultFreeFunc, err = lib.Prep("ct2_generation_result_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_generation_result_free: %w", err)
	}

	if ct2GeneratorCreateFunc, err = lib.Prep("ct2_generator_create", &ffi.TypePointer, &ffi.TypePointer, &FFITypeCt2modelconfig); err != nil {
		return fmt.Errorf("ct2_generator_create: %w", err)
	}

	if ct2GeneratorGenerateFunc, err = lib.Prep("ct2_generator_generate", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64, &FFITypeCt2generationoptions, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_generator_generate: %w", err)
	}

	if ct2GeneratorGenerateBatchFunc, err = lib.Prep("ct2_generator_generate_batch", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint64, &FFITypeCt2generationoptions, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_generator_generate_batch: %w", err)
	}

	if ct2GeneratorFreeFunc, err = lib.Prep("ct2_generator_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_generator_free: %w", err)
	}

	if ct2VersionFunc, err = lib.Prep("ct2_version", &ffi.TypePointer); err != nil {
		return fmt.Errorf("ct2_version: %w", err)
	}

	if ct2CudaAvailableFunc, err = lib.Prep("ct2_cuda_available", &ffi.TypeUint8); err != nil {
		return fmt.Errorf("ct2_cuda_available: %w", err)
	}

	if ct2CudaDeviceCountFunc, err = lib.Prep("ct2_cuda_device_count", &ffi.TypeSint32); err != nil {
		return fmt.Errorf("ct2_cuda_device_count: %w", err)
	}

	return nil
}

func Ct2GetLastError() string {
	var resultPtr *byte
	ct2GetLastErrorFunc.Call(unsafe.Pointer(&resultPtr))
	if resultPtr == nil {
		return ""
	}
	return unix.BytePtrToString(resultPtr)
}

func Ct2ClearError() {
	ct2ClearErrorFunc.Call(nil)
}

func Ct2ModelConfigDefault() Ct2modelconfig {
	var result Ct2modelconfig
	ct2ModelConfigDefaultFunc.Call(unsafe.Pointer(&result))
	return result
}

func Ct2StringsFree(arr *Ct2stringarray) {
	ct2StringsFreeFunc.Call(nil, unsafe.Pointer(&arr))
}

func Ct2FloatsFree(arr *Ct2floatarray) {
	ct2FloatsFreeFunc.Call(nil, unsafe.Pointer(&arr))
}

func Ct2StorageCreateFloat(data uintptr, shape uintptr, ndims uint64, device Ct2device) Ct2storageview {
	var result Ct2storageview
	ct2StorageCreateFloatFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&data), unsafe.Pointer(&shape), unsafe.Pointer(&ndims), unsafe.Pointer(&device))
	return result
}

func Ct2StorageGetShape(storage Ct2storageview, shape uintptr, ndims uintptr) int32 {
	var result ffi.Arg
	ct2StorageGetShapeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&storage), unsafe.Pointer(&shape), unsafe.Pointer(&ndims))
	return int32(result)
}

func Ct2StorageSize(storage Ct2storageview) int64 {
	var result int64
	ct2StorageSizeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&storage))
	return result
}

func Ct2StorageToFloat(storage Ct2storageview, buffer uintptr) int32 {
	var result ffi.Arg
	ct2StorageToFloatFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&storage), unsafe.Pointer(&buffer))
	return int32(result)
}

func Ct2StorageFree(storage Ct2storageview) {
	ct2StorageFreeFunc.Call(nil, unsafe.Pointer(&storage))
}

func Ct2WhisperOptionsDefault() Ct2whisperoptions {
	var result Ct2whisperoptions
	ct2WhisperOptionsDefaultFunc.Call(unsafe.Pointer(&result))
	return result
}

func Ct2WhisperResultFree(result *Ct2whisperresult) {
	ct2WhisperResultFreeFunc.Call(nil, unsafe.Pointer(&result))
}

func Ct2WhisperCreate(modelPath string, config Ct2modelconfig) Ct2whisper {
	modelPathPtr, _ := unix.BytePtrFromString(modelPath)
	var result Ct2whisper
	ct2WhisperCreateFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&modelPathPtr), &config)
	return result
}

func Ct2WhisperIsMultilingual(whisper Ct2whisper) bool {
	var result ffi.Arg
	ct2WhisperIsMultilingualFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&whisper))
	return result.Bool()
}

func Ct2WhisperNMels(whisper Ct2whisper) uint64 {
	var result uint64
	ct2WhisperNMelsFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&whisper))
	return result
}

func Ct2WhisperNumLanguages(whisper Ct2whisper) uint64 {
	var result uint64
	ct2WhisperNumLanguagesFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&whisper))
	return result
}

func Ct2WhisperGenerate(whisper Ct2whisper, features Ct2storageview, prompts uintptr, numPrompts uint64, options Ct2whisperoptions, resultOut *Ct2whisperresult) int32 {
	var ret ffi.Arg
	ct2WhisperGenerateFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&whisper), unsafe.Pointer(&features), unsafe.Pointer(&prompts), unsafe.Pointer(&numPrompts), &options, unsafe.Pointer(&resultOut))
	return int32(ret)
}

func Ct2WhisperDetectLanguage(whisper Ct2whisper, features Ct2storageview, languages *Ct2stringarray, probabilities *Ct2floatarray) int32 {
	var result ffi.Arg
	ct2WhisperDetectLanguageFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&whisper), unsafe.Pointer(&features), unsafe.Pointer(&languages), unsafe.Pointer(&probabilities))
	return int32(result)
}

func Ct2WhisperEncode(whisper Ct2whisper, features Ct2storageview, toCpu bool) Ct2storageview {
	var result Ct2storageview
	ct2WhisperEncodeFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&whisper), unsafe.Pointer(&features), unsafe.Pointer(&toCpu))
	return result
}

func Ct2WhisperFree(whisper Ct2whisper) {
	ct2WhisperFreeFunc.Call(nil, unsafe.Pointer(&whisper))
}

func Ct2TranslationOptionsDefault() Ct2translationoptions {
	var result Ct2translationoptions
	ct2TranslationOptionsDefaultFunc.Call(unsafe.Pointer(&result))
	return result
}

func Ct2TranslationResultFree(result *Ct2translationresult) {
	ct2TranslationResultFreeFunc.Call(nil, unsafe.Pointer(&result))
}

func Ct2TranslatorCreate(modelPath string, config Ct2modelconfig) Ct2translator {
	modelPathPtr, _ := unix.BytePtrFromString(modelPath)
	var result Ct2translator
	ct2TranslatorCreateFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&modelPathPtr), &config)
	return result
}

func Ct2TranslatorTranslateBatch(translator Ct2translator, sources string, sourceLengths uintptr, numSources uint64, options Ct2translationoptions, results *Ct2translationresult) int32 {
	sourcesPtr, _ := unix.BytePtrFromString(sources)
	var result ffi.Arg
	ct2TranslatorTranslateBatchFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&translator), unsafe.Pointer(&sourcesPtr), unsafe.Pointer(&sourceLengths), unsafe.Pointer(&numSources), &options, unsafe.Pointer(&results))
	return int32(result)
}

func Ct2TranslatorTranslate(translator Ct2translator, source uintptr, sourceLength uint64, options Ct2translationoptions, resultOut *Ct2translationresult) int32 {
	var ret ffi.Arg
	ct2TranslatorTranslateFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&translator), unsafe.Pointer(&source), unsafe.Pointer(&sourceLength), &options, unsafe.Pointer(&resultOut))
	return int32(ret)
}

func Ct2TranslatorFree(translator Ct2translator) {
	ct2TranslatorFreeFunc.Call(nil, unsafe.Pointer(&translator))
}

func Ct2GenerationOptionsDefault() Ct2generationoptions {
	var result Ct2generationoptions
	ct2GenerationOptionsDefaultFunc.Call(unsafe.Pointer(&result))
	return result
}

func Ct2GenerationResultFree(result *Ct2generationresult) {
	ct2GenerationResultFreeFunc.Call(nil, unsafe.Pointer(&result))
}

func Ct2GeneratorCreate(modelPath string, config Ct2modelconfig) Ct2generator {
	modelPathPtr, _ := unix.BytePtrFromString(modelPath)
	var result Ct2generator
	ct2GeneratorCreateFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&modelPathPtr), &config)
	return result
}

func Ct2GeneratorGenerate(generator Ct2generator, prompt uintptr, promptLength uint64, options Ct2generationoptions, resultOut *Ct2generationresult) int32 {
	var ret ffi.Arg
	ct2GeneratorGenerateFunc.Call(unsafe.Pointer(&ret), unsafe.Pointer(&generator), unsafe.Pointer(&prompt), unsafe.Pointer(&promptLength), &options, unsafe.Pointer(&resultOut))
	return int32(ret)
}

func Ct2GeneratorGenerateBatch(generator Ct2generator, prompts string, promptLengths uintptr, numPrompts uint64, options Ct2generationoptions, results *Ct2generationresult) int32 {
	promptsPtr, _ := unix.BytePtrFromString(prompts)
	var result ffi.Arg
	ct2GeneratorGenerateBatchFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&generator), unsafe.Pointer(&promptsPtr), unsafe.Pointer(&promptLengths), unsafe.Pointer(&numPrompts), &options, unsafe.Pointer(&results))
	return int32(result)
}

func Ct2GeneratorFree(generator Ct2generator) {
	ct2GeneratorFreeFunc.Call(nil, unsafe.Pointer(&generator))
}

func Ct2Version() string {
	var resultPtr *byte
	ct2VersionFunc.Call(unsafe.Pointer(&resultPtr))
	if resultPtr == nil {
		return ""
	}
	return unix.BytePtrToString(resultPtr)
}

func Ct2CudaAvailable() bool {
	var result ffi.Arg
	ct2CudaAvailableFunc.Call(unsafe.Pointer(&result))
	return result.Bool()
}

func Ct2CudaDeviceCount() int32 {
	var result ffi.Arg
	ct2CudaDeviceCountFunc.Call(unsafe.Pointer(&result))
	return int32(result)
}
