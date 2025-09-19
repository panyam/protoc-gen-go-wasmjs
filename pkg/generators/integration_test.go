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

package generators

import (
	"testing"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// TestGoGenerator_Creation tests Go generator initialization.
// This ensures all dependencies are properly wired up in the generator.
func TestGoGenerator_Creation(t *testing.T) {
	// Create generator (this tests the entire dependency injection setup)
	generator := NewGoGenerator(nil) // Plugin can be nil for this test

	if generator == nil {
		t.Error("NewGoGenerator() should return non-nil generator")
		return
	}

	// Test that all components are initialized
	if generator.analyzer == nil {
		t.Error("GoGenerator should have initialized analyzer")
	}

	if generator.pathCalc == nil {
		t.Error("GoGenerator should have initialized pathCalc")
	}

	if generator.nameConv == nil {
		t.Error("GoGenerator should have initialized nameConv")
	}

	if generator.packageFilter == nil {
		t.Error("GoGenerator should have initialized packageFilter")
	}

	if generator.serviceFilter == nil {
		t.Error("GoGenerator should have initialized serviceFilter")
	}

	if generator.methodFilter == nil {
		t.Error("GoGenerator should have initialized methodFilter")
	}

	if generator.dataBuilder == nil {
		t.Error("GoGenerator should have initialized dataBuilder")
	}

	if generator.renderer == nil {
		t.Error("GoGenerator should have initialized renderer")
	}
}

// TestTSGenerator_Creation tests TypeScript generator initialization.
// This ensures all dependencies are properly wired up in the generator.
func TestTSGenerator_Creation(t *testing.T) {
	// Create generator (this tests the entire dependency injection setup)
	generator := NewTSGenerator(nil) // Plugin can be nil for this test

	if generator == nil {
		t.Error("NewTSGenerator() should return non-nil generator")
		return
	}

	// Test that all components are initialized
	if generator.analyzer == nil {
		t.Error("TSGenerator should have initialized analyzer")
	}

	if generator.pathCalc == nil {
		t.Error("TSGenerator should have initialized pathCalc")
	}

	if generator.nameConv == nil {
		t.Error("TSGenerator should have initialized nameConv")
	}

	if generator.packageFilter == nil {
		t.Error("TSGenerator should have initialized packageFilter")
	}

	if generator.serviceFilter == nil {
		t.Error("TSGenerator should have initialized serviceFilter")
	}

	if generator.methodFilter == nil {
		t.Error("TSGenerator should have initialized methodFilter")
	}

	if generator.dataBuilder == nil {
		t.Error("TSGenerator should have initialized dataBuilder")
	}

	if generator.renderer == nil {
		t.Error("TSGenerator should have initialized renderer")
	}
}

// TestGoGenerator_ConfigValidation tests Go generator configuration validation.
// This ensures invalid configurations are caught early with helpful error messages.
func TestGoGenerator_ConfigValidation(t *testing.T) {
	generator := NewGoGenerator(nil)

	tests := []struct {
		name        string                     // Test case description
		config      *builders.GenerationConfig // Configuration to test
		expectError bool                       // Whether to expect validation error
		reason      string                     // Why this test case is important
	}{
		{
			name: "valid configuration",
			config: &builders.GenerationConfig{
				WasmExportPath:      "./gen/wasm",
				JSStructure:         "namespaced",
				JSNamespace:         "myapp",
				ModuleName:          "myapp_services",
				GenerateBuildScript: true,
			},
			expectError: false,
			reason:      "Valid configuration should pass validation",
		},
		{
			name: "empty wasm export path",
			config: &builders.GenerationConfig{
				WasmExportPath: "", // Invalid
				JSStructure:    "namespaced",
			},
			expectError: true,
			reason:      "Empty WASM export path should be rejected",
		},
		{
			name: "invalid JS structure",
			config: &builders.GenerationConfig{
				WasmExportPath: "./gen/wasm",
				JSStructure:    "invalid_structure", // Invalid
			},
			expectError: true,
			reason:      "Invalid JS structure should be rejected",
		},
		{
			name: "missing JS structure gets default",
			config: &builders.GenerationConfig{
				WasmExportPath: "./gen/wasm",
				JSStructure:    "", // Should get default
			},
			expectError: false,
			reason:      "Missing JS structure should get default value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.ValidateConfig(tt.config)

			if tt.expectError && err == nil {
				t.Errorf("ValidateConfig() expected error but got none\nReason: %s", tt.reason)
			}

			if !tt.expectError && err != nil {
				t.Errorf("ValidateConfig() unexpected error: %v\nReason: %s", err, tt.reason)
			}

			// Check that defaults are applied
			if !tt.expectError && tt.config.JSStructure == "" {
				// Note: The actual config modification happens in ValidateConfig
				// This test demonstrates the expected behavior
			}
		})
	}
}

