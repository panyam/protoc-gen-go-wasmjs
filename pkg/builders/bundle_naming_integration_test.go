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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBundleNamingIntegration reads the actual generated file and checks if the bundle name is correct
// This test confirms that the current behavior is broken before we fix it
func TestBundleNamingIntegration(t *testing.T) {
	// Path to the generated file in browser-callbacks example
	generatedFilePath := filepath.Join("..", "..", "examples", "browser-callbacks", 
		"web", "src", "generated", "presenter", "v1", "presenterServiceClient.ts")
	
	// Check if the file exists
	if _, err := os.Stat(generatedFilePath); os.IsNotExist(err) {
		t.Skip("Generated file not found - run buf generate first")
	}

	// Read the generated file
	content, err := os.ReadFile(generatedFilePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	fileContent := string(content)

	// Check current (broken) behavior - should show wrong bundle name
	if !strings.Contains(fileContent, "export class Presenter_v1Bundle") {
		t.Error("Expected to find 'Presenter_v1Bundle' (current broken behavior), but didn't")
	}

	if !strings.Contains(fileContent, "moduleName: 'presenter_v1'") {
		t.Error("Expected to find moduleName: 'presenter_v1' (current broken behavior), but didn't")
	}

	// These should NOT be present (they're what we want after the fix)
	if strings.Contains(fileContent, "Browser_callbacksBundle") {
		t.Error("Found 'Browser_callbacksBundle' which suggests the fix is already applied")
	}

	if strings.Contains(fileContent, "moduleName: 'browser_callbacks'") {
		t.Error("Found moduleName: 'browser_callbacks' which suggests the fix is already applied")
	}

	t.Log("✅ Test confirms current broken behavior: bundle named after package instead of module_name")
}
