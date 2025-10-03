
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema, BaseSchemaRegistry } from "@protoc-gen-go-wasmjs/runtime";


/**
 * Schema for LoadUserDataRequest message
 */
export const LoadUserDataRequestSchema: MessageSchema = {
  name: "LoadUserDataRequest",
  fields: [
    {
      name: "userId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for LoadUserDataResponse message
 */
export const LoadUserDataResponseSchema: MessageSchema = {
  name: "LoadUserDataResponse",
  fields: [
    {
      name: "username",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "email",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "permissions",
      type: FieldType.REPEATED,
      id: 3,
      repeated: true,
    },
    {
      name: "fromCache",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "google.protobuf.Timestamp",
    },
  ],
};


/**
 * Schema for StateUpdateRequest message
 */
export const StateUpdateRequestSchema: MessageSchema = {
  name: "StateUpdateRequest",
  fields: [
    {
      name: "action",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "params",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "updateMask",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "google.protobuf.FieldMask",
    },
  ],
};


/**
 * Schema for UIUpdate message
 */
export const UIUpdateSchema: MessageSchema = {
  name: "UIUpdate",
  fields: [
    {
      name: "component",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "action",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "data",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for TestRecord message
 */
export const TestRecordSchema: MessageSchema = {
  name: "TestRecord",
  fields: [
    {
      name: "helperRecord",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "utils.v1.HelperUtilType",
    },
    {
      name: "nestedHelper",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "utils.v1.ParentUtilMessage.NestedUtilType",
    },
  ],
};


/**
 * Schema for PreferencesRequest message
 */
export const PreferencesRequestSchema: MessageSchema = {
  name: "PreferencesRequest",
  fields: [
    {
      name: "preferences",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for PreferencesResponse message
 */
export const PreferencesResponseSchema: MessageSchema = {
  name: "PreferencesResponse",
  fields: [
    {
      name: "saved",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "itemsSaved",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for CallbackDemoRequest message
 */
export const CallbackDemoRequestSchema: MessageSchema = {
  name: "CallbackDemoRequest",
  fields: [
    {
      name: "demoName",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for CallbackDemoResponse message
 */
export const CallbackDemoResponseSchema: MessageSchema = {
  name: "CallbackDemoResponse",
  fields: [
    {
      name: "collectedInputs",
      type: FieldType.REPEATED,
      id: 1,
      repeated: true,
    },
    {
      name: "completed",
      type: FieldType.BOOLEAN,
      id: 2,
    },
  ],
};



/**
 * Package-scoped schema registry for presenter.v1
 */
export const presenter_v1SchemaRegistry: Record<string, MessageSchema> = {
  "presenter.v1.LoadUserDataRequest": LoadUserDataRequestSchema,
  "presenter.v1.LoadUserDataResponse": LoadUserDataResponseSchema,
  "presenter.v1.StateUpdateRequest": StateUpdateRequestSchema,
  "presenter.v1.UIUpdate": UIUpdateSchema,
  "presenter.v1.TestRecord": TestRecordSchema,
  "presenter.v1.PreferencesRequest": PreferencesRequestSchema,
  "presenter.v1.PreferencesResponse": PreferencesResponseSchema,
  "presenter.v1.CallbackDemoRequest": CallbackDemoRequestSchema,
  "presenter.v1.CallbackDemoResponse": CallbackDemoResponseSchema,
};

/**
 * Schema registry instance for presenter.v1 package with utility methods
 * Extends BaseSchemaRegistry with package-specific schema data
 */
// Schema utility functions (now inherited from BaseSchemaRegistry in runtime package)
// Creating instance with package-specific schema registry
const registryInstance = new BaseSchemaRegistry(presenter_v1SchemaRegistry);

export const getSchema = registryInstance.getSchema.bind(registryInstance);
export const getFieldSchema = registryInstance.getFieldSchema.bind(registryInstance);
export const getFieldSchemaById = registryInstance.getFieldSchemaById.bind(registryInstance);
export const isOneofField = registryInstance.isOneofField.bind(registryInstance);
export const getOneofFields = registryInstance.getOneofFields.bind(registryInstance);