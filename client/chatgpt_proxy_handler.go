package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func init() {
	setHandler("gpt_proxy", chatgptProxyHandler)
}

func chatgptProxyHandler(inputMsg string) (string, error) {
	// Create a map[string]interface{} for the body
	bodyData := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	// Convert the map to JSON
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodGet, "https://chatgpt-api.shn.hk/v1/", bytes.NewBuffer(jsonData))
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

	type error struct {
		Message string
		Type    string
	}
	type reply struct {
		Error error
	}

	var r reply
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}
	if r.Error.Message != "" {
		return "", errors.New(r.Error.Message)
	}
	return "nothing to say", nil
}
