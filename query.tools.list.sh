#!/bin/bash 
SERVICE_URL="http://localhost:8080"
#SERVICE_URL="https://mcp.amphipod.local:8080"

curl --no-buffer ${SERVICE_URL}/tools/list 

#curl -H "Authorization: Bearer shrimpsarebeautiful" --no-buffer ${SERVICE_URL}/tools/list 

echo ""