//go:build js && wasm

package main

import (
	"fmt"

	v1 "github.com/panyam/protoc-gen-go-wasmjs/example/gen/wasm/go/browser/v1"
	pv1 "github.com/panyam/protoc-gen-go-wasmjs/example/gen/wasm/go/presenter/v1"
	"github.com/panyam/protoc-gen-go-wasmjs/example/services"
)

func main() {
	fmt.Println("Example - Initializing WASM module...")

	// Create the browser API client (will be generated)
	browserAPI := v1.NewBrowserAPIClient()

	// Create the presenter service with browser API access
	presenterService := services.NewPresenterService(browserAPI)

	// Create exports with dependency injection
	browserExports := &v1.Browser_v1ServicesExports{
		// Browser-provided services (clients)
		BrowserAPI: browserAPI,
	}

	// Regular services (implemented in WASM)
	presenterExports := &pv1.Presenter_v1ServiceExports{
		PresenterService: presenterService,
	}

	// Register the JavaScript API
	presenterExports.RegisterAPI()
	browserExports.RegisterAPI()

	fmt.Println("Example - WASM module ready!")

	// Keep the program running
	select {}
}
