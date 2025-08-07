
// Generated TypeScript schemas from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldType, FieldSchema, MessageSchema } from "./deserializer_schemas";


/**
 * Schema for BookStats message
 */
export const BookStatsSchema: MessageSchema = {
  name: "BookStats",
  fields: [
    {
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage",
    },
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "viewCount",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "checkoutCount",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "averageRating",
      type: FieldType.NUMBER,
      id: 5,
    },
    {
      name: "reviewCount",
      type: FieldType.NUMBER,
      id: 6,
    },
  ],
};


/**
 * Schema for UserActivity message
 */
export const UserActivitySchema: MessageSchema = {
  name: "UserActivity",
  fields: [
    {
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage",
    },
    {
      name: "userId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "activityType",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "description",
      type: FieldType.STRING,
      id: 5,
    },
  ],
};


/**
 * Schema for GetBookStatsRequest message
 */
export const GetBookStatsRequestSchema: MessageSchema = {
  name: "GetBookStatsRequest",
  fields: [
    {
      name: "metadata",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.Metadata",
    },
    {
      name: "bookId",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "dateRange",
      type: FieldType.STRING,
      id: 3,
    },
  ],
};


/**
 * Schema for GetBookStatsResponse message
 */
export const GetBookStatsResponseSchema: MessageSchema = {
  name: "GetBookStatsResponse",
  fields: [
    {
      name: "metadata",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.Metadata",
    },
    {
      name: "stats",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "library.v2.BookStats",
    },
    {
      name: "error",
      type: FieldType.MESSAGE,
      id: 3,
      messageType: "library.common.ErrorInfo",
    },
    {
      name: "fieldErrors",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "library.v2.FieldErrorsEntry",
    },
  ],
};


/**
 * Schema for Pagination message
 */
export const PaginationSchema: MessageSchema = {
  name: "Pagination",
  fields: [
    {
      name: "pageKey",
      type: FieldType.STRING,
      id: 1,
    },
    {
      name: "pageOffset",
      type: FieldType.NUMBER,
      id: 2,
    },
    {
      name: "pageSize",
      type: FieldType.NUMBER,
      id: 3,
    },
  ],
};


/**
 * Schema for PaginationResponse message
 */
export const PaginationResponseSchema: MessageSchema = {
  name: "PaginationResponse",
  fields: [
    {
      name: "nextPageKey",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "nextPageOffset",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "hasMore",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "totalResults",
      type: FieldType.NUMBER,
      id: 5,
    },
  ],
};


/**
 * Schema for Book message
 */
export const BookSchema: MessageSchema = {
  name: "Book",
  fields: [
    {
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage",
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
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 8,
      repeated: true,
    },
    {
      name: "rating",
      type: FieldType.NUMBER,
      id: 9,
    },
    {
      name: "createdAt",
      type: FieldType.MESSAGE,
      id: 10,
      messageType: "google.protobuf.Timestamp",
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
      name: "base",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.BaseMessage",
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
    {
      name: "phone",
      type: FieldType.STRING,
      id: 4,
    },
    {
      name: "preferences",
      type: FieldType.REPEATED,
      id: 5,
      repeated: true,
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
      name: "metadata",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.Metadata",
    },
    {
      name: "query",
      type: FieldType.STRING,
      id: 2,
    },
    {
      name: "genre",
      type: FieldType.STRING,
      id: 3,
    },
    {
      name: "limit",
      type: FieldType.NUMBER,
      id: 4,
    },
    {
      name: "availableOnly",
      type: FieldType.BOOLEAN,
      id: 5,
    },
    {
      name: "tags",
      type: FieldType.REPEATED,
      id: 6,
      repeated: true,
    },
    {
      name: "minRating",
      type: FieldType.NUMBER,
      id: 7,
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
      name: "metadata",
      type: FieldType.MESSAGE,
      id: 1,
      messageType: "library.common.Metadata",
    },
    {
      name: "books",
      type: FieldType.MESSAGE,
      id: 2,
      messageType: "library.v2.Book",
      repeated: true,
    },
    {
      name: "totalCount",
      type: FieldType.NUMBER,
      id: 3,
    },
    {
      name: "hasMore",
      type: FieldType.BOOLEAN,
      id: 4,
    },
    {
      name: "pagination",
      type: FieldType.MESSAGE,
      id: 5,
      messageType: "library.v2.Pagination",
    },
    {
      name: "fieldErrors",
      type: FieldType.MESSAGE,
      id: 6,
      messageType: "library.v2.FieldErrorsEntry",
    },
  ],
};



/**
 * Package-scoped schema registry for library.v2
 */
export const LibraryV2SchemaRegistry: Record<string, MessageSchema> = {
  "library.v2.BookStats": BookStatsSchema,
  "library.v2.UserActivity": UserActivitySchema,
  "library.v2.GetBookStatsRequest": GetBookStatsRequestSchema,
  "library.v2.GetBookStatsResponse": GetBookStatsResponseSchema,
  "library.v2.Pagination": PaginationSchema,
  "library.v2.PaginationResponse": PaginationResponseSchema,
  "library.v2.Book": BookSchema,
  "library.v2.User": UserSchema,
  "library.v2.FindBooksRequest": FindBooksRequestSchema,
  "library.v2.FindBooksResponse": FindBooksResponseSchema,
};

/**
 * Get schema for a message type from library.v2 package
 */
export function getSchema(messageType: string): MessageSchema | undefined {
  return LibraryV2SchemaRegistry[messageType];
}

/**
 * Get field schema by name from library.v2 package
 */
export function getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.name === fieldName);
}

/**
 * Get field schema by proto field ID from library.v2 package
 */
export function getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
  const schema = getSchema(messageType);
  return schema?.fields.find(field => field.id === fieldId);
}

/**
 * Check if field is part of a oneof group in library.v2 package
 */
export function isOneofField(messageType: string, fieldName: string): boolean {
  const fieldSchema = getFieldSchema(messageType, fieldName);
  return fieldSchema?.oneofGroup !== undefined;
}

/**
 * Get all fields in a oneof group from library.v2 package
 */
export function getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
  const schema = getSchema(messageType);
  return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
}