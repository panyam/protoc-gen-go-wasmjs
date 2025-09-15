
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for FetchRequest message
 */
export const FetchRequestSchema: MessageSchema = {
  name: "FetchRequest",
  fields: [
    {
      name: "url",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "method",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "headers",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "browser.v1.HeadersEntry",
    },
    {
      name: "body",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for FetchResponse message
 */
export const FetchResponseSchema: MessageSchema = {
  name: "FetchResponse",
  fields: [
    {
      name: "status",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "statusText",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "headers",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "browser.v1.HeadersEntry",
    },
    {
      name: "body",
      type: FieldType.STRING,
      id: 4,
    },
  ],
};


/**
 * Schema for StorageKeyRequest message
 */
export const StorageKeyRequestSchema: MessageSchema = {
  name: "StorageKeyRequest",
  fields: [
    {
      name: "key",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for StorageValueResponse message
 */
export const StorageValueResponseSchema: MessageSchema = {
  name: "StorageValueResponse",
  fields: [
    {
      name: "value",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "exists",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};


/**
 * Schema for StorageSetRequest message
 */
export const StorageSetRequestSchema: MessageSchema = {
  name: "StorageSetRequest",
  fields: [
    {
      name: "key",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "value",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for StorageSetResponse message
 */
export const StorageSetResponseSchema: MessageSchema = {
  name: "StorageSetResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
  ],
};


/**
 * Schema for CookieRequest message
 */
export const CookieRequestSchema: MessageSchema = {
  name: "CookieRequest",
  fields: [
    {
      name: "name",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for CookieResponse message
 */
export const CookieResponseSchema: MessageSchema = {
  name: "CookieResponse",
  fields: [
    {
      name: "value",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "exists",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};


/**
 * Schema for AlertRequest message
 */
export const AlertRequestSchema: MessageSchema = {
  name: "AlertRequest",
  fields: [
    {
      name: "message",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for AlertResponse message
 */
export const AlertResponseSchema: MessageSchema = {
  name: "AlertResponse",
  fields: [
    {
      name: "shown",
      type: FieldType.BOOLEAN,
      id: 1,
    },
  ],
};



/**
 * Package-scoped schema registry for browser.v1
 */
export const BrowserV1SchemaRegistry: Record<string, MessageSchema> = {
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
};

/**
 * Get schema for a message type from browser.v1 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return BrowserV1SchemaRegistry[messageType];
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