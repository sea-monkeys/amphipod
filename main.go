package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

// Configuration struct to hold server settings
type Config struct {
	HTTPPort    string
	UseHTTPS    bool
	CertFile    string
	KeyFile     string
	AuthToken   string
	RequireAuth bool
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	return Config{
		HTTPPort:    getEnvOrDefault("HTTP_PORT", "8080"),
		UseHTTPS:    getEnvOrDefault("USE_HTTPS", "false") == "true",
		CertFile:    getEnvOrDefault("CERT_FILE", "cert.pem"),
		KeyFile:     getEnvOrDefault("KEY_FILE", "key.pem"),
		AuthToken:   getEnvOrDefault("AUTH_TOKEN", ""),
		RequireAuth: getEnvOrDefault("REQUIRE_AUTH", "false") == "true",
	}
}

// AuthMiddleware handles token-based authentication
func AuthMiddleware(next http.HandlerFunc, config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !config.RequireAuth || config.AuthToken == "" {
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token != config.AuthToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
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
	config := LoadConfig()

	ctx := context.Background()

	pluginConfig := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	//var httpPort = os.Getenv("HTTP_PORT")
	//if httpPort == "" {
	//	httpPort = "8080"
	//}

	// read the content of the mcp.list.json file and store it in a variable
	mcpTools, err := os.ReadFile("tools/mcp.list.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	mux := http.NewServeMux()

	// https://spec.modelcontextprotocol.io/specification/server/tools/#protocol-messages
	mux.HandleFunc("GET /tools/list", AuthMiddleware(func(response http.ResponseWriter, request *http.Request) {
		// Set headers for JSON and UTF-8
		response.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Set success status code (200 OK)
		response.WriteHeader(http.StatusOK)

		// Create and send JSON response
		response.Write(mcpTools)
	}, config))

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

	mux.HandleFunc("POST /tools/call", AuthMiddleware(func(response http.ResponseWriter, request *http.Request) {
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

		pluginInst, err := extism.NewPlugin(ctx, manifest, pluginConfig, nil) // new
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

	}, config))

	// Start the HTTP server
	//var errListening error
	//log.Println("🌍 http server is listening on: " + httpPort)
	//errListening = http.ListenAndServe(":"+httpPort, mux)

	//log.Fatal(errListening)

	// Start the server
	log.Printf("🌍 Server is listening on port %s (HTTPS: %v, Auth Required: %v)\n",
		config.HTTPPort, config.UseHTTPS, config.RequireAuth)

	var errListening error
	if config.UseHTTPS {
		errListening = http.ListenAndServeTLS(":"+config.HTTPPort, config.CertFile, config.KeyFile, mux)
	} else {
		errListening = http.ListenAndServe(":"+config.HTTPPort, mux)
	}

	log.Fatal(errListening)

}
