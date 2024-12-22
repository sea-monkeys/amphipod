#!/bin/bash 
#echo -n "🤖 tool: $1 params: $2"

name=$(echo $2 | jq -r '.name')

case "$1" in
    "say_hello")
        echo -n "👋🙂 hello $name"
        ;;
    "say_goodbye")
        echo -n "✋😢 goodbye $name"
        ;;
    *)
        echo "Unknown tool: $1"
        exit 1
        ;;
esac
