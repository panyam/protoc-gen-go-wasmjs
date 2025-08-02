
import { Book as BookInterface, User as UserInterface, FindBooksRequest as FindBooksRequestInterface, FindBooksResponse as FindBooksResponseInterface, CheckoutBookRequest as CheckoutBookRequestInterface, CheckoutBookResponse as CheckoutBookResponseInterface, ReturnBookRequest as ReturnBookRequestInterface, ReturnBookResponse as ReturnBookResponseInterface, GetUserBooksRequest as GetUserBooksRequestInterface, GetUserBooksResponse as GetUserBooksResponseInterface, GetUserRequest as GetUserRequestInterface, GetUserResponse as GetUserResponseInterface, CreateUserRequest as CreateUserRequestInterface, CreateUserResponse as CreateUserResponseInterface } from "./library_interfaces";


import { Book as ConcreteBook, User as ConcreteUser, FindBooksRequest as ConcreteFindBooksRequest, FindBooksResponse as ConcreteFindBooksResponse, CheckoutBookRequest as ConcreteCheckoutBookRequest, CheckoutBookResponse as ConcreteCheckoutBookResponse, ReturnBookRequest as ConcreteReturnBookRequest, ReturnBookResponse as ConcreteReturnBookResponse, GetUserBooksRequest as ConcreteGetUserBooksRequest, GetUserBooksResponse as ConcreteGetUserBooksResponse, GetUserRequest as ConcreteGetUserRequest, GetUserResponse as ConcreteGetUserResponse, CreateUserRequest as ConcreteCreateUserRequest, CreateUserResponse as ConcreteCreateUserResponse } from "./library_models";


export class LibraryV1Factory {
  newBook = (data?: any): BookInterface => {
    const out = new ConcreteBook();
    if (data) {
      out.id = data.id ?? "";
      out.title = data.title ?? "";
      out.author = data.author ?? "";
      out.isbn = data.isbn ?? "";
      out.year = data.year ?? 0;
      out.genre = data.genre ?? "";
      out.available = data.available ?? false;
    }
    return out;
  }

  newUser = (data?: any): UserInterface => {
    const out = new ConcreteUser();
    if (data) {
      out.id = data.id ?? "";
      out.name = data.name ?? "";
      out.email = data.email ?? "";
    }
    return out;
  }

  newFindBooksRequest = (data?: any): FindBooksRequestInterface => {
    const out = new ConcreteFindBooksRequest();
    if (data) {
      out.query = data.query ?? "";
      out.genre = data.genre ?? "";
      out.limit = data.limit ?? 0;
      out.availableOnly = data.availableOnly ?? false;
    }
    return out;
  }

  newFindBooksResponse = (data?: any): FindBooksResponseInterface => {
    const out = new ConcreteFindBooksResponse();
    if (data) {
      if (data.books && Array.isArray(data.books)) {
        out.books = data.books.map((item: any) => this.newBook(item));
      }
      out.totalCount = data.totalCount ?? 0;
    }
    return out;
  }

  newCheckoutBookRequest = (data?: any): CheckoutBookRequestInterface => {
    const out = new ConcreteCheckoutBookRequest();
    if (data) {
      out.bookId = data.bookId ?? "";
      out.userId = data.userId ?? "";
    }
    return out;
  }

  newCheckoutBookResponse = (data?: any): CheckoutBookResponseInterface => {
    const out = new ConcreteCheckoutBookResponse();
    if (data) {
      out.success = data.success ?? false;
      out.message = data.message ?? "";
      out.dueDate = data.dueDate ?? "";
    }
    return out;
  }

  newReturnBookRequest = (data?: any): ReturnBookRequestInterface => {
    const out = new ConcreteReturnBookRequest();
    if (data) {
      out.bookId = data.bookId ?? "";
      out.userId = data.userId ?? "";
    }
    return out;
  }

  newReturnBookResponse = (data?: any): ReturnBookResponseInterface => {
    const out = new ConcreteReturnBookResponse();
    if (data) {
      out.success = data.success ?? false;
      out.message = data.message ?? "";
    }
    return out;
  }

  newGetUserBooksRequest = (data?: any): GetUserBooksRequestInterface => {
    const out = new ConcreteGetUserBooksRequest();
    if (data) {
      out.userId = data.userId ?? "";
    }
    return out;
  }

  newGetUserBooksResponse = (data?: any): GetUserBooksResponseInterface => {
    const out = new ConcreteGetUserBooksResponse();
    if (data) {
      if (data.books && Array.isArray(data.books)) {
        out.books = data.books.map((item: any) => this.newBook(item));
      }
    }
    return out;
  }

  newGetUserRequest = (data?: any): GetUserRequestInterface => {
    const out = new ConcreteGetUserRequest();
    if (data) {
      out.userId = data.userId ?? "";
    }
    return out;
  }

  newGetUserResponse = (data?: any): GetUserResponseInterface => {
    const out = new ConcreteGetUserResponse();
    if (data) {
      if (data.user) out.user = this.newUser(data.user);
    }
    return out;
  }

  newCreateUserRequest = (data?: any): CreateUserRequestInterface => {
    const out = new ConcreteCreateUserRequest();
    if (data) {
      out.name = data.name ?? "";
      out.email = data.email ?? "";
    }
    return out;
  }

  newCreateUserResponse = (data?: any): CreateUserResponseInterface => {
    const out = new ConcreteCreateUserResponse();
    if (data) {
      if (data.user) out.user = this.newUser(data.user);
    }
    return out;
  }

}