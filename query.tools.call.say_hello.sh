#!/bin/bash 
SERVICE_URL="http://localhost:8080"
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
# --silent

curl --no-buffer ${SERVICE_URL}/tools/call \
    -H "Content-Type: application/json" \
    -d "${DATA}" 

echo ""