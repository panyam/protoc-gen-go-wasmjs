package main

import (
	"context"
	"fmt"
	"time"

	// Import the generated packages
	browserv1 "github.com/panyam/protoc-gen-go-wasmjs/example/gen/go/browser/v1"
	presenterv1 "github.com/panyam/protoc-gen-go-wasmjs/example/gen/go/presenter/v1"
	browserwasmgen "github.com/panyam/protoc-gen-go-wasmjs/example/gen/wasm/go/browser/v1"
	presenterwasmgen "github.com/panyam/protoc-gen-go-wasmjs/example/gen/wasm/go/presenter/v1"
	"google.golang.org/grpc"
)

// PresenterServiceImpl implements the PresenterService
type PresenterServiceImpl struct {
	presenterv1.UnimplementedPresenterServiceServer
	browserClient *browserwasmgen.BrowserAPIClient
}

// LoadUserData fetches user data from API and stores it locally
func (s *PresenterServiceImpl) LoadUserData(ctx context.Context, req *presenterv1.LoadUserDataRequest) (*presenterv1.LoadUserDataResponse, error) {
	fmt.Printf("LoadUserData called for user: %s\n", req.UserId)

	// TODO: Make actual API call via browser service
	// For now, return mock data
	return &presenterv1.LoadUserDataResponse{
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

// RunCallbackDemo demonstrates browser callbacks by prompting user 3 times
func (s *PresenterServiceImpl) RunCallbackDemo(ctx context.Context, req *presenterv1.CallbackDemoRequest) (*presenterv1.CallbackDemoResponse, error) {
	fmt.Printf("RunCallbackDemo called with demo: %s\n", req.DemoName)

	// Log start of demo
	_, err := s.browserClient.LogToWindow(ctx, &browserv1.LogRequest{
		Message: fmt.Sprintf("Starting callback demo: %s", req.DemoName),
		Level:   "info",
	})
	if err != nil {
		fmt.Printf("Error logging to window: %v\n", err)
		panic(err)
	}

	// Collect 3 inputs from the user
	var collectedInputs []string
	prompts := []string{
		"Enter your favorite color:",
		"Enter your favorite animal:",
		"Enter your favorite number:",
	}

	for i, prompt := range prompts {
		// Log the prompt
		_, err := s.browserClient.LogToWindow(ctx, &browserv1.LogRequest{
			Message: fmt.Sprintf("Prompting user %d/3: %s", i+1, prompt),
			Level:   "info",
		})
		if err != nil {
			fmt.Printf("Error logging prompt: %v\n", err)
		}

		// Prompt the user
		promptResp, err := s.browserClient.PromptUser(ctx, &browserv1.PromptRequest{
			Message:      prompt,
			DefaultValue: "",
		})
		if err != nil {
			// Log error
			s.browserClient.LogToWindow(ctx, &browserv1.LogRequest{
				Message: fmt.Sprintf("Error prompting user: %v", err),
				Level:   "error",
			})
			return nil, fmt.Errorf("failed to prompt user: %w", err)
		}

		if promptResp.Cancelled {
			// User cancelled
			s.browserClient.LogToWindow(ctx, &browserv1.LogRequest{
				Message: "User cancelled the demo",
				Level:   "warning",
			})
			return &presenterv1.CallbackDemoResponse{
				CollectedInputs: collectedInputs,
				Completed:       false,
			}, nil
		}

		// Log the response
		_, err = s.browserClient.LogToWindow(ctx, &browserv1.LogRequest{
			Message: fmt.Sprintf("User entered: %s", promptResp.Value),
			Level:   "success",
		})
		if err != nil {
			fmt.Printf("Error logging response: %v\n", err)
		}

		collectedInputs = append(collectedInputs, promptResp.Value)

		// Small delay between prompts
		time.Sleep(500 * time.Millisecond)
	}

	// Log completion
	_, err = s.browserClient.LogToWindow(ctx, &browserv1.LogRequest{
		Message: fmt.Sprintf("Demo completed! Collected %d inputs: %v", len(collectedInputs), collectedInputs),
		Level:   "success",
	})
	if err != nil {
		fmt.Printf("Error logging completion: %v\n", err)
	}

	return &presenterv1.CallbackDemoResponse{
		CollectedInputs: collectedInputs,
		Completed:       true,
	}, nil
}

func main() {
	// Create browser API client using the generated client
	browserClient := browserwasmgen.NewBrowserAPIClient()

	// Initialize browser exports (creates empty namespace for browser_v1)
	browserExports := &browserwasmgen.Browser_v1ServicesExports{
		BrowserAPI: browserClient,
	}
	browserExports.RegisterAPI()

	// Initialize presenter service implementations
	presenterExports := &presenterwasmgen.Presenter_v1ServicesExports{
		PresenterService: &PresenterServiceImpl{
			browserClient: browserClient,
		},
	}
	presenterExports.RegisterAPI()

	fmt.Println("Browser-callbacks example WASM module ready!")

	// Keep the WASM module running
	select {}
}
