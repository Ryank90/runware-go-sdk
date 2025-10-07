package runware

import "github.com/google/uuid"

// OutputType specifies the format of the output image
type OutputType string

const (
	OutputTypeURL        OutputType = "URL"
	OutputTypeBase64Data OutputType = "base64Data"
	OutputTypeDataURI    OutputType = "dataURI"
)

// OutputFormat specifies the image file format
type OutputFormat string

const (
	OutputFormatJPG  OutputFormat = "jpg"
	OutputFormatPNG  OutputFormat = "png"
	OutputFormatWEBP OutputFormat = "webp"
)

// DeliveryMethod specifies how the result should be delivered
type DeliveryMethod string

const (
	DeliveryMethodStream DeliveryMethod = "stream"
	DeliveryMethodPOST   DeliveryMethod = "post"
)

// SafetyMode specifies the content safety check mode
type SafetyMode string

const (
	SafetyModeStrict   SafetyMode = "strict"
	SafetyModeModerate SafetyMode = "moderate"
	SafetyModeRelaxed  SafetyMode = "relaxed"
)

// ControlMode specifies the ControlNet control mode
type ControlMode string

const (
	ControlModeBalanced ControlMode = "balanced"
	ControlModePrompt   ControlMode = "prompt"
	ControlModeControl  ControlMode = "control"
)

// PromptWeighting specifies the prompt weighting algorithm
type PromptWeighting string

const (
	PromptWeightingCompel PromptWeighting = "compel"
	PromptWeightingSD     PromptWeighting = "sd"
)

// TaskType constants for API operations
const (
	TaskTypeImageInference         = "imageInference"
	TaskTypeVideoInference         = "videoInference"
	TaskTypePromptEnhance          = "promptEnhance"
	TaskTypeImageCaption           = "imageCaption"
	TaskTypeImageUpload            = "imageUpload"
	TaskTypeUpscaleGan             = "upscaleGan"
	TaskTypeImageBackgroundRemoval = "imageBackgroundRemoval"
	TaskTypeGetResponse            = "getResponse"
)

// Acceleration specifies the acceleration mode
type Acceleration string

const (
	AccelerationTurbo Acceleration = "turbo"
)

// Scheduler specifies the sampling scheduler
type Scheduler string

const (
	SchedulerEuler          Scheduler = "euler"
	SchedulerEulerA         Scheduler = "euler_a"
	SchedulerDPMPP2M        Scheduler = "dpmpp_2m"
	SchedulerDPMPP2MKarras  Scheduler = "dpmpp_2m_karras"
	SchedulerDPMPPSDE       Scheduler = "dpmpp_sde"
	SchedulerDPMPPSDEKarras Scheduler = "dpmpp_sde_karras"
	SchedulerLMS            Scheduler = "lms"
	SchedulerLMSKarras      Scheduler = "lms_karras"
	SchedulerHeun           Scheduler = "heun"
	SchedulerDDIM           Scheduler = "ddim"
	SchedulerPNDM           Scheduler = "pndm"
)

// Safety contains content safety check parameters
type Safety struct {
	CheckContent bool       `json:"checkContent,omitempty"`
	Mode         SafetyMode `json:"mode,omitempty"`
}

// Outpaint defines the outpainting parameters
type Outpaint struct {
	Top    int `json:"top,omitempty"`
	Right  int `json:"right,omitempty"`
	Bottom int `json:"bottom,omitempty"`
	Left   int `json:"left,omitempty"`
	Blur   int `json:"blur,omitempty"`
}

// AcceleratorOptions contains advanced acceleration parameters
type AcceleratorOptions struct {
	TeaCache                 *bool    `json:"teaCache,omitempty"`
	TeaCacheDistance         *int     `json:"teaCacheDistance,omitempty"`
	DeepCache                *bool    `json:"deepCache,omitempty"`
	DeepCacheInterval        *int     `json:"deepCacheInterval,omitempty"`
	DeepCacheBranchID        *string  `json:"deepCacheBranchId,omitempty"`
	CacheStartStep           *int     `json:"cacheStartStep,omitempty"`
	CacheStartStepPercentage *float64 `json:"cacheStartStepPercentage,omitempty"`
	CacheEndStep             *int     `json:"cacheEndStep,omitempty"`
	CacheEndStepPercentage   *float64 `json:"cacheEndStepPercentage,omitempty"`
	CacheMaxConsecutiveSteps *int     `json:"cacheMaxConsecutiveSteps,omitempty"`
}

