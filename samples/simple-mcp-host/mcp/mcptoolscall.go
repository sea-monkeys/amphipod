package mcp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func ToolsCall(serviceURL string, params map[string]interface{}) (string, error) {

	// Create the request data
	/*
	   data := map[string]interface{}{
	       "name": "say_goodbye",
	       "arguments": map[string]string{
	           "name": "Jane Doe",
	       },
	   }
	*/

	// Convert data to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	//fmt.Printf("Sending: %s on %s\n\n", string(jsonData), serviceURL)

	// Create the request
	req, err := http.NewRequest("POST", serviceURL+"/tools/call", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}
