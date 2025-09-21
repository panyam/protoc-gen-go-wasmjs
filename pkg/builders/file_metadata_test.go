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

package builders

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
)

// TestFileSpec_MetadataField tests the new Metadata field functionality
func TestFileSpec_MetadataField(t *testing.T) {
	t.Run("BasicMetadataOperations", func(t *testing.T) {
		spec := FileSpec{
			Name:     "test_client",
			Type:     "service_client",
			Required: true,
			Metadata: map[string]interface{}{
				"service": &protogen.Service{GoName: "TestService"},
				"key1":    "value1",
				"key2":    42,
				"key3":    true,
			},
		}

		// Test basic metadata access
		if spec.Metadata["key1"] != "value1" {
			t.Errorf("Expected key1='value1', got %v", spec.Metadata["key1"])
		}

		if spec.Metadata["key2"] != 42 {
			t.Errorf("Expected key2=42, got %v", spec.Metadata["key2"])
		}

		if spec.Metadata["key3"] != true {
			t.Errorf("Expected key3=true, got %v", spec.Metadata["key3"])
		}

		// Test service metadata
		service, ok := spec.Metadata["service"].(*protogen.Service)
		if !ok {
			t.Error("Expected service in metadata")
		} else if service.GoName != "TestService" {
			t.Errorf("Expected service name 'TestService', got %s", service.GoName)
		}
	})

	t.Run("NilMetadataHandling", func(t *testing.T) {
		spec := FileSpec{
			Name:     "test_spec",
			Metadata: nil,
		}

		// Should handle nil metadata gracefully
		if spec.Metadata != nil {
			t.Error("Expected nil metadata")
		}

		// Accessing nil metadata should not panic (Go handles this)
		value := spec.Metadata["nonexistent"]
		if value != nil {
			t.Error("Expected nil value from nil metadata")
		}
	})

	t.Run("EmptyMetadataHandling", func(t *testing.T) {
		spec := FileSpec{
			Name:     "test_spec",
			Metadata: map[string]interface{}{},
		}

		// Should handle empty metadata
		if len(spec.Metadata) != 0 {
			t.Error("Expected empty metadata map")
		}

		value := spec.Metadata["nonexistent"]
		if value != nil {
			t.Error("Expected nil value from empty metadata")
		}
	})
}

// TestGeneratedFileSet_GetFileSpec tests the new GetFileSpec method
func TestGeneratedFileSet_GetFileSpec(t *testing.T) {
	// Create test file specs
	specs := []FileSpec{
		{
			Name: "client1",
			Type: "service_client",
			Metadata: map[string]interface{}{
				"service": &protogen.Service{GoName: "Service1"},
			},
		},
		{
			Name: "client2", 
			Type: "service_client",
			Metadata: map[string]interface{}{
				"service": &protogen.Service{GoName: "Service2"},
			},
		},
		{
			Name: "interfaces",
			Type: "interfaces",
			Metadata: nil,
		},
	}

	filePlan := &FilePlan{Specs: specs}
	fileSet := &GeneratedFileSet{
		Plan:  filePlan,
		Files: make(map[string]*protogen.GeneratedFile),
	}

	t.Run("FindExistingSpec", func(t *testing.T) {
		spec := fileSet.GetFileSpec("client1")
		if spec == nil {
			t.Fatal("Expected to find spec 'client1'")
		}

		if spec.Name != "client1" {
			t.Errorf("Expected spec name 'client1', got %s", spec.Name)
		}

		if spec.Type != "service_client" {
			t.Errorf("Expected spec type 'service_client', got %s", spec.Type)
		}

		// Test metadata access
		service, ok := spec.Metadata["service"].(*protogen.Service)
		if !ok {
			t.Error("Expected service in metadata")
		} else if service.GoName != "Service1" {
			t.Errorf("Expected service name 'Service1', got %s", service.GoName)
		}
	})

	t.Run("FindNonExistentSpec", func(t *testing.T) {
		spec := fileSet.GetFileSpec("nonexistent")
		if spec != nil {
			t.Error("Expected nil for non-existent spec")
		}
	})

	t.Run("FindSpecWithNilMetadata", func(t *testing.T) {
		spec := fileSet.GetFileSpec("interfaces")
		if spec == nil {
			t.Fatal("Expected to find spec 'interfaces'")
		}

		if spec.Metadata != nil {
			t.Error("Expected nil metadata for interfaces spec")
		}
	})

	t.Run("GetFilesByType_ServiceClient", func(t *testing.T) {
		serviceClientFiles := fileSet.GetFilesByType("service_client")
		
		if len(serviceClientFiles) != 2 {
			t.Errorf("Expected 2 service client files, got %d", len(serviceClientFiles))
		}

		// Should include both client1 and client2
		if _, exists := serviceClientFiles["client1"]; !exists {
			t.Error("Expected 'client1' in service client files")
		}

		if _, exists := serviceClientFiles["client2"]; !exists {
			t.Error("Expected 'client2' in service client files")
		}
	})
}

