import { StatefulOptions as ConcreteStatefulOptions, StatefulMethodOptions as ConcreteStatefulMethodOptions, MessagePatch as ConcreteMessagePatch, PatchBatch as ConcretePatchBatch, PatchResponse as ConcretePatchResponse } from "./models";
/**
 * Enhanced factory with context-aware object construction
 */
export class WasmjsV1Factory {
    constructor() {
        /**
         * Enhanced factory method for StatefulOptions
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newStatefulOptions = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteStatefulOptions();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for StatefulMethodOptions
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newStatefulMethodOptions = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteStatefulMethodOptions();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for MessagePatch
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newMessagePatch = (parent, attributeName, attributeKey, data) => {
            const out = new ConcreteMessagePatch();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for PatchBatch
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newPatchBatch = (parent, attributeName, attributeKey, data) => {
            const out = new ConcretePatchBatch();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Enhanced factory method for PatchResponse
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw data to potentially populate from
         * @returns Factory result with instance and population status
         */
        this.newPatchResponse = (parent, attributeName, attributeKey, data) => {
            const out = new ConcretePatchResponse();
            // Factory does not populate by default - let deserializer handle it
            return { instance: out, fullyLoaded: false };
        };
        /**
         * Factory method for converting protobuf Timestamp data to native Date
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw protobuf timestamp data
         * @returns Factory result with Date instance
         */
        this.newTimestamp = (parent, attributeName, attributeKey, data) => {
            if (!data) {
                return { instance: new Date(), fullyLoaded: true };
            }
            let date;
            if (typeof data === 'string') {
                // Handle ISO string format
                date = new Date(data);
            }
            else if (data.seconds !== undefined) {
                // Handle protobuf format with seconds/nanos
                const seconds = typeof data.seconds === 'string'
                    ? parseInt(data.seconds, 10)
                    : data.seconds;
                const nanos = data.nanos || 0;
                date = new Date(seconds * 1000 + Math.floor(nanos / 1000000));
            }
            else {
                date = new Date();
            }
            return { instance: date, fullyLoaded: true };
        };
        /**
         * Factory method for converting protobuf FieldMask data to native string array
         * @param parent Parent object containing this field
         * @param attributeName Field name in parent object
         * @param attributeKey Array index, map key, or union tag (for containers)
         * @param data Raw protobuf field mask data
         * @returns Factory result with string array instance
         */
        this.newFieldMask = (parent, attributeName, attributeKey, data) => {
            if (!data) {
                return { instance: [], fullyLoaded: true };
            }
            let paths;
            if (Array.isArray(data)) {
                paths = data;
            }
            else if (data.paths && Array.isArray(data.paths)) {
                paths = data.paths;
            }
            else {
                paths = [];
            }
            return { instance: paths, fullyLoaded: true };
        };
    }
    /**
     * Get factory method for a fully qualified message type
     * Enables cross-package factory delegation
     */
    getFactoryMethod(messageType) {
        // Extract package from message type (e.g., "library.common.BaseMessage" -> "library.common")
        const parts = messageType.split('.');
        if (parts.length < 2) {
            return undefined;
        }
        const packageName = parts.slice(0, -1).join('.');
        const typeName = parts[parts.length - 1];
        const methodName = 'new' + typeName;
        // Check if this is our own package first
        const currentPackage = "wasmjs.v1";
        if (packageName === currentPackage) {
            return this[methodName];
        }
        // Check external type factory mappings
        const externalFactory = this.externalTypeFactories()[messageType];
        if (externalFactory) {
            return externalFactory;
        }
        // Delegate to appropriate dependency factory
        return undefined;
    }
    /**
     * Generic object deserializer that respects factory decisions
     */
    deserializeObject(instance, data) {
        if (!data || typeof data !== 'object')
            return instance;
        for (const [key, value] of Object.entries(data)) {
            if (value !== null && value !== undefined) {
                instance[key] = value;
            }
        }
        return instance;
    }
    // External type conversion methods
    /**
     * Mapping of external types to their factory methods
     */
    externalTypeFactories() {
        return {
            "google.protobuf.Timestamp": this.newTimestamp,
            "google.protobuf.FieldMask": this.newFieldMask,
        };
    }
    ;
    /**
     * Convert native Date to protobuf Timestamp format for serialization
     */
    serializeTimestamp(date) {
        if (!date)
            return null;
        return {
            seconds: Math.floor(date.getTime() / 1000).toString(),
            nanos: (date.getTime() % 1000) * 1000000
        };
    }
    /**
     * Convert native string array to protobuf FieldMask format for serialization
     */
    serializeFieldMask(paths) {
        if (!paths || !Array.isArray(paths))
            return null;
        return { paths };
    }
}
