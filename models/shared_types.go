package models

// Shared enums and basic types used by multiple domains

// OutputType specifies the format of the output image/audio/video reference
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
	DeliveryMethodAsync  DeliveryMethod = "async"
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

// Task types (string values used on the wire)
const (
	TaskTypeImageInference         = "imageInference"
	TaskTypeVideoInference         = "videoInference"
	TaskTypeAudioInference         = "audioInference"
	TaskTypePromptEnhance          = "promptEnhance"
	TaskTypeImageCaption           = "imageCaption"
	TaskTypeImageUpload            = "imageUpload"
	TaskTypeUpscaleGan             = "imageUpscale"
	TaskTypeImageBackgroundRemoval = "imageBackgroundRemoval"
	TaskTypeGetResponse            = "getResponse"
)

// Acceleration specifies acceleration mode
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

// Provider-specific (image)
type BFLProviderSettings struct {
	PromptUpsampling *bool `json:"promptUpsampling,omitempty"`
	SafetyTolerance  *int  `json:"safetyTolerance,omitempty"`
	Raw              *bool `json:"raw,omitempty"`
}

type ByteDanceProviderSettings struct {
	MaxSequentialImages *int `json:"maxSequentialImages,omitempty"`
}

// Ideogram settings
type IdeogramRenderingSpeed string

const (
	IdeogramRenderingSpeedFast IdeogramRenderingSpeed = "fast"
	IdeogramRenderingSpeedSlow IdeogramRenderingSpeed = "slow"
)

type IdeogramStyleType string

const (
	IdeogramStyleTypeAuto      IdeogramStyleType = "auto"
	IdeogramStyleTypeGeneral   IdeogramStyleType = "general"
	IdeogramStyleTypeRealistic IdeogramStyleType = "realistic"
	IdeogramStyleTypeDesign    IdeogramStyleType = "design"
	IdeogramStyleTypeRender3D  IdeogramStyleType = "render_3d"
	IdeogramStyleTypeAnime     IdeogramStyleType = "anime"
)

type IdeogramStylePreset string

const (
	IdeogramStylePresetGeneral   IdeogramStylePreset = "general"
	IdeogramStylePresetRealistic IdeogramStylePreset = "realistic"
	IdeogramStylePresetDesign    IdeogramStylePreset = "design"
	IdeogramStylePresetAnime     IdeogramStylePreset = "anime"
	IdeogramStylePresetRender3D  IdeogramStylePreset = "render_3d"
)

type ColorPaletteMember struct {
	ColorHex    string   `json:"colorHex"`
	ColorWeight *float64 `json:"colorWeight,omitempty"`
}
type ColorPalette struct {
	Name    *string              `json:"name,omitempty"`
	Members []ColorPaletteMember `json:"members,omitempty"`
}

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

type ProviderSettings struct {
	BFL       *BFLProviderSettings       `json:"bfl,omitempty"`
	ByteDance *ByteDanceProviderSettings `json:"bytedance,omitempty"`
	Ideogram  *IdeogramProviderSettings  `json:"ideogram,omitempty"`
}

// Async task status
type TaskStatus string

const (
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusError      TaskStatus = "error"
)

// Generic API envelope
type ErrorResponse struct {
	Error    string `json:"error"`
	ErrorID  string `json:"errorId,omitempty"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message,omitempty"`
	TaskUUID string `json:"taskUUID,omitempty"`
	TaskType string `json:"taskType,omitempty"`
}

type APIResponse struct {
	Data  []interface{}  `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}
