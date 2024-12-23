#!/bin/bash 
SERVICE_URL="http://localhost:8080"
#SERVICE_URL="https://mcp.amphipod.local:8080"

read -r -d '' DATA <<- EOM
{
  "name":"add_numbers",
  "arguments": {
    "number1":28,
    "number2":14
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