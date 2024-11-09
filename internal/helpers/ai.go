package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/huggingface"
)

type GeminiReqPayload struct {
	Contents []Content `json:"contents"`
}

type GeminiResp struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

func ParseGeminiOutput(input string) string {
	var result strings.Builder

	processed := strings.ReplaceAll(input, "\\`", "`")
	lines := strings.Split(processed, "\n")

	lastLineIndex := len(lines) - 1
	if strings.TrimSpace(lines[lastLineIndex]) == "```" {
		lines = lines[:lastLineIndex]
	}

	for i, line := range lines {
		result.WriteString(line)

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func QueryGemini(apiKey string, prompt string) (string, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", apiKey)

	payload := GeminiReqPayload{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(body))
	}

	var data GeminiResp
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	parsedOutput := ParseGeminiOutput(data.Candidates[0].Content.Parts[0].Text)
	return parsedOutput, nil
}

func QueryHuggingFace(apiKey, modelName, prompt string) (string, error) {
	llm, err := huggingface.New(
		huggingface.WithToken(apiKey),
		huggingface.WithModel(modelName),
	)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		prompt,
		llms.WithTemperature(0.1),
	)
	if err != nil {
		return "", err
	}

	return completion, err
}
