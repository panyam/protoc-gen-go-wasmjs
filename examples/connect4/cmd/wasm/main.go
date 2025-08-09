//go:build js && wasm
// +build js,wasm

// Main.go for multiplayer Connect4 WASM module
// This creates the actual WASM module that registers service implementations

package main

import (
	// Import the generated WASM package
	"github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/gen/wasm"

	// Import our actual service implementation
	"github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/services"
)

func main() {
	// Initialize service implementations with our actual Connect4Service
	exports := &wasm.Multiplayer_connect4ServicesExports{
		Connect4Service: services.NewConnect4Service(),
	}

	// Register the JavaScript API
	exports.RegisterAPI()

	// Keep the WASM module running
	select {}
}
