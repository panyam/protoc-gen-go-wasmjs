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
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

// FilterResult represents the outcome of a filtering operation.
// This provides rich information about why something was included or excluded,
// enabling better debugging and validation of filtering behavior.
type FilterResult struct {
	Include bool   // Whether the item should be included
	Reason  string // Human-readable reason for the decision
}

// Included creates a FilterResult for items that should be included.
func Included(reason string) FilterResult {
	return FilterResult{Include: true, Reason: reason}
}

// Excluded creates a FilterResult for items that should be excluded.
func Excluded(reason string) FilterResult {
	return FilterResult{Include: false, Reason: reason}
}

// ServiceFilterResult contains the result of service filtering with additional metadata.
// This is used when filtering services and provides context about the service type.
type ServiceFilterResult struct {
	FilterResult
	IsBrowserProvided bool   // Whether this service is browser-provided
	CustomName        string // Custom JavaScript name (if any)
}

// MethodFilterResult contains the result of method filtering with additional metadata.
// This provides information about method transformations and special handling.
type MethodFilterResult struct {
	FilterResult
	CustomJSName      string // Custom JavaScript method name (if any)
	IsAsync           bool   // Whether method requires async/callback handling
	IsServerStreaming bool   // Whether method uses server-side streaming
}

// PackageFilterResult contains the result of package filtering.
// This helps understand why certain packages are included or excluded from generation.
type PackageFilterResult struct {
	FilterResult
	HasServices bool // Whether package contains services
	HasMessages bool // Whether package contains messages
	HasEnums    bool // Whether package contains enums
}

// CollectionResult represents the result of collecting items (messages, enums) from proto files.
// This provides insight into what was found and what filtering was applied.
type CollectionResult[T any] struct {
	Items         []T    // Collected items that passed filtering
	TotalFound    int    // Total items found before filtering
	FilesScanned  int    // Number of files examined
	SkippedReason string // Reason for skipped items (if any)
}

// NewCollectionResult creates a new collection result with the given items.
func NewCollectionResult[T any](items []T, totalFound, filesScanned int) CollectionResult[T] {
	return CollectionResult[T]{
		Items:        items,
		TotalFound:   totalFound,
		FilesScanned: filesScanned,
	}
}

// FilterStats provides statistics about filtering operations.
// This is useful for debugging filtering behavior and understanding generation scope.
type FilterStats struct {
	// Service filtering stats
	ServicesTotal    int // Total services found
	ServicesIncluded int // Services included after filtering
	ServicesExcluded int // Services excluded by filters

	// Method filtering stats
	MethodsTotal    int // Total methods found
	MethodsIncluded int // Methods included after filtering
	MethodsExcluded int // Methods excluded by filters

	// Content collection stats
	MessagesTotal int // Total messages found
	EnumsTotal    int // Total enums found
	PackagesTotal int // Total packages processed
}

// NewFilterStats creates an empty FilterStats instance.
func NewFilterStats() *FilterStats {
	return &FilterStats{}
}

// AddServiceResult updates statistics with a service filtering result.
func (fs *FilterStats) AddServiceResult(result ServiceFilterResult) {
	fs.ServicesTotal++
	if result.Include {
		fs.ServicesIncluded++
	} else {
		fs.ServicesExcluded++
	}
}

// AddMethodResult updates statistics with a method filtering result.
func (fs *FilterStats) AddMethodResult(result MethodFilterResult) {
	fs.MethodsTotal++
	if result.Include {
		fs.MethodsIncluded++
	} else {
		fs.MethodsExcluded++
	}
}

// AddCollectionStats updates statistics with collection results.
func (fs *FilterStats) AddCollectionStats(messagesFound, enumsFound, packagesScanned int) {
	fs.MessagesTotal += messagesFound
	fs.EnumsTotal += enumsFound
	fs.PackagesTotal += packagesScanned
}

// Summary returns a human-readable summary of filtering statistics.
// This is useful for debugging and understanding what was generated.
func (fs *FilterStats) Summary() string {
	return fmt.Sprintf(
		"Filtering Summary: %d/%d services, %d/%d methods, %d messages, %d enums from %d packages",
		fs.ServicesIncluded, fs.ServicesTotal,
		fs.MethodsIncluded, fs.MethodsTotal,
		fs.MessagesTotal, fs.EnumsTotal, fs.PackagesTotal,
	)
}

// FilterContext provides context for filtering operations.
// This carries information about the current state of generation that filters may need.
type FilterContext struct {
	CurrentPackage string           // Package currently being processed
	AllFiles       []*protogen.File // All files in the generation request
	PackageFiles   []*protogen.File // Files in the current package
	Stats          *FilterStats     // Statistics collector
}

// NewFilterContext creates a new filter context for the given package and files.
func NewFilterContext(currentPackage string, packageFiles, allFiles []*protogen.File) *FilterContext {
	return &FilterContext{
		CurrentPackage: currentPackage,
		AllFiles:       allFiles,
		PackageFiles:   packageFiles,
		Stats:          NewFilterStats(),
	}
}
