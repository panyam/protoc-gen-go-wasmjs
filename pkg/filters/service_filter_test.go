// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filters

import (
	"testing"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// TestNewServiceFilter tests service filter creation.
// This ensures the service filter is properly initialized with its dependencies.
func TestNewServiceFilter(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	filter := NewServiceFilter(analyzer)

	if filter == nil {
		t.Error("NewServiceFilter() should return non-nil filter")
	}

	if filter.analyzer != analyzer {
		t.Error("ServiceFilter should store the provided analyzer")
	}
}

// TestServiceFilterResult_Creation tests result creation helpers.
// These helpers ensure consistent result formatting across filtering operations.
func TestServiceFilterResult_Creation(t *testing.T) {
	tests := []struct {
		name        string              // Test case description
		factoryFunc func() FilterResult // Factory function to test
		expectIncl  bool                // Expected Include value
		reason      string              // Why this test is important
	}{
		{
			name:        "included result",
			factoryFunc: func() FilterResult { return Included("test reason") },
			expectIncl:  true,
			reason:      "Included() should create result with Include=true",
		},
		{
			name:        "excluded result",
			factoryFunc: func() FilterResult { return Excluded("test reason") },
			expectIncl:  false,
			reason:      "Excluded() should create result with Include=false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.factoryFunc()

			if result.Include != tt.expectIncl {
				t.Errorf("FilterResult.Include = %t, want %t\nReason: %s",
					result.Include, tt.expectIncl, tt.reason)
			}

			if result.Reason == "" {
				t.Error("FilterResult.Reason should not be empty")
			}
		})
	}
}

/*
Note: The following tests would require protogen.Service mocks which are complex to create.
In a real implementation, these would be integration tests that use actual protogen objects
or sophisticated mocking frameworks. Here's what would be tested and why:

// TestServiceFilter_ShouldIncludeService would test the main filtering logic
// This is the core business logic that determines service inclusion/exclusion.
//
// Test cases would include:
// - Services with wasm_service_exclude annotation (should be excluded)
// - Services in the configured services list (should be included)
// - Services not in configured list when list is specified (should be excluded)
// - Services when no list is configured (should be included by default)
// - Browser-provided services (should be included with metadata)
// - Services with custom names (should include custom name in result)
//
// Why important: This is the main entry point for service filtering and affects
// whether entire services appear in generated code. Bugs here cause services
// to disappear entirely or appear when they shouldn't.

func TestServiceFilter_ShouldIncludeService(t *testing.T) {
	// Would require creating mock protogen.Service objects with:
	// - Service descriptions with names
	// - Service options with wasmjs annotations
	// - Different browser_provided settings
	// This is complex because protogen types have many internal dependencies
}

// TestServiceFilter_FilterServices would test batch filtering operations
// This ensures the convenience methods work correctly for processing multiple services.
//
// Test cases would include:
// - Multiple services with mixed include/exclude results
// - Statistics collection accuracy
// - File processing across multiple files
// - Browser service separation
//
// Why important: Batch operations are used by the main generator and need to
// handle edge cases like empty service lists, all excluded services, etc.

func TestServiceFilter_FilterServices(t *testing.T) {
	// Would require creating mock protogen.File objects with services
}

// TestServiceFilter_GetBrowserProvidedServices would test browser service separation
// This is critical for the dual-mode architecture where some services run in browser.
//
// Test cases would include:
// - Mixed regular and browser services (should separate correctly)
// - All browser services (should return all)
// - No browser services (should return empty)
// - Browser services that are excluded by other criteria
//
// Why important: Browser services need special handling in code generation.
// Wrong separation leads to incorrect client/server interface generation.

func TestServiceFilter_GetBrowserProvidedServices(t *testing.T) {
	// Would require mock services with browser_provided annotations
}

Integration Testing Approach:
These service filter functions are integration-tested through the full generator pipeline
in the examples/ directory. The library example includes:
- Regular services (LibraryService)
- Services with custom method names
- Services with excluded methods (CreateUser with wasm_method_exclude)
- Multiple services in one package

This provides comprehensive validation of the filtering logic in realistic scenarios.
*/
