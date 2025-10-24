package models

type ImageInferenceRequest struct {
	TaskType           string              `json:"taskType"`
	TaskUUID           string              `json:"taskUUID"`
	OutputType         *OutputType         `json:"outputType,omitempty"`
	OutputFormat       *OutputFormat       `json:"outputFormat,omitempty"`
	OutputQuality      *int                `json:"outputQuality,omitempty"`
	WebhookURL         *string             `json:"webhookURL,omitempty"`
	DeliveryMethod     *DeliveryMethod     `json:"deliveryMethod,omitempty"`
	UploadEndpoint     *string             `json:"uploadEndpoint,omitempty"`
	Safety             *Safety             `json:"safety,omitempty"`
	TTL                *int                `json:"ttl,omitempty"`
	IncludeCost        *bool               `json:"includeCost,omitempty"`
	PositivePrompt     string              `json:"positivePrompt"`
	NegativePrompt     *string             `json:"negativePrompt,omitempty"`
	SeedImage          *string             `json:"seedImage,omitempty"`
	MaskImage          *string             `json:"maskImage,omitempty"`
	MaskMargin         *int                `json:"maskMargin,omitempty"`
	Strength           *float64            `json:"strength,omitempty"`
	ReferenceImages    []string            `json:"referenceImages,omitempty"`
	Outpaint           *Outpaint           `json:"outpaint,omitempty"`
	Height             int                 `json:"height"`
	Width              int                 `json:"width"`
	Model              string              `json:"model"`
	VAE                *string             `json:"vae,omitempty"`
	Steps              *int                `json:"steps,omitempty"`
	Scheduler          *Scheduler          `json:"scheduler,omitempty"`
	Seed               *int64              `json:"seed,omitempty"`
	CFGScale           *float64            `json:"CFGScale,omitempty"`
	ClipSkip           *int                `json:"clipSkip,omitempty"`
	PromptWeighting    *PromptWeighting    `json:"promptWeighting,omitempty"`
	NumberResults      *int                `json:"numberResults,omitempty"`
	Acceleration       *Acceleration       `json:"acceleration,omitempty"`
	AdvancedFeatures   *AdvancedFeatures   `json:"advancedFeatures,omitempty"`
	AcceleratorOptions *AcceleratorOptions `json:"acceleratorOptions,omitempty"`
	PuLID              *PuLID              `json:"puLID,omitempty"`
	ACEPlusPlus        *ACEPlusPlus        `json:"acePlusPlus,omitempty"`
	Refiner            *Refiner            `json:"refiner,omitempty"`
	Embeddings         []Embedding         `json:"embeddings,omitempty"`
	ControlNet         []ControlNet        `json:"controlNet,omitempty"`
	LoRA               []LoRA              `json:"lora,omitempty"`
	IPAdapters         []IPAdapter         `json:"ipAdapters,omitempty"`
	ProviderSettings   *ProviderSettings   `json:"providerSettings,omitempty"`
}

type ImageInferenceResponse struct {
	TaskType        string   `json:"taskType"`
	TaskUUID        string   `json:"taskUUID"`
	ImageUUID       string   `json:"imageUUID"`
	ImageURL        *string  `json:"imageURL,omitempty"`
	ImageBase64Data *string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    *string  `json:"imageDataURI,omitempty"`
	Seed            *int64   `json:"seed,omitempty"`
	NSFWContent     *bool    `json:"NSFWContent,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
}

type UploadImageRequest struct {
	TaskType     string  `json:"taskType"`
	TaskUUID     string  `json:"taskUUID"`
	ImageBase64  *string `json:"imageBase64,omitempty"`
	ImageDataURI *string `json:"imageDataURI,omitempty"`
	ImageURL     *string `json:"imageURL,omitempty"`
}

type UploadImageResponse struct {
	TaskType  string `json:"taskType"`
	TaskUUID  string `json:"taskUUID"`
	ImageUUID string `json:"imageUUID"`
}

type UpscaleGanRequest struct {
	TaskType       string          `json:"taskType"`
	TaskUUID       string          `json:"taskUUID"`
	InputImage     string          `json:"inputImage"`
	UpscaleFactor  int             `json:"upscaleFactor"`
	OutputType     *OutputType     `json:"outputType,omitempty"`
	OutputFormat   *OutputFormat   `json:"outputFormat,omitempty"`
	OutputQuality  *int            `json:"outputQuality,omitempty"`
	WebhookURL     *string         `json:"webhookURL,omitempty"`
	DeliveryMethod *DeliveryMethod `json:"deliveryMethod,omitempty"`
	UploadEndpoint *string         `json:"uploadEndpoint,omitempty"`
	IncludeCost    *bool           `json:"includeCost,omitempty"`
}

type UpscaleGanResponse struct {
	TaskType        string   `json:"taskType"`
	TaskUUID        string   `json:"taskUUID"`
	ImageUUID       string   `json:"imageUUID"`
	ImageURL        *string  `json:"imageURL,omitempty"`
	ImageBase64Data *string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    *string  `json:"imageDataURI,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
}

type RemoveImageBackgroundRequest struct {
	TaskType       string          `json:"taskType"`
	TaskUUID       string          `json:"taskUUID"`
	InputImage     string          `json:"inputImage"`
	OutputType     *OutputType     `json:"outputType,omitempty"`
	OutputFormat   *OutputFormat   `json:"outputFormat,omitempty"`
	OutputQuality  *int            `json:"outputQuality,omitempty"`
	WebhookURL     *string         `json:"webhookURL,omitempty"`
	DeliveryMethod *DeliveryMethod `json:"deliveryMethod,omitempty"`
	UploadEndpoint *string         `json:"uploadEndpoint,omitempty"`
	IncludeCost    *bool           `json:"includeCost,omitempty"`
	Rgba           []int           `json:"rgba,omitempty"`
}

type RemoveImageBackgroundResponse struct {
	TaskType        string   `json:"taskType"`
	TaskUUID        string   `json:"taskUUID"`
	ImageUUID       string   `json:"imageUUID"`
	ImageURL        *string  `json:"imageURL,omitempty"`
	ImageBase64Data *string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    *string  `json:"imageDataURI,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
}

type EnhancePromptRequest struct {
	TaskType        string `json:"taskType"`
	TaskUUID        string `json:"taskUUID"`
	Prompt          string `json:"prompt"`
	PromptMaxLength *int   `json:"promptMaxLength,omitempty"`
	PromptVersions  *int   `json:"promptVersions,omitempty"`
	IncludeCost     *bool  `json:"includeCost,omitempty"`
}

type EnhancePromptResponse struct {
	TaskType string   `json:"taskType"`
	TaskUUID string   `json:"taskUUID"`
	Text     string   `json:"text"`
	Cost     *float64 `json:"cost,omitempty"`
}

type ImageCaptionRequest struct {
	TaskType    string `json:"taskType"`
	TaskUUID    string `json:"taskUUID"`
	InputImage  string `json:"inputImage"`
	IncludeCost *bool  `json:"includeCost,omitempty"`
}

type ImageCaptionResponse struct {
	TaskType string   `json:"taskType"`
	TaskUUID string   `json:"taskUUID"`
	Text     string   `json:"text"`
	Cost     *float64 `json:"cost,omitempty"`
}
