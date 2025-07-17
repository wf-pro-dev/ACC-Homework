package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
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

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	//GenerateTranscript()

	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: viper.GetString("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	var content string
	file, err := os.Open("transcript.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	content = string(byteContent)

	summaryPrompt := `
	Generate 4 fields based on the transcript:
	1. Title
	2. Full name 
	3. Key Points
	4. Summary

	Title should be a short title for the note (Maximum 5 words).
	Full name should be the full name for the note (Minimum 10 words).
	Keywords should be a list of key points from the transcript (Maximum 5 keywords).
	
	Summary should be a markdown summary of the transcript:
	Use a visual representation when you can.
	Format into well defined section. 
	Use markdown formatting.
	Use bold when you need to highlight something.
	Use italic when you need to emphasize something.
	Use code blocks when you need to show code.
	Use links when you need to reference something.
	Use images when you need to show a diagram.
	

	The transcript is:
	` + content

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(summaryPrompt),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"title":    {Type: genai.TypeString},
					"name":     {Type: genai.TypeString},
					"keywords": {Type: genai.TypeString},
					"summary":  {Type: genai.TypeString},
				},
				PropertyOrdering: []string{"title", "name", "keywords", "summary"},
			},
		},
	)
	if err != nil {
		log.Fatal("Error generating summary: ", err)
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal([]byte(result.Text()), &resultMap)
	if err != nil {
		log.Fatal("Error unmarshalling summary result: ", err)
	}
	fmt.Println(resultMap["title"])
	fmt.Println(resultMap["name"])
	fmt.Println(resultMap["keywords"])

	// Write the result to transcript.txt
	outputFile, err := os.Create("summary.md")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(resultMap["summary"].(string))
	if err != nil {
		log.Fatal("Error writing summary : ", err)
	}

}

func GenerateTranscript() {
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
	outputFile, err := os.Create("transcript.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(result.Text())
	if err != nil {
		log.Fatal(err)
	}
}
