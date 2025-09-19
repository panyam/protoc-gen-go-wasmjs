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

package filters

import (
	"testing"
)

// TestNewFilterCriteria tests the creation of filter criteria with default values.
// This ensures sensible defaults that maintain backward compatibility with existing generator behavior.
func TestNewFilterCriteria(t *testing.T) {
	criteria := NewFilterCriteria()

	// Test default values match current generator behavior
	if criteria.ServicesSet == nil {
		t.Error("ServicesSet should be initialized, not nil")
	}

	if len(criteria.ServicesSet) != 0 {
		t.Error("ServicesSet should be empty by default (include all services)")
	}

	if len(criteria.MethodIncludes) != 0 {
		t.Error("MethodIncludes should be empty by default")
	}

	if len(criteria.MethodExcludes) != 0 {
		t.Error("MethodExcludes should be empty by default")
	}

	if criteria.MethodRenames == nil {
		t.Error("MethodRenames should be initialized, not nil")
	}

	// Test defaults that ensure proper filtering behavior
	if !criteria.ExcludeAnnotationPackages {
		t.Error("ExcludeAnnotationPackages should be true by default (skip wasmjs.v1)")
	}

	if !criteria.ExcludeEmptyPackages {
		t.Error("ExcludeEmptyPackages should be true by default (skip packages with no content)")
	}

	if !criteria.ExcludeMapEntries {
		t.Error("ExcludeMapEntries should be true by default (skip synthetic map messages)")
	}

	if criteria.ExcludeNestedMessages {
		t.Error("ExcludeNestedMessages should be false by default (include nested messages)")
	}

	if criteria.ExcludeNestedEnums {
		t.Error("ExcludeNestedEnums should be false by default (include nested enums)")
	}
}

// TestFilterCriteria_HasServiceFilter tests detection of service filtering configuration.
// This is important for determining whether to apply service-specific filtering logic.
func TestFilterCriteria_HasServiceFilter(t *testing.T) {
	tests := []struct {
		name        string          // Test case description
		servicesSet map[string]bool // Services configuration
		expected    bool            // Expected result
		reason      string          // Why this test case is important
	}{
		{
			name:        "no services configured",
			servicesSet: map[string]bool{},
			expected:    false,
			reason:      "Empty services set means include all services (no filtering)",
		},
		{
			name:        "specific services configured",
			servicesSet: map[string]bool{"UserService": true, "LibraryService": true},
			expected:    true,
			reason:      "Non-empty services set means apply service filtering",
		},
		{
			name:        "nil services set",
			servicesSet: nil,
			expected:    false,
			reason:      "Nil services set should be treated as no filtering",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria()
			criteria.ServicesSet = tt.servicesSet

			result := criteria.HasServiceFilter()
			if result != tt.expected {
				t.Errorf("HasServiceFilter() = %t, want %t\nReason: %s",
					result, tt.expected, tt.reason)
			}
		})
	}
}

// TestFilterCriteria_GetMethodRename tests method rename functionality.
// This ensures custom method names are applied correctly during generation.
func TestFilterCriteria_GetMethodRename(t *testing.T) {
	tests := []struct {
		name         string            // Test case description
		renames      map[string]string // Method renames configuration
		methodName   string            // Method name to look up
		expectedName string            // Expected result
		reason       string            // Why this test case is important
	}{
		{
			name:         "method with custom rename",
			renames:      map[string]string{"FindBooks": "searchBooks", "GetUser": "fetchUser"},
			methodName:   "FindBooks",
			expectedName: "searchBooks",
			reason:       "Configured renames should be applied",
		},
		{
			name:         "method without rename",
			renames:      map[string]string{"FindBooks": "searchBooks"},
			methodName:   "GetUser",
			expectedName: "GetUser",
			reason:       "Methods without renames should return original name",
		},
		{
			name:         "empty renames map",
			renames:      map[string]string{},
			methodName:   "FindBooks",
			expectedName: "FindBooks",
			reason:       "Empty renames should return original names",
		},
		{
			name:         "nil renames map",
			renames:      nil,
			methodName:   "FindBooks",
			expectedName: "FindBooks",
			reason:       "Nil renames should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria := NewFilterCriteria()
			criteria.MethodRenames = tt.renames

			result := criteria.GetMethodRename(tt.methodName)
			if result != tt.expectedName {
				t.Errorf("GetMethodRename(%s) = %s, want %s\nReason: %s",
					tt.methodName, result, tt.expectedName, tt.reason)
			}
		})
	}
}

