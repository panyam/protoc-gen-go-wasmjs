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

// TestNewMethodFilter tests method filter creation.
// This ensures the method filter is properly initialized with its dependencies.
func TestNewMethodFilter(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	filter := NewMethodFilter(analyzer)

	if filter == nil {
		t.Error("NewMethodFilter() should return non-nil filter")
	}

	if filter.analyzer != analyzer {
		t.Error("MethodFilter should store the provided analyzer")
	}
}

// TestMethodFilter_ValidateMethodPatterns tests glob pattern validation.
// This is critical because invalid patterns would cause runtime panics during filtering.
// Testing this pure function ensures we catch configuration errors early.
func TestMethodFilter_ValidateMethodPatterns(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	filter := NewMethodFilter(analyzer)

	tests := []struct {
		name        string   // Test case description
		includes    []string // Method include patterns
		excludes    []string // Method exclude patterns
		expectError bool     // Whether to expect validation error
		reason      string   // Why this test case is important
	}{
		{
			name:        "valid glob patterns",
			includes:    []string{"Get*", "Find*", "Create*"},
			excludes:    []string{"*Internal", "*Debug", "*Test"},
			expectError: false,
			reason:      "Standard glob patterns should validate successfully",
		},
		{
			name:        "empty patterns",
			includes:    []string{},
			excludes:    []string{},
			expectError: false,
			reason:      "Empty pattern lists should be valid (no filtering)",
		},
		{
			name:        "single character wildcards",
			includes:    []string{"Get?", "Find?User"},
			excludes:    []string{"?Internal"},
			expectError: false,
			reason:      "Single character wildcards (?) should be valid",
		},
		{
			name:        "literal method names",
			includes:    []string{"GetUser", "FindBooks", "CreateLibrary"},
			excludes:    []string{"DeleteUser", "RemoveBook"},
			expectError: false,
			reason:      "Literal method names (no wildcards) should be valid",
		},
		{
			name:        "invalid glob pattern in includes",
			includes:    []string{"Get*", "[invalid"}, // Unclosed bracket
			excludes:    []string{"*Internal"},
			expectError: true,
			reason:      "Invalid glob syntax should be caught during validation",
		},
		{
			name:        "invalid glob pattern in excludes",
			includes:    []string{"Get*"},
			excludes:    []string{"*Internal", "[unclosed"}, // Unclosed bracket
			expectError: true,
			reason:      "Invalid exclude patterns should be caught during validation",
		},
		{
			name:        "complex valid patterns",
			includes:    []string{"[Gg]et*", "*[Uu]ser*", "Find[A-Z]*"},
			excludes:    []string{"*[Ii]nternal*", "*[Dd]ebug*"},
			expectError: false,
			reason:      "Complex but valid glob patterns should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria()
			criteria.MethodIncludes = tt.includes
			criteria.MethodExcludes = tt.excludes

			err := filter.ValidateMethodPatterns(criteria)

			if tt.expectError && err == nil {
				t.Errorf("ValidateMethodPatterns() expected error but got none\nReason: %s", tt.reason)
			}

			if !tt.expectError && err != nil {
				t.Errorf("ValidateMethodPatterns() unexpected error: %v\nReason: %s", err, tt.reason)
			}
		})
	}
}

