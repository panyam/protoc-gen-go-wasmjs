// Generated TypeScript factory + deserializer (annotation-based)
// DO NOT EDIT - This file is auto-generated


import { MessageTypeConstructor, BaseDeserializer, FactoryInterface } from "@protoc-gen-go-wasmjs/runtime";


import { test_one_package_v1SchemaRegistry } from "./schemas";




import {SampleRequest as SampleRequestInterface,SampleResponse as SampleResponseInterface} from "./models/interfaces";

import {SecondRequest as SecondRequestInterface,SecondResponse as SecondResponseInterface} from "./models2/interfaces";

import {SampleRequest as ConcreteSampleRequest,SampleResponse as ConcreteSampleResponse} from "./models/models";

import {SecondRequest as ConcreteSecondRequest,SecondResponse as ConcreteSecondResponse} from "./models2/models";




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
export class Test_one_package_v1Factory {


  /**
   * Enhanced factory method for SampleRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newSampleRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<SampleRequestInterface> => {
    const out = new ConcreteSampleRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for SampleResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newSampleResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<SampleResponseInterface> => {
    const out = new ConcreteSampleResponse();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for SecondRequest
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newSecondRequest = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<SecondRequestInterface> => {
    const out = new ConcreteSecondRequest();
    
    // Factory does not populate by default - let deserializer handle it
    return { instance: out, fullyLoaded: false };
  }

  /**
   * Enhanced factory method for SecondResponse
   * @param parent Parent object containing this field
   * @param attributeName Field name in parent object
   * @param attributeKey Array index, map key, or union tag (for containers)
   * @param data Raw data to potentially populate from
   * @returns Factory result with instance and population status
   */
  newSecondResponse = (
    parent?: any,
    attributeName?: string,
    attributeKey?: string | number,
    data?: any
  ): FactoryResult<SecondResponseInterface> => {
    const out = new ConcreteSecondResponse();
    
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
    const currentPackage = "test_one_package.v1";
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

// Shared factory instance to avoid creating new instances on every deserializer construction
const DEFAULT_FACTORY = new Test_one_package_v1Factory();

/**
 * Schema-aware deserializer for test_one_package.v1 package
 * Extends BaseDeserializer with package-specific configuration
 */
export class Test_one_package_v1Deserializer extends BaseDeserializer {
  constructor(
    schemaRegistry = test_one_package_v1SchemaRegistry,
    factory: FactoryInterface = DEFAULT_FACTORY
  ) {
    super(schemaRegistry, factory);
  }

  /**
   * Static utility method to create and deserialize a message without needing a deserializer instance
   * @param messageType Fully qualified message type (use Class.MESSAGE_TYPE)
   * @param data Raw data to deserialize
   * @returns Deserialized instance or null if creation failed
   */
  static fromMsgType<T>(messageType: string, data: any) {
    const deserializer = new Test_one_package_v1Deserializer(); // Uses default factory and schema registry
    return deserializer.createAndDeserialize<T>(messageType, data);
  }


  /**
   * Static utility method - infers messageType from type parameter
   * Type-safe convenience method
   */
  static from<T>(
    typeConstructor: MessageTypeConstructor<T>,
    data: any
  ): T | null {
    const deserializer = new Test_one_package_v1Deserializer();
    return deserializer.createAndDeserialize<T>(typeConstructor.MESSAGE_TYPE, data);
  }

  // Deserialize if data is already a partial instance
  static fromPartial<T extends { __MESSAGE_TYPE: string }>(data: T): T {
    const messageType = data.__MESSAGE_TYPE;
    return this.fromMsgType<T>(messageType, data);
  }
}
