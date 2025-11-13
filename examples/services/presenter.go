package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	browserv1 "github.com/panyam/protoc-gen-go-wasmjs/example/gen/go/browser/v1"
	presenterv1 "github.com/panyam/protoc-gen-go-wasmjs/example/gen/go/presenter/v1"
)

// PresenterService implements the presenter logic that runs in WASM
type PresenterService struct {
	presenterv1.UnimplementedPresenterServiceServer
	browserAPI BrowserAPIInterface // Injected browser API client interface
}

// BrowserAPIInterface defines the methods we need from the browser API
type BrowserAPIInterface interface {
	GetLocalStorage(ctx context.Context, req *browserv1.StorageKeyRequest) (*browserv1.StorageValueResponse, error)
	SetLocalStorage(ctx context.Context, req *browserv1.StorageSetRequest) (*browserv1.StorageSetResponse, error)
	Fetch(ctx context.Context, req *browserv1.FetchRequest) (*browserv1.FetchResponse, error)
	Alert(ctx context.Context, req *browserv1.AlertRequest) (*browserv1.AlertResponse, error)
}

// NewPresenterService creates a new presenter service with browser API access
func NewPresenterService(browserAPI BrowserAPIInterface) *PresenterService {
	return &PresenterService{
		browserAPI: browserAPI,
	}
}

// LoadUserData fetches user data from API and stores it locally
func (s *PresenterService) LoadUserData(ctx context.Context, req *presenterv1.LoadUserRequest) (*presenterv1.LoadUserResponse, error) {
	// First check localStorage for cached data
	cacheKey := fmt.Sprintf("user_data_%s", req.UserId)
	storageResp, err := s.browserAPI.GetLocalStorage(ctx, &browserv1.StorageKeyRequest{
		Key: cacheKey,
	})
	if err == nil && storageResp.Exists {
		// Parse cached data
		var cachedUser presenterv1.LoadUserResponse
		if err := json.Unmarshal([]byte(storageResp.Value), &cachedUser); err == nil {
			cachedUser.FromCache = true
			return &cachedUser, nil
		}
	}

	// Fetch fresh data from API
	apiURL := fmt.Sprintf("https://api.example.com/users/%s", req.UserId)
	fetchResp, err := s.browserAPI.Fetch(ctx, &browserv1.FetchRequest{
		Url:    apiURL,
		Method: "GET",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %w", err)
	}

	if fetchResp.Status != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", fetchResp.Status, fetchResp.StatusText)
	}

	// Parse API response
	var apiData struct {
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
	}
	if err := json.Unmarshal([]byte(fetchResp.Body), &apiData); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Create response
	response := &presenterv1.LoadUserResponse{
		Username:    apiData.Username,
		Email:       apiData.Email,
		Permissions: apiData.Permissions,
		FromCache:   false,
	}

	// Cache the data in localStorage
	cacheData, _ := json.Marshal(response)
	_, _ = s.browserAPI.SetLocalStorage(ctx, &browserv1.StorageSetRequest{
		Key:   cacheKey,
		Value: string(cacheData),
	})

	return response, nil
}

// UpdateUIState processes state changes and returns UI updates
func (s *PresenterService) UpdateUIState(req *presenterv1.StateUpdateRequest, stream presenterv1.PresenterService_UpdateUIStateServer) error {
	// Simulate processing state updates and sending UI commands
	ctx := stream.Context()

	// Send initial loading state
	if err := stream.Send(&presenterv1.UIUpdate{
		Component: "loading",
		Action:    "show",
		Data: map[string]string{
			"message": "Processing " + req.Action,
		},
	}); err != nil {
		return err
	}

	// Simulate some processing
	time.Sleep(500 * time.Millisecond)

	// Based on action, send different UI updates
	switch req.Action {
	case "refresh":
		// Send multiple UI updates
		updates := []presenterv1.UIUpdate{
			{
				Component: "header",
				Action:    "update",
				Data: map[string]string{
					"timestamp": time.Now().Format(time.RFC3339),
				},
			},
			{
				Component: "content",
				Action:    "refresh",
				Data:      req.Params,
			},
			{
				Component: "loading",
				Action:    "hide",
				Data:      map[string]string{},
			},
		}

		for _, update := range updates {
			if err := stream.Send(&update); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond) // Simulate gradual updates
		}

	case "navigate":
		// Show alert about navigation
		_, _ = s.browserAPI.Alert(ctx, &browserv1.AlertRequest{
			Message: fmt.Sprintf("Navigating to: %s", req.Params["page"]),
		})

		// Send navigation UI update
		if err := stream.Send(&presenterv1.UIUpdate{
			Component: "router",
			Action:    "navigate",
			Data:      req.Params,
		}); err != nil {
			return err
		}

		// Hide loading
		if err := stream.Send(&presenterv1.UIUpdate{
			Component: "loading",
			Action:    "hide",
			Data:      map[string]string{},
		}); err != nil {
			return err
		}

	default:
		// Unknown action
		if err := stream.Send(&presenterv1.UIUpdate{
			Component: "error",
			Action:    "show",
			Data: map[string]string{
				"message": "Unknown action: " + req.Action,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

// SavePreferences saves user preferences to localStorage
func (s *PresenterService) SavePreferences(ctx context.Context, req *presenterv1.PreferencesRequest) (*presenterv1.PreferencesResponse, error) {
	savedCount := 0

	for key, value := range req.Preferences {
		prefKey := fmt.Sprintf("pref_%s", key)
		resp, err := s.browserAPI.SetLocalStorage(ctx, &browserv1.StorageSetRequest{
			Key:   prefKey,
			Value: value,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save preference %s: %w", key, err)
		}
		if resp.Success {
			savedCount++
		}
	}

	// Show success message
	if savedCount > 0 {
		_, _ = s.browserAPI.Alert(ctx, &browserv1.AlertRequest{
			Message: fmt.Sprintf("Successfully saved %d preferences", savedCount),
		})
	}

	return &presenterv1.PreferencesResponse{
		Saved:      savedCount > 0,
		ItemsSaved: int32(savedCount),
	}, nil
}
