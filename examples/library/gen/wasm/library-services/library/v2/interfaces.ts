// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { Metadata } from "../common/interfaces";
import { ErrorInfo } from "../common/interfaces";
import { BaseMessage } from "../common/interfaces";



/**
 * BookStats provides analytics for books
 */
export interface BookStats {
  base?: BaseMessage;
  bookId: string;
  viewCount: number;
  checkoutCount: number;
  averageRating: number;
  reviewCount: number;
}


/**
 * UserActivity tracks user behavior
 */
export interface UserActivity {
  base?: BaseMessage;
  userId: string;
  activityType: string;
  bookId: string;
  description: string;
}


/**
 * GetBookStatsRequest for analytics
 */
export interface GetBookStatsRequest {
  metadata?: Metadata;
  bookId: string;
  dateRange: string;
}


/**
 * GetBookStatsResponse for analytics
 */
export interface GetBookStatsResponse {
  metadata?: Metadata;
  stats?: BookStats;
  error?: ErrorInfo;
  fieldErrors?: Map<string, string>;
}



export interface Pagination {
  /** *
 Instead of an offset an abstract  "page" key is provided that offers
 an opaque "pointer" into some offset in a result set. */
  pageKey: string;
  /** *
 If a pagekey is not supported we can also support a direct integer offset
 for cases where it makes sense. */
  pageOffset: number;
  /** *
 Number of results to return. */
  pageSize: number;
}



export interface PaginationResponse {
  /** *
 The key/pointer string that subsequent List requests should pass to
 continue the pagination. */
  nextPageKey: string;
  /** *
 Also support an integer offset if possible */
  nextPageOffset: number;
  /** *
 Whether theere are more results. */
  hasMore: boolean;
  /** *
 Total number of results. */
  totalResults: number;
}


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
  createdAt?: Date;
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
  pagination?: Pagination;
  fieldErrors?: Map<string, string>;
}

