
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for FetchRequest message
 */
export const FetchRequestSchema: MessageSchema = {
  name: "FetchRequest",
  fields: [
  ],
};


/**
 * Schema for FetchResponse message
 */
export const FetchResponseSchema: MessageSchema = {
  name: "FetchResponse",
  fields: [
  ],
};


/**
 * Schema for StorageKeyRequest message
 */
export const StorageKeyRequestSchema: MessageSchema = {
  name: "StorageKeyRequest",
  fields: [
  ],
};


/**
 * Schema for StorageValueResponse message
 */
export const StorageValueResponseSchema: MessageSchema = {
  name: "StorageValueResponse",
  fields: [
  ],
};


/**
 * Schema for StorageSetRequest message
 */
export const StorageSetRequestSchema: MessageSchema = {
  name: "StorageSetRequest",
  fields: [
  ],
};


/**
 * Schema for StorageSetResponse message
 */
export const StorageSetResponseSchema: MessageSchema = {
  name: "StorageSetResponse",
  fields: [
  ],
};


/**
 * Schema for CookieRequest message
 */
export const CookieRequestSchema: MessageSchema = {
  name: "CookieRequest",
  fields: [
  ],
};


/**
 * Schema for CookieResponse message
 */
export const CookieResponseSchema: MessageSchema = {
  name: "CookieResponse",
  fields: [
  ],
};


/**
 * Schema for AlertRequest message
 */
export const AlertRequestSchema: MessageSchema = {
  name: "AlertRequest",
  fields: [
  ],
};


/**
 * Schema for AlertResponse message
 */
export const AlertResponseSchema: MessageSchema = {
  name: "AlertResponse",
  fields: [
  ],
};


/**
 * Schema for PromptRequest message
 */
export const PromptRequestSchema: MessageSchema = {
  name: "PromptRequest",
  fields: [
  ],
};


/**
 * Schema for PromptResponse message
 */
export const PromptResponseSchema: MessageSchema = {
  name: "PromptResponse",
  fields: [
  ],
};


/**
 * Schema for LogRequest message
 */
export const LogRequestSchema: MessageSchema = {
  name: "LogRequest",
  fields: [
  ],
};


/**
 * Schema for LogResponse message
 */
export const LogResponseSchema: MessageSchema = {
  name: "LogResponse",
  fields: [
  ],
};



/**
 * Package-scoped schema registry for browser.v1
 */
export const browser_v1SchemaRegistry: Record<string, MessageSchema> = {
  "browser.v1.FetchRequest": FetchRequestSchema,
  "browser.v1.FetchResponse": FetchResponseSchema,
  "browser.v1.StorageKeyRequest": StorageKeyRequestSchema,
  "browser.v1.StorageValueResponse": StorageValueResponseSchema,
  "browser.v1.StorageSetRequest": StorageSetRequestSchema,
  "browser.v1.StorageSetResponse": StorageSetResponseSchema,
  "browser.v1.CookieRequest": CookieRequestSchema,
  "browser.v1.CookieResponse": CookieResponseSchema,
  "browser.v1.AlertRequest": AlertRequestSchema,
  "browser.v1.AlertResponse": AlertResponseSchema,
  "browser.v1.PromptRequest": PromptRequestSchema,
  "browser.v1.PromptResponse": PromptResponseSchema,
  "browser.v1.LogRequest": LogRequestSchema,
  "browser.v1.LogResponse": LogResponseSchema,
};

/**
 * Get schema for a message type from browser.v1 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return browser_v1SchemaRegistry[messageType];
}

/**
 * Get field schema by name from browser.v1 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from browser.v1 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in browser.v1 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from browser.v1 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}