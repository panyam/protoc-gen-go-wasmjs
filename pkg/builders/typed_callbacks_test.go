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
	"testing"
)

// TestTypedCallbacks_ImportCollection tests that TypeScript type imports
// are properly collected from service method signatures
func TestTypedCallbacks_ImportCollection(t *testing.T) {
	// Create test service data with various method types
	services := []ServiceData{
		{
			Name:   "TestService",
			JSName: "testService",
			Methods: []MethodData{
				{
					Name:           "SyncMethod",
					JSName:         "syncMethod",
					RequestTSType:  "SyncRequest",
					ResponseTSType: "SyncResponse",
					IsAsync:        false,
					ShouldGenerate: true,
				},
				{
					Name:           "AsyncMethod",
					JSName:         "asyncMethod",
					RequestTSType:  "AsyncRequest",
					ResponseTSType: "AsyncResponse",
					IsAsync:        true,
					ShouldGenerate: true,
				},
				{
					Name:           "StreamingMethod",
					JSName:         "streamingMethod",
					RequestTSType:  "StreamRequest",
					ResponseTSType: "StreamResponse",
					IsServerStreaming: true,
					ShouldGenerate: true,
				},
				{
					Name:           "FilteredMethod",
					JSName:         "filteredMethod",
					RequestTSType:  "FilteredRequest",
					ResponseTSType: "FilteredResponse",
					ShouldGenerate: false, // Should be excluded
				},
			},
		},
	}

	dataBuilder := &TSDataBuilder{}
	imports := dataBuilder.collectServiceTypeImports(services)

	// Should collect only types from methods that ShouldGenerate=true
	expectedImports := map[string]bool{
		"SyncRequest":     true,
		"SyncResponse":    true,
		"AsyncRequest":    true,
		"AsyncResponse":   true,
		"StreamRequest":   true,
		"StreamResponse":  true,
		// FilteredRequest/FilteredResponse should be excluded
	}

	// Check that we got the right number of imports
	if len(imports) != len(expectedImports) {
		t.Errorf("Expected %d imports, got %d: %v", len(expectedImports), len(imports), imports)
	}

	// Check that all expected imports are present
	importSet := make(map[string]bool)
	for _, imp := range imports {
		importSet[imp] = true
	}

	for expectedImport := range expectedImports {
		if !importSet[expectedImport] {
			t.Errorf("Expected import %s not found in collected imports", expectedImport)
		}
	}

	// Check that filtered types are not included
	if importSet["FilteredRequest"] || importSet["FilteredResponse"] {
		t.Error("Filtered method types should not be included in imports")
	}
}

// TestTypedCallbacks_DuplicateImportHandling tests that duplicate type imports
// are properly deduplicated
func TestTypedCallbacks_DuplicateImportHandling(t *testing.T) {
	// Create service data with duplicate types
	services := []ServiceData{
		{
			Name:   "Service1",
			JSName: "service1",
			Methods: []MethodData{
				{
					Name:           "Method1",
					RequestTSType:  "CommonRequest",  // Duplicate
					ResponseTSType: "CommonResponse", // Duplicate
					ShouldGenerate: true,
				},
				{
					Name:           "Method2", 
					RequestTSType:  "CommonRequest",  // Duplicate
					ResponseTSType: "UniqueResponse",
					ShouldGenerate: true,
				},
			},
		},
		{
			Name:   "Service2",
			JSName: "service2", 
			Methods: []MethodData{
				{
					Name:           "Method3",
					RequestTSType:  "CommonRequest",  // Duplicate across services
					ResponseTSType: "CommonResponse", // Duplicate across services
					ShouldGenerate: true,
				},
			},
		},
	}

	dataBuilder := &TSDataBuilder{}
	imports := dataBuilder.collectServiceTypeImports(services)

	// Should deduplicate - only unique types
	expectedUniqueTypes := []string{
		"CommonRequest",
		"CommonResponse", 
		"UniqueResponse",
	}

	if len(imports) != len(expectedUniqueTypes) {
		t.Errorf("Expected %d unique imports, got %d: %v", len(expectedUniqueTypes), len(imports), imports)
	}

	// Verify each expected type appears exactly once
	typeCount := make(map[string]int)
	for _, imp := range imports {
		typeCount[imp]++
	}

	for _, expectedType := range expectedUniqueTypes {
		if typeCount[expectedType] != 1 {
			t.Errorf("Expected type %s to appear once, appeared %d times", expectedType, typeCount[expectedType])
		}
	}
}

// TestTypedCallbacks_EmptyTypes tests handling of empty or missing type information
func TestTypedCallbacks_EmptyTypes(t *testing.T) {
	services := []ServiceData{
		{
			Name:   "TestService",
			JSName: "testService",
			Methods: []MethodData{
				{
					Name:           "MethodWithEmptyTypes",
					RequestTSType:  "", // Empty
					ResponseTSType: "", // Empty
					ShouldGenerate: true,
				},
				{
					Name:           "MethodWithOneType",
					RequestTSType:  "ValidRequest",
					ResponseTSType: "", // Empty response
					ShouldGenerate: true,
				},
			},
		},
	}

	dataBuilder := &TSDataBuilder{}
	imports := dataBuilder.collectServiceTypeImports(services)

	// Should only include non-empty types
	expectedImports := []string{"ValidRequest"}

	if len(imports) != len(expectedImports) {
		t.Errorf("Expected %d imports, got %d: %v", len(expectedImports), len(imports), imports)
	}

	for _, expectedImport := range expectedImports {
		found := false
		for _, imp := range imports {
			if imp == expectedImport {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected import %s not found", expectedImport)
		}
	}
}

// TestTypedCallbacks_NoServices tests handling when no services are provided
func TestTypedCallbacks_NoServices(t *testing.T) {
	dataBuilder := &TSDataBuilder{}
	imports := dataBuilder.collectServiceTypeImports([]ServiceData{})

	if len(imports) != 0 {
		t.Errorf("Expected no imports for empty services, got: %v", imports)
	}
}
