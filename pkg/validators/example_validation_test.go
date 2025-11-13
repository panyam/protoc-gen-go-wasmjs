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
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestExamples_Basic validates that the example
// works correctly with our framework (without triggering build constraint issues)
func TestExamples_Basic(t *testing.T) {
	// Get path to example (avoiding build constraint issues by not running go test there)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	exampleDir := filepath.Join(wd, "../../example")
	if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
		t.Skip("example not found")
		return
	}

	t.Run("PerServiceFilesGenerated", func(t *testing.T) {
		// Validate that per-service generation worked
		expectedFiles := []string{
			"web/src/generated/presenter/v1/presenterServiceClient.ts",
			"web/src/generated/browser/v1/browserAPIClient.ts",
		}

		for _, file := range expectedFiles {
			fullPath := filepath.Join(exampleDir, file)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("Per-service file not generated: %s", file)
			} else {
				t.Logf("✅ Per-service file exists: %s", file)
			}
		}
	})

	t.Run("TypeScriptCompilation", func(t *testing.T) {
		// Test TypeScript compilation works with per-service clients
		cmd := exec.Command("pnpm", "typecheck")
		cmd.Dir = filepath.Join(exampleDir, "web")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Example TypeScript compilation failed: %v\nOutput: %s", err, string(output))
		} else {
			t.Logf("✅ Example TypeScript compiles correctly")
		}
	})

	t.Run("WASMBuilds", func(t *testing.T) {
		// Test WASM builds successfully
		cmd := exec.Command("make", "wasm")
		cmd.Dir = exampleDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Example WASM build failed: %v\nOutput: %s", err, string(output))
		} else {
			// Check WASM file was created
			wasmPath := filepath.Join(exampleDir, "web/public/browser_example.wasm")
			if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
				t.Error("WASM file not created")
			} else {
				t.Logf("✅ Example WASM builds successfully")
			}
		}
	})

	t.Run("MainJSUsesPerServiceClients", func(t *testing.T) {
		// Validate that main.js is properly updated to use per-service clients
		mainPath := filepath.Join(exampleDir, "web/src/main.ts")
		content, err := os.ReadFile(mainPath)
		if err != nil {
			t.Skipf("main.ts not found: %v", err)
			return
		}

		contentStr := string(content)

		// Should import per-service client bundle
		if !strings.Contains(contentStr, "from './generated/presenter/v1/presenterServiceClient'") {
			t.Error("main.ts should import per-service client bundle")
		}

		// Should use correct bundle class (new architecture)
		if !strings.Contains(contentStr, "new ExampleBundle()") {
			t.Error("main.ts should use ExampleBundle")
		}

		// Should access service via composition pattern
		if !strings.Contains(contentStr, "presenterService.") && !strings.Contains(contentStr, "PresenterService") {
			t.Error("main.ts should access presenterService via composition pattern")
		}

		t.Log("Example main.js properly uses per-service clients with composition pattern")
	})

	t.Run("AsyncMethodProperlyConfigured", func(t *testing.T) {
		// Validate async method configuration
		protoPath := filepath.Join(exampleDir, "proto/presenter/v1/presenter.proto")
		content, err := os.ReadFile(protoPath)
		if err != nil {
			t.Skipf("Proto file not found: %v", err)
			return
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "option (wasmjs.v1.async_method)") {
			t.Error("Example should use async_method annotation")
		}

		// Check generated TypeScript has typed callback signature
		clientPath := filepath.Join(exampleDir, "web/src/generated/presenter/v1/presenterServiceClient.ts")
		if clientContent, err := os.ReadFile(clientPath); err == nil {
			clientStr := string(clientContent)

			// Check for typed callback signature (updated for Phase 2)
			if !strings.Contains(clientStr, "callback: (response: CallbackDemoResponse, error?: string) => void") {
				t.Error("Generated client should have typed callback signature for async method")
			} else {
				t.Logf("✅ Example async method properly configured with typed callbacks")
			}

			// Also check that it imports the response type
			if !strings.Contains(clientStr, "CallbackDemoResponse,") {
				t.Error("Generated client should import CallbackDemoResponse type")
			}
		}
	})
}

// TestExamples_DirectoryStructures validates that all examples follow
// the correct per-service directory structure
func TestExamples_DirectoryStructures(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	examplesDir := filepath.Join(wd, "../../examples")
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		t.Skip("examples directory not found")
		return
	}

	// Find all examples
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		exampleName := entry.Name()
		if exampleName == "." || exampleName == ".." {
			continue
		}

		t.Run(exampleName, func(t *testing.T) {
			examplePath := filepath.Join(examplesDir, exampleName)
			generatedPath := filepath.Join(examplePath, "web/src/generated")

			// Skip if no web/src/generated directory (not all examples have web UIs)
			if _, err := os.Stat(generatedPath); os.IsNotExist(err) {
				t.Skipf("Example %s has no web UI", exampleName)
				return
			}

			// Check that generated directory follows proto structure
			// (not specific files, just that the structure makes sense)
			entries, err := os.ReadDir(generatedPath)
			if err != nil {
				t.Errorf("Failed to read generated directory for %s: %v", exampleName, err)
				return
			}

			// Should have at least one subdirectory (proto package structure)
			hasSubdirs := false
			for _, entry := range entries {
				if entry.IsDir() {
					hasSubdirs = true
					t.Logf("✅ Example %s has proto package directory: %s", exampleName, entry.Name())
				}
			}

			if !hasSubdirs {
				t.Errorf("Example %s should have proto package subdirectories", exampleName)
			}
		})
	}
}
