/**
 * Factory result interface
 */
interface FactoryResult<T> {
    instance: T;
    fullyLoaded: boolean;
}
/**
 * Factory method type for creating instances
 */
type FactoryMethod<T = any> = (parent?: any, attributeName?: string, attributeKey?: string | number, data?: any) => FactoryResult<T>;
/**
 * Factory interface that deserializer expects
 */
interface FactoryInterface {
    /**
     * Get factory method for a fully qualified message type
     * This enables cross-package factory delegation
     */
    getFactoryMethod?(messageType: string): FactoryMethod | undefined;
}

/**
 * Patch operations for modifying protobuf message fields
 */
declare enum PatchOperation {
    SET = "SET",
    INSERT_LIST = "INSERT_LIST",
    REMOVE_LIST = "REMOVE_LIST",
    MOVE_LIST = "MOVE_LIST",
    INSERT_MAP = "INSERT_MAP",
    REMOVE_MAP = "REMOVE_MAP",
    CLEAR_LIST = "CLEAR_LIST",
    CLEAR_MAP = "CLEAR_MAP"
}
/**
 * A single patch operation on a protobuf message field
 */
interface MessagePatch {
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
interface PatchBatch {
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
declare enum PatchSource {
    LOCAL = "LOCAL",
    REMOTE = "REMOTE",
    SERVER = "SERVER",
    STORAGE = "STORAGE"
}
/**
 * Response message for methods that return patches
 */
interface PatchResponse {
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
interface ChangeTransport {
    /** Register callback for incoming changes */
    onChanges(callback: (batches: PatchBatch[]) => void): void;
    /** Disconnect from the transport */
    disconnect(): void;
}

export { type ChangeTransport, type FactoryInterface, type FactoryMethod, type FactoryResult, type MessagePatch, type PatchBatch, PatchOperation, type PatchResponse, PatchSource };
