package models

import (
	"encoding/json"
	"testing"
)

func TestOutputTypeConstants(t *testing.T) {
	tests := []struct {
		name  string
		value OutputType
		want  string
	}{
		{"URL", OutputTypeURL, "URL"},
		{"Base64Data", OutputTypeBase64Data, "base64Data"},
		{"DataURI", OutputTypeDataURI, "dataURI"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.want {
				t.Errorf("OutputType = %v, want %v", tt.value, tt.want)
			}
		})
	}
}

func TestSchedulerConstants(t *testing.T) {
	schedulers := []Scheduler{
		SchedulerEuler,
		SchedulerEulerA,
		SchedulerDPMPP2M,
		SchedulerDPMPP2MKarras,
		SchedulerDPMPPSDE,
		SchedulerDPMPPSDEKarras,
		SchedulerLMS,
		SchedulerLMSKarras,
		SchedulerHeun,
		SchedulerDDIM,
		SchedulerPNDM,
	}

	for _, s := range schedulers {
		if s == "" {
			t.Errorf("Scheduler constant is empty")
		}
	}
}

func TestImageInferenceRequestJSON(t *testing.T) {
	req := NewImageInferenceRequest("test prompt", "test-model", 512, 512)

	steps := 30
	req.Steps = &steps

	cfg := 7.5
	req.CFGScale = &cfg

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var decoded ImageInferenceRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if decoded.PositivePrompt != req.PositivePrompt {
		t.Errorf("Decoded PositivePrompt = %v, want %v", decoded.PositivePrompt, req.PositivePrompt)
	}

	if decoded.Steps == nil || *decoded.Steps != *req.Steps {
		t.Errorf("Decoded Steps = %v, want %v", decoded.Steps, req.Steps)
	}

	if decoded.CFGScale == nil || *decoded.CFGScale != *req.CFGScale {
		t.Errorf("Decoded CFGScale = %v, want %v", decoded.CFGScale, req.CFGScale)
	}
}

func TestControlNetJSON(t *testing.T) {
	weight := 0.8
	cn := ControlNet{
		Model:      "test-model",
		GuideImage: "test-uuid",
		Weight:     &weight,
	}

	data, err := json.Marshal(cn)
	if err != nil {
		t.Fatalf("Failed to marshal ControlNet: %v", err)
	}

	var decoded ControlNet
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ControlNet: %v", err)
	}

	if decoded.Model != cn.Model {
		t.Errorf("Decoded Model = %v, want %v", decoded.Model, cn.Model)
	}

	if decoded.GuideImage != cn.GuideImage {
		t.Errorf("Decoded GuideImage = %v, want %v", decoded.GuideImage, cn.GuideImage)
	}

	if decoded.Weight == nil || *decoded.Weight != *cn.Weight {
		t.Errorf("Decoded Weight = %v, want %v", decoded.Weight, cn.Weight)
	}
}

func TestLoRAJSON(t *testing.T) {
	weight := 0.95
	lora := LoRA{
		Model:  "test-lora",
		Weight: &weight,
	}

	data, err := json.Marshal(lora)
	if err != nil {
		t.Fatalf("Failed to marshal LoRA: %v", err)
	}

	var decoded LoRA
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal LoRA: %v", err)
	}

	if decoded.Model != lora.Model {
		t.Errorf("Decoded Model = %v, want %v", decoded.Model, lora.Model)
	}

	if decoded.Weight == nil || *decoded.Weight != *lora.Weight {
		t.Errorf("Decoded Weight = %v, want %v", decoded.Weight, lora.Weight)
	}
}

func TestOutpaintJSON(t *testing.T) {
	outpaint := Outpaint{
		Top:    128,
		Right:  64,
		Bottom: 128,
		Left:   64,
		Blur:   10,
	}

	data, err := json.Marshal(outpaint)
	if err != nil {
		t.Fatalf("Failed to marshal Outpaint: %v", err)
	}

	var decoded Outpaint
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Outpaint: %v", err)
	}

	if decoded.Top != outpaint.Top {
		t.Errorf("Decoded Top = %v, want %v", decoded.Top, outpaint.Top)
	}

	if decoded.Right != outpaint.Right {
		t.Errorf("Decoded Right = %v, want %v", decoded.Right, outpaint.Right)
	}
}

func TestColorPaletteJSON(t *testing.T) {
	weight1 := 1.0
	weight2 := 0.7

	palette := ColorPalette{
		Members: []ColorPaletteMember{
			{ColorHex: "#FF5733", ColorWeight: &weight1},
			{ColorHex: "#C70039", ColorWeight: &weight2},
		},
	}

	data, err := json.Marshal(palette)
	if err != nil {
		t.Fatalf("Failed to marshal ColorPalette: %v", err)
	}

	var decoded ColorPalette
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ColorPalette: %v", err)
	}

	if len(decoded.Members) != len(palette.Members) {
		t.Errorf("Decoded Members length = %v, want %v", len(decoded.Members), len(palette.Members))
	}

	if decoded.Members[0].ColorHex != palette.Members[0].ColorHex {
		t.Errorf("Decoded ColorHex = %v, want %v", decoded.Members[0].ColorHex, palette.Members[0].ColorHex)
	}
}
