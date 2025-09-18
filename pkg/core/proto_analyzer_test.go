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

package core

import (
	"testing"
)

// TestProtoAnalyzer_ExtractPackageName tests the extraction of package names from
// fully qualified message types. This is critical for cross-package imports and
// dependency resolution in TypeScript generation.
func TestProtoAnalyzer_ExtractPackageName(t *testing.T) {
	analyzer := NewProtoAnalyzer()

	tests := []struct {
		name            string // Test case description
		fullMessageType string // Input: fully qualified message type
		expectedPackage string // Expected output: package name
		reason          string // Why this test case is important
	}{
		{
			name:            "standard package with version",
			fullMessageType: "library.v1.Book",
			expectedPackage: "library.v1",
			reason:          "Most common case - package with version needs correct extraction",
		},
		{
			name:            "deeply nested package",
			fullMessageType: "company.product.module.v2.ComplexMessage",
			expectedPackage: "company.product.module.v2",
			reason:          "Deep package hierarchies should be handled correctly",
		},
		{
			name:            "single level package",
			fullMessageType: "common.Message",
			expectedPackage: "common",
			reason:          "Simple packages without versions should work",
		},
		{
			name:            "no package (just message)",
			fullMessageType: "Message",
			expectedPackage: "",
			reason:          "Messages without packages should return empty string",
		},
		{
			name:            "empty input",
			fullMessageType: "",
			expectedPackage: "",
			reason:          "Edge case - empty input should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.ExtractPackageName(tt.fullMessageType)
			if result != tt.expectedPackage {
				t.Errorf("ExtractPackageName(%s) = %s, want %s\nReason: %s",
					tt.fullMessageType, result, tt.expectedPackage, tt.reason)
			}
		})
	}
}

// TestProtoAnalyzer_ExtractMessageName tests the extraction of message names from
// fully qualified types. This is used for generating TypeScript class and interface names.
func TestProtoAnalyzer_ExtractMessageName(t *testing.T) {
	analyzer := NewProtoAnalyzer()

	tests := []struct {
		name            string // Test case description
		fullMessageType string // Input: fully qualified message type
		expectedMessage string // Expected output: message name
		reason          string // Why this test case is important
	}{
		{
			name:            "standard qualified message",
			fullMessageType: "library.v1.Book",
			expectedMessage: "Book",
			reason:          "Standard case - should extract just the message name",
		},
		{
			name:            "deeply nested message",
			fullMessageType: "company.product.module.v2.ComplexMessage",
			expectedMessage: "ComplexMessage",
			reason:          "Should extract message name regardless of package depth",
		},
		{
			name:            "unqualified message",
			fullMessageType: "SimpleMessage",
			expectedMessage: "SimpleMessage",
			reason:          "Messages without packages should return the whole string",
		},
		{
			name:            "empty input",
			fullMessageType: "",
			expectedMessage: "",
			reason:          "Edge case - empty input should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.ExtractMessageName(tt.fullMessageType)
			if result != tt.expectedMessage {
				t.Errorf("ExtractMessageName(%s) = %s, want %s\nReason: %s",
					tt.fullMessageType, result, tt.expectedMessage, tt.reason)
			}
		})
	}
}

