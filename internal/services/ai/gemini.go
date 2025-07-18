package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"google.golang.org/genai"
)

type Notes struct {
	Title    string
	Keywords string
}

type File struct {
	Path string
	Text string
}

func GenerateSummary(transcript string, files []File) *Notes {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: viper.GetString("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	summaryPrompt := fmt.Sprintf(`
	Generate 3 fields based on the transcript:
	1. Title
	2. Key Points
	3. Summary

	Title should be a short title for the note (Maximum 5 words).
	Key Points should be a list of key points from the transcript (Maximum 5 keywords).

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
	%s
`, transcript)
	if len(files) > 0 {
		summaryPrompt += "this is supporting files for the transcript:\n"
		for _, file := range files {
			summaryPrompt += fmt.Sprintf(`
			The file %s is:
			%s
		`, file.Path, file.Text)
		}
	}

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
					"keywords": {Type: genai.TypeString},
					"summary":  {Type: genai.TypeString},
				},
				PropertyOrdering: []string{"title", "keywords", "summary"},
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

	return &Notes{
		Title:    resultMap["title"].(string),
		Keywords: resultMap["keywords"].(string),
	}

}
