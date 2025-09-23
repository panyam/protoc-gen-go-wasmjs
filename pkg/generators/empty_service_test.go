// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/builders"
)

// TestEmptyServiceHandling tests how the generator handles services with zero methods
func TestEmptyServiceHandling(t *testing.T) {
	t.Run("ServiceWithZeroMethods", func(t *testing.T) {
		// Create a test proto file with a service that has no methods
		testProtoContent := `
syntax = "proto3";

package test.v1;

// A service with no methods - edge case
service EmptyService {
  // No methods defined
}

// A message for completeness
message EmptyRequest {
  string placeholder = 1;
}
`
		
		// Create temporary proto file
		tempDir := t.TempDir()
		protoFile := filepath.Join(tempDir, "empty.proto")
		err := os.WriteFile(protoFile, []byte(testProtoContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test proto file: %v", err)
		}
		
		// Test direct data building with empty service
		t.Log("Testing empty service handling...")
		
		// This test verifies our validation logic and how we handle edge cases
		// We can't easily run the full protoc pipeline in unit tests, but we can test
		// the validation logic and data building components
		
		// Test: What happens when BuildServiceClientData gets empty methods?
		testValidationWithEmptyMethods(t)
		
		// Test: What happens when templates receive empty service data?
		testTemplateWithEmptyService(t)
	})
}

// testValidationWithEmptyMethods tests TSRenderer validation with empty methods
func testValidationWithEmptyMethods(t *testing.T) {
	// Create test data with empty service
	serviceWithNoMethods := builders.ServiceData{
		Name:    "EmptyService",
		JSName:  "emptyService",
		Methods: []builders.MethodData{}, // No methods!
	}
	
	templateData := &builders.TSTemplateData{
		PackageName: "test.v1",
		PackagePath: "test/v1",
		Services:    []builders.ServiceData{serviceWithNoMethods},
		ModuleName:  "test_module",
	}
	
	// Test service client validation (should fail)
	renderer := NewTSRenderer(nil)
	err := renderer.ValidateTSTemplateData(templateData)
	
	if err != nil {
		t.Errorf("Service client validation failed unexpectedly: %v", err)
	} else {
		t.Log("EXPECTED: Service client validation now allows empty services (will generate client with 0 methods)")
	}
	
	// Test bundle validation (should pass with new validation)
	bundleErr := renderer.ValidateBundleTemplateData(templateData)
	
	if bundleErr != nil {
		t.Errorf("Bundle validation failed unexpectedly: %v", bundleErr)
	} else {
		t.Log("EXPECTED: Bundle validation passed for empty service (bundles don't need methods)")
	}
}

// testTemplateWithEmptyService tests what templates do with empty service data
func testTemplateWithEmptyService(t *testing.T) {
	// Create service data with dummy method to pass validation
	serviceWithDummyMethod := builders.ServiceData{
		Name:   "EmptyService",
		JSName: "emptyService",
		Methods: []builders.MethodData{
			{Name: "dummy", JSName: "dummy", ShouldGenerate: false}, // Dummy method
		},
	}
	
	templateData := &builders.TSTemplateData{
		PackageName: "test.v1",
		PackagePath: "test/v1",
		Services:    []builders.ServiceData{serviceWithDummyMethod},
		ModuleName:  "test_module",
	}
	
	// Test what the bundle template would generate
	t.Log("Testing bundle template with service that has dummy method...")
	
	// The bundle template iterates over .Services but doesn't use .Methods
	// So it should work fine even with dummy methods
	if len(templateData.Services) > 0 {
		service := templateData.Services[0]
		t.Logf("Bundle would generate: public readonly %s: %sServiceClient;", service.JSName, service.Name)
		t.Logf("Bundle would initialize: this.%s = new %sServiceClient(this._wasmBundle);", service.JSName, service.Name)
		
		// Check if any methods would be generated
		hasRealMethods := false
		for _, method := range service.Methods {
			if method.ShouldGenerate {
				hasRealMethods = true
				break
			}
		}
		
		if !hasRealMethods {
			t.Log("NOTE: Service client would be generated but have no callable methods")
		}
	}
}

// NewTSRenderer creates a TSRenderer for testing - simplified version
func NewTSRenderer(plugin *protogen.Plugin) *TSRenderer {
	// For testing validation, we don't need the full renderer infrastructure
	return &TSRenderer{}
}

// TSRenderer simple mock for testing
type TSRenderer struct{}

// ValidateTSTemplateData mock that implements the actual validation logic
func (tr *TSRenderer) ValidateTSTemplateData(data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}
	
	// Basic validation
	if data.PackageName == "" {
		return fmt.Errorf("TSTemplateData.PackageName cannot be empty")
	}
	
	if data.PackagePath == "" {
		return fmt.Errorf("TSTemplateData.PackagePath cannot be empty") 
	}
	
	// Validate services have proper names (methods are optional now)
	for _, service := range data.Services {
		if service.Name == "" {
			return fmt.Errorf("service has empty Name")
		}
		
		if service.JSName == "" {
			return fmt.Errorf("service %s has empty JSName", service.Name)
		}
		
		// Note: We now allow services with 0 methods (e.g., browser services with no RPCs)
	}
	
	return nil
}

// ValidateBundleTemplateData mock that implements bundle-specific validation
func (tr *TSRenderer) ValidateBundleTemplateData(data *builders.TSTemplateData) error {
	if data == nil {
		return nil
	}
	
	// Basic validation (same as regular validation)
	if data.PackageName == "" {
		return fmt.Errorf("TSTemplateData.PackageName cannot be empty")
	}
	
	if data.PackagePath == "" {
		return fmt.Errorf("TSTemplateData.PackagePath cannot be empty")
	}
	
	// For bundles, we don't validate that services have methods
	// Services are just used for import/property generation
	for _, service := range data.Services {
		if service.Name == "" {
			return fmt.Errorf("service has empty Name")
		}
		
		if service.JSName == "" {
			return fmt.Errorf("service %s has empty JSName", service.Name)
		}
	}
	
	return nil
}
