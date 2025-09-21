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

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// TestGetModuleName validates that getModuleName uses the configured module_name
// instead of generating from package names
func TestGetModuleName(t *testing.T) {
	tests := []struct {
		name              string
		packageName       string
		configModuleName  string
		expectedModuleName string
	}{
		{
			name:              "uses configured module name",
			packageName:       "presenter.v1",
			configModuleName:  "browser_callbacks",
			expectedModuleName: "browser_callbacks",
		},
		{
			name:              "falls back to package name when no config",
			packageName:       "presenter.v1", 
			configModuleName:  "",
			expectedModuleName: "presenter_v1_services",
		},
		{
			name:              "handles underscores in config",
			packageName:       "library.v2",
			configModuleName:  "my_awesome_module",
			expectedModuleName: "my_awesome_module",
		},
		{
			name:              "second package should use same configured module",
			packageName:       "browser.v1",
			configModuleName:  "browser_callbacks", 
			expectedModuleName: "browser_callbacks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create minimal components needed
			nameConverter := core.NewNameConverter()
			
			// Create builder with minimal dependencies
			builder := &TSDataBuilder{
				nameConv: nameConverter,
			}

			// Create config
			config := &GenerationConfig{
				ModuleName: tt.configModuleName,
			}

			// Test the getModuleName function
			actualModuleName := builder.getModuleName(tt.packageName, config)
			
			if actualModuleName != tt.expectedModuleName {
				t.Errorf("Expected ModuleName=%s, got=%s", tt.expectedModuleName, actualModuleName)
			}
		})
	}
}
