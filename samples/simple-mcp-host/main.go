package main

import (
	"mcp-host/mcp"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

var (
	FALSE = false
	TRUE  = true
)

func main() {
	ctx := context.Background()

	var ollamaRawUrl string
	if ollamaRawUrl = os.Getenv("OLLAMA_HOST"); ollamaRawUrl == "" {
		ollamaRawUrl = "http://localhost:11434"
	}

	var mcpServerUrl string
	if mcpServerUrl = os.Getenv("MCP_SERVER_HOST"); mcpServerUrl == "" {
		mcpServerUrl = "http://localhost:8080"
	}

	var toolsLlm string
	if toolsLlm = os.Getenv("TOOLS_LLM"); toolsLlm == "" {
		toolsLlm = "allenporter/xlam:1b"
	}

	strToolsList, _ := mcp.ToolsList(mcpServerUrl)
	strOllamaToolsList, _ := mcp.TransformToOllamaToolsFormat(strToolsList)
	//fmt.Println("🔧", strOllamaToolsList)

	url, _ := url.Parse(ollamaRawUrl)

	client := api.NewClient(url, http.DefaultClient)

	// transform strOllamaToolsList to api.Tools
	var toolsList api.Tools
	jsonErr := json.Unmarshal([]byte(strOllamaToolsList), &toolsList)
	if jsonErr != nil {
		log.Fatalln("😡", jsonErr)
	}


	// Prompt construction
	messages := []api.Message{
		{Role: "user", Content: "Say hello to Bob Morane"},
		{Role: "user", Content: "add 28 to 12"},
		{Role: "user", Content: "Say goodbye to Sarah Connor"},
		{Role: "user", Content: "Say goodbye to Jane Doe"},
		{Role: "user", Content: "Say hello to John Doe"},
		{Role: "user", Content: "add 35 to 7"},
	}

	req := &api.ChatRequest{
		Model:    toolsLlm,
		Messages: messages,
		Options: map[string]interface{}{
			"temperature":   0.0,
			"repeat_last_n": 2,
		},
		Tools:  toolsList,
		Stream: &FALSE,
	}

	
	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		
		for _, toolCall := range resp.Message.ToolCalls {
			fmt.Println("🔧", toolCall.Function.Name, toolCall.Function.Arguments)

			data := map[string]interface{}{
				"name": toolCall.Function.Name,
				"arguments": toolCall.Function.Arguments,
			}

			mcpResp, err := mcp.ToolsCall(mcpServerUrl, data)
			if err != nil {
				log.Fatalln("😡", err)
			}
			fmt.Println("🎉", mcpResp)

		}

		return nil
	})

	if err != nil {
		log.Fatalln("😡", err)
	}
	
}
