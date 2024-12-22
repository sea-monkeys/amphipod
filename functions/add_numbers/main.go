package main

import (
	"encoding/json"
	"strconv"

	"github.com/extism/go-pdk"
)

type Arguments struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

//export handle
func handle() {
	arguments := pdk.InputString()

	var args Arguments
	json.Unmarshal([]byte(arguments), &args)
	res := args.Number1 + args.Number2

	pdk.OutputString("🤖 result = " + strconv.Itoa(res))

}

func main() {}
