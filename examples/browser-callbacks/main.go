package main

import (
	"context"
	"fmt"

	// Import the generated WASM package
	presenterv1 "github.com/panyam/protoc-gen-go-wasmjs/examples/browser-callbacks/gen/go/presenter/v1"
	wasm "github.com/panyam/protoc-gen-go-wasmjs/examples/browser-callbacks/gen/wasm"
	"google.golang.org/grpc"
)

// PresenterServiceImpl implements the PresenterService
type PresenterServiceImpl struct {
	presenterv1.UnimplementedPresenterServiceServer
}

// LoadUserData fetches user data from API and stores it locally
func (s *PresenterServiceImpl) LoadUserData(ctx context.Context, req *presenterv1.LoadUserRequest) (*presenterv1.LoadUserResponse, error) {
	fmt.Printf("LoadUserData called for user: %s\n", req.UserId)

	// TODO: Make actual API call via browser service
	// For now, return mock data
	return &presenterv1.LoadUserResponse{
		Username:    "john_doe",
		Email:       "john@example.com",
		Permissions: []string{"read", "write"},
		FromCache:   false,
	}, nil
}

// UpdateUIState processes state changes and returns UI updates
func (s *PresenterServiceImpl) UpdateUIState(req *presenterv1.StateUpdateRequest, stream grpc.ServerStreamingServer[presenterv1.UIUpdate]) error {
	fmt.Printf("UpdateUIState called with action: %s\n", req.Action)

	// Use component from params if provided, otherwise default
	component := "default"
	if comp, ok := req.Params["component"]; ok {
		component = comp
	}

	// Send a series of UI updates
	updates := []*presenterv1.UIUpdate{
		{
			Component: component,
			Action:    "show",
			Data:      map[string]string{"message": "Loading..."},
		},
		{
			Component: component,
			Action:    "update",
			Data:      map[string]string{"message": "Processing..."},
		},
		{
			Component: component,
			Action:    "complete",
			Data:      map[string]string{"message": "Done!", "status": "success"},
		},
	}

	for _, update := range updates {
		if err := stream.Send(update); err != nil {
			return err
		}
	}

	return nil
}

// SavePreferences saves user preferences to localStorage
func (s *PresenterServiceImpl) SavePreferences(ctx context.Context, req *presenterv1.PreferencesRequest) (*presenterv1.PreferencesResponse, error) {
	fmt.Printf("SavePreferences called with %d preferences\n", len(req.Preferences))

	// TODO: Actually save to localStorage via browser service
	// For now, just return success
	return &presenterv1.PreferencesResponse{
		Saved:      true,
		ItemsSaved: int32(len(req.Preferences)),
	}, nil
}

func main() {
	// Initialize service implementations
	exports := &wasm.Browser_exampleServicesExports{
		PresenterService: &PresenterServiceImpl{},
	}

	// Register the JavaScript API
	exports.RegisterAPI()

	fmt.Println("Browser-callbacks example WASM module ready!")

	// Keep the WASM module running
	select {}
}
