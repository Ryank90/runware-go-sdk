package models

// Internal lightweight interfaces for performance (used by transport)
type TaskIdentifiable interface {
	GetTaskUUID() string
	GetTaskType() string
}
type ResultCountProvider interface{ GetNumberResults() *int }

// Implementations on request types
func (r *ImageInferenceRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *ImageInferenceRequest) GetTaskType() string    { return r.TaskType }
func (r *ImageInferenceRequest) GetNumberResults() *int { return r.NumberResults }

func (r *UploadImageRequest) GetTaskUUID() string { return r.TaskUUID }
func (r *UploadImageRequest) GetTaskType() string { return r.TaskType }

func (r *UpscaleGanRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *UpscaleGanRequest) GetTaskType() string    { return r.TaskType }
func (r *UpscaleGanRequest) GetNumberResults() *int { return nil }

func (r *RemoveImageBackgroundRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *RemoveImageBackgroundRequest) GetTaskType() string    { return r.TaskType }
func (r *RemoveImageBackgroundRequest) GetNumberResults() *int { return nil }

func (r *EnhancePromptRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *EnhancePromptRequest) GetTaskType() string    { return r.TaskType }
func (r *EnhancePromptRequest) GetNumberResults() *int { return nil }

func (r *ImageCaptionRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *ImageCaptionRequest) GetTaskType() string    { return r.TaskType }
func (r *ImageCaptionRequest) GetNumberResults() *int { return nil }

func (r *VideoInferenceRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *VideoInferenceRequest) GetTaskType() string    { return r.TaskType }
func (r *VideoInferenceRequest) GetNumberResults() *int { return r.NumberResults }

func (r *AudioInferenceRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *AudioInferenceRequest) GetTaskType() string    { return r.TaskType }
func (r *AudioInferenceRequest) GetNumberResults() *int { return r.NumberResults }

func (r *GetResponseRequest) GetTaskUUID() string    { return r.TaskUUID }
func (r *GetResponseRequest) GetTaskType() string    { return r.TaskType }
func (r *GetResponseRequest) GetNumberResults() *int { return nil }
