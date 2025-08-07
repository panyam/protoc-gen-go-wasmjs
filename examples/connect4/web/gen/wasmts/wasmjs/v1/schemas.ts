
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


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
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "wasmjs.v1.MetadataEntry",
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
export const WasmjsV1SchemaRegistry: Record<string, MessageSchema> = {
  "wasmjs.v1.StatefulOptions": StatefulOptionsSchema,
  "wasmjs.v1.StatefulMethodOptions": StatefulMethodOptionsSchema,
  "wasmjs.v1.MessagePatch": MessagePatchSchema,
  "wasmjs.v1.PatchBatch": PatchBatchSchema,
  "wasmjs.v1.PatchResponse": PatchResponseSchema,
};

/**
 * Get schema for a message type from wasmjs.v1 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return WasmjsV1SchemaRegistry[messageType];
}

/**
 * Get field schema by name from wasmjs.v1 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from wasmjs.v1 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in wasmjs.v1 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from wasmjs.v1 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}