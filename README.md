# AmphiPod

> Model Context Protocol HTTP server (Quick and dirty POC implementation of MCP on HTTP for experimentation)

Why: avoid implementation of a client for every language and platform

## Architecture Overview

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


## Workflow Overview

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
