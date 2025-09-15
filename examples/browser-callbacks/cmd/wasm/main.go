//go:build js && wasm

package main

import (
	"fmt"

	"github.com/panyam/protoc-gen-go-wasmjs/examples/browser-callbacks/gen/wasm"
	"github.com/panyam/protoc-gen-go-wasmjs/examples/browser-callbacks/services"
)

func main() {
	fmt.Println("Browser Callbacks Example - Initializing WASM module...")

	// Create the browser API client (will be generated)
	browserAPI := wasm.NewBrowserAPIClient()

	// Create the presenter service with browser API access
	presenterService := services.NewPresenterService(browserAPI)

	// Create exports with dependency injection
	exports := &wasm.Browser_exampleServicesExports{
		// Regular services (implemented in WASM)
		PresenterService: presenterService,

		// Browser-provided services (clients)
		BrowserAPI: browserAPI,
	}

	// Register the JavaScript API
	exports.RegisterAPI()

	fmt.Println("Browser Callbacks Example - WASM module ready!")

	// Keep the program running
	select {}
}
