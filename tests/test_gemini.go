package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/williamfotso/acc/internal/services/ai"
	"google.golang.org/genai"
)

var (
	transcriptPrompt = `
	Generate a class transcript for a college course on a Mathematics topic. 
	I only want the transcript, no annotations, no schema.
	Include only words that are spoken by the professor.
	Don't add annotations or anything else. Just the transcript.
	The class is 2 hours long and the professor is speaking for the entire time at a pace of 100 words per minute.
	`
)

func GenerateTranscript() string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: viper.GetString("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(transcriptPrompt),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Write the result to transcript.txt
	return result.Text()
}

func main() {
	text := GenerateTranscript()
	ai.GenerateSummary(text, nil)

}
