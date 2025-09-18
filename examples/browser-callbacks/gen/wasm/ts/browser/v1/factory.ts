

import { FetchRequest as FetchRequestInterface, FetchResponse as FetchResponseInterface, StorageKeyRequest as StorageKeyRequestInterface, StorageValueResponse as StorageValueResponseInterface, StorageSetRequest as StorageSetRequestInterface, StorageSetResponse as StorageSetResponseInterface, CookieRequest as CookieRequestInterface, CookieResponse as CookieResponseInterface, AlertRequest as AlertRequestInterface, AlertResponse as AlertResponseInterface, PromptRequest as PromptRequestInterface, PromptResponse as PromptResponseInterface, LogRequest as LogRequestInterface, LogResponse as LogResponseInterface } from "./interfaces";


import { FetchRequest as ConcreteFetchRequest, FetchResponse as ConcreteFetchResponse, StorageKeyRequest as ConcreteStorageKeyRequest, StorageValueResponse as ConcreteStorageValueResponse, StorageSetRequest as ConcreteStorageSetRequest, StorageSetResponse as ConcreteStorageSetResponse, CookieRequest as ConcreteCookieRequest, CookieResponse as ConcreteCookieResponse, AlertRequest as ConcreteAlertRequest, AlertResponse as ConcreteAlertResponse, PromptRequest as ConcretePromptRequest, PromptResponse as ConcretePromptResponse, LogRequest as ConcreteLogRequest, LogResponse as ConcreteLogResponse } from "./models";



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
export class BrowserV1Factory {


  /**
   * Enhanced factory method for FetchRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newFetchRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<FetchRequestInterface> => {
    const out = new ConcreteFetchRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for FetchResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newFetchResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<FetchResponseInterface> => {
    const out = new ConcreteFetchResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for StorageKeyRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newStorageKeyRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<StorageKeyRequestInterface> => {
    const out = new ConcreteStorageKeyRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for StorageValueResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newStorageValueResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<StorageValueResponseInterface> => {
    const out = new ConcreteStorageValueResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for StorageSetRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newStorageSetRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<StorageSetRequestInterface> => {
    const out = new ConcreteStorageSetRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for StorageSetResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newStorageSetResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<StorageSetResponseInterface> => {
    const out = new ConcreteStorageSetResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CookieRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCookieRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CookieRequestInterface> => {
    const out = new ConcreteCookieRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for CookieResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newCookieResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<CookieResponseInterface> => {
    const out = new ConcreteCookieResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for AlertRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newAlertRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<AlertRequestInterface> => {
    const out = new ConcreteAlertRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for AlertResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newAlertResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<AlertResponseInterface> => {
    const out = new ConcreteAlertResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for PromptRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPromptRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PromptRequestInterface> => {
    const out = new ConcretePromptRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for PromptResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newPromptResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<PromptResponseInterface> => {
    const out = new ConcretePromptResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for LogRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newLogRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<LogRequestInterface> => {
    const out = new ConcreteLogRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for LogResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newLogResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<LogResponseInterface> => {
    const out = new ConcreteLogResponse();
    
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
    const currentPackage = "browser.v1";
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
