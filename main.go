package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

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

		//Here call the function that will return the result
		cmd := exec.Command("bash", "-c", "./tools/runner.sh "+functionName+" '"+string(jsonArgumentsString)+ "'")
		output, err := cmd.Output()


		if err != nil {
			// Handle error case
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]any{
				"result": map[string]any{
					"content": "Error executing the command",
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
