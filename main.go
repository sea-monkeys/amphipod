package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

/*
WASM Plugins
store all your plugins in a normal Go hash map, protected by a Mutex
(reproduce something like the node.js event loop)
to avoid "memory collision 💥"
*/
var m sync.Mutex
var plugins = make(map[string]*extism.Plugin)

func StorePlugin(plugin *extism.Plugin) {
	plugins["code"] = plugin
}

func GetPlugin() (extism.Plugin, error) {
	if plugin, ok := plugins["code"]; ok {
		return *plugin, nil
	} else {
		return extism.Plugin{}, errors.New("🔴 no plugin")
	}
}

/*
GetBytesBody returns the body of an HTTP request as a []byte.
  - It takes a pointer to an http.Request as a parameter.
  - It returns a []byte.
*/
func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {

	ctx := context.Background()

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// read the content of the mcp.list.json file and store it in a variable
	mcpTools, err := os.ReadFile("tools/mcp.list.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	mux := http.NewServeMux()

	// https://spec.modelcontextprotocol.io/specification/server/tools/#protocol-messages
	mux.HandleFunc("GET /tools/list", func(response http.ResponseWriter, request *http.Request) {
		// Set headers for JSON and UTF-8
		response.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Set success status code (200 OK)
		response.WriteHeader(http.StatusOK)

		// Create and send JSON response
		response.Write(mcpTools)
	})

	/*
		Request
			{
				"name": "get_weather",
				"arguments": {
					"location": "New York"
				}
			}
		Response:
			{
				"result": {
					"content": [{
						"type": "text",
						"text": "Current weather in New York:\nTemperature: 72°F\nConditions: Partly cloudy"
					}],
					"isError": false
				}
			}
	*/

	mux.HandleFunc("POST /tools/call", func(response http.ResponseWriter, request *http.Request) {
		// Set headers for JSON and UTF-8
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		body := GetBytesBody(request)
		var data map[string]any

		if err := json.Unmarshal(body, &data); err != nil {
			// Handle error case
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]any{
				"result": map[string]any{
					"content": "Invalid JSON payload",
					"isError": true,
				},
			})
			return
		}

		fmt.Println("🚀", data["name"], "parameters:", data["arguments"])

		functionName := data["name"].(string)
		jsonArguments := data["arguments"]

		// convert jsonArguments to json string
		jsonArgumentsString, err := json.Marshal(jsonArguments)

		if err != nil {
			// Handle error case
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]any{
				"result": map[string]any{
					"content": "Error converting arguments to JSON",
					"isError": true,
				},
			})
			return
		}

		manifest := extism.Manifest{
			Wasm: []extism.Wasm{
				extism.WasmFile{
					Path: "functions/" + functionName + "/plugin.wasm"},
			},
			AllowedHosts: []string{"*"},
			Config:       map[string]string{},
		}

		pluginInst, err := extism.NewPlugin(ctx, manifest, config, nil) // new
		if err != nil {
			// Handle error case
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]any{
				"result": map[string]any{
					"content": "Error when loading the WASM plugin",
					"isError": true,
				},
			})
			return
		}

		_, output, err := pluginInst.Call("handle", jsonArgumentsString)

		if err != nil {
			// Handle error case
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]any{
				"result": map[string]any{
					"content": "Error executing the WASM fumction",
					"isError": true,
				},
			})
			return
		}

		// Set success status code (200 OK)
		response.WriteHeader(http.StatusOK)
		// Create and send JSON response
		json.NewEncoder(response).Encode(map[string]any{
			"result": map[string]any{
				"content": []map[string]string{
					{
						"type": "text",
						"text": string(output),
					},
				},
				"isError": false,
			},
		})

	})

	// Start the HTTP server
	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
