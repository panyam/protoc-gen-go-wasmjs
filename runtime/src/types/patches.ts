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
 * Patch operations for modifying protobuf message fields
 */
export enum PatchOperation {
  SET = 'SET',
  INSERT_LIST = 'INSERT_LIST',
  REMOVE_LIST = 'REMOVE_LIST',
  MOVE_LIST = 'MOVE_LIST',
  INSERT_MAP = 'INSERT_MAP',
  REMOVE_MAP = 'REMOVE_MAP',
  CLEAR_LIST = 'CLEAR_LIST',
  CLEAR_MAP = 'CLEAR_MAP',
}

/**
 * A single patch operation on a protobuf message field
 */
export interface MessagePatch {
  /** The type of operation to perform */
  operation: PatchOperation;
  
  /** Path to the field being modified (e.g., "players[2].name", "places['tile_123'].latitude") */
  fieldPath: string;
  
  /** The new value to set (for SET, INSERT_LIST, INSERT_MAP operations) */
  value?: any;
  
  /** Index for list operations (INSERT_LIST, REMOVE_LIST, MOVE_LIST) */
  index?: number;
  
  /** Map key for map operations (INSERT_MAP, REMOVE_MAP) */
  key?: string;
  
  /** Source index for MOVE_LIST operations */
  oldIndex?: number;
  
  /** Monotonically increasing change number for ordering */
  changeNumber: number;
  
  /** Timestamp when the change was created (microseconds since epoch) */
  timestamp: number;
  
  /** Optional user ID who made the change (for conflict resolution) */
  userId?: string;
  
  /** Optional transaction ID to group related patches */
  transactionId?: string;
}

/**
 * A batch of patches applied to a single entity
 */
export interface PatchBatch {
  /** The fully qualified protobuf message type (e.g., "example.Game") */
  messageType: string;
  
  /** The unique identifier of the entity being modified */
  entityId: string;
  
  /** List of patches to apply in order */
  patches: MessagePatch[];
  
  /** The highest change number in this batch */
  changeNumber: number;
  
  /** Source of the changes */
  source: PatchSource;
  
  /** Optional metadata about the batch */
  metadata?: Record<string, string>;
}

/**
 * Source of patch changes
 */
export enum PatchSource {
  LOCAL = 'LOCAL',
  REMOTE = 'REMOTE', 
  SERVER = 'SERVER',
  STORAGE = 'STORAGE',
}

/**
 * Response message for methods that return patches
 */
export interface PatchResponse {
  /** The patches to apply */
  patchBatches: PatchBatch[];
  
  /** Success status */
  success: boolean;
  
  /** Error message if success is false */
  errorMessage?: string;
  
  /** The new change number after applying these patches */
  newChangeNumber: number;
}

/**
 * Transport interface for receiving patch updates
 */
export interface ChangeTransport {
  /** Register callback for incoming changes */
  onChanges(callback: (batches: PatchBatch[]) => void): void;
  
  /** Disconnect from the transport */
  disconnect(): void;
}
