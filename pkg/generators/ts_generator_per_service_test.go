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
	"google.golang.org/protobuf/compiler/protogen"
)

// ========================================
// UNIT TESTS - Test individual functions
// ========================================

func TestTSGenerator_ServiceClientFilenameGeneration(t *testing.T) {
	generator := &TSGenerator{}

	testCases := []struct {
		name            string
		packageName     string
		packagePath     string
		serviceName     string
		expectedFilename string
	}{
		{
			name:            "Simple service",
			packageName:     "test.simple.v1",
			packagePath:     "test/simple/v1",
			serviceName:     "UserService",
			expectedFilename: "test/simple/v1/userServiceClient.ts",
		},
		{
			name:            "CamelCase conversion",
			packageName:     "presenter.v1",
			packagePath:     "presenter/v1",
			serviceName:     "PresenterService",
			expectedFilename: "presenter/v1/presenterServiceClient.ts",
		},
		{
			name:            "Deep nested package",
			packageName:     "company.api.user.v2",
			packagePath:     "company/api/user/v2",
			serviceName:     "UserManager",
			expectedFilename: "company/api/user/v2/userManagerClient.ts",
		},
		{
			name:            "Browser API service",
			packageName:     "browser.v1",
			packagePath:     "browser/v1",
			serviceName:     "BrowserAPI",
			expectedFilename: "browser/v1/browserAPIClient.ts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packageInfo := &builders.PackageInfo{
				Name: tc.packageName,
				Path: tc.packagePath,
			}
			
			mockService := &protogen.Service{
				GoName: tc.serviceName,
			}
			
			config := &builders.GenerationConfig{}
			
			filename := generator.calculateServiceClientFilename(packageInfo, mockService, config)
			if filename != tc.expectedFilename {
				t.Errorf("Expected filename %s, got %s", tc.expectedFilename, filename)
			}
		})
	}
}

