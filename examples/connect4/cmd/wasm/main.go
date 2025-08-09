//go:build js && wasm
// +build js,wasm

// Main.go for multiplayer Connect4 WASM module
// This creates the actual WASM module that registers service implementations

package main

import (
	"syscall/js"

	// Import the generated WASM package
	multiplayer_connect4 "github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/gen/wasm"

	// Import our actual service implementation
	"github.com/panyam/protoc-gen-go-wasmjs/examples/connect4/services"
)

func main() {
	// Initialize service implementations with our actual Connect4Service
	service := services.NewConnect4Service()
	exports := &multiplayer_connect4.Multiplayer_connect4ServicesExports{
		Connect4Service: service,
	}

	// Register the JavaScript API
	exports.RegisterAPI()

	// Expose storage callback configuration functions to browser
	setupStorageCallbacks(service)

	// Keep the WASM module running
	select {}
}

// setupStorageCallbacks exposes functions for browser to configure storage operations
func setupStorageCallbacks(service *services.Connect4Service) {
	// Browser calls this to set up storage callbacks
	js.Global().Set("setWasmStorageCallbacks", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 3 {
			return map[string]any{
				"success": false,
				"message": "Expected 3 callback functions: save, load, poll",
			}
		}

		saveFunc := args[0]
		loadFunc := args[1]
		pollFunc := args[2]

		// Set the callbacks in the service
		service.SetStorageCallbacks(saveFunc, loadFunc, pollFunc)

		return map[string]any{
			"success": true,
			"message": "Storage callbacks configured successfully",
		}
	}))

	// Browser calls this when external storage changes are detected
	js.Global().Set("wasmOnExternalStorageChange", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 2 {
			return map[string]any{
				"success": false,
				"message": "Expected gameId and gameState JSON",
			}
		}

		gameId := args[0].String()
		gameStateJson := args[1].String()

		// Notify service of external change
		service.HandleExternalStorageChange(gameId, gameStateJson)

		return map[string]any{
			"success": true,
			"message": "External storage change processed",
		}
	}))
}
