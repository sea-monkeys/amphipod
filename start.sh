#!/bin/bash 
export HTTP_PORT=8080                      # Port to listen on
export USE_HTTPS=true                      # Enable HTTPS
export CERT_FILE=mcp.amphipod.local.crt    # Path to SSL certificate
export KEY_FILE=mcp.amphipod.local.key     # Path to SSL private key
export AUTH_TOKEN=shrimpsarebeautiful      # Authentication token
export REQUIRE_AUTH=true                   # Enable authentication
go run main.go