// TestParseFromConfig tests parsing filter criteria from string configuration.
// This ensures the bridge between existing config system and new filter layer works correctly.
func TestParseFromConfig(t *testing.T) {
	tests := []struct {
		name             string            // Test case description
		servicesStr      string            // Services configuration string
		methodIncludeStr string            // Method includes configuration
		methodExcludeStr string            // Method excludes configuration
		methodRenameStr  string            // Method renames configuration
		expectedServices []string          // Expected services in set
		expectedIncludes []string          // Expected include patterns
		expectedExcludes []string          // Expected exclude patterns
		expectedRenames  map[string]string // Expected renames
		expectError      bool              // Whether to expect an error
		reason           string            // Why this test case is important
	}{
		{
			name:             "complete configuration",
			servicesStr:      "UserService,LibraryService",
			methodIncludeStr: "Get*,Find*",
			methodExcludeStr: "*Internal,*Debug",
			methodRenameStr:  "FindBooks:searchBooks,GetUser:fetchUser",
			expectedServices: []string{"UserService", "LibraryService"},
			expectedIncludes: []string{"Get*", "Find*"},
			expectedExcludes: []string{"*Internal", "*Debug"},
			expectedRenames:  map[string]string{"FindBooks": "searchBooks", "GetUser": "fetchUser"},
			expectError:      false,
			reason:           "Complete configuration should parse all components correctly",
		},
		{
			name:             "empty configuration",
			servicesStr:      "",
			methodIncludeStr: "",
			methodExcludeStr: "",
			methodRenameStr:  "",
			expectedServices: []string{},
			expectedIncludes: []string{},
			expectedExcludes: []string{},
			expectedRenames:  map[string]string{},
			expectError:      false,
			reason:           "Empty configuration should result in no filtering (include all)",
		},
		{
			name:             "whitespace handling",
			servicesStr:      " UserService , LibraryService ",
			methodIncludeStr: " Get* , Find* ",
			methodExcludeStr: " *Internal , *Debug ",
			methodRenameStr:  " FindBooks : searchBooks , GetUser : fetchUser ",
			expectedServices: []string{"UserService", "LibraryService"},
			expectedIncludes: []string{"Get*", "Find*"},
			expectedExcludes: []string{"*Internal", "*Debug"},
			expectedRenames:  map[string]string{"FindBooks": "searchBooks", "GetUser": "fetchUser"},
			expectError:      false,
			reason:           "Whitespace should be trimmed from all configuration values",
		},
		{
			name:            "invalid rename format",
			servicesStr:     "UserService",
			methodRenameStr: "FindBooks->searchBooks", // Wrong separator
			expectError:     true,
			reason:          "Invalid rename format should return error",
		},
		{
			name:            "empty rename parts",
			servicesStr:     "UserService",
			methodRenameStr: "FindBooks:", // Empty new name
			expectError:     true,
			reason:          "Empty rename parts should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			criteria, err := ParseFromConfig(tt.servicesStr, tt.methodIncludeStr, tt.methodExcludeStr, tt.methodRenameStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseFromConfig() expected error but got none\nReason: %s", tt.reason)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseFromConfig() unexpected error: %v\nReason: %s", err, tt.reason)
				return
			}

			// Validate services
			if len(criteria.ServicesSet) != len(tt.expectedServices) {
				t.Errorf("ServicesSet length = %d, want %d", len(criteria.ServicesSet), len(tt.expectedServices))
			}
			for _, service := range tt.expectedServices {
				if !criteria.ServicesSet[service] {
					t.Errorf("Service %s not found in ServicesSet", service)
				}
			}

			// Validate includes
			if len(criteria.MethodIncludes) != len(tt.expectedIncludes) {
				t.Errorf("MethodIncludes length = %d, want %d", len(criteria.MethodIncludes), len(tt.expectedIncludes))
			}
			for i, expected := range tt.expectedIncludes {
				if i >= len(criteria.MethodIncludes) || criteria.MethodIncludes[i] != expected {
					t.Errorf("MethodIncludes[%d] = %s, want %s", i, criteria.MethodIncludes[i], expected)
				}
			}

			// Validate excludes
			if len(criteria.MethodExcludes) != len(tt.expectedExcludes) {
				t.Errorf("MethodExcludes length = %d, want %d", len(criteria.MethodExcludes), len(tt.expectedExcludes))
			}
			for i, expected := range tt.expectedExcludes {
				if i >= len(criteria.MethodExcludes) || criteria.MethodExcludes[i] != expected {
					t.Errorf("MethodExcludes[%d] = %s, want %s", i, criteria.MethodExcludes[i], expected)
				}
			}

			// Validate renames
			if len(criteria.MethodRenames) != len(tt.expectedRenames) {
				t.Errorf("MethodRenames length = %d, want %d", len(criteria.MethodRenames), len(tt.expectedRenames))
			}
			for oldName, expectedNewName := range tt.expectedRenames {
				if newName, exists := criteria.MethodRenames[oldName]; !exists || newName != expectedNewName {
					t.Errorf("MethodRenames[%s] = %s, want %s", oldName, newName, expectedNewName)
				}
			}
		})
	}
}

