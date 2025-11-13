
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for SecondRequest message
 */
export const SecondRequestSchema: MessageSchema = {
  name: "SecondRequest",
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
 * Schema for SecondResponse message
 */
export const SecondResponseSchema: MessageSchema = {
  name: "SecondResponse",
  fields: [
    {
      name: "x",
      type: FieldType.NUMBER,
      id: 1,
    },
  ],
};



/**
 * Package-scoped schema registry for test_one_package.v1
 */
export const test_one_package_v1SchemaRegistry: Record<string, MessageSchema> = {
  "test_one_package.v1.SecondRequest": SecondRequestSchema,
  "test_one_package.v1.SecondResponse": SecondResponseSchema,
};

/**
 * Schema registry instance for test_one_package.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(test_one_package_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);