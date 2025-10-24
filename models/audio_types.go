package models

type AudioOutputFormat string

const (
	AudioOutputFormatMP3 AudioOutputFormat = "MP3"
)

type AudioSettings struct {
	SampleRate *int `json:"sampleRate,omitempty"`
	Bitrate    *int `json:"bitrate,omitempty"`
}

type ElevenLabsMusicSettings struct {
	PromptInfluence *float64 `json:"promptInfluence,omitempty"`
}
type ElevenLabsAudioSettings struct {
	Music *ElevenLabsMusicSettings `json:"music,omitempty"`
}

type AudioProviderSettings struct {
	ElevenLabs *ElevenLabsAudioSettings `json:"elevenlabs,omitempty"`
}

type AudioInferenceRequest struct {
	TaskType         string                 `json:"taskType"`
	TaskUUID         string                 `json:"taskUUID"`
	OutputType       *OutputType            `json:"outputType,omitempty"`
	OutputFormat     *AudioOutputFormat     `json:"outputFormat,omitempty"`
	WebhookURL       *string                `json:"webhookURL,omitempty"`
	DeliveryMethod   *DeliveryMethod        `json:"deliveryMethod,omitempty"`
	UploadEndpoint   *string                `json:"uploadEndpoint,omitempty"`
	TTL              *int                   `json:"ttl,omitempty"`
	IncludeCost      *bool                  `json:"includeCost,omitempty"`
	PositivePrompt   string                 `json:"positivePrompt"`
	Model            string                 `json:"model"`
	Duration         *int                   `json:"duration,omitempty"`
	NumberResults    *int                   `json:"numberResults,omitempty"`
	AudioSettings    *AudioSettings         `json:"audioSettings,omitempty"`
	ProviderSettings *AudioProviderSettings `json:"providerSettings,omitempty"`
}

type AudioInferenceResponse struct {
	TaskType        string     `json:"taskType"`
	TaskUUID        string     `json:"taskUUID"`
	Status          TaskStatus `json:"status,omitempty"`
	AudioUUID       string     `json:"audioUUID,omitempty"`
	AudioURL        *string    `json:"audioURL,omitempty"`
	AudioBase64Data *string    `json:"audioBase64Data,omitempty"`
	AudioDataURI    *string    `json:"audioDataURI,omitempty"`
	Cost            *float64   `json:"cost,omitempty"`
}