// TestProtoAnalyzer_GetBaseFileName tests extraction of base filenames from proto paths.
// This is essential for generating consistent TypeScript file names that match proto files.
func TestProtoAnalyzer_GetBaseFileName(t *testing.T) {
	analyzer := NewProtoAnalyzer()

	tests := []struct {
		name         string // Test case description
		protoFile    string // Input: proto file path
		expectedBase string // Expected output: base filename
		reason       string // Why this test case is important
	}{
		{
			name:         "standard proto file path",
			protoFile:    "proto/library/v1/library.proto",
			expectedBase: "library",
			reason:       "Standard case - extract filename without .proto extension",
		},
		{
			name:         "deeply nested path",
			protoFile:    "src/proto/company/product/v2/messages.proto",
			expectedBase: "messages",
			reason:       "Should work regardless of path depth",
		},
		{
			name:         "filename only",
			protoFile:    "service.proto",
			expectedBase: "service",
			reason:       "Should handle simple filenames without directories",
		},
		{
			name:         "no extension",
			protoFile:    "proto/common/types",
			expectedBase: "types",
			reason:       "Should handle files without .proto extension",
		},
		{
			name:         "complex filename",
			protoFile:    "proto/user-auth/v1/user_auth_service.proto",
			expectedBase: "user_auth_service",
			reason:       "Should preserve complex filenames with underscores and hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.GetBaseFileName(tt.protoFile)
			if result != tt.expectedBase {
				t.Errorf("GetBaseFileName(%s) = %s, want %s\nReason: %s",
					tt.protoFile, result, tt.expectedBase, tt.reason)
			}
		})
	}
}

// Note: The following tests would require protogen.Field and protogen.Message mocks
// which would make this test file much more complex. In a real implementation,
// these would be integration tests that use actual protogen objects or sophisticated
// mocks. For now, I'm including the test structure with comments about what would
// be tested and why.

/*
// TestProtoAnalyzer_GetProtoFieldType would test field type extraction
// This is critical because incorrect type mapping leads to runtime errors
// in the generated TypeScript code.
//
// Test cases would include:
// - All primitive types (string, int32, bool, etc.)
// - Message types (should return message name)
// - Enum types (should return enum name)
// - Repeated and optional fields
//
// Why important: Type mapping errors cause TypeScript compilation errors
// or runtime failures when WASM and TypeScript try to exchange data.

func TestProtoAnalyzer_GetProtoFieldType(t *testing.T) {
	// This would require creating mock protogen.Field objects
	// with different field types to test the type extraction logic
}

// TestProtoAnalyzer_IsMapField would test map field detection
// This is important because map fields need special handling in TypeScript
// generation (they become Map<K, V> types instead of arrays).
//
// Test cases would include:
// - Regular message fields (should return false)
// - Actual map fields (should return true)
// - Repeated fields that aren't maps (should return false)
//
// Why important: Incorrectly identifying map fields leads to wrong
// TypeScript types and runtime errors during JSON serialization.

func TestProtoAnalyzer_IsMapField(t *testing.T) {
	// Would require mock protogen.Field objects with map entry messages
}

// TestProtoAnalyzer_GetMapKeyValueTypes would test map type extraction
// This ensures Map<string, number> is generated correctly for map<string, int32>.
//
// Test cases would include:
// - string->message maps
// - number->string maps
// - complex key/value combinations
// - non-map fields (should return "any", "any")
//
// Why important: Wrong map types cause TypeScript compilation errors
// and prevent proper type checking in client code.

func TestProtoAnalyzer_GetMapKeyValueTypes(t *testing.T) {
	// Would require mock map entry message structures
}

// TestProtoAnalyzer_IsBrowserProvidedService would test annotation detection
// This is critical for the browser service architecture where some services
// are implemented in JavaScript rather than Go WASM.
//
// Test cases would include:
// - Services with browser_provided=true annotation
// - Services without the annotation (should return false)
// - Services with malformed annotations
//
// Why important: Incorrect detection leads to wrong code generation -
// either missing browser service interfaces or incorrect WASM bindings.

func TestProtoAnalyzer_IsBrowserProvidedService(t *testing.T) {
	// Would require mock protogen.Service objects with wasmjs annotations
}

// Additional tests would cover:
// - GetCustomMethodName: Ensures custom JS method names are extracted correctly
// - IsAsyncMethod: Critical for preventing WASM deadlocks with browser APIs
// - IsMethodExcluded/IsServiceExcluded: Ensures filtered code isn't generated
// - GetCustomServiceName: For custom service naming in JS namespaces
// - GetOneofGroups: Essential for proper oneof handling in TypeScript
// - IsNestedMessage: Affects import generation and type references
*/