// AdvancedFeatures contains advanced feature flags
type AdvancedFeatures struct {
	LayerDiffuse *bool `json:"layerDiffuse,omitempty"`
}

// PuLID contains PuLID identity preservation parameters
type PuLID struct {
	InputImages            []string `json:"inputImages"`
	IDWeight               *float64 `json:"idWeight,omitempty"`
	TrueCFGScale           *float64 `json:"trueCFGScale,omitempty"`
	CFGStartStep           *int     `json:"CFGStartStep,omitempty"`
	CFGStartStepPercentage *float64 `json:"CFGStartStepPercentage,omitempty"`
}

// ACEPlusPlusType specifies the ACE++ operation type
type ACEPlusPlusType string

const (
	ACEPlusPlusTypeFaceSwap   ACEPlusPlusType = "faceSwap"
	ACEPlusPlusTypeRepainting ACEPlusPlusType = "repainting"
)

// ACEPlusPlus contains ACE++ parameters
type ACEPlusPlus struct {
	Type            ACEPlusPlusType `json:"type"`
	InputImages     []string        `json:"inputImages"`
	InputMasks      []string        `json:"inputMasks,omitempty"`
	RepaintingScale *float64        `json:"repaintingScale,omitempty"`
}

// Refiner contains refiner model parameters
type Refiner struct {
	Model               string   `json:"model"`
	StartStep           *int     `json:"startStep,omitempty"`
	StartStepPercentage *float64 `json:"startStepPercentage,omitempty"`
}

// Embedding contains embedding parameters
type Embedding struct {
	Model  string   `json:"model"`
	Weight *float64 `json:"weight,omitempty"`
}

// ControlNet contains ControlNet parameters
type ControlNet struct {
	Model               string       `json:"model"`
	GuideImage          string       `json:"guideImage"`
	Weight              *float64     `json:"weight,omitempty"`
	StartStep           *int         `json:"startStep,omitempty"`
	StartStepPercentage *float64     `json:"startStepPercentage,omitempty"`
	EndStep             *int         `json:"endStep,omitempty"`
	EndStepPercentage   *float64     `json:"endStepPercentage,omitempty"`
	ControlMode         *ControlMode `json:"controlMode,omitempty"`
}

// LoRA contains LoRA parameters
type LoRA struct {
	Model  string   `json:"model"`
	Weight *float64 `json:"weight,omitempty"`
}

// IPAdapter contains IP-Adapter parameters
type IPAdapter struct {
	Model      string   `json:"model"`
	GuideImage string   `json:"guideImage"`
	Weight     *float64 `json:"weight,omitempty"`
}

// BFLProviderSettings contains Black Forest Labs provider settings
type BFLProviderSettings struct {
	PromptUpsampling *bool `json:"promptUpsampling,omitempty"`
	SafetyTolerance  *int  `json:"safetyTolerance,omitempty"`
	Raw              *bool `json:"raw,omitempty"`
}

// ByteDanceProviderSettings contains ByteDance provider settings
type ByteDanceProviderSettings struct {
	MaxSequentialImages *int `json:"maxSequentialImages,omitempty"`
}

// IdeogramRenderingSpeed specifies the rendering speed for Ideogram
type IdeogramRenderingSpeed string

const (
	IdeogramRenderingSpeedFast IdeogramRenderingSpeed = "fast"
	IdeogramRenderingSpeedSlow IdeogramRenderingSpeed = "slow"
)

// IdeogramStyleType specifies the style type for Ideogram
type IdeogramStyleType string

const (
	IdeogramStyleTypeAuto      IdeogramStyleType = "auto"
	IdeogramStyleTypeGeneral   IdeogramStyleType = "general"
	IdeogramStyleTypeRealistic IdeogramStyleType = "realistic"
	IdeogramStyleTypeDesign    IdeogramStyleType = "design"
	IdeogramStyleTypeRender3D  IdeogramStyleType = "render_3d"
	IdeogramStyleTypeAnime     IdeogramStyleType = "anime"
)

