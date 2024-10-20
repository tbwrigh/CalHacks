package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GPTContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type GPTMessage struct {
	Role    string       `json:"role"`
	Content []GPTContent `json:"content"`
}

type GPTRequest struct {
	Messages         []GPTMessage `json:"messages"`
	Model            string       `json:"model"`
	Temperature      float64      `json:"temperature"`
	MaxTokens        int          `json:"max_tokens"`
	TopP             float64      `json:"top_p"`
	FrequencyPenalty float64      `json:"frequency_penalty"`
	PresencePenalty  float64      `json:"presence_penalty"`
	ResponseFormat   struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

type GPTResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GPTResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Index   int                `json:"index"`
		Message GPTResponseMessage `json:"message"`
	}
}

func BuildDefGPTRequest(messages []GPTMessage, model string) GPTRequest {
	return GPTRequest{
		Messages:         messages,
		Model:            model,
		Temperature:      1,
		MaxTokens:        2048,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		ResponseFormat: struct {
			Type string `json:"type"`
		}{
			Type: "text",
		},
	}
}

func MakeGPTRequest(messages []GPTMessage, model string) ([]GPTResponseMessage, error) {
	request := BuildDefGPTRequest(messages, model)

	// Make the request
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Set the Authorization header
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client
	client := &http.Client{}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Unmarshal the response body into a struct
	var response GPTResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(response)

	return []GPTResponseMessage{response.Choices[0].Message}, nil
}

func GPTFixSecurity(startLine, endLine int, description, code string) (string, error) {
	messages := []GPTMessage{
		{
			Role: "system",
			Content: []GPTContent{
				{
					Type: "text",
					Text: "You are a security expert who is editing vulnerable code. You will receive code with a vulnerability in it. Additionally you will get a description of the vulnerability and lines the vulnerability are on. You must return the entire input code with the vulnerability secured. You should only return that.",
				},
			},
		},
		{
			Role: "user",
			Content: []GPTContent{
				{
					Type: "text",
					Text: "I have a security issue described as follows: " + description +
						". The code snippet is as follows: \n```" + code +
						"\n```\n The issue is on lines " + fmt.Sprint(startLine) + " to " + fmt.Sprint(endLine) + ".",
				},
			},
		},
	}

	response, err := MakeGPTRequest(messages, "gpt-4o")
	if err != nil {
		return "", err
	}

	return response[0].Content, nil
}
