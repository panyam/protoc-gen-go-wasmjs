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

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// TestNewMessageCollector tests message collector creation.
// This ensures the collector is properly initialized with its dependencies.
func TestNewMessageCollector(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	collector := NewMessageCollector(analyzer)

	if collector == nil {
		t.Error("NewMessageCollector() should return non-nil collector")
	}

	if collector.analyzer != analyzer {
		t.Error("MessageCollector should store the provided analyzer")
	}
}

// TestNewCollectionResult tests collection result creation.
// This tests the generic result structure used by all collectors.
func TestNewCollectionResult(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	totalFound := 5
	filesScanned := 3

	result := NewCollectionResult(items, totalFound, filesScanned)

	if len(result.Items) != len(items) {
		t.Errorf("Items length = %d, want %d", len(result.Items), len(items))
	}

	for i, expected := range items {
		if result.Items[i] != expected {
			t.Errorf("Items[%d] = %s, want %s", i, result.Items[i], expected)
		}
	}

	if result.TotalFound != totalFound {
		t.Errorf("TotalFound = %d, want %d", result.TotalFound, totalFound)
	}

	if result.FilesScanned != filesScanned {
		t.Errorf("FilesScanned = %d, want %d", result.FilesScanned, filesScanned)
	}
}

// TestMessageCollector_GetMessageNames tests name extraction from collection results.
// This is a pure function that extracts just the names for template usage.
func TestMessageCollector_GetMessageNames(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	collector := NewMessageCollector(analyzer)

	// Create test message info
	messages := []MessageInfo{
		{Name: "Book", FullyQualifiedName: "library.v1.Book"},
		{Name: "User", FullyQualifiedName: "library.v1.User"},
		{Name: "Checkout", FullyQualifiedName: "library.v1.Checkout"},
	}

	result := NewCollectionResult(messages, len(messages), 1)
	names := collector.GetMessageNames(result)

	expectedNames := []string{"Book", "User", "Checkout"}

	if len(names) != len(expectedNames) {
		t.Errorf("GetMessageNames() length = %d, want %d", len(names), len(expectedNames))
		return
	}

	for i, expected := range expectedNames {
		if names[i] != expected {
			t.Errorf("GetMessageNames()[%d] = %s, want %s", i, names[i], expected)
		}
	}
}

// TestMessageCollector_GetMessagesByFile tests grouping messages by source file.
// This is important for generating accurate TypeScript imports that reference the correct files.
func TestMessageCollector_GetMessagesByFile(t *testing.T) {
	analyzer := core.NewProtoAnalyzer()
	collector := NewMessageCollector(analyzer)

	// Create test messages from different files
	messages := []MessageInfo{
		{Name: "Book", ProtoFile: "proto/library/v1/book.proto"},
		{Name: "User", ProtoFile: "proto/library/v1/user.proto"},
		{Name: "Checkout", ProtoFile: "proto/library/v1/book.proto"}, // Same file as Book
	}

	result := NewCollectionResult(messages, len(messages), 2)
	messagesByFile := collector.GetMessagesByFile(result)

	// Should have 2 files
	if len(messagesByFile) != 2 {
		t.Errorf("GetMessagesByFile() returned %d files, want 2", len(messagesByFile))
		return
	}

	// Check book.proto file should have 2 messages
	bookMessages := messagesByFile["proto/library/v1/book.proto"]
	if len(bookMessages) != 2 {
		t.Errorf("book.proto should have 2 messages, got %d", len(bookMessages))
	}

	// Check user.proto file should have 1 message
	userMessages := messagesByFile["proto/library/v1/user.proto"]
	if len(userMessages) != 1 {
		t.Errorf("user.proto should have 1 message, got %d", len(userMessages))
	}

	// Verify specific messages are in correct files
	found := false
	for _, msg := range bookMessages {
		if msg.Name == "Book" || msg.Name == "Checkout" {
			found = true
		}
	}
	if !found {
		t.Error("book.proto should contain Book and Checkout messages")
	}

	if userMessages[0].Name != "User" {
		t.Errorf("user.proto should contain User message, got %s", userMessages[0].Name)
	}
}

/*
Note: The following tests would require protogen.File and protogen.Message mocks
which are complex to create. In a real implementation, these would be integration
tests. Here's what would be tested and why:

// TestMessageCollector_CollectMessages would test the main collection logic
// This is the core functionality that finds and filters messages across packages.
//
// Test cases would include:
// - Files with top-level messages only
// - Files with nested messages (should include if not excluded)
// - Files with map entry messages (should exclude if configured)
// - Mixed files with various message types
// - Empty files (should handle gracefully)
// - Deeply nested message hierarchies
//
// Why important: This is the main entry point for message collection and affects
// what TypeScript interfaces and classes are generated. Bugs here cause missing
// types or incorrect type definitions.

func TestMessageCollector_CollectMessages(t *testing.T) {
	// Would require creating mock protogen.File objects with:
	// - File descriptors with package names
	// - Message arrays with nested structures
	// - Map entry message detection
	// - Comment extraction
}

// TestMessageCollector_HasAnyMessages would test early termination logic
// This is used for performance optimization when packages have no messages.
//
// Test cases would include:
// - Files with messages (should return true)
// - Files with only map entries and map entries excluded (should return false)
// - Empty files (should return false)
// - Files with only nested messages when nested excluded (should return false)
//
// Why important: Prevents unnecessary processing of packages that won't generate
// TypeScript code, improving performance and reducing generated code size.

func TestMessageCollector_HasAnyMessages(t *testing.T) {
	// Would require mock files with various message configurations
}

// TestMessageCollector_CollectTopLevelMessages would test top-level only collection
// This specialized method is used when nested messages should be ignored.
//
// Test cases would include:
// - Files with both top-level and nested messages (should return only top-level)
// - Files with only nested messages (should return empty)
// - Files with only top-level messages (should return all)
//
// Why important: Some use cases need only top-level types for simpler generation.

func TestMessageCollector_CollectTopLevelMessages(t *testing.T) {
	// Would require distinguishing top-level vs nested messages
}

Integration Testing Approach:
The message collection logic is integration-tested through examples:

examples/library/proto/library/v1/library.proto includes:
- Top-level messages: Book, User, FindBooksRequest, etc.
- Various field types: strings, numbers, booleans, repeated fields
- Messages with comments and metadata

This validates:
- Correct message discovery and metadata extraction
- Proper handling of different message types and structures
- Cross-package message dependencies (if present)
*/
