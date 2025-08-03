// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Enhanced Book with common base
 */
export interface Book {
  base?: BaseMessage;
  title: string;
  author: string;
  isbn: string;
  year: number;
  genre: string;
  available: boolean;
  tags: string[];
  rating: number;
}


/**
 * Enhanced User with common base
 */
export interface User {
  base?: BaseMessage;
  name: string;
  email: string;
  phone: string;
  preferences: string[];
}


/**
 * Enhanced request with metadata
 */
export interface FindBooksRequest {
  metadata?: Metadata;
  query: string;
  genre: string;
  limit: number;
  availableOnly: boolean;
  tags: string[];
  minRating: number;
}


/**
 * Enhanced response with metadata
 */
export interface FindBooksResponse {
  metadata?: Metadata;
  books?: Book[];
  totalCount: number;
  hasMore: boolean;
}

