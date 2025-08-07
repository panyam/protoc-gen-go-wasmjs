
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for BaseMessage message
 */
export const BaseMessageSchema: MessageSchema = {
  name: "BaseMessage",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "timestamp",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "version",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for Metadata message
 */
export const MetadataSchema: MessageSchema = {
  name: "Metadata",
  fields: [
    {
      name: "requestId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "userAgent",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "headers",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "library.common.HeadersEntry",
    },
  ],
};


/**
 * Schema for ErrorInfo message
 */
export const ErrorInfoSchema: MessageSchema = {
  name: "ErrorInfo",
  fields: [
    {
      name: "code",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "message",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "details",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
  ],
};



/**
 * Package-scoped schema registry for library.common
 */
export const LibraryCommonSchemaRegistry: Record<string, MessageSchema> = {
  "library.common.BaseMessage": BaseMessageSchema,
  "library.common.Metadata": MetadataSchema,
  "library.common.ErrorInfo": ErrorInfoSchema,
};

/**
 * Get schema for a message type from library.common package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return LibraryCommonSchemaRegistry[messageType];
}

/**
 * Get field schema by name from library.common package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from library.common package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in library.common package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from library.common package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}