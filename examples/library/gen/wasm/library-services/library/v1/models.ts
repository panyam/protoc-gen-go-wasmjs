import { Book as BookInterface, User as UserInterface, FindBooksRequest as FindBooksRequestInterface, FindBooksResponse as FindBooksResponseInterface, CheckoutBookRequest as CheckoutBookRequestInterface, CheckoutBookResponse as CheckoutBookResponseInterface, ReturnBookRequest as ReturnBookRequestInterface, ReturnBookResponse as ReturnBookResponseInterface, GetUserBooksRequest as GetUserBooksRequestInterface, GetUserBooksResponse as GetUserBooksResponseInterface, GetUserRequest as GetUserRequestInterface, GetUserResponse as GetUserResponseInterface, CreateUserRequest as CreateUserRequestInterface, CreateUserResponse as CreateUserResponseInterface } from "./interfaces";
import { LibraryV1Deserializer } from "./deserializer";


/**
 * Book represents a book in the library
 */
export class Book implements BookInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.Book";

  id: string = "";
  title: string = "";
  author: string = "";
  isbn: string = "";
  year: number = 0;
  genre: string = "";
  available: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized Book instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<Book>(Book.MESSAGE_TYPE, data);
  }
}


/**
 * User represents a library user
 */
export class User implements UserInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.User";

  id: string = "";
  name: string = "";
  email: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized User instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<User>(User.MESSAGE_TYPE, data);
  }
}


/**
 * FindBooksRequest is the request for finding books
 */
export class FindBooksRequest implements FindBooksRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.FindBooksRequest";

  /** Search query (title, author, or ISBN) */
  query: string = "";
  /** Genre filter */
  genre: string = "";
  /** Maximum number of results */
  limit: number = 0;
  /** Only show available books */
  availableOnly: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized FindBooksRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<FindBooksRequest>(FindBooksRequest.MESSAGE_TYPE, data);
  }
}


/**
 * FindBooksResponse is the response for finding books
 */
export class FindBooksResponse implements FindBooksResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.FindBooksResponse";

  books: Book[] = [];
  totalCount: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized FindBooksResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<FindBooksResponse>(FindBooksResponse.MESSAGE_TYPE, data);
  }
}


/**
 * CheckoutBookRequest is the request for checking out a book
 */
export class CheckoutBookRequest implements CheckoutBookRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.CheckoutBookRequest";

  bookId: string = "";
  userId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CheckoutBookRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<CheckoutBookRequest>(CheckoutBookRequest.MESSAGE_TYPE, data);
  }
}


/**
 * CheckoutBookResponse is the response for checking out a book
 */
export class CheckoutBookResponse implements CheckoutBookResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.CheckoutBookResponse";

  success: boolean = false;
  message: string = "";
  dueDate: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CheckoutBookResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<CheckoutBookResponse>(CheckoutBookResponse.MESSAGE_TYPE, data);
  }
}


/**
 * ReturnBookRequest is the request for returning a book
 */
export class ReturnBookRequest implements ReturnBookRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.ReturnBookRequest";

  bookId: string = "";
  userId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized ReturnBookRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<ReturnBookRequest>(ReturnBookRequest.MESSAGE_TYPE, data);
  }
}


/**
 * ReturnBookResponse is the response for returning a book
 */
export class ReturnBookResponse implements ReturnBookResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.ReturnBookResponse";

  success: boolean = false;
  message: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized ReturnBookResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<ReturnBookResponse>(ReturnBookResponse.MESSAGE_TYPE, data);
  }
}


/**
 * GetUserBooksRequest is the request for getting user's checked out books
 */
export class GetUserBooksRequest implements GetUserBooksRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.GetUserBooksRequest";

  userId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetUserBooksRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<GetUserBooksRequest>(GetUserBooksRequest.MESSAGE_TYPE, data);
  }
}


/**
 * GetUserBooksResponse is the response for getting user's checked out books
 */
export class GetUserBooksResponse implements GetUserBooksResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.GetUserBooksResponse";

  books: Book[] = [];

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetUserBooksResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<GetUserBooksResponse>(GetUserBooksResponse.MESSAGE_TYPE, data);
  }
}


/**
 * GetUserRequest is the request for getting user info
 */
export class GetUserRequest implements GetUserRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.GetUserRequest";

  userId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetUserRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<GetUserRequest>(GetUserRequest.MESSAGE_TYPE, data);
  }
}


/**
 * GetUserResponse is the response for getting user info
 */
export class GetUserResponse implements GetUserResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.GetUserResponse";

  user?: User;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized GetUserResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<GetUserResponse>(GetUserResponse.MESSAGE_TYPE, data);
  }
}


/**
 * CreateUserRequest is the request for creating a user
 */
export class CreateUserRequest implements CreateUserRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.CreateUserRequest";

  name: string = "";
  email: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CreateUserRequest instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<CreateUserRequest>(CreateUserRequest.MESSAGE_TYPE, data);
  }
}


/**
 * CreateUserResponse is the response for creating a user
 */
export class CreateUserResponse implements CreateUserResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "library.v1.CreateUserResponse";

  user?: User;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CreateUserResponse instance or null if creation failed
   */
  static from(data: any) {
    return LibraryV1Deserializer.from<CreateUserResponse>(CreateUserResponse.MESSAGE_TYPE, data);
  }
}