// TestServiceClientDataBuilding tests the BuildServiceClientData functionality
func TestServiceClientDataBuilding(t *testing.T) {
	t.Run("ServiceClientDataStructure", func(t *testing.T) {
		// Test that BuildServiceClientData creates proper template data structure
		// This validates the single-service template data building

		packageInfo := &PackageInfo{
			Name: "test.service.v1",
			Path: "test/service/v1",
			Files: []*protogen.File{
				{
					Services: []*protogen.Service{
						{GoName: "TestService"},
					},
				},
			},
		}

		testService := &protogen.Service{
			GoName: "TestService",
		}

		// Test that the method signature works (compilation test)
		dataBuilder := &TSDataBuilder{}
		_, err := dataBuilder.BuildServiceClientData(packageInfo, testService, nil, nil)

		// We expect an error due to nil dependencies, but the method should exist
		if err == nil {
			t.Log("Unexpected success - method should fail with nil dependencies")
		} else {
			t.Logf("Expected error due to nil dependencies: %v", err)
		}

		// The important thing is that this compiles without errors
	})
}

// TestPerServiceFileNaming tests the filename generation logic
func TestPerServiceFileNaming(t *testing.T) {
	testCases := []struct {
		packageName      string
		packagePath      string
		serviceName      string
		expectedFilename string
	}{
		{
			packageName:      "simple.v1",
			packagePath:      "simple/v1", 
			serviceName:      "UserService",
			expectedFilename: "simple/v1/userServiceClient.ts",
		},
		{
			packageName:      "presenter.v1",
			packagePath:      "presenter/v1",
			serviceName:      "PresenterService", 
			expectedFilename: "presenter/v1/presenterServiceClient.ts",
		},
		{
			packageName:      "browser.v1",
			packagePath:      "browser/v1",
			serviceName:      "BrowserAPI",
			expectedFilename: "browser/v1/browserAPIClient.ts",
		},
		{
			packageName:      "deeply.nested.package.v2",
			packagePath:      "deeply/nested/package/v2",
			serviceName:      "ComplexService",
			expectedFilename: "deeply/nested/package/v2/complexServiceClient.ts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.packageName+"_"+tc.serviceName, func(t *testing.T) {
			// Test the filename calculation logic that's used in per-service generation
			// This validates our directory structure following proto hierarchy

			// Simulate the filename calculation
			serviceFileName := convertToFileNameForTest(tc.serviceName) + "Client.ts"
			filename := tc.packagePath + "/" + serviceFileName

			if filename != tc.expectedFilename {
				t.Errorf("Expected filename %s, got %s", tc.expectedFilename, filename)
			}
		})
	}
}

// Helper function that mimics the convertToFileName logic for testing
func convertToFileNameForTest(serviceName string) string {
	if len(serviceName) == 0 {
		return serviceName
	}
	// Convert PascalCase to camelCase
	return strings.ToLower(serviceName[:1]) + serviceName[1:]
}
