package models

import "github.com/google/uuid"

// Constructors that set defaults and generate task UUIDs

func NewImageInferenceRequest(prompt, model string, width, height int) *ImageInferenceRequest {
	return &ImageInferenceRequest{TaskType: TaskTypeImageInference, TaskUUID: uuid.New().String(), PositivePrompt: prompt, Model: model, Width: width, Height: height}
}

func NewUploadImageRequest() *UploadImageRequest {
	return &UploadImageRequest{TaskType: TaskTypeImageUpload, TaskUUID: uuid.New().String()}
}

func NewUpscaleGanRequest(inputImage string, upscaleFactor int) *UpscaleGanRequest {
	return &UpscaleGanRequest{TaskType: TaskTypeUpscaleGan, TaskUUID: uuid.New().String(), InputImage: inputImage, UpscaleFactor: upscaleFactor}
}

func NewRemoveImageBackgroundRequest(inputImage string) *RemoveImageBackgroundRequest {
	return &RemoveImageBackgroundRequest{TaskType: TaskTypeImageBackgroundRemoval, TaskUUID: uuid.New().String(), InputImage: inputImage}
}

func NewEnhancePromptRequest(prompt string) *EnhancePromptRequest {
	return &EnhancePromptRequest{TaskType: TaskTypePromptEnhance, TaskUUID: uuid.New().String(), Prompt: prompt}
}

func NewImageCaptionRequest(inputImage string) *ImageCaptionRequest {
	return &ImageCaptionRequest{TaskType: TaskTypeImageCaption, TaskUUID: uuid.New().String(), InputImage: inputImage}
}

func NewVideoInferenceRequest(prompt, model string) *VideoInferenceRequest {
	width, height, fps := 1920, 1080, 30
	numberResults := 1
	outputType := OutputTypeURL
	delivery := DeliveryMethodAsync
	return &VideoInferenceRequest{TaskType: TaskTypeVideoInference, TaskUUID: uuid.New().String(), PositivePrompt: prompt, Model: model, Width: &width, Height: &height, FPS: &fps, NumberResults: &numberResults, OutputType: &outputType, DeliveryMethod: &delivery}
}

func NewAudioInferenceRequest(prompt, model string, duration int) *AudioInferenceRequest {
	outputType := OutputTypeURL
	delivery := DeliveryMethodAsync
	format := AudioOutputFormatMP3
	numberResults := 1
	return &AudioInferenceRequest{TaskType: TaskTypeAudioInference, TaskUUID: uuid.New().String(), PositivePrompt: prompt, Model: model, Duration: &duration, OutputType: &outputType, OutputFormat: &format, DeliveryMethod: &delivery, NumberResults: &numberResults}
}

type GetResponseRequest struct {
	TaskType string `json:"taskType"`
	TaskUUID string `json:"taskUUID"`
}

func NewGetResponseRequest(taskUUID string) *GetResponseRequest {
	return &GetResponseRequest{TaskType: TaskTypeGetResponse, TaskUUID: taskUUID}
}
