
import { BaseMessage } from "../common/interfaces";
import { Metadata } from "../common/interfaces";
import { ErrorInfo } from "../common/interfaces";



import { BookStats as BookStatsInterface, UserActivity as UserActivityInterface, GetBookStatsRequest as GetBookStatsRequestInterface, GetBookStatsResponse as GetBookStatsResponseInterface, Pagination as PaginationInterface, PaginationResponse as PaginationResponseInterface, Book as BookInterface, User as UserInterface, FindBooksRequest as FindBooksRequestInterface, FindBooksResponse as FindBooksResponseInterface } from "./interfaces";


import { BookStats as ConcreteBookStats, UserActivity as ConcreteUserActivity, GetBookStatsRequest as ConcreteGetBookStatsRequest, GetBookStatsResponse as ConcreteGetBookStatsResponse, Pagination as ConcretePagination, PaginationResponse as ConcretePaginationResponse, Book as ConcreteBook, User as ConcreteUser, FindBooksRequest as ConcreteFindBooksRequest, FindBooksResponse as ConcreteFindBooksResponse } from "./models";

import { LibraryCommonFactory } from "../common/factory";


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
export class LibraryV2Factory {
  // Dependency factory for library.common package
  private commonFactory = new LibraryCommonFactory();


  /**
   * Enhanced factory method for BookStats
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newBookStats = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<BookStatsInterface> => {
    const out = new ConcreteBookStats();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for UserActivity
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newUserActivity = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<UserActivityInterface> => {
    const out = new ConcreteUserActivity();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetBookStatsRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetBookStatsRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetBookStatsRequestInterface> => {
    const out = new ConcreteGetBookStatsRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for GetBookStatsResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newGetBookStatsResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<GetBookStatsResponseInterface> => {
    const out = new ConcreteGetBookStatsResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for Pagination
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPagination = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PaginationInterface> => {
    const out = new ConcretePagination();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for PaginationResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPaginationResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PaginationResponseInterface> => {
    const out = new ConcretePaginationResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

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
    const currentPackage = "library.v2";
    if (packageName === currentPackage) {
      return (this as any)[methodName];
    }
    
    // Check external type factory mappings
    const externalFactory = this.externalTypeFactories()[messageType];
    if (externalFactory) {
      return externalFactory;
    }
    
    // Delegate to appropriate dependency factory
    if (packageName === "library.common") {
      return (this.commonFactory as any)[methodName];
    }

    
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
