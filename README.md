# AmphiPod 🦐

AmphiPod is a lightweight, HTTP-based **Model Context Protocol** (MCP) server implementation written in Go. It simplifies the integration of AI tools by providing an HTTP interface to the MCP specification and executing tools through WebAssembly plugins.

## 🌟 Key Features

- **HTTP-Based MCP Server**: Implements Anthropic's MCP specification using HTTP, making it accessible from any programming language
- **WebAssembly Tool Execution**: Leverages the Extism Framework for running tool plugins in WebAssembly
- **Language Agnostic**: Can be used with any programming language that supports HTTP requests
- **Lightweight & Fast**: Written in Go for optimal performance
- **Easy Integration**: Simple HTTP API for tool registration and execution

## 🎯 Why AmphiPod?

While Anthropic's MCP specification is excellent, it requires implementing MCP clients in each programming language. Currently, official SDKs are only available for Python, TypeScript, and Kotlin. AmphiPod solves this limitation by:

1. Using HTTP as the transport protocol, enabling any language with HTTP capabilities to interact with the MCP server
2. Simplifying the integration process for host AI applications
3. Providing a consistent interface across different programming languages
4. Executing tools through WebAssembly plugins for enhanced security and portability

## 🚀 Getting Started

> 🚧 This is a work in progress. Please check back later for installation instructions.
<!--
### Prerequisites

### Installation

### Basic Usage
-->

## 🔧 API Endpoints

- `GET /tools/list`: Get available tools
- `POST /tools/call`: Execute a tool

## 📦 Tool Development

AmphiPod uses **WebAssembly** plugins powered by the [**Extism** Framework](https://extism.org/). 

Look at the `./functions` directory for examples of WebAssembly plugins. 

## 🔄 Architecture Overview

```mermaid
flowchart TD
    subgraph Client["Host GenAI Application"]
        GenAI["Host GenAI application with MCP Client"]
    end

    subgraph Servers["MCP Servers"]
        ServerA["HTTP MCP Server A"]
        ServerB["HTTP MCP Server B"]
        ServerC["HTTP MCP Server C"]
    end

    subgraph LocalResources["Local Resources"]
        ResourceA["Local Data Source and Resource A"]
        ResourceB["Local Data Source and Resource B"]
    end

    subgraph RemoteServices["Remote Services"]
        ServiceC["Remote Services on Internet C"]
    end

    %% Client to Server connections
    GenAI <--> |HTTP| ServerA
    GenAI <--> |HTTP| ServerB
    GenAI <--> |HTTP| ServerC

    %% Server to Resource connections
    ServerA --> ResourceA
    ServerB --> ResourceB
    
    %% Server to Remote Service connection
    ServerC --> ServiceC

    %% Styling
    classDef default fill:#f9f9f9,stroke:#333,stroke-width:2px;
    classDef client fill:#e1f3d8,stroke:#333,stroke-width:2px;
    classDef server fill:#dae8fc,stroke:#333,stroke-width:2px;
    classDef resource fill:#ffe6cc,stroke:#333,stroke-width:2px;
    classDef remote fill:#fff2cc,stroke:#333,stroke-width:2px;

    class GenAI client;
    class ServerA,ServerB,ServerC server;
    class ResourceA,ResourceB resource;
    class ServiceC remote;
```

## 🔄 Workflow Overview

The MCP Client will be a simple HTTP client run by the Host GenAI application. The MCP Client will make HTTP requests to the MCP Server to get the list of tools and to make tool calls. The MCP Server will respond with the list of tools and the output of the tool calls.


```mermaid
sequenceDiagram
    participant LLM as Tools LLM
    participant Host as Host GenAI App
    participant Client as MCP Client
    participant Server as MCP Server

    Host->>Client: Initialize HTTP Client
    Client->>Server: GET /tools/list
    Server-->>Client: Return tools JSON
    Client-->>Host: Return tools list
    
    Host->>Host: Convert JSON to tools list
    Host->>Host: Generate LLM prompt
    Host->>LLM: Send prompt
    LLM-->>Host: Generate tool response
    
    Host->>Client: Send tool call request
    Client->>Server: POST /tools/call
    Server-->>Client: Return tool output
    Client-->>Host: Return tool response
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit pull requests, report bugs, and suggest features.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Anthropic](https://www.anthropic.com) for the MCP specification
- The [Extism](https://extism.org/) team for their excellent WebAssembly framework

## 📚 Further Reading

- [Introducing the Model Context Protocol](https://www.anthropic.com/news/model-context-protocol)
- [Model Context Protocol Specification](https://modelcontextprotocol.io/introduction)
- [Extism Framework Documentation](https://extism.org/)
- [WebAssembly Documentation](https://webassembly.org/)
- [WASI Documentation](https://wasi.dev/)

