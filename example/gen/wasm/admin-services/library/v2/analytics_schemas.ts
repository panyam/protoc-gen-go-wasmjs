
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

/**
 * Field type enumeration for proto field types
 */
export enum FieldType {
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
export interface FieldSchema {
  name: string;
  type: FieldType;
  id: number; // Proto field number (e.g., text_query = 1)
  messageType?: string; // For MESSAGE type fields
  repeated?: boolean; // For array fields
  mapKeyType?: FieldType; // For MAP type fields
  mapValueType?: FieldType | string; // For MAP type fields
  oneofGroup?: string; // For ONEOF fields
  optional?: boolean;
}

/**
 * Message schema interface
 */
export interface MessageSchema {
  name: string;
  fields: FieldSchema[];
  oneofGroups?: string[]; // List of oneof group names
}


/**
 * Schema for BookStats message
 */
export const BookStatsSchema: MessageSchema = {
  name: "BookStats",
  fields: [
    {
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage",
    },
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "viewCount",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "checkoutCount",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "averageRating",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "reviewCount",
      type: FieldType.NUMBER,
      id: 6,
    },
  ],
};


/**
 * Schema for UserActivity message
 */
export const UserActivitySchema: MessageSchema = {
  name: "UserActivity",
  fields: [
    {
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage",
    },
    {
      name: "userId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "activityType",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for GetBookStatsRequest message
 */
export const GetBookStatsRequestSchema: MessageSchema = {
  name: "GetBookStatsRequest",
  fields: [
    {
      name: "metadata",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.Metadata",
    },
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "dateRange",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for GetBookStatsResponse message
 */
export const GetBookStatsResponseSchema: MessageSchema = {
  name: "GetBookStatsResponse",
  fields: [
    {
      name: "metadata",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.Metadata",
    },
    {
      name: "stats",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "library.v2.BookStats",
    },
    {
      name: "error",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "library.common.ErrorInfo",
    },
  ],
};



/**
 * Package-scoped schema registry for library.v2
 */
export const LibraryV2SchemaRegistry: Record<string, MessageSchema> = {
  "library.v2.BookStats": BookStatsSchema,
  "library.v2.UserActivity": UserActivitySchema,
  "library.v2.GetBookStatsRequest": GetBookStatsRequestSchema,
  "library.v2.GetBookStatsResponse": GetBookStatsResponseSchema,
};

/**
 * Get schema for a message type from library.v2 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return LibraryV2SchemaRegistry[messageType];
}

/**
 * Get field schema by name from library.v2 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from library.v2 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in library.v2 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from library.v2 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}