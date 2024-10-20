package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Choice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		QueueTime        float64 `json:"queue_time"`
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func MakeRequest(messages []Message, model string) ([]Message, error) {
	request := Request{
		Messages: messages,
		Model:    model,
	}

	// Make the request
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create a new HTTP request
	url := "https://api.groq.com/openai/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	groq_key := os.Getenv("GROQ_API_KEY")

	// Set the API key in the Authorization header
	req.Header.Set("Authorization", "Bearer "+groq_key)

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Making request to Groq")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(string(body))

	// Unmarshal the response body into a struct
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(response)

	return []Message{response.Choices[0].Message}, nil
}

func FixSecurity(startLine int, endLine int, description string, code string) (string, error) {
	prompt := []Message{
		{
			Role:    "assistant",
			Content: "You are a security expert who is editing vulnerable code. You will receive code with a vulnerability in it. Additionally you will get a description of the vulnerability and lines the vulnerability are on. You must return the entire input code with the vulnerability secured. You should only return that.",
		},
		{
			Role: "user",
			Content: "I have a security issue described as follows: " + description +
				". The code snippet is as follows: \n```" + code +
				"```\n The issue is on lines " + fmt.Sprint(startLine) + " to " + fmt.Sprint(endLine) + ".",
		},
	}

	messages, err := MakeRequest(prompt, "llama-3.2-90b-text-preview")

	if err != nil {
		return "", err
	}

	if len(messages) == 0 {
		return "", fmt.Errorf("No response from Groq")
	}

	fix := messages[len(messages)-1].Content

	return fix, nil

}
