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

/**
 * Field type enumeration for proto field types
 */
export enum FieldType {
  STRING = "string",
  NUMBER = "number", 
  BOOLEAN = "boolean",
  MESSAGE = "message",
  REPEATED = "repeated",
  MAP = "map",
  ONEOF = "oneof"
}

/**
 * Schema interface for field definitions
 */
export interface FieldSchema {
  name: string;
  type: FieldType;
  id: number; // Proto field number (e.g., text_query = 1)
  messageType?: string; // For MESSAGE type fields
  repeated?: boolean; // For array fields
  mapKeyType?: FieldType; // For MAP type fields
  mapValueType?: FieldType | string; // For MAP type fields
  oneofGroup?: string; // For ONEOF fields
  optional?: boolean;
}

/**
 * Message schema interface
 */
export interface MessageSchema {
  name: string;
  fields: FieldSchema[];
  oneofGroups?: string[]; // List of oneof group names
}
