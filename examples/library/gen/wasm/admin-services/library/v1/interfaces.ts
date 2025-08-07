// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Book represents a book in the library
 */
export interface Book {
  id: string;
  title: string;
  author: string;
  isbn: string;
  year: number;
  genre: string;
  available: boolean;
}


/**
 * User represents a library user
 */
export interface User {
  id: string;
  name: string;
  email: string;
}


/**
 * FindBooksRequest is the request for finding books
 */
export interface FindBooksRequest {
  /** Search query (title, author, or ISBN) */
  query: string;
  /** Genre filter */
  genre: string;
  /** Maximum number of results */
  limit: number;
  /** Only show available books */
  availableOnly: boolean;
}


/**
 * FindBooksResponse is the response for finding books
 */
export interface FindBooksResponse {
  books?: Book[];
  totalCount: number;
}


/**
 * CheckoutBookRequest is the request for checking out a book
 */
export interface CheckoutBookRequest {
  bookId: string;
  userId: string;
}


/**
 * CheckoutBookResponse is the response for checking out a book
 */
export interface CheckoutBookResponse {
  success: boolean;
  message: string;
  dueDate: string;
}


/**
 * ReturnBookRequest is the request for returning a book
 */
export interface ReturnBookRequest {
  bookId: string;
  userId: string;
}


/**
 * ReturnBookResponse is the response for returning a book
 */
export interface ReturnBookResponse {
  success: boolean;
  message: string;
}


/**
 * GetUserBooksRequest is the request for getting user's checked out books
 */
export interface GetUserBooksRequest {
  userId: string;
}


/**
 * GetUserBooksResponse is the response for getting user's checked out books
 */
export interface GetUserBooksResponse {
  books?: Book[];
}


/**
 * GetUserRequest is the request for getting user info
 */
export interface GetUserRequest {
  userId: string;
}


/**
 * GetUserResponse is the response for getting user info
 */
export interface GetUserResponse {
  user?: User;
}


/**
 * CreateUserRequest is the request for creating a user
 */
export interface CreateUserRequest {
  name: string;
  email: string;
}


/**
 * CreateUserResponse is the response for creating a user
 */
export interface CreateUserResponse {
  user?: User;
}

