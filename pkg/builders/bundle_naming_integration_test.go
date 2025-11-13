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

// TestBundleNamingIntegration reads the actual generated bundle file and checks if the bundle name is correct
// This test verifies that bundles use the configured module_name instead of package names
func TestBundleNamingIntegration(t *testing.T) {
	// Path to the generated bundle file in example
	generatedFilePath := filepath.Join("..", "..", "example",
		"web", "src", "generated", "index.ts")

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

	// Check for CORRECT behavior - simple base bundle with configured module_name
	if !strings.Contains(fileContent, "export class ExampleBundle extends WASMBundle") {
		t.Error("Expected to find 'ExampleBundle extends WASMBundle' (base bundle class), but didn't")
	}

	if !strings.Contains(fileContent, "moduleName: 'example'") {
		t.Error("Expected to find moduleName: 'example' (using configured module_name), but didn't")
	}

	// Verify it's the simple base bundle, not the old complex bundle
	if strings.Contains(fileContent, "public readonly presenterService") {
		t.Error("Found service properties in bundle - should be simple base bundle without services")
	}

	if strings.Contains(fileContent, "new PresenterServiceServiceClient") {
		t.Error("Found service instantiation in bundle - should be simple base bundle")
	}

	t.Log("Bundle naming and architecture updated: simple base bundle with configured module_name")
}
