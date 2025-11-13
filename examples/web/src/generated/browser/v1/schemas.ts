
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


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
      type: FieldType.STRING,
      id: 3,
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
      type: FieldType.STRING,
      id: 3,
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
 * Schema for PromptRequest message
 */
export const PromptRequestSchema: MessageSchema = {
  name: "PromptRequest",
  fields: [
    {
      name: "message",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "defaultValue",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for PromptResponse message
 */
export const PromptResponseSchema: MessageSchema = {
  name: "PromptResponse",
  fields: [
    {
      name: "value",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "cancelled",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};


/**
 * Schema for LogRequest message
 */
export const LogRequestSchema: MessageSchema = {
  name: "LogRequest",
  fields: [
    {
      name: "message",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "level",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for LogResponse message
 */
export const LogResponseSchema: MessageSchema = {
  name: "LogResponse",
  fields: [
    {
      name: "logged",
      type: FieldType.BOOLEAN,
      id: 1,
    },
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
 * Schema registry instance for browser.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(browser_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);