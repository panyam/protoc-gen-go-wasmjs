
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for ModelWithOptionalAndDefaults message
 */
export const ModelWithOptionalAndDefaultsSchema: MessageSchema = {
  name: "ModelWithOptionalAndDefaults",
  fields: [
    {
      name: "neededValue",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "optionalString",
      type: FieldType.STRING,
      id: 2,
      oneofGroup: "_optional_string",
      optional: true,
    },
  ],
  oneofGroups: ["_optional_string"],
};


/**
 * Schema for HelperUtilType message
 */
export const HelperUtilTypeSchema: MessageSchema = {
  name: "HelperUtilType",
  fields: [
    {
      name: "value1",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "updateMask",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "google.protobuf.FieldMask",
    },
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.Timestamp",
    },
  ],
};


/**
 * Schema for NestedUtilType message
 */
export const NestedUtilTypeSchema: MessageSchema = {
  name: "NestedUtilType",
  fields: [
    {
      name: "topLevelCount",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "topLevelValue",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for ParentUtilMessage message
 */
export const ParentUtilMessageSchema: MessageSchema = {
  name: "ParentUtilMessage",
  fields: [
    {
      name: "parentValue",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "nested",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "utils.v1.ParentUtilMessage.NestedUtilType",
    },
  ],
};


/**
 * Schema for ParentUtilMessage_NestedUtilType message
 */
export const ParentUtilMessage_NestedUtilTypeSchema: MessageSchema = {
  name: "ParentUtilMessage_NestedUtilType",
  fields: [
    {
      name: "nestedValue",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "nestedCount",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};



/**
 * Package-scoped schema registry for utils.v1
 */
export const utils_v1SchemaRegistry: Record<string, MessageSchema> = {
  "utils.v1.ModelWithOptionalAndDefaults": ModelWithOptionalAndDefaultsSchema,
  "utils.v1.HelperUtilType": HelperUtilTypeSchema,
  "utils.v1.NestedUtilType": NestedUtilTypeSchema,
  "utils.v1.ParentUtilMessage": ParentUtilMessageSchema,
  "utils.v1.ParentUtilMessage_NestedUtilType": ParentUtilMessage_NestedUtilTypeSchema,
};

/**
 * Schema registry instance for utils.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(utils_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);