// IdeogramStylePreset specifies style presets for Ideogram
type IdeogramStylePreset string

const (
	IdeogramStylePresetGeneral   IdeogramStylePreset = "general"
	IdeogramStylePresetRealistic IdeogramStylePreset = "realistic"
	IdeogramStylePresetDesign    IdeogramStylePreset = "design"
	IdeogramStylePresetAnime     IdeogramStylePreset = "anime"
	IdeogramStylePresetRender3D  IdeogramStylePreset = "render_3d"
)

// ColorPaletteMember represents a color in a palette
type ColorPaletteMember struct {
	ColorHex    string   `json:"colorHex"`
	ColorWeight *float64 `json:"colorWeight,omitempty"`
}

// ColorPalette represents a color palette
type ColorPalette struct {
	Name    *string              `json:"name,omitempty"`
	Members []ColorPaletteMember `json:"members,omitempty"`
}

// IdeogramProviderSettings contains Ideogram provider settings
type IdeogramProviderSettings struct {
	RenderingSpeed       *IdeogramRenderingSpeed `json:"renderingSpeed,omitempty"`
	MagicPrompt          *bool                   `json:"magicPrompt,omitempty"`
	StyleType            *IdeogramStyleType      `json:"styleType,omitempty"`
	StyleReferenceImages []string                `json:"styleReferenceImages,omitempty"`
	RemixStrength        *float64                `json:"remixStrength,omitempty"`
	StylePreset          *IdeogramStylePreset    `json:"stylePreset,omitempty"`
	StyleCode            *string                 `json:"styleCode,omitempty"`
	ColorPalette         *ColorPalette           `json:"colorPalette,omitempty"`
}

// ProviderSettings contains provider-specific settings
type ProviderSettings struct {
	BFL       *BFLProviderSettings       `json:"bfl,omitempty"`
	ByteDance *ByteDanceProviderSettings `json:"bytedance,omitempty"`
	Ideogram  *IdeogramProviderSettings  `json:"ideogram,omitempty"`
}

// ImageInferenceRequest represents a request for image inference
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

// ImageInferenceResponse represents a response from image inference
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

// UploadImageRequest represents a request to upload an image
type UploadImageRequest struct {
	TaskType     string  `json:"taskType"`
	TaskUUID     string  `json:"taskUUID"`
	ImageBase64  *string `json:"imageBase64,omitempty"`
	ImageDataURI *string `json:"imageDataURI,omitempty"`
	ImageURL     *string `json:"imageURL,omitempty"`
}

// UploadImageResponse represents a response from image upload
type UploadImageResponse struct {
	TaskType  string `json:"taskType"`
	TaskUUID  string `json:"taskUUID"`
	ImageUUID string `json:"imageUUID"`
}

// UpscaleGanRequest represents a request for GAN-based upscaling
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

// UpscaleGanResponse represents a response from GAN upscaling
type UpscaleGanResponse struct {
	TaskType        string   `json:"taskType"`
	TaskUUID        string   `json:"taskUUID"`
	ImageUUID       string   `json:"imageUUID"`
	ImageURL        *string  `json:"imageURL,omitempty"`
	ImageBase64Data *string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    *string  `json:"imageDataURI,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
}

// RemoveImageBackgroundRequest represents a request to remove image background
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

// RemoveImageBackgroundResponse represents a response from background removal
type RemoveImageBackgroundResponse struct {
	TaskType        string   `json:"taskType"`
	TaskUUID        string   `json:"taskUUID"`
	ImageUUID       string   `json:"imageUUID"`
	ImageURL        *string  `json:"imageURL,omitempty"`
	ImageBase64Data *string  `json:"imageBase64Data,omitempty"`
	ImageDataURI    *string  `json:"imageDataURI,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
}

