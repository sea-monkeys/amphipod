#!/bin/bash 
SERVICE_URL="http://localhost:8080"
read -r -d '' DATA <<- EOM
{
  "name":"say_goodbye",
  "arguments": {
    "name":"Jane Doe"
  }
}
EOM

echo "Sending: ${DATA} on ${SERVICE_URL}"
echo ""


curl --no-buffer ${SERVICE_URL}/tools/call \
    -H "Content-Type: application/json" \
    -d "${DATA}" 

echo ""