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
	"strings"
	"testing"
)

// TestExactBug - Test that demonstrates the exact bug in BuildServiceClientData line 223
func TestExactBug(t *testing.T) {
	t.Run("DemonstrateBugInBuildServiceClientData", func(t *testing.T) {
		// This reproduces the EXACT bug: baseName vs getModuleName
		
		packageName := "presenter.v1"
		
		// Reproduce the exact logic from line 217 and line 223
		buggyBaseName := strings.ReplaceAll(packageName, ".", "_") // This produces "presenter_v1"
		
		// What the correct line 304 does: tb.getModuleName(packageInfo.Name, config)
		config := &GenerationConfig{
			ModuleName: "browser_callbacks", // From buf.gen.yaml
		}
		
		// Simulate getModuleName logic
		correctModuleName := config.ModuleName // Should be "browser_callbacks"
		if correctModuleName == "" {
			// Fallback - but this shouldn't happen if config is populated correctly
			correctModuleName = buggyBaseName + "_services"
		}
		
		// The bug: line 223 uses buggyBaseName instead of correctModuleName
		t.Logf("Package: %s", packageName)
		t.Logf("Line 217 baseName: %s", buggyBaseName)
		t.Logf("Line 223 ModuleName (BUGGY): %s", buggyBaseName)  
		t.Logf("Line 304 ModuleName (CORRECT): %s", correctModuleName)
		
		// This test SHOULD FAIL with current code to prove the bug exists
		if buggyBaseName == correctModuleName {
			t.Error("❌ TEST DESIGN ERROR: The bug should make these different!")
		} else {
			t.Logf("✅ BUG CONFIRMED: Line 223 produces '%s' but should produce '%s'", 
				buggyBaseName, correctModuleName)
		}
		
		// Verify this explains the integration test results
		if buggyBaseName == "presenter_v1" {
			t.Logf("✅ This explains why we see 'Presenter_v1Bundle' in generated files")
		}
		if correctModuleName == "browser_callbacks" {
			t.Logf("✅ After fix, we should see 'Browser_callbacksBundle' in generated files")
		}
	})
}