// EnhancePromptRequest represents a request to enhance a prompt
type EnhancePromptRequest struct {
	TaskType        string `json:"taskType"`
	TaskUUID        string `json:"taskUUID"`
	Prompt          string `json:"prompt"`
	PromptMaxLength *int   `json:"promptMaxLength,omitempty"`
	PromptVersions  *int   `json:"promptVersions,omitempty"`
	IncludeCost     *bool  `json:"includeCost,omitempty"`
}

// EnhancePromptResponse represents a response from prompt enhancement
type EnhancePromptResponse struct {
	TaskType string   `json:"taskType"`
	TaskUUID string   `json:"taskUUID"`
	Text     string   `json:"text"`
	Cost     *float64 `json:"cost,omitempty"`
}

// ImageCaptionRequest represents a request for image captioning
type ImageCaptionRequest struct {
	TaskType    string `json:"taskType"`
	TaskUUID    string `json:"taskUUID"`
	InputImage  string `json:"inputImage"`
	IncludeCost *bool  `json:"includeCost,omitempty"`
}

// ImageCaptionResponse represents a response from image captioning
type ImageCaptionResponse struct {
	TaskType string   `json:"taskType"`
	TaskUUID string   `json:"taskUUID"`
	Text     string   `json:"text"`
	Cost     *float64 `json:"cost,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error    string `json:"error"`
	ErrorID  string `json:"errorId,omitempty"`
	TaskUUID string `json:"taskUUID,omitempty"`
	TaskType string `json:"taskType,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Data  []interface{}  `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

// NewImageInferenceRequest creates a new image inference request with required fields
func NewImageInferenceRequest(prompt, model string, width, height int) *ImageInferenceRequest {
	return &ImageInferenceRequest{
		TaskType:       "imageInference",
		TaskUUID:       uuid.New().String(),
		PositivePrompt: prompt,
		Model:          model,
		Width:          width,
		Height:         height,
	}
}

// NewUploadImageRequest creates a new image upload request
func NewUploadImageRequest() *UploadImageRequest {
	return &UploadImageRequest{
		TaskType: "imageUpload",
		TaskUUID: uuid.New().String(),
	}
}

// NewUpscaleGanRequest creates a new upscale request
func NewUpscaleGanRequest(inputImage string, upscaleFactor int) *UpscaleGanRequest {
	return &UpscaleGanRequest{
		TaskType:      "upscaleGan",
		TaskUUID:      uuid.New().String(),
		InputImage:    inputImage,
		UpscaleFactor: upscaleFactor,
	}
}

// NewRemoveImageBackgroundRequest creates a new background removal request
func NewRemoveImageBackgroundRequest(inputImage string) *RemoveImageBackgroundRequest {
	return &RemoveImageBackgroundRequest{
		TaskType:   "imageBackgroundRemoval",
		TaskUUID:   uuid.New().String(),
		InputImage: inputImage,
	}
}

// NewEnhancePromptRequest creates a new prompt enhancement request
func NewEnhancePromptRequest(prompt string) *EnhancePromptRequest {
	return &EnhancePromptRequest{
		TaskType: "promptEnhance",
		TaskUUID: uuid.New().String(),
		Prompt:   prompt,
	}
}

// NewImageCaptionRequest creates a new image caption request
func NewImageCaptionRequest(inputImage string) *ImageCaptionRequest {
	return &ImageCaptionRequest{
		TaskType:   "imageCaption",
		TaskUUID:   uuid.New().String(),
		InputImage: inputImage,
	}
}

// ========================================
// Video Inference Types
// ========================================

// VideoOutputFormat specifies the video file format
type VideoOutputFormat string

const (
	VideoOutputFormatMP4  VideoOutputFormat = "MP4"
	VideoOutputFormatWEBM VideoOutputFormat = "WEBM"
)

// FramePosition specifies the position of a frame in the video
type FramePosition string

const (
	FramePositionFirst FramePosition = "first"
	FramePositionLast  FramePosition = "last"
)

// TaskStatus represents the status of an async task
type TaskStatus string

const (
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusError      TaskStatus = "error"
)

// FrameImage defines an image constraint for a specific frame position
type FrameImage struct {
	InputImage string        `json:"inputImage"`
	Frame      FramePosition `json:"frame"`
}

// ReferenceImage defines a reference image for video generation
type ReferenceImage struct {
	InputImage string `json:"inputImage"`
}

// ReferenceVideo defines a reference video for video generation
type ReferenceVideo struct {
	InputVideo string `json:"inputVideo"`
}

// InputAudio defines an audio input for video generation
type InputAudio struct {
	InputAudio string `json:"inputAudio"`
}

// Speech defines text-to-speech parameters for video generation
type Speech struct {
	Voice string `json:"voice"` // Voice ID for speech synthesis
	Text  string `json:"text"`  // Text to convert to speech
}

// GoogleVideoSettings contains Google (Veo) provider-specific settings
type GoogleVideoSettings struct {
	EnhancePrompt *bool `json:"enhancePrompt,omitempty"`
	GenerateAudio *bool `json:"generateAudio,omitempty"`
}

// ByteDanceVideoSettings contains ByteDance provider-specific settings
type ByteDanceVideoSettings struct {
	CameraFixed *bool `json:"cameraFixed,omitempty"`
}

// MiniMaxVideoSettings contains MiniMax provider-specific settings
type MiniMaxVideoSettings struct {
	PromptOptimizer *bool `json:"promptOptimizer,omitempty"`
}

// PixVerseVideoSettings contains PixVerse provider-specific settings
type PixVerseVideoSettings struct {
	Style              *string `json:"style,omitempty"`              // "realistic", "anime", "3d"
	Effect             *string `json:"effect,omitempty"`             // Effect template
	CameraMovement     *string `json:"cameraMovement,omitempty"`     // Camera movement type
	MotionMode         *string `json:"motionMode,omitempty"`         // Motion mode setting
	SoundEffectSwitch  *bool   `json:"soundEffectSwitch,omitempty"`  // Enable sound effects
	SoundEffectContent *string `json:"soundEffectContent,omitempty"` // Sound effect content
}

// ViduVideoSettings contains Vidu provider-specific settings
type ViduVideoSettings struct {
	MovementAmplitude *string `json:"movementAmplitude,omitempty"` // "auto", "small", "medium", "large"
	BGM               *bool   `json:"bgm,omitempty"`               // Background music (4s videos only)
	Style             *string `json:"style,omitempty"`             // "general", "anime"
}

// VideoProviderSettings contains provider-specific settings for video generation
type VideoProviderSettings struct {
	Google    *GoogleVideoSettings    `json:"google,omitempty"`
	ByteDance *ByteDanceVideoSettings `json:"bytedance,omitempty"`
	MiniMax   *MiniMaxVideoSettings   `json:"minimax,omitempty"`
	PixVerse  *PixVerseVideoSettings  `json:"pixverse,omitempty"`
	Vidu      *ViduVideoSettings      `json:"vidu,omitempty"`
}

// VideoInferenceRequest represents a request for video inference
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

// VideoInferenceResponse represents a response from video inference
type VideoInferenceResponse struct {
	TaskType  string     `json:"taskType"`
	TaskUUID  string     `json:"taskUUID"`
	Status    TaskStatus `json:"status,omitempty"`
	VideoUUID string     `json:"videoUUID,omitempty"`
	VideoURL  *string    `json:"videoURL,omitempty"`
	Seed      *int64     `json:"seed,omitempty"`
	Cost      *float64   `json:"cost,omitempty"`
}

// GetResponseRequest represents a request to get the status/result of an async task
type GetResponseRequest struct {
	TaskType string `json:"taskType"`
	TaskUUID string `json:"taskUUID"`
}

// NewVideoInferenceRequest creates a new video inference request with required fields
func NewVideoInferenceRequest(prompt, model string) *VideoInferenceRequest {
	asyncMethod := DeliveryMethodPOST // Default to async for video
	return &VideoInferenceRequest{
		TaskType:       "videoInference",
		TaskUUID:       uuid.New().String(),
		PositivePrompt: prompt,
		Model:          model,
		DeliveryMethod: &asyncMethod,
	}
}

// NewGetResponseRequest creates a new get response request
func NewGetResponseRequest(taskUUID string) *GetResponseRequest {
	return &GetResponseRequest{
		TaskType: "getResponse",
		TaskUUID: taskUUID,
	}
}
