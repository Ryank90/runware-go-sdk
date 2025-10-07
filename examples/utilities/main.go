package main

import (
	"context"
	"fmt"
	"log"

	runware "github.com/Ryank90/runware-go-sdk"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	client, err := runware.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to Runware API")

	// Example 1: Enhance a prompt
	fmt.Println("\n=== Prompt Enhancement ===")
	enhanceReq := runware.NewEnhancePromptRequest("a beautiful sunset")
	enhanceResp, err := client.EnhancePrompt(ctx, enhanceReq)
	if err != nil {
		log.Printf("Failed to enhance prompt: %v", err)
	} else {
		fmt.Printf("Original: a beautiful sunset\n")
		fmt.Printf("Enhanced: %s\n", enhanceResp.Text)
	}

	// Example 2: Upload and process an image
	fmt.Println("\n=== Image Upload ===")
	uploadResp, err := client.UploadImageFromURL(ctx, "https://example.com/image.jpg")
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
	} else {
		fmt.Printf("Image uploaded: %s\n", uploadResp.ImageUUID)

		// Example 3: Remove background
		fmt.Println("\n=== Background Removal ===")
		bgRemovalReq := runware.NewRemoveImageBackgroundRequest(uploadResp.ImageUUID)
		includeCost := true
		bgRemovalReq.IncludeCost = &includeCost

		bgRemovalResp, err := client.RemoveBackground(ctx, bgRemovalReq)
		if err != nil {
			log.Printf("Failed to remove background: %v", err)
		} else {
			fmt.Printf("Background removed: %s\n", bgRemovalResp.ImageUUID)
			if bgRemovalResp.ImageURL != nil {
				fmt.Printf("Result URL: %s\n", *bgRemovalResp.ImageURL)
			}
		}

		// Example 4: Upscale image
		fmt.Println("\n=== Image Upscaling ===")
		upscaleReq := runware.NewUpscaleGanRequest(uploadResp.ImageUUID, 4)
		upscaleReq.IncludeCost = &includeCost

		upscaleResp, err := client.UpscaleImage(ctx, upscaleReq)
		if err != nil {
			log.Printf("Failed to upscale image: %v", err)
		} else {
			fmt.Printf("Image upscaled: %s\n", upscaleResp.ImageUUID)
			if upscaleResp.ImageURL != nil {
				fmt.Printf("Result URL: %s\n", *upscaleResp.ImageURL)
			}
		}

		// Example 5: Caption image
		fmt.Println("\n=== Image Captioning ===")
		captionReq := runware.NewImageCaptionRequest(uploadResp.ImageUUID)
		captionResp, err := client.CaptionImage(ctx, captionReq)
		if err != nil {
			log.Printf("Failed to caption image: %v", err)
		} else {
			fmt.Printf("Image caption: %s\n", captionResp.Text)
		}
	}
}
