package mcp

import (
	"encoding/json"
)

// Input JSON structure
type InputJSON struct {
	Result struct {
		Tools []InputTool `json:"tools"`
	} `json:"result"`
}

type InputTool struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	InputSchema JSONSchema `json:"inputSchema"`
}

type JSONSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]Property    `json:"properties"`
	Required   []string              `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Output JSON structure
type OutputTool struct {
	Type     string        `json:"type"`
	Function FunctionSpec  `json:"function"`
}

type FunctionSpec struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  JSONSchema  `json:"parameters"`
}

func TransformToOllamaToolsFormat(inputJSON string) (string, error) {
	// decode input JSON
	var input InputJSON
	if err := json.Unmarshal([]byte(inputJSON), &input); err != nil {
		return "", err
	}

	// Create output tools slice
	var outputTools []OutputTool

	// transform every tool
	for _, tool := range input.Result.Tools {
		outputTool := OutputTool{
			Type: "function",
			Function: FunctionSpec{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters: JSONSchema{
					Type:       tool.InputSchema.Type,
					Properties: tool.InputSchema.Properties,
					Required:   tool.InputSchema.Required,
				},
			},
		}
		outputTools = append(outputTools, outputTool)
	}

	// Encode output JSON with indentation
	output, err := json.MarshalIndent(outputTools, "", "  ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}
