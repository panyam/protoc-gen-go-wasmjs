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

