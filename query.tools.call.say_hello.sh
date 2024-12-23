#!/bin/bash 
SERVICE_URL="http://localhost:8080"
#SERVICE_URL="https://mcp.amphipod.local:8080"

read -r -d '' DATA <<- EOM
{
  "name":"say_hello",
  "arguments": {
    "name":"John Doe"
  }
}
EOM

echo "Sending: ${DATA} on ${SERVICE_URL}"
echo ""

curl --no-buffer ${SERVICE_URL}/tools/call \
    -H "Content-Type: application/json" \
    -d "${DATA}" 

#curl -H "Authorization: Bearer shrimpsarebeautiful" --no-buffer ${SERVICE_URL}/tools/call \
#    -H "Content-Type: application/json" \
#    -d "${DATA}" 

echo ""