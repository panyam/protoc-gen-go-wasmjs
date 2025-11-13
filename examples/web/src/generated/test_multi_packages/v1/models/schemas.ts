
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for SampleRequest message
 */
export const SampleRequestSchema: MessageSchema = {
  name: "SampleRequest",
  fields: [
    {
      name: "a",
      type: FieldType.NUMBER,
      id: 1,
    },
    {
      name: "b",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for SampleResponse message
 */
export const SampleResponseSchema: MessageSchema = {
  name: "SampleResponse",
  fields: [
    {
      name: "x",
      type: FieldType.NUMBER,
      id: 1,
    },
  ],
};



/**
 * Package-scoped schema registry for test_multi_packages.v1.models
 */
export const test_multi_packages_v1_modelsSchemaRegistry: Record<string, MessageSchema> = {
  "test_multi_packages.v1.models.SampleRequest": SampleRequestSchema,
  "test_multi_packages.v1.models.SampleResponse": SampleResponseSchema,
};

/**
 * Schema registry instance for test_multi_packages.v1.models package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(test_multi_packages_v1_modelsSchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);