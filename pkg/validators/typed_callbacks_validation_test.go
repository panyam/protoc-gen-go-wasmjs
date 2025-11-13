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

package validators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTypedCallbacks_GeneratedImports validates that the generated TypeScript clients
// properly import and use typed interfaces for method signatures
func TestTypedCallbacks_GeneratedImports(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Test with example (has various method types)
	exampleDir := filepath.Join(wd, "../../example")
	clientPath := filepath.Join(exampleDir, "web/src/generated/presenter/v1/presenterServiceClient.ts")

	content, err := os.ReadFile(clientPath)
	if err != nil {
		t.Skipf("Generated client not found: %v", err)
		return
	}

	contentStr := string(content)

	t.Run("TypeScriptImportsGenerated", func(t *testing.T) {
		// Should import TypeScript types
		if !strings.Contains(contentStr, "import {") {
			t.Error("Generated client should import TypeScript types")
		}

		if !strings.Contains(contentStr, "} from './interfaces';") {
			t.Error("Generated client should import from interfaces file")
		}

		// Should import specific types we expect
		expectedImports := []string{
			"LoadUserRequest",
			"LoadUserResponse",
			"CallbackDemoRequest",
			"CallbackDemoResponse",
			"StateUpdateRequest",
			"UIUpdate",
			"PreferencesRequest",
			"PreferencesResponse",
		}

		for _, expectedImport := range expectedImports {
			if !strings.Contains(contentStr, expectedImport+",") {
				t.Errorf("Expected import %s not found in client", expectedImport)
			}
		}

		t.Logf("✅ TypeScript type imports properly generated")
	})

	t.Run("TypedMethodSignatures", func(t *testing.T) {
		// Test that method signatures use proper types instead of 'any'

		// Sync method should be fully typed
		expectedSyncSignature := "loadUserData(request: LoadUserRequest): Promise<LoadUserResponse>"
		if !strings.Contains(contentStr, expectedSyncSignature) {
			t.Error("Sync method should have fully typed signature")
		}

		// Async method should have typed callback
		expectedAsyncSignature := "runCallbackDemo(request: CallbackDemoRequest, callback: (response: CallbackDemoResponse, error?: string) => void)"
		if !strings.Contains(contentStr, expectedAsyncSignature) {
			t.Error("Async method should have typed callback signature")
		}

		// Streaming method should have typed callback
		expectedStreamSignature := "updateUIState(request: StateUpdateRequest, callback: (response: UIUpdate | null, error: string | null, done: boolean) => boolean)"
		if !strings.Contains(contentStr, expectedStreamSignature) {
			t.Error("Streaming method should have typed callback signature")
		}

		// Should NOT contain 'any' types in method signatures
		if strings.Contains(contentStr, "request: any") {
			t.Error("Generated client should not use 'any' for request types")
		}

		if strings.Contains(contentStr, "Promise<any>") {
			t.Error("Generated client should not use 'any' for response types")
		}

		t.Logf("✅ Method signatures are fully typed")
	})

	t.Run("InterfaceMethodSignatures", func(t *testing.T) {
		// Test that service interface also uses typed signatures
		interfacePattern := "export interface PresenterServiceMethods"
		if !strings.Contains(contentStr, interfacePattern) {
			t.Error("Should generate typed service interface")
		}

		// Check that interface methods are typed
		interfaceStart := strings.Index(contentStr, interfacePattern)
		if interfaceStart == -1 {
			t.Fatal("Interface not found in generated content")
		}

		interfaceEnd := strings.Index(contentStr[interfaceStart:], "}")
		if interfaceEnd == -1 {
			t.Fatal("Interface end not found")
		}

		interfaceContent := contentStr[interfaceStart : interfaceStart+interfaceEnd]

		// Interface should have typed methods
		if strings.Contains(interfaceContent, "request: any") {
			t.Error("Service interface should not use 'any' for request types")
		}

		if strings.Contains(interfaceContent, "response: any") {
			t.Error("Service interface should not use 'any' for response types")
		}

		t.Logf("✅ Service interface is fully typed")
	})
}
