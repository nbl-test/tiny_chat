package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	setHandler("gpt_proxy", chatgptProxyHandler)
}

// {"id":"chatcmpl-6viCY5M91KBSHAwYTMTSCB8sho5la","object":"chat.completion","created":1679212626,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":13,"completion_tokens":33,"total_tokens":46},"choices":[{"message":{"role":"assistant","content":"\n\nAs an AI language model, I do not have emotions or feelings, but I am functioning well. Thank you for asking. How can I assist you today?"},"finish_reason":"stop","index":0}]}
type ChatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func chatgptProxyHandler(inputMsg string) (string, error) {
	// Create a map[string]interface{} for the body
	bodyData := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []interface{}{
			map[string]interface{}{"role": "user", "content": inputMsg},
		},
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, "https://chatgpt-api.shn.hk/v1/", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set the content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r ChatCompletion
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyBytes = []byte(strings.TrimPrefix(string(bodyBytes), "OpenAI API responded:"))
	fmt.Println("response is", string(bodyBytes))
	err = json.Unmarshal(bodyBytes, &r)
	if err != nil {
		return "", err
	}
	if r.Error.Message != "" {
		return "", errors.New(r.Error.Message)
	}
	if len(r.Choices) == 0 {
		return "nothing to say", nil
	}
	return r.Choices[0].Message.Content, nil
}
