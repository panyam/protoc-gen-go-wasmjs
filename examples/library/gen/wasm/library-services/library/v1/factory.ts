

import { Book as BookInterface, User as UserInterface, FindBooksRequest as FindBooksRequestInterface, FindBooksResponse as FindBooksResponseInterface, CheckoutBookRequest as CheckoutBookRequestInterface, CheckoutBookResponse as CheckoutBookResponseInterface, ReturnBookRequest as ReturnBookRequestInterface, ReturnBookResponse as ReturnBookResponseInterface, GetUserBooksRequest as GetUserBooksRequestInterface, GetUserBooksResponse as GetUserBooksResponseInterface, GetUserRequest as GetUserRequestInterface, GetUserResponse as GetUserResponseInterface, CreateUserRequest as CreateUserRequestInterface, CreateUserResponse as CreateUserResponseInterface } from "./interfaces";


import { Book as ConcreteBook, User as ConcreteUser, FindBooksRequest as ConcreteFindBooksRequest, FindBooksResponse as ConcreteFindBooksResponse, CheckoutBookRequest as ConcreteCheckoutBookRequest, CheckoutBookResponse as ConcreteCheckoutBookResponse, ReturnBookRequest as ConcreteReturnBookRequest, ReturnBookResponse as ConcreteReturnBookResponse, GetUserBooksRequest as ConcreteGetUserBooksRequest, GetUserBooksResponse as ConcreteGetUserBooksResponse, GetUserRequest as ConcreteGetUserRequest, GetUserResponse as ConcreteGetUserResponse, CreateUserRequest as ConcreteCreateUserRequest, CreateUserResponse as ConcreteCreateUserResponse } from "./models";



/**
 * Factory result interface for enhanced factory methods
 */
export interface FactoryResult<T> {
  instance: T;
  fullyLoaded: boolean;
}

/**
 * Enhanced factory with context-aware object construction
 */
export class LibraryV1Factory {


  /**
   * Enhanced factory method for Book
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newBook = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<BookInterface> => {
    const out = new ConcreteBook();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for User
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newUser = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<UserInterface> => {
    const out = new ConcreteUser();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for FindBooksRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newFindBooksRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<FindBooksRequestInterface> => {
    const out = new ConcreteFindBooksRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for FindBooksResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newFindBooksResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<FindBooksResponseInterface> => {
    const out = new ConcreteFindBooksResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CheckoutBookRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCheckoutBookRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CheckoutBookRequestInterface> => {
    const out = new ConcreteCheckoutBookRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CheckoutBookResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCheckoutBookResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CheckoutBookResponseInterface> => {
    const out = new ConcreteCheckoutBookResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for ReturnBookRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newReturnBookRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<ReturnBookRequestInterface> => {
    const out = new ConcreteReturnBookRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for ReturnBookResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newReturnBookResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<ReturnBookResponseInterface> => {
    const out = new ConcreteReturnBookResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetUserBooksRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetUserBooksRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetUserBooksRequestInterface> => {
    const out = new ConcreteGetUserBooksRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetUserBooksResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetUserBooksResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetUserBooksResponseInterface> => {
    const out = new ConcreteGetUserBooksResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetUserRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetUserRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetUserRequestInterface> => {
    const out = new ConcreteGetUserRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetUserResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetUserResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetUserResponseInterface> => {
    const out = new ConcreteGetUserResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CreateUserRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCreateUserRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CreateUserRequestInterface> => {
    const out = new ConcreteCreateUserRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CreateUserResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCreateUserResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CreateUserResponseInterface> => {
    const out = new ConcreteCreateUserResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }



  /**
   * Get factory method for a fully qualified message type
   * Enables cross-package factory delegation
   */
  getFactoryMethod(messageType: string): ((parent?: any, attributeName?: string, attributeKey?: string | number, data?: any) => FactoryResult<any>) | undefined {
    // Extract package from message type (e.g., "library.common.BaseMessage" -> "library.common")
    const parts = messageType.split('.');
    if (parts.length < 2) {
      return undefined;
    }
    
    const packageName = parts.slice(0, -1).join('.');
    const typeName = parts[parts.length - 1];
    const methodName = 'new' + typeName;
    
    // Check if this is our own package first
    const currentPackage = "library.v1";
    if (packageName === currentPackage) {
      return (this as any)[methodName];
    }
    
    // Check external type factory mappings
    const externalFactory = this.externalTypeFactories()[messageType];
    if (externalFactory) {
      return externalFactory;
    }
    
    // Delegate to appropriate dependency factory

    
    return undefined;
  }



  /**
   * Generic object deserializer that respects factory decisions
   */
  protected deserializeObject(instance: any, data: any): any {
    if (!data || typeof data !== 'object') return instance;
    
    for (const [key, value] of Object.entries(data)) {
      if (value !== null && value !== undefined) {
        instance[key] = value;
      }
    }
    return instance;
  }

  // External type conversion methods

  /**
   * Mapping of external types to their factory methods
   */
  private externalTypeFactories(): Record<string, (parent?: any, attributeName?: string, attributeKey?: string | number, data?: any) => FactoryResult<any>> { 
      return {
          "google.protobuf.Timestamp": this.newTimestamp,
          "google.protobuf.FieldMask": this.newFieldMask,
      }
  };

  /**
   * Convert native Date to protobuf Timestamp format for serialization
   */
  serializeTimestamp(date: Date): any {
    if (!date) return null;
    return {
      seconds: Math.floor(date.getTime() / 1000).toString(),
      nanos: (date.getTime() % 1000) * 1000000
    };
  }

  /**
   * Factory method for converting protobuf Timestamp data to native Date
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object  
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw protobuf timestamp data
   * @returns Factory result with Date instance
   */
  newTimestamp = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<Date> => {
    if (!data) {
      return { instance: new Date(), fullyLoaded: true };
    }
    
    let date: Date;
    if (typeof data === 'string') {
      // Handle ISO string format
      date = new Date(data);
    } else if (data.seconds !== undefined) {
      // Handle protobuf format with seconds/nanos
      const seconds = typeof data.seconds === 'string' 
        ? parseInt(data.seconds, 10) 
        : data.seconds;
      const nanos = data.nanos || 0;
      date = new Date(seconds * 1000 + Math.floor(nanos / 1000000));
    } else {
      date = new Date();
    }
    
    return { instance: date, fullyLoaded: true };
  }

  /**
   * Convert native string array to protobuf FieldMask format for serialization
   */
  serializeFieldMask(paths: string[]): any {
    if (!paths || !Array.isArray(paths)) return null;
    return { paths };
  }

  /**
   * Factory method for converting protobuf FieldMask data to native string array
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw protobuf field mask data
   * @returns Factory result with string array instance
   */
  newFieldMask = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<string[]> => {
    if (!data) {
      return { instance: [], fullyLoaded: true };
    }
    
    let paths: string[];
    if (Array.isArray(data)) {
      paths = data;
    } else if (data.paths && Array.isArray(data.paths)) {
      paths = data.paths;
    } else {
      paths = [];
    }
    
    return { instance: paths, fullyLoaded: true };
  }
}