// TestMethodFilter_GetMethodJSName tests JavaScript method name generation.
// This is a pure function that can be tested without protogen mocks by testing
// the name conversion logic directly.
func TestMethodFilter_GetMethodJSName(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	filter := NewMethodFilter(analyzer)
	nameConverter := core.NewNameConverter()

	tests := []struct {
		name         string            // Test case description
		methodName   string            // Original method name
		renames      map[string]string // Method renames configuration
		expectedName string            // Expected JavaScript name
		reason       string            // Why this test case is important
	}{
		{
			name:         "method with configured rename",
			methodName:   "FindBooks",
			renames:      map[string]string{"FindBooks": "searchBooks"},
			expectedName: "searchBooks",
			reason:       "Configured renames should override default camelCase conversion",
		},
		{
			name:         "method without rename",
			methodName:   "GetUser",
			renames:      map[string]string{"FindBooks": "searchBooks"},
			expectedName: "getUser",
			reason:       "Methods without renames should use camelCase conversion",
		},
		{
			name:         "PascalCase to camelCase",
			methodName:   "CreateLibraryItem",
			renames:      map[string]string{},
			expectedName: "createLibraryItem",
			reason:       "Default behavior should convert PascalCase to camelCase",
		},
		{
			name:         "already camelCase",
			methodName:   "login",
			renames:      map[string]string{},
			expectedName: "login",
			reason:       "Already camelCase names should remain unchanged",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria()
			criteria.MethodRenames = tt.renames

			// Note: This test simulates the GetMethodJSName behavior without annotations
			// since we can't easily mock protogen.Method objects. The annotation logic
			// is tested through integration tests with real proto files.

			// Test the rename logic directly
			result := criteria.GetMethodRename(tt.methodName)
			if result != tt.methodName && tt.renames[tt.methodName] != "" {
				// If there was a rename, use it
				if result != tt.expectedName {
					t.Errorf("Method rename: GetMethodRename(%s) = %s, want %s\nReason: %s",
						tt.methodName, result, tt.expectedName, tt.reason)
				}
			} else {
				// No rename, should use camelCase conversion
				camelCase := nameConverter.ToCamelCase(tt.methodName)
				if camelCase != tt.expectedName {
					t.Errorf("CamelCase conversion: ToCamelCase(%s) = %s, want %s\nReason: %s",
						tt.methodName, camelCase, tt.expectedName, tt.reason)
				}
			}
		})
	}

	// Suppress unused variable warning for filter
	_ = filter
}

/*
Note: The following tests would require protogen.Method and protogen.Service mocks
which are complex to create. In a real implementation, these would be integration
tests that use actual protogen objects. Here's what would be tested and why:

// TestMethodFilter_ShouldIncludeMethod would test the main method filtering logic
// This is the core business logic that determines method inclusion/exclusion.
//
// Test cases would include:
// - Methods with wasm_method_exclude annotation (should be excluded)
// - Client streaming methods (should be excluded - not supported)
// - Methods matching exclude patterns (should be excluded)
// - Methods matching include patterns when includes are configured (should be included)
// - Methods not matching include patterns when includes are configured (should be excluded)
// - Methods when no patterns are configured (should be included by default)
// - Async methods (should be included with async=true metadata)
// - Server streaming methods (should be included with streaming=true metadata)
// - Methods with custom names (should include custom name in result)
//
// Why important: This determines which methods appear in generated JavaScript APIs.
// Bugs here cause methods to disappear or appear incorrectly, breaking client code.

func TestMethodFilter_ShouldIncludeMethod(t *testing.T) {
	// Would require creating mock protogen.Method objects with:
	// - Method descriptions with names
	// - Method options with wasmjs annotations
	// - Streaming settings (client/server/unary)
	// - Comments and metadata
}

// TestMethodFilter_FilterMethods would test batch method filtering
// This ensures statistics collection and batch processing work correctly.
//
// Test cases would include:
// - Service with all methods included (should return all)
// - Service with all methods excluded (should return empty)
// - Service with mixed include/exclude results
// - Statistics accuracy for included/excluded counts
// - Empty service (no methods)
//
// Why important: Batch filtering is used by the main generator and needs to
// handle edge cases correctly while maintaining accurate statistics.

func TestMethodFilter_FilterMethods(t *testing.T) {
	// Would require mock protogen.Service with protogen.Method arrays
}

// TestMethodFilter_HasAnyMethods would test early termination logic
// This is used to determine if a service is worth generating at all.
//
// Test cases would include:
// - Service with at least one included method (should return true)
// - Service with all methods excluded (should return false)
// - Empty service (should return false)
// - Service with only client streaming methods (should return false)
//
// Why important: Prevents generation of empty service wrappers when all
// methods are filtered out, reducing generated code size and complexity.

func TestMethodFilter_HasAnyMethods(t *testing.T) {
	// Would require mock services with various method configurations
}

Integration Testing Approach:
The method filtering logic is comprehensively tested through the examples:

examples/library/proto/library/v1/library.proto includes:
- FindBooks with custom name annotation: option (wasmjs.v1.wasm_method_name) = "searchBooks"
- CreateUser with exclusion: option (wasmjs.v1.wasm_method_exclude) = true
- Regular methods without annotations
- Multiple services (LibraryService, UserService)

This provides realistic testing of:
- Annotation-based filtering and renaming
- Service filtering and method filtering interactions
- Complex filtering scenarios with real protobuf definitions
*/
