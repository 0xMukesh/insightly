package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GeminiReqPayload struct {
	Contents []struct {
		Parts []struct {
			Text string
		}
	}
}

type GeminiResp struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string
			}
		}
	}
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

func SendReqToGemini(apiKey string, prompt string) (string, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", apiKey)

	payload := fmt.Sprintf(`{"contents":[{"parts":[{"text":"%s"}]}]}`, prompt)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
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
