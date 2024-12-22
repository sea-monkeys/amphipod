package mcp

import (
	"io"
	"net/http"
)

func ToolsList(serviceURL string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", serviceURL+"/tools/list", nil)
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
		return "", err
	}

	return string(body), nil
}