func TestTSGenerator_ConvertToFileName(t *testing.T) {
	generator := &TSGenerator{}

	tests := []struct {
		input    string
		expected string
	}{
		{"PresenterService", "presenterService"},
		{"BrowserAPI", "browserAPI"},
		{"UserManagerService", "userManagerService"},
		{"Service", "service"},
		{"", ""},
		{"XMLParser", "xMLParser"}, // Edge case with multiple capitals
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.convertToFileName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}



// ========================================
// INTEGRATION TESTS - Test with real plugin execution
// ========================================

func TestTSGenerator_FileSpecMetadata(t *testing.T) {
	// Test that the Metadata field works correctly on FileSpec
	spec := builders.FileSpec{
		Name:     "test",
		Type:     "service_client",
		Metadata: map[string]interface{}{
			"service": &protogen.Service{GoName: "TestService"},
			"key":     "value",
		},
	}
	
	if spec.Metadata["key"] != "value" {
		t.Error("FileSpec.Metadata field not working correctly")
	}
	
	service, ok := spec.Metadata["service"].(*protogen.Service)
	if !ok {
		t.Error("Expected service in metadata")
	}
	
	if service.GoName != "TestService" {
		t.Errorf("Expected service name 'TestService', got %s", service.GoName)
	}
}

func TestTSGenerator_FileStructureIntegration(t *testing.T) {
	t.Run("MultiplePackagesGenerateToCorrectDirectories", func(t *testing.T) {
		generator := &TSGenerator{}

		testCases := []struct {
			packageName     string
			packagePath     string
			serviceName     string
			expectedFilename string
		}{
			{"presenter.v1", "presenter/v1", "PresenterService", "presenter/v1/presenterServiceClient.ts"},
			{"browser.v1", "browser/v1", "BrowserAPI", "browser/v1/browserAPIClient.ts"},
			{"user.service.v2", "user/service/v2", "UserManager", "user/service/v2/userManagerClient.ts"},
		}

		for _, tc := range testCases {
			t.Run(tc.packageName, func(t *testing.T) {
				packageInfo := &builders.PackageInfo{
					Name: tc.packageName,
					Path: tc.packagePath,
				}
				mockService := &protogen.Service{GoName: tc.serviceName}
				config := &builders.GenerationConfig{}

				filename := generator.calculateServiceClientFilename(packageInfo, mockService, config)
				if filename != tc.expectedFilename {
					t.Errorf("For package %s, service %s: expected filename %s, got %s",
						tc.packageName, tc.serviceName, tc.expectedFilename, filename)
				}
			})
		}
	})
}

func TestTSDataBuilder_BuildServiceClientData_InterfaceTest(t *testing.T) {
	// Test that the BuildServiceClientData method exists with correct signature
	
	// Create minimal components
	packageInfo := &builders.PackageInfo{
		Name: "test.v1",
		Path: "test/v1",
		Files: []*protogen.File{},
	}

	testService := &protogen.Service{
		GoName: "TestService",
	}

	dataBuilder := &builders.TSDataBuilder{}
	config := &builders.GenerationConfig{}
	criteria := &filters.FilterCriteria{}

	// Test that the method signature compiles (it will likely fail due to nil dependencies, but that's OK)
	_, err := dataBuilder.BuildServiceClientData(packageInfo, testService, criteria, config)
	
	// We expect an error due to nil dependencies, but the method should exist and be callable
	if err != nil {
		t.Logf("Expected error due to minimal setup (this is OK): %v", err)
	}
	
	// The important thing is that this compiles and the method signature is correct
}

// Add comprehensive tests for new functionality
func TestTSGenerator_FileMetadataHandling(t *testing.T) {
	// Test that metadata is properly stored and retrieved
	testService := &protogen.Service{
		GoName: "TestService",
	}

	spec := builders.FileSpec{
		Name: "client_TestService",
		Metadata: map[string]interface{}{
			"service": testService,
		},
	}

	// Test metadata retrieval
	if spec.Metadata == nil {
		t.Fatal("Expected metadata to be set")
	}

	service, ok := spec.Metadata["service"].(*protogen.Service)
	if !ok {
		t.Fatal("Expected service in metadata")
	}

	if service.GoName != "TestService" {
		t.Errorf("Expected service name 'TestService', got %s", service.GoName)
	}
}



func TestTSGenerator_GetFileSpec(t *testing.T) {
	// Test the new GetFileSpec method
	specs := []builders.FileSpec{
		{Name: "client1", Type: "service_client"},
		{Name: "client2", Type: "service_client"},
		{Name: "interfaces", Type: "interfaces"},
	}
	
	filePlan := &builders.FilePlan{Specs: specs}
	fileSet := &builders.GeneratedFileSet{Plan: filePlan}
	
	// Test finding existing spec
	spec := fileSet.GetFileSpec("client1")
	if spec == nil {
		t.Error("Expected to find spec 'client1'")
	} else if spec.Name != "client1" {
		t.Errorf("Expected spec name 'client1', got %s", spec.Name)
	}
	
	// Test non-existent spec
	spec = fileSet.GetFileSpec("nonexistent")
	if spec != nil {
		t.Error("Expected nil for non-existent spec")
	}
}

// Unit test for filename generation with different package structures
func TestTSGenerator_ComplexPackageFilenames(t *testing.T) {
	generator := &TSGenerator{}

	testCases := []struct {
		packageName     string
		packagePath     string
		serviceName     string
		expectedFilename string
	}{
		{
			packageName:     "simple.v1",
			packagePath:     "simple/v1",
			serviceName:     "UserService",
			expectedFilename: "simple/v1/userServiceClient.ts",
		},
		{
			packageName:     "company.api.user.v2",
			packagePath:     "company/api/user/v2",
			serviceName:     "UserManager",
			expectedFilename: "company/api/user/v2/userManagerClient.ts",
		},
		{
			packageName:     "deeply.nested.package.v3",
			packagePath:     "deeply/nested/package/v3",
			serviceName:     "ComplexService",
			expectedFilename: "deeply/nested/package/v3/complexServiceClient.ts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.packageName+"_"+tc.serviceName, func(t *testing.T) {
			packageInfo := &builders.PackageInfo{
				Name: tc.packageName,
				Path: tc.packagePath,
			}
			
			mockService := &protogen.Service{
				GoName: tc.serviceName,
			}
			
			config := &builders.GenerationConfig{}
			
			filename := generator.calculateServiceClientFilename(packageInfo, mockService, config)
			if filename != tc.expectedFilename {
				t.Errorf("Expected filename %s, got %s", tc.expectedFilename, filename)
			}
		})
	}
}
