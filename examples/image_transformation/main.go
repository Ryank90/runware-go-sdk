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

	// First, upload an image
	fmt.Println("Uploading image...")
	uploadResp, err := client.UploadImageFromURL(ctx, "https://example.com/image.jpg")
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}

	fmt.Printf("Image uploaded with UUID: %s\n", uploadResp.ImageUUID)

	// Transform the image with a prompt
	prompt := "a watercolor painting style, soft brushstrokes, artistic interpretation"
	model := "civitai:139562@297320"
	strength := 0.7

	fmt.Println("Transforming image...")
	response, err := client.ImageToImage(ctx, prompt, model, uploadResp.ImageUUID, 1024, 1024, strength)
	if err != nil {
		log.Fatalf("Failed to transform image: %v", err)
	}

	fmt.Printf("Image transformed successfully!\n")
	fmt.Printf("Result UUID: %s\n", response.ImageUUID)
	if response.ImageURL != nil {
		fmt.Printf("Result URL: %s\n", *response.ImageURL)
	}
}
