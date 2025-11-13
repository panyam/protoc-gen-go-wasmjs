
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for StatefulOptions message
 */
export const StatefulOptionsSchema: MessageSchema = {
  name: "StatefulOptions",
  fields: [
    {
      name: "enabled",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "stateMessageType",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "conflictResolution",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for StatefulMethodOptions message
 */
export const StatefulMethodOptionsSchema: MessageSchema = {
  name: "StatefulMethodOptions",
  fields: [
    {
      name: "returnsPatches",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "broadcasts",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};


/**
 * Schema for AsyncMethodOptions message
 */
export const AsyncMethodOptionsSchema: MessageSchema = {
  name: "AsyncMethodOptions",
  fields: [
    {
      name: "isAsync",
      type: FieldType.BOOLEAN,
      id: 1,
    },
  ],
};


/**
 * Schema for MessagePatch message
 */
export const MessagePatchSchema: MessageSchema = {
  name: "MessagePatch",
  fields: [
    {
      name: "operation",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "fieldPath",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "valueJson",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "index",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "key",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "oldIndex",
      type: FieldType.NUMBER,
      id: 6,
    },
    {
      name: "changeNumber",
      type: FieldType.NUMBER,
      id: 7,
    },
    {
      name: "timestamp",
      type: FieldType.NUMBER,
      id: 8,
    },
    {
      name: "userId",
      type: FieldType.STRING,
      id: 9,
    },
    {
      name: "transactionId",
      type: FieldType.STRING,
      id: 10,
    },
  ],
};


/**
 * Schema for PatchBatch message
 */
export const PatchBatchSchema: MessageSchema = {
  name: "PatchBatch",
  fields: [
    {
      name: "messageType",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "entityId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "patches",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "wasmjs.v1.MessagePatch",
      repeated: true,
    },
    {
      name: "changeNumber",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "source",
      type: FieldType.STRING,
      id: 5,
    },
    {
      name: "metadata",
      type: FieldType.STRING,
      id: 6,
    },
  ],
};


/**
 * Schema for PatchResponse message
 */
export const PatchResponseSchema: MessageSchema = {
  name: "PatchResponse",
  fields: [
    {
      name: "patchBatches",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "wasmjs.v1.PatchBatch",
      repeated: true,
    },
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 2,
    },
    {
      name: "errorMessage",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "newChangeNumber",
      type: FieldType.NUMBER,
      id: 4,
    },
  ],
};



/**
 * Package-scoped schema registry for wasmjs.v1
 */
export const wasmjs_v1SchemaRegistry: Record<string, MessageSchema> = {
  "wasmjs.v1.StatefulOptions": StatefulOptionsSchema,
  "wasmjs.v1.StatefulMethodOptions": StatefulMethodOptionsSchema,
  "wasmjs.v1.AsyncMethodOptions": AsyncMethodOptionsSchema,
  "wasmjs.v1.MessagePatch": MessagePatchSchema,
  "wasmjs.v1.PatchBatch": PatchBatchSchema,
  "wasmjs.v1.PatchResponse": PatchResponseSchema,
};

/**
 * Schema registry instance for wasmjs.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(wasmjs_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);