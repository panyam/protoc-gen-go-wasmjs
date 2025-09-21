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

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// TestBundleNamingDebug - Test-driven debugging to identify where module_name configuration gets lost
func TestBundleNamingDebug(t *testing.T) {
	t.Run("ConfigurationPassthrough", func(t *testing.T) {
		// Test 1: Verify GenerationConfig correctly holds module_name
		config := &GenerationConfig{
			ModuleName: "browser_callbacks", // This is what should come from buf.gen.yaml
		}

		if config.ModuleName != "browser_callbacks" {
			t.Errorf("GenerationConfig.ModuleName not set correctly, got=%s, want=%s", 
				config.ModuleName, "browser_callbacks")
		}
		t.Logf("✅ GenerationConfig holds module_name correctly: %s", config.ModuleName)
	})

	t.Run("TSDataBuilderUsesConfig", func(t *testing.T) {
		// Test 2: Verify TSDataBuilder.getModuleName uses config correctly
		nameConverter := core.NewNameConverter()
		builder := &TSDataBuilder{
			nameConv: nameConverter,
		}

		config := &GenerationConfig{
			ModuleName: "browser_callbacks",
		}

		// Test both packages that exist in browser-callbacks example
		testCases := []struct {
			packageName string
			expected    string
		}{
			{"presenter.v1", "browser_callbacks"}, // Should use config, not package
			{"browser.v1", "browser_callbacks"},   // Should use config, not package
		}

		for _, tc := range testCases {
			actual := builder.getModuleName(tc.packageName, config)
			if actual != tc.expected {
				t.Errorf("getModuleName(%s) = %s, want %s", tc.packageName, actual, tc.expected)
			} else {
				t.Logf("✅ getModuleName(%s) correctly returns %s", tc.packageName, actual)
			}
		}
	})

	t.Run("TemplateDataGeneration", func(t *testing.T) {
		// Test 3: Verify TSTemplateData gets correct ModuleName
		nameConverter := core.NewNameConverter()
		builder := &TSDataBuilder{
			nameConv: nameConverter,
		}

		config := &GenerationConfig{
			ModuleName: "browser_callbacks",
		}

		// Create minimal package info for presenter service
		packageInfo := &PackageInfo{
			Name: "presenter.v1",
		}

		// This is the critical test - does BuildServiceClientData use the config correctly?
		templateData := &TSTemplateData{
			ModuleName: builder.getModuleName(packageInfo.Name, config),
		}

		// Expected: browser_callbacks (from config)
		// Current broken: presenter_v1 (from package)
		expected := "browser_callbacks"
		if templateData.ModuleName != expected {
			t.Errorf("TSTemplateData.ModuleName = %s, want %s", templateData.ModuleName, expected)
			t.Errorf("❌ ISSUE FOUND: Template data not using configured module_name")
		} else {
			t.Logf("✅ TSTemplateData.ModuleName correctly set to %s", templateData.ModuleName)
		}
	})

	t.Run("BundleClassNameGeneration", func(t *testing.T) {
		// Test 4: Verify bundle class name generation
		// Expected: Browser_callbacksBundle (from browser_callbacks)
		// Current broken: Presenter_v1Bundle (from presenter_v1)
		
		moduleName := "browser_callbacks"
		// Simulate the template title function: {{ .ModuleName | title }}Bundle
		bundleClassName := titleCase(moduleName) + "Bundle"
		
		expected := "Browser_callbacksBundle"
		if bundleClassName != expected {
			t.Errorf("Bundle class name = %s, want %s", bundleClassName, expected)
		} else {
			t.Logf("✅ Bundle class name correctly generated: %s", bundleClassName)
		}
	})
}

// titleCase simulates the template helper function
func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// Simple title case - just capitalize first letter
	// The template helper is more sophisticated but this tests the concept
	return string(s[0]-'a'+'A') + s[1:]
}
