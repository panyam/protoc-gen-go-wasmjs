// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

/**
 * Conflict resolution strategies for stateful objects
 */
export enum ConflictResolution {
  /** Use change numbers to determine order (default) */
  CHANGE_NUMBER_BASED = 0,
  /** Use timestamps for ordering (requires synchronized clocks) */
  TIMESTAMP_BASED = 1,
  /** Last writer wins (simple but may lose data) */
  LAST_WRITER_WINS = 2,
}

/**
 * Patch operations for modifying protobuf message fields
 */
export enum PatchOperation {
  /** Set a field value (scalar, message, or replace entire field) */
  SET = 0,
  /** Insert an item into a repeated field at a specific index */
  INSERT_LIST = 1,
  /** Remove an item from a repeated field at a specific index */
  REMOVE_LIST = 2,
  /** Move an item within a repeated field from one index to another */
  MOVE_LIST = 3,
  /** Insert or update a key-value pair in a map field */
  INSERT_MAP = 4,
  /** Remove a key-value pair from a map field */
  REMOVE_MAP = 5,
  /** Clear all items from a repeated field */
  CLEAR_LIST = 6,
  /** Clear all key-value pairs from a map field */
  CLEAR_MAP = 7,
}

/**
 * Source of patch changes
 */
export enum PatchSource {
  /** Changes originating from local user actions */
  LOCAL = 0,
  /** Changes from remote users via real-time sync */
  REMOTE = 1,
  /** Authoritative changes from the server */
  SERVER = 2,
  /** Changes loaded from persistent storage */
  STORAGE = 3,
}


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
 * Configuration for async methods
 */
export interface AsyncMethodOptions {
  /** Whether this method should be generated as async with callback parameter */
  isAsync: boolean;
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
  metadata: Record<string, string>;
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

