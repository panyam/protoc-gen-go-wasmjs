// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Configuration for stateful services
 */
export interface StatefulOptions {
  /** Whether stateful proxy generation is enabled */
  enabled: boolean;
  /** The fully qualified name of the message type that represents the state
 (e.g., "example.Game", "library.Map") */
  stateMessageType: string;
  /** Strategy for resolving conflicts when multiple changes occur */
  conflictResolution: ConflictResolution;
}


/**
 * Configuration for stateful methods
 */
export interface StatefulMethodOptions {
  /** Whether this method returns patch operations instead of full objects */
  returnsPatches: boolean;
  /** Whether changes from this method should be broadcast to other clients */
  broadcasts: boolean;
}


/**
 * A single patch operation on a protobuf message field
 */
export interface MessagePatch {
  /** The type of operation to perform */
  operation: PatchOperation;
  /** Path to the field being modified (e.g., "players[2].name", "places['tile_123'].latitude") */
  fieldPath: string;
  /** The new value to set (for SET, INSERT_LIST, INSERT_MAP operations)
 Encoded as JSON for type flexibility */
  valueJson: string;
  /** Index for list operations (INSERT_LIST, REMOVE_LIST, MOVE_LIST) */
  index: number;
  /** Map key for map operations (INSERT_MAP, REMOVE_MAP) */
  key: string;
  /** Source index for MOVE_LIST operations */
  oldIndex: number;
  /** Monotonically increasing change number for ordering */
  changeNumber: number;
  /** Timestamp when the change was created (microseconds since epoch) */
  timestamp: number;
  /** Optional user ID who made the change (for conflict resolution) */
  userId: string;
  /** Optional transaction ID to group related patches */
  transactionId: string;
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
  patches?: MessagePatch[];
  /** The highest change number in this batch */
  changeNumber: number;
  /** Source of the changes */
  source: PatchSource;
  /** Optional metadata about the batch */
  metadata?: Map<string, string>;
}


/**
 * Response message for methods that return patches
 */
export interface PatchResponse {
  /** The patches to apply */
  patchBatches?: PatchBatch[];
  /** Success status */
  success: boolean;
  /** Error message if success is false */
  errorMessage: string;
  /** The new change number after applying these patches */
  newChangeNumber: number;
}

