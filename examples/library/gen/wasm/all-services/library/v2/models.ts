import { BaseMessage } from "../common/interfaces";
import { Metadata } from "../common/interfaces";
import { ErrorInfo } from "../common/interfaces";


import { BookStats as BookStatsInterface, UserActivity as UserActivityInterface, GetBookStatsRequest as GetBookStatsRequestInterface, GetBookStatsResponse as GetBookStatsResponseInterface, Pagination as PaginationInterface, PaginationResponse as PaginationResponseInterface, Book as BookInterface, User as UserInterface, FindBooksRequest as FindBooksRequestInterface, FindBooksResponse as FindBooksResponseInterface } from "./interfaces";
import { LibraryV2Deserializer } from "./deserializer";


/**
 * BookStats provides analytics for books
 */
export class BookStats implements BookStatsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.BookStats";

  base?: BaseMessage;
  bookId: string = "";
  viewCount: number = 0;
  checkoutCount: number = 0;
  averageRating: number = 0;
  reviewCount: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized BookStats instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<BookStats>(BookStats.MESSAGE_TYPE, data);
  }
}


/**
 * UserActivity tracks user behavior
 */
export class UserActivity implements UserActivityInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.UserActivity";

  base?: BaseMessage;
  userId: string = "";
  activityType: string = "";
  bookId: string = "";
  description: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized UserActivity instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<UserActivity>(UserActivity.MESSAGE_TYPE, data);
  }
}


/**
 * GetBookStatsRequest for analytics
 */
export class GetBookStatsRequest implements GetBookStatsRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.GetBookStatsRequest";

  metadata?: Metadata;
  bookId: string = "";
  dateRange: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetBookStatsRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<GetBookStatsRequest>(GetBookStatsRequest.MESSAGE_TYPE, data);
  }
}


/**
 * GetBookStatsResponse for analytics
 */
export class GetBookStatsResponse implements GetBookStatsResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.GetBookStatsResponse";

  metadata?: Metadata;
  stats?: BookStats;
  error?: ErrorInfo;
  fieldErrors?: Map<string, string>;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetBookStatsResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<GetBookStatsResponse>(GetBookStatsResponse.MESSAGE_TYPE, data);
  }
}



export class Pagination implements PaginationInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.Pagination";

  /** *
 Instead of an offset an abstract  "page" key is provided that offers
 an opaque "pointer" into some offset in a result set. */
  pageKey: string = "";
  /** *
 If a pagekey is not supported we can also support a direct integer offset
 for cases where it makes sense. */
  pageOffset: number = 0;
  /** *
 Number of results to return. */
  pageSize: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized Pagination instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<Pagination>(Pagination.MESSAGE_TYPE, data);
  }
}



export class PaginationResponse implements PaginationResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.PaginationResponse";

  /** *
 The key/pointer string that subsequent List requests should pass to
 continue the pagination. */
  nextPageKey: string = "";
  /** *
 Also support an integer offset if possible */
  nextPageOffset: number = 0;
  /** *
 Whether theere are more results. */
  hasMore: boolean = false;
  /** *
 Total number of results. */
  totalResults: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PaginationResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<PaginationResponse>(PaginationResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Enhanced Book with common base
 */
export class Book implements BookInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.Book";

  base?: BaseMessage;
  title: string = "";
  author: string = "";
  isbn: string = "";
  year: number = 0;
  genre: string = "";
  available: boolean = false;
  tags: string[] = [];
  rating: number = 0;
  createdAt?: Date;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized Book instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<Book>(Book.MESSAGE_TYPE, data);
  }
}


/**
 * Enhanced User with common base
 */
export class User implements UserInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.User";

  base?: BaseMessage;
  name: string = "";
  email: string = "";
  phone: string = "";
  preferences: string[] = [];

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized User instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<User>(User.MESSAGE_TYPE, data);
  }
}


/**
 * Enhanced request with metadata
 */
export class FindBooksRequest implements FindBooksRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.FindBooksRequest";

  metadata?: Metadata;
  query: string = "";
  genre: string = "";
  limit: number = 0;
  availableOnly: boolean = false;
  tags: string[] = [];
  minRating: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized FindBooksRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<FindBooksRequest>(FindBooksRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Enhanced response with metadata
 */
export class FindBooksResponse implements FindBooksResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v2.FindBooksResponse";

  metadata?: Metadata;
  books: Book[] = [];
  totalCount: number = 0;
  hasMore: boolean = false;
  pagination?: Pagination;
  fieldErrors?: Map<string, string>;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized FindBooksResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV2Deserializer.from<FindBooksResponse>(FindBooksResponse.MESSAGE_TYPE, data);
  }
}


