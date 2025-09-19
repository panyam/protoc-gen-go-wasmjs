
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for LoadUserRequest message
 */
export const LoadUserRequestSchema: MessageSchema = {
  name: "LoadUserRequest",
  fields: [
  ],
};


/**
 * Schema for LoadUserResponse message
 */
export const LoadUserResponseSchema: MessageSchema = {
  name: "LoadUserResponse",
  fields: [
  ],
};


/**
 * Schema for StateUpdateRequest message
 */
export const StateUpdateRequestSchema: MessageSchema = {
  name: "StateUpdateRequest",
  fields: [
  ],
};


/**
 * Schema for UIUpdate message
 */
export const UIUpdateSchema: MessageSchema = {
  name: "UIUpdate",
  fields: [
  ],
};


/**
 * Schema for PreferencesRequest message
 */
export const PreferencesRequestSchema: MessageSchema = {
  name: "PreferencesRequest",
  fields: [
  ],
};


/**
 * Schema for PreferencesResponse message
 */
export const PreferencesResponseSchema: MessageSchema = {
  name: "PreferencesResponse",
  fields: [
  ],
};


/**
 * Schema for CallbackDemoRequest message
 */
export const CallbackDemoRequestSchema: MessageSchema = {
  name: "CallbackDemoRequest",
  fields: [
  ],
};


/**
 * Schema for CallbackDemoResponse message
 */
export const CallbackDemoResponseSchema: MessageSchema = {
  name: "CallbackDemoResponse",
  fields: [
  ],
};



/**
 * Package-scoped schema registry for presenter.v1
 */
export const presenter_v1SchemaRegistry: Record<string, MessageSchema> = {
  "presenter.v1.LoadUserRequest": LoadUserRequestSchema,
  "presenter.v1.LoadUserResponse": LoadUserResponseSchema,
  "presenter.v1.StateUpdateRequest": StateUpdateRequestSchema,
  "presenter.v1.UIUpdate": UIUpdateSchema,
  "presenter.v1.PreferencesRequest": PreferencesRequestSchema,
  "presenter.v1.PreferencesResponse": PreferencesResponseSchema,
  "presenter.v1.CallbackDemoRequest": CallbackDemoRequestSchema,
  "presenter.v1.CallbackDemoResponse": CallbackDemoResponseSchema,
};

/**
 * Get schema for a message type from presenter.v1 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return presenter_v1SchemaRegistry[messageType];
}

/**
 * Get field schema by name from presenter.v1 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from presenter.v1 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in presenter.v1 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from presenter.v1 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}