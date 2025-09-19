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
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"github.com/panyam/protoc-gen-go-wasmjs/pkg/core"
)

// MessageInfo represents metadata about a protobuf message for generation.
// This is a simplified version that focuses on the information needed by filters.
type MessageInfo struct {
	Name               string // Message name (e.g., "Book")
	FullyQualifiedName string // Fully qualified name (e.g., "library.v1.Book")
	PackageName        string // Proto package name
	ProtoFile          string // Source proto file path
	IsNested           bool   // Whether this is a nested message
	IsMapEntry         bool   // Whether this is a synthetic map entry message
	Comment            string // Leading comment from proto
}

// MessageCollector provides logic for collecting and filtering protobuf messages.
// This handles the complex logic of finding all messages across packages while
// applying appropriate filtering criteria.
type MessageCollector struct {
	analyzer *core.ProtoAnalyzer
}

// NewMessageCollector creates a new message collector with the given analyzer.
func NewMessageCollector(analyzer *core.ProtoAnalyzer) *MessageCollector {
	return &MessageCollector{
		analyzer: analyzer,
	}
}

// CollectMessages collects all messages from the given files that meet the criteria.
// This handles recursive collection of nested messages and applies filtering rules.
// Returns CollectionResult with detailed statistics about what was found and collected.
func (mc *MessageCollector) CollectMessages(files []*protogen.File, criteria *FilterCriteria) CollectionResult[MessageInfo] {
	var messages []MessageInfo
	totalFound := 0

	for _, file := range files {
		// Collect top-level messages
		for _, message := range file.Messages {
			totalFound++

			// Apply map entry filtering
			if criteria.ExcludeMapEntries && message.Desc.IsMapEntry() {
				continue // Skip synthetic map entry messages
			}

			messageInfo := mc.buildMessageInfo(message, file, false)
			messages = append(messages, messageInfo)

			// Collect nested messages if not excluded
			if !criteria.ExcludeNestedMessages {
				nestedMessages, nestedCount := mc.collectNestedMessages(message, file, criteria)
				messages = append(messages, nestedMessages...)
				totalFound += nestedCount
			}
		}
	}

	return NewCollectionResult(messages, totalFound, len(files))
}

// collectNestedMessages recursively collects nested message definitions.
// Returns the collected messages and the total count (including filtered ones).
func (mc *MessageCollector) collectNestedMessages(message *protogen.Message, file *protogen.File, criteria *FilterCriteria) ([]MessageInfo, int) {
	var nestedMessages []MessageInfo
	totalCount := 0

	for _, nested := range message.Messages {
		totalCount++

		// Apply map entry filtering
		if criteria.ExcludeMapEntries && nested.Desc.IsMapEntry() {
			continue // Skip synthetic map entry messages
		}

		nestedInfo := mc.buildMessageInfo(nested, file, true)
		nestedMessages = append(nestedMessages, nestedInfo)

		// Recursively collect deeply nested messages
		if !criteria.ExcludeNestedMessages {
			deeplyNested, deepCount := mc.collectNestedMessages(nested, file, criteria)
			nestedMessages = append(nestedMessages, deeplyNested...)
			totalCount += deepCount
		}
	}

	return nestedMessages, totalCount
}

// buildMessageInfo constructs MessageInfo from a protogen.Message.
// This extracts the essential metadata needed for filtering and generation decisions.
func (mc *MessageCollector) buildMessageInfo(message *protogen.Message, file *protogen.File, isNested bool) MessageInfo {
	messageName := string(message.Desc.Name())
	packageName := string(file.Desc.Package())

	// Build fully qualified name (e.g., "library.v1.Book")
	fullyQualifiedName := packageName + "." + messageName

	return MessageInfo{
		Name:               messageName,
		FullyQualifiedName: fullyQualifiedName,
		PackageName:        packageName,
		ProtoFile:          file.Desc.Path(),
		IsNested:           isNested,
		IsMapEntry:         message.Desc.IsMapEntry(),
		Comment:            strings.TrimSpace(string(message.Comments.Leading)),
	}
}

// HasAnyMessages checks if any messages would be collected with the given criteria.
// This is useful for early termination when no messages would be generated.
func (mc *MessageCollector) HasAnyMessages(files []*protogen.File, criteria *FilterCriteria) bool {
	for _, file := range files {
		for _, message := range file.Messages {
			// Apply map entry filtering
			if criteria.ExcludeMapEntries && message.Desc.IsMapEntry() {
				continue
			}

			// Found at least one message that would be included
			return true
		}
	}
	return false
}

// CollectTopLevelMessages collects only top-level messages (no nesting).
// This is a specialized method for cases where nested messages should be ignored.
func (mc *MessageCollector) CollectTopLevelMessages(files []*protogen.File, criteria *FilterCriteria) CollectionResult[MessageInfo] {
	var messages []MessageInfo
	totalFound := 0

	for _, file := range files {
		for _, message := range file.Messages {
			totalFound++

			// Apply map entry filtering
			if criteria.ExcludeMapEntries && message.Desc.IsMapEntry() {
				continue
			}

			messageInfo := mc.buildMessageInfo(message, file, false)
			messages = append(messages, messageInfo)
		}
	}

	return NewCollectionResult(messages, totalFound, len(files))
}

// CollectMessagesByPackage collects messages grouped by package name.
// This is useful for generating package-specific TypeScript files.
func (mc *MessageCollector) CollectMessagesByPackage(files []*protogen.File, criteria *FilterCriteria) map[string]CollectionResult[MessageInfo] {
	packageMessages := make(map[string]CollectionResult[MessageInfo])

	// Group files by package
	packageFiles := make(map[string][]*protogen.File)
	for _, file := range files {
		packageName := string(file.Desc.Package())
		packageFiles[packageName] = append(packageFiles[packageName], file)
	}

	// Collect messages for each package
	for packageName, pkgFiles := range packageFiles {
		result := mc.CollectMessages(pkgFiles, criteria)
		packageMessages[packageName] = result
	}

	return packageMessages
}

// GetMessageNames extracts just the message names from a collection result.
// This is useful for template generation where only names are needed.
func (mc *MessageCollector) GetMessageNames(result CollectionResult[MessageInfo]) []string {
	var names []string
	for _, msg := range result.Items {
		names = append(names, msg.Name)
	}
	return names
}

// GetMessagesByFile groups messages by their source proto file.
// This is useful for generating file-specific imports and references.
func (mc *MessageCollector) GetMessagesByFile(result CollectionResult[MessageInfo]) map[string][]MessageInfo {
	messagesByFile := make(map[string][]MessageInfo)

	for _, msg := range result.Items {
		messagesByFile[msg.ProtoFile] = append(messagesByFile[msg.ProtoFile], msg)
	}

	return messagesByFile
}
