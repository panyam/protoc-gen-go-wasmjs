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

package builders

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/filters"
)

// TestBrowserServicePanic reproduces the exact nil pointer panic scenario:
// service TestService {
//    option (wasmjs.v1.browser_provided) = true;
// }
// (Service with browser_provided option but no RPC methods)
func TestBrowserServicePanic(t *testing.T) {
	t.Run("BrowserServiceWithNoMethods", func(t *testing.T) {
		// Create a mock browser service with no methods (the exact scenario from the panic)
		mockService := &protogen.Service{
			GoName:  "TestService",
			Methods: []*protogen.Method{}, // No methods - this causes the panic
		}
		
		// Create minimal package info
		packageInfo := &PackageInfo{
			Name: "test.v1",
			Path: "test/v1",
		}
		
		// Create TSDataBuilder
		dataBuilder := createMockTSDataBuilder()
		
		// Create filter criteria
		criteria := &filters.FilterCriteria{}
		
		// Create generation config
		config := &GenerationConfig{
			JSStructure: "namespaced",
			JSNamespace: "testApp",
			ModuleName:  "test_module",
		}
		
		// This should reproduce the panic: BuildServiceClientData calls buildServiceDataForTS
		// which returns nil for services with no methods, then line 206 does *serviceData (panic!)
		defer func() {
			if r := recover(); r != nil {
				t.Logf("REPRODUCED PANIC: %v", r)
				t.Log("This confirms the nil pointer dereference bug when buildServiceDataForTS returns nil")
			}
		}()
		
		// Call the method that causes the panic
		result, err := dataBuilder.BuildServiceClientData(packageInfo, mockService, criteria, config)
		
		if err != nil {
			t.Logf("BuildServiceClientData failed with error: %v", err)
		} else if result == nil {
			t.Log("BuildServiceClientData returned nil (service had no methods)")
		} else {
			t.Log("BuildServiceClientData succeeded unexpectedly")
		}
	})
}

// TestBrowserServicePanicFixed tests the fix for the nil pointer panic
func TestBrowserServicePanicFixed(t *testing.T) {
	t.Run("BrowserServiceWithNoMethodsFixed", func(t *testing.T) {
		// Same setup as panic test
		mockService := &protogen.Service{
			GoName:  "TestService", 
			Methods: []*protogen.Method{},
		}
		
		packageInfo := &PackageInfo{
			Name: "test.v1",
			Path: "test/v1",
		}
		
		dataBuilder := createMockTSDataBuilder()
		criteria := &filters.FilterCriteria{}
		config := &GenerationConfig{
			JSStructure: "namespaced",
			JSNamespace: "testApp", 
			ModuleName:  "test_module",
		}
		
		// This should NOT panic after the fix
		result, err := dataBuilder.BuildServiceClientData(packageInfo, mockService, criteria, config)
		
		if err != nil {
			t.Logf("BuildServiceClientData failed gracefully: %v", err)
		} else if result == nil {
			t.Log("BuildServiceClientData returned nil gracefully (no methods to generate)")
		} else {
			t.Errorf("Unexpected success for service with no methods")
		}
		
		t.Log("SUCCESS: No panic occurred - fix is working")
	})
}

// createMockTSDataBuilder creates a minimal TSDataBuilder for testing
func createMockTSDataBuilder() *TSDataBuilder {
	// For testing, we need to provide the minimum required dependencies
	// In practice, these would be properly initialized
	return &TSDataBuilder{
		// These would normally be properly initialized, but for testing the panic,
		// we just need the BuildServiceClientData method to not crash
		methodFilter: &filters.MethodFilter{}, // Mock method filter
	}
}
