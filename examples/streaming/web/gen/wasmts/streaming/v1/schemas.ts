
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for TickRequest message
 */
export const TickRequestSchema: MessageSchema = {
  name: "TickRequest",
  fields: [
    {
      name: "count",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "intervalMs",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "message",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for TickResponse message
 */
export const TickResponseSchema: MessageSchema = {
  name: "TickResponse",
  fields: [
    {
      name: "tickNumber",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "timestamp",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "message",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "isFinal",
      type: FieldType.BOOLEAN,
      id: 4,
    },
  ],
};



/**
 * Package-scoped schema registry for streaming.v1
 */
export const StreamingV1SchemaRegistry: Record<string, MessageSchema> = {
  "streaming.v1.TickRequest": TickRequestSchema,
  "streaming.v1.TickResponse": TickResponseSchema,
};

/**
 * Get schema for a message type from streaming.v1 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return StreamingV1SchemaRegistry[messageType];
}

/**
 * Get field schema by name from streaming.v1 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from streaming.v1 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in streaming.v1 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from streaming.v1 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}