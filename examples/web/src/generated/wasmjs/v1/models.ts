import { StatefulOptions as StatefulOptionsInterface, StatefulMethodOptions as StatefulMethodOptionsInterface, AsyncMethodOptions as AsyncMethodOptionsInterface, MessagePatch as MessagePatchInterface, PatchBatch as PatchBatchInterface, PatchResponse as PatchResponseInterface, ConflictResolution, PatchOperation, PatchSource } from "./interfaces";




/**
 * Configuration for stateful services
 */
export class StatefulOptions implements StatefulOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "wasmjs.v1.StatefulOptions";
  readonly __MESSAGE_TYPE = StatefulOptions.MESSAGE_TYPE;

  /** Whether stateful proxy generation is enabled */
  enabled: boolean = false;
  /** The fully qualified name of the message type that represents the state
 (e.g., "example.Game", "library.Map") */
  stateMessageType: string = "";
  /** Strategy for resolving conflicts when multiple changes occur */
  conflictResolution: ConflictResolution = ConflictResolution.CHANGE_NUMBER_BASED;

  
}


/**
 * Configuration for stateful methods
 */
export class StatefulMethodOptions implements StatefulMethodOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "wasmjs.v1.StatefulMethodOptions";
  readonly __MESSAGE_TYPE = StatefulMethodOptions.MESSAGE_TYPE;

  /** Whether this method returns patch operations instead of full objects */
  returnsPatches: boolean = false;
  /** Whether changes from this method should be broadcast to other clients */
  broadcasts: boolean = false;

  
}


/**
 * Configuration for async methods
 */
export class AsyncMethodOptions implements AsyncMethodOptionsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "wasmjs.v1.AsyncMethodOptions";
  readonly __MESSAGE_TYPE = AsyncMethodOptions.MESSAGE_TYPE;

  /** Whether this method should be generated as async with callback parameter */
  isAsync: boolean = false;

  
}


/**
 * A single patch operation on a protobuf message field
 */
export class MessagePatch implements MessagePatchInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "wasmjs.v1.MessagePatch";
  readonly __MESSAGE_TYPE = MessagePatch.MESSAGE_TYPE;

  /** The type of operation to perform */
  operation: PatchOperation = PatchOperation.SET;
  /** Path to the field being modified (e.g., "players[2].name", "places['tile_123'].latitude") */
  fieldPath: string = "";
  /** The new value to set (for SET, INSERT_LIST, INSERT_MAP operations)
 Encoded as JSON for type flexibility */
  valueJson: string = "";
  /** Index for list operations (INSERT_LIST, REMOVE_LIST, MOVE_LIST) */
  index: number = 0;
  /** Map key for map operations (INSERT_MAP, REMOVE_MAP) */
  key: string = "";
  /** Source index for MOVE_LIST operations */
  oldIndex: number = 0;
  /** Monotonically increasing change number for ordering */
  changeNumber: number = 0;
  /** Timestamp when the change was created (microseconds since epoch) */
  timestamp: number = 0;
  /** Optional user ID who made the change (for conflict resolution) */
  userId: string = "";
  /** Optional transaction ID to group related patches */
  transactionId: string = "";

  
}


/**
 * A batch of patches applied to a single entity
 */
export class PatchBatch implements PatchBatchInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "wasmjs.v1.PatchBatch";
  readonly __MESSAGE_TYPE = PatchBatch.MESSAGE_TYPE;

  /** The fully qualified protobuf message type (e.g., "example.Game") */
  messageType: string = "";
  /** The unique identifier of the entity being modified */
  entityId: string = "";
  /** List of patches to apply in order */
  patches: MessagePatch[] = [];
  /** The highest change number in this batch */
  changeNumber: number = 0;
  /** Source of the changes */
  source: PatchSource = PatchSource.LOCAL;
  /** Optional metadata about the batch */
  metadata: Record<string, string> = {};

  
}


/**
 * Response message for methods that return patches
 */
export class PatchResponse implements PatchResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "wasmjs.v1.PatchResponse";
  readonly __MESSAGE_TYPE = PatchResponse.MESSAGE_TYPE;

  /** The patches to apply */
  patchBatches: PatchBatch[] = [];
  /** Success status */
  success: boolean = false;
  /** Error message if success is false */
  errorMessage: string = "";
  /** The new change number after applying these patches */
  newChangeNumber: number = 0;

  
}


