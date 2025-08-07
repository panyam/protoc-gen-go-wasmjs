
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for Book message
 */
export const BookSchema: MessageSchema = {
  name: "Book",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "title",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "author",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "isbn",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "year",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "genre",
      type: FieldType.STRING,
      id: 6,
    },
    {
      name: "available",
      type: FieldType.BOOLEAN,
      id: 7,
    },
  ],
};


/**
 * Schema for User message
 */
export const UserSchema: MessageSchema = {
  name: "User",
  fields: [
    {
      name: "id",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "name",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "email",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for FindBooksRequest message
 */
export const FindBooksRequestSchema: MessageSchema = {
  name: "FindBooksRequest",
  fields: [
    {
      name: "query",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "genre",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "limit",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "availableOnly",
      type: FieldType.BOOLEAN,
      id: 4,
    },
  ],
};


/**
 * Schema for FindBooksResponse message
 */
export const FindBooksResponseSchema: MessageSchema = {
  name: "FindBooksResponse",
  fields: [
    {
      name: "books",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.v1.Book",
      repeated: true,
    },
    {
      name: "totalCount",
      type: FieldType.NUMBER,
      id: 2,
    },
  ],
};


/**
 * Schema for CheckoutBookRequest message
 */
export const CheckoutBookRequestSchema: MessageSchema = {
  name: "CheckoutBookRequest",
  fields: [
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "userId",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for CheckoutBookResponse message
 */
export const CheckoutBookResponseSchema: MessageSchema = {
  name: "CheckoutBookResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "message",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "dueDate",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for ReturnBookRequest message
 */
export const ReturnBookRequestSchema: MessageSchema = {
  name: "ReturnBookRequest",
  fields: [
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "userId",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for ReturnBookResponse message
 */
export const ReturnBookResponseSchema: MessageSchema = {
  name: "ReturnBookResponse",
  fields: [
    {
      name: "success",
      type: FieldType.BOOLEAN,
      id: 1,
    },
    {
      name: "message",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for GetUserBooksRequest message
 */
export const GetUserBooksRequestSchema: MessageSchema = {
  name: "GetUserBooksRequest",
  fields: [
    {
      name: "userId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for GetUserBooksResponse message
 */
export const GetUserBooksResponseSchema: MessageSchema = {
  name: "GetUserBooksResponse",
  fields: [
    {
      name: "books",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.v1.Book",
      repeated: true,
    },
  ],
};


/**
 * Schema for GetUserRequest message
 */
export const GetUserRequestSchema: MessageSchema = {
  name: "GetUserRequest",
  fields: [
    {
      name: "userId",
      type: FieldType.STRING,
      id: 1,
    },
  ],
};


/**
 * Schema for GetUserResponse message
 */
export const GetUserResponseSchema: MessageSchema = {
  name: "GetUserResponse",
  fields: [
    {
      name: "user",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.v1.User",
    },
  ],
};


/**
 * Schema for CreateUserRequest message
 */
export const CreateUserRequestSchema: MessageSchema = {
  name: "CreateUserRequest",
  fields: [
    {
      name: "name",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "email",
      type: FieldType.STRING,
      id: 2,
    },
  ],
};


/**
 * Schema for CreateUserResponse message
 */
export const CreateUserResponseSchema: MessageSchema = {
  name: "CreateUserResponse",
  fields: [
    {
      name: "user",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.v1.User",
    },
  ],
};



/**
 * Package-scoped schema registry for library.v1
 */
export const LibraryV1SchemaRegistry: Record<string, MessageSchema> = {
  "library.v1.Book": BookSchema,
  "library.v1.User": UserSchema,
  "library.v1.FindBooksRequest": FindBooksRequestSchema,
  "library.v1.FindBooksResponse": FindBooksResponseSchema,
  "library.v1.CheckoutBookRequest": CheckoutBookRequestSchema,
  "library.v1.CheckoutBookResponse": CheckoutBookResponseSchema,
  "library.v1.ReturnBookRequest": ReturnBookRequestSchema,
  "library.v1.ReturnBookResponse": ReturnBookResponseSchema,
  "library.v1.GetUserBooksRequest": GetUserBooksRequestSchema,
  "library.v1.GetUserBooksResponse": GetUserBooksResponseSchema,
  "library.v1.GetUserRequest": GetUserRequestSchema,
  "library.v1.GetUserResponse": GetUserResponseSchema,
  "library.v1.CreateUserRequest": CreateUserRequestSchema,
  "library.v1.CreateUserResponse": CreateUserResponseSchema,
};

/**
 * Get schema for a message type from library.v1 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return LibraryV1SchemaRegistry[messageType];
}

/**
 * Get field schema by name from library.v1 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from library.v1 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in library.v1 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from library.v1 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}