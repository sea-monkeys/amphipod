# Simple MCP Host

First Start the MCP Server, at the root of the project:

```bash
go run main.go
```

Then, in another terminal, you can run the following commands to test the MCP Server endpoints, in the current directory:

```bash
export OLLAMA_HOST=http://localhost:11434
export MCP_SERVER_HOST=http://localhost:8080
export TOOLS_LLM=qwen2.5:0.5b
go run main.go
```

