package models

type VideoOutputFormat string

const (
	VideoOutputFormatMP4  VideoOutputFormat = "MP4"
	VideoOutputFormatWEBM VideoOutputFormat = "WEBM"
)

type FramePosition string

const (
	FramePositionFirst FramePosition = "first"
	FramePositionLast  FramePosition = "last"
)

type FrameImage struct {
	InputImage string        `json:"inputImage"`
	Frame      FramePosition `json:"frame"`
}
type ReferenceImage struct {
	InputImage string `json:"inputImage"`
}
type ReferenceVideo struct {
	InputVideo string `json:"inputVideo"`
}
type InputAudio struct {
	InputAudio string `json:"inputAudio"`
}

type Speech struct {
	Voice string `json:"voice"`
	Text  string `json:"text"`
}

type GoogleVideoSettings struct {
	EnhancePrompt *bool `json:"enhancePrompt,omitempty"`
	GenerateAudio *bool `json:"generateAudio,omitempty"`
}
type ByteDanceVideoSettings struct {
	CameraFixed *bool `json:"cameraFixed,omitempty"`
}
type MiniMaxVideoSettings struct {
	PromptOptimizer *bool `json:"promptOptimizer,omitempty"`
}

type PixVerseVideoSettings struct {
	Style              *string `json:"style,omitempty"`
	Effect             *string `json:"effect,omitempty"`
	CameraMovement     *string `json:"cameraMovement,omitempty"`
	MotionMode         *string `json:"motionMode,omitempty"`
	SoundEffectSwitch  *bool   `json:"soundEffectSwitch,omitempty"`
	SoundEffectContent *string `json:"soundEffectContent,omitempty"`
}

type ViduVideoSettings struct {
	MovementAmplitude *string `json:"movementAmplitude,omitempty"`
	BGM               *bool   `json:"bgm,omitempty"`
	Style             *string `json:"style,omitempty"`
}

type VideoProviderSettings struct {
	Google    *GoogleVideoSettings    `json:"google,omitempty"`
	ByteDance *ByteDanceVideoSettings `json:"bytedance,omitempty"`
	MiniMax   *MiniMaxVideoSettings   `json:"minimax,omitempty"`
	PixVerse  *PixVerseVideoSettings  `json:"pixverse,omitempty"`
	Vidu      *ViduVideoSettings      `json:"vidu,omitempty"`
}

type VideoInferenceRequest struct {
	TaskType           string                 `json:"taskType"`
	TaskUUID           string                 `json:"taskUUID"`
	OutputType         *OutputType            `json:"outputType,omitempty"`
	OutputFormat       *VideoOutputFormat     `json:"outputFormat,omitempty"`
	OutputQuality      *int                   `json:"outputQuality,omitempty"`
	WebhookURL         *string                `json:"webhookURL,omitempty"`
	DeliveryMethod     *DeliveryMethod        `json:"deliveryMethod,omitempty"`
	UploadEndpoint     *string                `json:"uploadEndpoint,omitempty"`
	Safety             *Safety                `json:"safety,omitempty"`
	TTL                *int                   `json:"ttl,omitempty"`
	IncludeCost        *bool                  `json:"includeCost,omitempty"`
	PositivePrompt     string                 `json:"positivePrompt"`
	NegativePrompt     *string                `json:"negativePrompt,omitempty"`
	FrameImages        []FrameImage           `json:"frameImages,omitempty"`
	ReferenceImages    []ReferenceImage       `json:"referenceImages,omitempty"`
	ReferenceVideos    []ReferenceVideo       `json:"referenceVideos,omitempty"`
	InputAudios        []InputAudio           `json:"inputAudios,omitempty"`
	Width              *int                   `json:"width,omitempty"`
	Height             *int                   `json:"height,omitempty"`
	Model              string                 `json:"model"`
	Duration           *int                   `json:"duration,omitempty"`
	FPS                *int                   `json:"fps,omitempty"`
	Steps              *int                   `json:"steps,omitempty"`
	Seed               *int64                 `json:"seed,omitempty"`
	CFGScale           *float64               `json:"CFGScale,omitempty"`
	Speech             *Speech                `json:"speech,omitempty"`
	NumberResults      *int                   `json:"numberResults,omitempty"`
	Acceleration       *Acceleration          `json:"acceleration,omitempty"`
	AdvancedFeatures   *AdvancedFeatures      `json:"advancedFeatures,omitempty"`
	AcceleratorOptions *AcceleratorOptions    `json:"acceleratorOptions,omitempty"`
	LoRA               []LoRA                 `json:"lora,omitempty"`
	ProviderSettings   *VideoProviderSettings `json:"providerSettings,omitempty"`
}

type VideoInferenceResponse struct {
	TaskType     string     `json:"taskType"`
	TaskUUID     string     `json:"taskUUID"`
	Status       TaskStatus `json:"status,omitempty"`
	VideoUUID    string     `json:"videoUUID,omitempty"`
	VideoURL     *string    `json:"videoURL,omitempty"`
	ThumbnailURL *string    `json:"thumbnailURL,omitempty"`
	Seed         *int64     `json:"seed,omitempty"`
	Cost         *float64   `json:"cost,omitempty"`
}
