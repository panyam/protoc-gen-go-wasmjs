/**
 * Field type enumeration for proto field types
 */
declare enum FieldType {
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
interface FieldSchema {
    name: string;
    type: FieldType;
    id: number;
    messageType?: string;
    repeated?: boolean;
    mapKeyType?: FieldType;
    mapValueType?: FieldType | string;
    oneofGroup?: string;
    optional?: boolean;
}
/**
 * Message schema interface
 */
interface MessageSchema {
    name: string;
    fields: FieldSchema[];
    oneofGroups?: string[];
}

export { type FieldSchema, FieldType, type MessageSchema };