// TestServiceFilterCriteria tests the service-specific criteria factory.
// This ensures the specialized factory creates correct criteria for service filtering only.
func TestServiceFilterCriteria(t *testing.T) {
	servicesSet := map[string]bool{
		"UserService":    true,
		"LibraryService": true,
	}

	criteria := ServiceFilterCriteria(servicesSet)

	// Should have the provided services set
	if len(criteria.ServicesSet) != len(servicesSet) {
		t.Errorf("ServicesSet length = %d, want %d", len(criteria.ServicesSet), len(servicesSet))
	}

	for service := range servicesSet {
		if !criteria.ServicesSet[service] {
			t.Errorf("Service %s not found in ServicesSet", service)
		}
	}

	// Should have default values for other fields
	if len(criteria.MethodIncludes) != 0 {
		t.Error("MethodIncludes should be empty for service-only criteria")
	}

	if len(criteria.MethodExcludes) != 0 {
		t.Error("MethodExcludes should be empty for service-only criteria")
	}
}

// TestMethodFilterCriteria tests the method-specific criteria factory.
// This ensures the specialized factory creates correct criteria for method filtering only.
func TestMethodFilterCriteria(t *testing.T) {
	includes := []string{"Get*", "Find*"}
	excludes := []string{"*Internal", "*Debug"}
	renames := map[string]string{"FindBooks": "searchBooks"}

	criteria := MethodFilterCriteria(includes, excludes, renames)

	// Should have the provided method filters
	if len(criteria.MethodIncludes) != len(includes) {
		t.Errorf("MethodIncludes length = %d, want %d", len(criteria.MethodIncludes), len(includes))
	}

	if len(criteria.MethodExcludes) != len(excludes) {
		t.Errorf("MethodExcludes length = %d, want %d", len(criteria.MethodExcludes), len(excludes))
	}

	if len(criteria.MethodRenames) != len(renames) {
		t.Errorf("MethodRenames length = %d, want %d", len(criteria.MethodRenames), len(renames))
	}

	// Should have empty services set (no service filtering)
	if len(criteria.ServicesSet) != 0 {
		t.Error("ServicesSet should be empty for method-only criteria")
	}
}