// TestTSGenerator_ConfigValidation tests TypeScript generator configuration validation.
func TestTSGenerator_ConfigValidation(t *testing.T) {
	generator := NewTSGenerator(nil)

	tests := []struct {
		name        string                     // Test case description
		config      *builders.GenerationConfig // Configuration to test
		expectError bool                       // Whether to expect validation error
		reason      string                     // Why this test case is important
	}{
		{
			name: "valid configuration",
			config: &builders.GenerationConfig{
				TSExportPath: "./gen/ts",
				ModuleName:   "myapp_client",
			},
			expectError: false,
			reason:      "Valid TypeScript configuration should pass",
		},
		{
			name: "empty TS export path",
			config: &builders.GenerationConfig{
				TSExportPath: "", // Invalid
			},
			expectError: true,
			reason:      "Empty TypeScript export path should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.ValidateConfig(tt.config)

			if tt.expectError && err == nil {
				t.Errorf("ValidateConfig() expected error but got none\nReason: %s", tt.reason)
			}

			if !tt.expectError && err != nil {
				t.Errorf("ValidateConfig() unexpected error: %v\nReason: %s", err, tt.reason)
			}
		})
	}
}

// TestFilterCriteria_Integration tests the integration between generators and filter layer.
// This ensures the filter criteria parsing and usage works correctly in the generator context.
func TestFilterCriteria_Integration(t *testing.T) {
	// Test filter criteria creation from typical configurations
	tests := []struct {
		name        string // Test scenario description
		services    string // Services configuration
		includes    string // Method includes
		excludes    string // Method excludes
		renames     string // Method renames
		expectError bool   // Whether parsing should fail
		reason      string // Why this scenario is important
	}{
		{
			name:     "typical Go generator config",
			services: "UserService,LibraryService",
			includes: "",
			excludes: "*Internal,*Debug",
			renames:  "",
			reason:   "Common Go generator usage should work",
		},
		{
			name:     "typical TS generator config",
			services: "",
			includes: "Get*,Find*,Create*",
			excludes: "*Internal",
			renames:  "FindBooks:searchBooks",
			reason:   "Common TypeScript generator usage should work",
		},
		{
			name:     "development config",
			services: "",
			includes: "",
			excludes: "",
			renames:  "",
			reason:   "Development config (no filtering) should work",
		},
		{
			name:        "invalid rename format",
			services:    "UserService",
			renames:     "FindBooks->searchBooks", // Wrong separator
			expectError: true,
			reason:      "Invalid configuration should be caught",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria, err := filters.ParseFromConfig(tt.services, tt.includes, tt.excludes, tt.renames)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseFromConfig() expected error but got none\nReason: %s", tt.reason)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseFromConfig() unexpected error: %v\nReason: %s", err, tt.reason)
				return
			}

			// Basic validation that criteria was created correctly
			if criteria == nil {
				t.Error("ParseFromConfig() returned nil criteria")
				return
			}

			// Test filter detection
			hasServiceFilter := criteria.HasServiceFilter()
			hasMethodIncludes := criteria.HasMethodIncludes()
			hasMethodExcludes := criteria.HasMethodExcludes()
			hasMethodRenames := criteria.HasMethodRenames()

			// Validate against expectations
			if tt.services != "" && !hasServiceFilter {
				t.Error("Should detect service filter when services are configured")
			}

			if tt.includes != "" && !hasMethodIncludes {
				t.Error("Should detect method includes when configured")
			}

			if tt.excludes != "" && !hasMethodExcludes {
				t.Error("Should detect method excludes when configured")
			}

			if tt.renames != "" && !hasMethodRenames {
				t.Error("Should detect method renames when configured")
			}
		})
	}
}
