import { WasmjsV1Deserializer } from "./deserializer";
/**
 * Configuration for stateful services
 */
export class StatefulOptions {
    constructor() {
        /** Whether stateful proxy generation is enabled */
        this.enabled = false;
        /** The fully qualified name of the message type that represents the state
       (e.g., "example.Game", "library.Map") */
        this.stateMessageType = "";
        /** Strategy for resolving conflicts when multiple changes occur */
        this.conflictResolution = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized StatefulOptions instance or null if creation failed
     */
    static from(data) {
        return WasmjsV1Deserializer.from(StatefulOptions.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
StatefulOptions.MESSAGE_TYPE = "wasmjs.v1.StatefulOptions";
/**
 * Configuration for stateful methods
 */
export class StatefulMethodOptions {
    constructor() {
        /** Whether this method returns patch operations instead of full objects */
        this.returnsPatches = false;
        /** Whether changes from this method should be broadcast to other clients */
        this.broadcasts = false;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized StatefulMethodOptions instance or null if creation failed
     */
    static from(data) {
        return WasmjsV1Deserializer.from(StatefulMethodOptions.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
StatefulMethodOptions.MESSAGE_TYPE = "wasmjs.v1.StatefulMethodOptions";
/**
 * A single patch operation on a protobuf message field
 */
export class MessagePatch {
    constructor() {
        /** The type of operation to perform */
        this.operation = 0;
        /** Path to the field being modified (e.g., "players[2].name", "places['tile_123'].latitude") */
        this.fieldPath = "";
        /** The new value to set (for SET, INSERT_LIST, INSERT_MAP operations)
       Encoded as JSON for type flexibility */
        this.valueJson = "";
        /** Index for list operations (INSERT_LIST, REMOVE_LIST, MOVE_LIST) */
        this.index = 0;
        /** Map key for map operations (INSERT_MAP, REMOVE_MAP) */
        this.key = "";
        /** Source index for MOVE_LIST operations */
        this.oldIndex = 0;
        /** Monotonically increasing change number for ordering */
        this.changeNumber = 0;
        /** Timestamp when the change was created (microseconds since epoch) */
        this.timestamp = 0;
        /** Optional user ID who made the change (for conflict resolution) */
        this.userId = "";
        /** Optional transaction ID to group related patches */
        this.transactionId = "";
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized MessagePatch instance or null if creation failed
     */
    static from(data) {
        return WasmjsV1Deserializer.from(MessagePatch.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
MessagePatch.MESSAGE_TYPE = "wasmjs.v1.MessagePatch";
/**
 * A batch of patches applied to a single entity
 */
export class PatchBatch {
    constructor() {
        /** The fully qualified protobuf message type (e.g., "example.Game") */
        this.messageType = "";
        /** The unique identifier of the entity being modified */
        this.entityId = "";
        /** List of patches to apply in order */
        this.patches = [];
        /** The highest change number in this batch */
        this.changeNumber = 0;
        /** Source of the changes */
        this.source = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PatchBatch instance or null if creation failed
     */
    static from(data) {
        return WasmjsV1Deserializer.from(PatchBatch.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
PatchBatch.MESSAGE_TYPE = "wasmjs.v1.PatchBatch";
/**
 * Response message for methods that return patches
 */
export class PatchResponse {
    constructor() {
        /** The patches to apply */
        this.patchBatches = [];
        /** Success status */
        this.success = false;
        /** Error message if success is false */
        this.errorMessage = "";
        /** The new change number after applying these patches */
        this.newChangeNumber = 0;
    }
    /**
     * Create and deserialize an instance from raw data
     * @param data Raw data to deserialize
     * @returns Deserialized PatchResponse instance or null if creation failed
     */
    static from(data) {
        return WasmjsV1Deserializer.from(PatchResponse.MESSAGE_TYPE, data);
    }
}
/**
 * Fully qualified message type for schema resolution
 */
PatchResponse.MESSAGE_TYPE = "wasmjs.v1.PatchResponse";