// TestParseCommaSeparated tests comma-separated string parsing utility.
// This is critical for parsing all the list-based configuration options.
func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string   // Test case description
		input    string   // Input comma-separated string
		expected []string // Expected parsed result
		reason   string   // Why this test case is important
	}{
		{
			name:     "standard comma-separated list",
			input:    "UserService,LibraryService,GameService",
			expected: []string{"UserService", "LibraryService", "GameService"},
			reason:   "Standard case - multiple items separated by commas",
		},
		{
			name:     "single item",
			input:    "UserService",
			expected: []string{"UserService"},
			reason:   "Single item should work without commas",
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
			reason:   "Empty input should return empty slice",
		},
		{
			name:     "whitespace handling",
			input:    " UserService , LibraryService , GameService ",
			expected: []string{"UserService", "LibraryService", "GameService"},
			reason:   "Whitespace around items should be trimmed",
		},
		{
			name:     "empty items in list",
			input:    "UserService,,LibraryService,",
			expected: []string{"UserService", "LibraryService"},
			reason:   "Empty items should be filtered out",
		},
		{
			name:     "only whitespace and commas",
			input:    " , , ",
			expected: []string{},
			reason:   "String with only whitespace and commas should return empty slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommaSeparated(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("parseCommaSeparated(%s) length = %d, want %d\nReason: %s",
					tt.input, len(result), len(tt.expected), tt.reason)
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("parseCommaSeparated(%s)[%d] = %s, want %s\nReason: %s",
						tt.input, i, result[i], expected, tt.reason)
				}
			}
		})
	}
}

// TestParseMethodRenames tests method rename parsing from configuration strings.
// This ensures method renames are parsed correctly and invalid formats are caught.
func TestParseMethodRenames(t *testing.T) {
	tests := []struct {
		name        string            // Test case description
		input       string            // Input rename string
		expected    map[string]string // Expected parsed renames
		expectError bool              // Whether to expect an error
		reason      string            // Why this test case is important
	}{
		{
			name:     "standard renames",
			input:    "FindBooks:searchBooks,GetUser:fetchUser",
			expected: map[string]string{"FindBooks": "searchBooks", "GetUser": "fetchUser"},
			reason:   "Standard rename format should parse correctly",
		},
		{
			name:     "single rename",
			input:    "FindBooks:searchBooks",
			expected: map[string]string{"FindBooks": "searchBooks"},
			reason:   "Single rename should work",
		},
		{
			name:     "empty string",
			input:    "",
			expected: map[string]string{},
			reason:   "Empty input should return empty map",
		},
		{
			name:     "whitespace handling",
			input:    " FindBooks : searchBooks , GetUser : fetchUser ",
			expected: map[string]string{"FindBooks": "searchBooks", "GetUser": "fetchUser"},
			reason:   "Whitespace around colons and commas should be trimmed",
		},
		{
			name:        "invalid format - no colon",
			input:       "FindBooks->searchBooks",
			expectError: true,
			reason:      "Invalid separator should return error",
		},
		{
			name:        "invalid format - multiple colons",
			input:       "Find:Books:searchBooks",
			expectError: true,
			reason:      "Multiple colons should return error",
		},
		{
			name:        "empty old name",
			input:       ":searchBooks",
			expectError: true,
			reason:      "Empty old name should return error",
		},
		{
			name:        "empty new name",
			input:       "FindBooks:",
			expectError: true,
			reason:      "Empty new name should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseMethodRenames(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("parseMethodRenames(%s) expected error but got none\nReason: %s",
						tt.input, tt.reason)
				}
				return
			}

			if err != nil {
				t.Errorf("parseMethodRenames(%s) unexpected error: %v\nReason: %s",
					tt.input, err, tt.reason)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("parseMethodRenames(%s) length = %d, want %d\nReason: %s",
					tt.input, len(result), len(tt.expected), tt.reason)
				return
			}

			for oldName, expectedNewName := range tt.expected {
				if newName, exists := result[oldName]; !exists || newName != expectedNewName {
					t.Errorf("parseMethodRenames(%s)[%s] = %s, want %s\nReason: %s",
						tt.input, oldName, newName, expectedNewName, tt.reason)
				}
			}
		})
	}
}
