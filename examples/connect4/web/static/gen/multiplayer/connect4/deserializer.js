// Generated TypeScript schema-aware deserializer
// DO NOT EDIT - This file is auto-generated
import { FieldType } from "./deserializer_schemas";
import { Connect4Factory } from "./factory";
import { Connect4SchemaRegistry } from "./schemas";
// Shared factory instance to avoid creating new instances on every deserializer construction
const DEFAULT_FACTORY = new Connect4Factory();
/**
 * Schema-aware deserializer for connect4 package
 */
export class Connect4Deserializer {
    constructor(schemaRegistry = Connect4SchemaRegistry, factory = DEFAULT_FACTORY) {
        this.schemaRegistry = schemaRegistry;
        this.factory = factory;
    }
    /**
     * Deserialize an object using schema information
     * @param instance The target instance to populate
     * @param data The source data to deserialize from
     * @param messageType The fully qualified message type (e.g., "library.v1.Book")
     * @returns The populated instance
     */
    deserialize(instance, data, messageType) {
        if (!data || typeof data !== 'object') {
            return instance;
        }
        const schema = this.schemaRegistry[messageType];
        if (!schema) {
            // Fallback to simple property copying if no schema found
            return this.fallbackDeserialize(instance, data);
        }
        // Process each field according to its schema
        for (const fieldSchema of schema.fields) {
            const fieldValue = data[fieldSchema.name];
            if (fieldValue === null || fieldValue === undefined) {
                continue;
            }
            this.deserializeField(instance, fieldSchema, fieldValue);
        }
        return instance;
    }
    /**
     * Deserialize a single field based on its schema
     */
    deserializeField(instance, fieldSchema, fieldValue) {
        const fieldName = fieldSchema.name;
        switch (fieldSchema.type) {
            case FieldType.STRING:
            case FieldType.NUMBER:
            case FieldType.BOOLEAN:
                // Simple primitive types - direct assignment
                instance[fieldName] = fieldValue;
                break;
            case FieldType.MESSAGE:
                if (fieldSchema.repeated) {
                    // Handle repeated message fields (arrays)
                    instance[fieldName] = this.deserializeMessageArray(fieldValue, fieldSchema.messageType, instance, fieldName);
                }
                else {
                    // Handle single message field
                    instance[fieldName] = this.deserializeMessageField(fieldValue, fieldSchema.messageType, instance, fieldName);
                }
                break;
            case FieldType.REPEATED:
                // Handle repeated primitive fields
                if (Array.isArray(fieldValue)) {
                    instance[fieldName] = [...fieldValue]; // Simple copy for primitives
                }
                break;
            case FieldType.ONEOF:
                // Handle oneof fields (would need additional logic for union types)
                instance[fieldName] = fieldValue;
                break;
            case FieldType.MAP:
                // Handle map fields (would need additional schema info for key/value types)
                instance[fieldName] = Object.assign({}, fieldValue);
                break;
            default:
                // Fallback to direct assignment
                instance[fieldName] = fieldValue;
                break;
        }
    }
    /**
     * Deserialize a single message field
     */
    deserializeMessageField(fieldValue, messageType, parent, attributeName) {
        // Try to get factory method using cross-package delegation
        let factoryMethod;
        if (this.factory.getFactoryMethod) {
            factoryMethod = this.factory.getFactoryMethod(messageType);
        }
        else {
            // Fallback to simple method name lookup
            const factoryMethodName = this.getFactoryMethodName(messageType);
            factoryMethod = this.factory[factoryMethodName];
        }
        if (factoryMethod) {
            const result = factoryMethod(parent, attributeName, undefined, fieldValue);
            if (result.fullyLoaded) {
                return result.instance;
            }
            else {
                // Factory created instance but didn't populate - use deserializer
                return this.deserialize(result.instance, fieldValue, messageType);
            }
        }
        // No factory method found - fallback
        return this.fallbackDeserialize({}, fieldValue);
    }
    /**
     * Deserialize an array of message objects
     */
    deserializeMessageArray(fieldValue, messageType, parent, attributeName) {
        if (!Array.isArray(fieldValue)) {
            return [];
        }
        // Try to get factory method using cross-package delegation
        let factoryMethod;
        if (this.factory.getFactoryMethod) {
            factoryMethod = this.factory.getFactoryMethod(messageType);
        }
        else {
            // Fallback to simple method name lookup
            const factoryMethodName = this.getFactoryMethodName(messageType);
            factoryMethod = this.factory[factoryMethodName];
        }
        return fieldValue.map((item, index) => {
            if (factoryMethod) {
                const result = factoryMethod(parent, attributeName, index, item);
                if (result.fullyLoaded) {
                    return result.instance;
                }
                else {
                    // Factory created instance but didn't populate - use deserializer
                    return this.deserialize(result.instance, item, messageType);
                }
            }
            // No factory method found - fallback
            return this.fallbackDeserialize({}, item);
        });
    }
    /**
     * Convert message type to factory method name
     * "library.v1.Book" -> "newBook"
     */
    getFactoryMethodName(messageType) {
        const parts = messageType.split('.');
        const typeName = parts[parts.length - 1]; // Get last part (e.g., "Book")
        return 'new' + typeName;
    }
    /**
     * Fallback deserializer for when no schema is available
     */
    fallbackDeserialize(instance, data) {
        if (!data || typeof data !== 'object') {
            return instance;
        }
        for (const [key, value] of Object.entries(data)) {
            if (value !== null && value !== undefined) {
                instance[key] = value;
            }
        }
        return instance;
    }
    /**
     * Create and deserialize a new instance of a message type
     */
    createAndDeserialize(messageType, data) {
        // Try to get factory method using cross-package delegation
        let factoryMethod;
        if (this.factory.getFactoryMethod) {
            factoryMethod = this.factory.getFactoryMethod(messageType);
        }
        else {
            // Fallback to simple method name lookup
            const factoryMethodName = this.getFactoryMethodName(messageType);
            factoryMethod = this.factory[factoryMethodName];
        }
        if (!factoryMethod) {
            throw new Error(`Could not find factory method to deserialize: ${messageType}`);
        }
        const result = factoryMethod(undefined, undefined, undefined, data);
        if (result.fullyLoaded) {
            return result.instance;
        }
        else {
            return this.deserialize(result.instance, data, messageType);
        }
    }
    /**
     * Static utility method to create and deserialize a message without needing a deserializer instance
     * @param messageType Fully qualified message type (use Class.MESSAGE_TYPE)
     * @param data Raw data to deserialize
     * @returns Deserialized instance or null if creation failed
     */
    static from(messageType, data) {
        const deserializer = new Connect4Deserializer(); // Uses default factory and schema registry
        return deserializer.createAndDeserialize(messageType, data);
    }
}
