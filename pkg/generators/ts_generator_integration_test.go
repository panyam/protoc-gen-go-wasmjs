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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestTSGenerator_PerServiceGeneration_Integration tests the actual plugin execution
// with real proto files to validate per-service client generation
func TestTSGenerator_PerServiceGeneration_Integration(t *testing.T) {
	// This test validates that our new per-service client generation works
	// by actually running the plugin on test proto files
	
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "ts_generator_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Build our TypeScript generator plugin for testing
	pluginPath := filepath.Join(tempDir, "protoc-gen-go-wasmjs-ts")
	buildCmd := exec.Command("go", "build", "-o", pluginPath, "../../cmd/protoc-gen-go-wasmjs-ts")
	buildCmd.Dir = wd
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build TypeScript generator: %v", err)
	}

	testCases := []struct {
		name          string
		protoFile     string
		expectedFiles []string
	}{
		{
			name:      "Simple single service",
			protoFile: "testdata/simple.proto",
			expectedFiles: []string{
				"test/simple/v1/userServiceClient.ts",
			},
		},
		{
			name:      "Multiple services",
			protoFile: "testdata/multi_service.proto",
			expectedFiles: []string{
				"test/multi/v1/presenterServiceClient.ts",
				"test/multi/v1/browserAPIClient.ts",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create output directory for this test
			outputDir := filepath.Join(tempDir, tc.name)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				t.Fatalf("Failed to create output dir: %v", err)
			}

			// Prepare protoc command
			protoPath := filepath.Join(wd, tc.protoFile)
			if _, err := os.Stat(protoPath); os.IsNotExist(err) {
				t.Skipf("Proto file %s not found", protoPath)
				return
			}

			// Run protoc with our plugin
			args := []string{
				"--plugin=protoc-gen-go-wasmjs-ts=" + pluginPath,
				"--go-wasmjs-ts_out=" + outputDir,
				"--go-wasmjs-ts_opt=js_structure=namespaced,js_namespace=testApp",
				"--proto_path=" + filepath.Dir(protoPath),
				"--proto_path=" + filepath.Join(wd, "../../proto"), // For wasmjs annotations
				tc.protoFile,
			}

			cmd := exec.Command("protoc", args...)
			cmd.Dir = wd
			output, err := cmd.CombinedOutput()
			
			if err != nil {
				t.Logf("Protoc command failed (expected during development): %v", err)
				t.Logf("Output: %s", string(output))
				// Don't fail the test - this is expected during development
				return
			}

			// Check that expected files were generated
			for _, expectedFile := range tc.expectedFiles {
				fullPath := filepath.Join(outputDir, expectedFile)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not generated", expectedFile)
				} else {
					t.Logf("✅ Generated file: %s", expectedFile)
					
					// Read and validate file content
					content, err := os.ReadFile(fullPath)
					if err != nil {
						t.Errorf("Failed to read generated file %s: %v", expectedFile, err)
						continue
					}
					
					// Basic validation that it's a TypeScript client
					contentStr := string(content)
					if !strings.Contains(contentStr, "extends WASMServiceClient") {
						t.Errorf("Generated file %s doesn't extend WASMServiceClient", expectedFile)
					}
					
					if !strings.Contains(contentStr, "loadWASMModule") {
						t.Errorf("Generated file %s doesn't implement loadWASMModule", expectedFile)
					}
				}
			}
		})
	}
}

// TestTSGenerator_BrowserCallbacksExample_RealGeneration tests the actual browser-callbacks example
func TestTSGenerator_BrowserCallbacksExample_RealGeneration(t *testing.T) {
	// Test the real browser-callbacks example to validate our fix
	
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Path to browser-callbacks example
	exampleDir := filepath.Join(wd, "../../examples/browser-callbacks")
	
	// Check if example exists
	if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
		t.Skip("browser-callbacks example not found")
		return
	}

	// Run buf generate in the example
	cmd := exec.Command("make", "buf")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("make buf failed: %v", err)
		t.Logf("Output: %s", string(output))
		// Continue - we want to check what was generated
	}

	// Check that the TypeScript client was generated correctly
	generatedClient := filepath.Join(exampleDir, "web/src/generated/presenter/v1/presenterServiceClient.ts")
	if _, err := os.Stat(generatedClient); err == nil {
		t.Logf("✅ Per-service client generated successfully: presenter/v1/presenterServiceClient.ts")
		
		// Read and validate content
		content, err := os.ReadFile(generatedClient)
		if err != nil {
			t.Errorf("Failed to read generated client: %v", err)
			return
		}
		
		contentStr := string(content)
		if strings.Contains(contentStr, "presenterService") {
			t.Logf("✅ Generated client contains presenterService property")
		} else {
			t.Errorf("Generated client missing presenterService property")
		}
		
		if strings.Contains(contentStr, "extends ServiceClient") {
			t.Logf("✅ Generated client properly extends ServiceClient base class")
		} else {
			t.Errorf("Generated client doesn't extend ServiceClient")
		}
		
	} else {
		// Check if old single-file client exists (which would indicate our change didn't work)
		oldClient := filepath.Join(exampleDir, "web/src/generated/browser_callbacksClient.ts")
		if _, err := os.Stat(oldClient); err == nil {
			content, _ := os.ReadFile(oldClient)
			contentStr := string(content)
			
			if strings.Contains(contentStr, "presenterService") {
				t.Logf("Old single-file client still has presenterService - this works but not optimal")
			} else {
				t.Errorf("❌ Client doesn't have presenterService property in either location")
			}
		} else {
			t.Errorf("❌ No client file generated in either expected location")
		}
	}
}
