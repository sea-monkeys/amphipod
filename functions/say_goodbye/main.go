package main

import (
	"encoding/json"
	"github.com/extism/go-pdk"
)

type Arguments struct {
	Name string `json:"name"`
}

//export handle
func handle() {
	arguments := pdk.InputString()
	var args Arguments
	json.Unmarshal([]byte(arguments), &args)

	pdk.OutputString("🟣👋😢 Goodbye " + args.Name)
}

func main() {}
