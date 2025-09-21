// Copyright 2025 Sri Panyam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { FieldType, FieldSchema, MessageSchema } from './types.js';
import { FactoryInterface, FactoryResult } from '../types/factory.js';

/**
 * Base deserializer class containing all non-template-dependent logic
 */
export abstract class BaseDeserializer {
  protected constructor(
    protected schemaRegistry: Record<string, MessageSchema>,
    protected factory: FactoryInterface
  ) {}

  /**
   * Deserialize an object using schema information
   * @param instance The target instance to populate
   * @param data The source data to deserialize from
   * @param messageType The fully qualified message type (e.g., "library.v1.Book")
   * @returns The populated instance
   */
  deserialize<T>(instance: T, data: any, messageType: string): T {
    if (!data || typeof data !== 'object') {
      return instance;
    }

    const schema = this.schemaRegistry[messageType];
    if (!schema) {
      // Fallback to simple property copying if no schema found
      return this.fallbackDeserialize(instance, data);
    }

    // Process each field according to its schema
    for (const fieldSchema of schema.fields) {
      const fieldValue = data[fieldSchema.name];
      if (fieldValue === null || fieldValue === undefined) {
        continue;
      }

      this.deserializeField(instance, fieldSchema, fieldValue);
    }

    return instance;
  }

  /**
   * Deserialize a single field based on its schema
   */
  protected deserializeField(instance: any, fieldSchema: FieldSchema, fieldValue: any): void {
    const fieldName = fieldSchema.name;

    switch (fieldSchema.type) {
      case FieldType.STRING:
      case FieldType.NUMBER:
      case FieldType.BOOLEAN:
        // Simple primitive types - direct assignment
        instance[fieldName] = fieldValue;
        break;

      case FieldType.MESSAGE:
        if (fieldSchema.repeated) {
          // Handle repeated message fields (arrays)
          instance[fieldName] = this.deserializeMessageArray(
            fieldValue,
            fieldSchema.messageType!,
            instance,
            fieldName
          );
        } else {
          // Handle single message field
          instance[fieldName] = this.deserializeMessageField(
            fieldValue,
            fieldSchema.messageType!,
            instance,
            fieldName
          );
        }
        break;

      case FieldType.REPEATED:
        // Handle repeated primitive fields
        if (Array.isArray(fieldValue)) {
          instance[fieldName] = [...fieldValue]; // Simple copy for primitives
        }
        break;

      case FieldType.ONEOF:
        // Handle oneof fields (would need additional logic for union types)
        instance[fieldName] = fieldValue;
        break;

      case FieldType.MAP:
        // Handle map fields (would need additional schema info for key/value types)
        instance[fieldName] = { ...fieldValue };
        break;

      default:
        // Fallback to direct assignment
        instance[fieldName] = fieldValue;
        break;
    }
  }

  /**
   * Deserialize a single message field
   */
  protected deserializeMessageField(
    fieldValue: any,
    messageType: string,
    parent: any,
    attributeName: string
  ): any {
    // Try to get factory method using cross-package delegation
    let factoryMethod;
    
    if (this.factory.getFactoryMethod) {
      factoryMethod = this.factory.getFactoryMethod(messageType);
    } else {
      // Fallback to simple method name lookup
      const factoryMethodName = this.getFactoryMethodName(messageType);
      factoryMethod = (this.factory as any)[factoryMethodName];
    }

    if (factoryMethod) {
      const result = factoryMethod(parent, attributeName, undefined, fieldValue);
      if (result.fullyLoaded) {
        return result.instance;
      } else {
        // Factory created instance but didn't populate - use deserializer
        return this.deserialize(result.instance, fieldValue, messageType);
      }
    }

    // No factory method found - fallback
    return this.fallbackDeserialize({}, fieldValue);
  }

  /**
   * Deserialize an array of message objects
   */
  protected deserializeMessageArray(
    fieldValue: any[],
    messageType: string,
    parent: any,
    attributeName: string
  ): any[] {
    if (!Array.isArray(fieldValue)) {
      return [];
    }

    // Try to get factory method using cross-package delegation
    let factoryMethod;
    
    if (this.factory.getFactoryMethod) {
      factoryMethod = this.factory.getFactoryMethod(messageType);
    } else {
      // Fallback to simple method name lookup
      const factoryMethodName = this.getFactoryMethodName(messageType);
      factoryMethod = (this.factory as any)[factoryMethodName];
    }

    return fieldValue.map((item, index) => {
      if (factoryMethod) {
        const result = factoryMethod(parent, attributeName, index, item);
        if (result.fullyLoaded) {
          return result.instance;
        } else {
          // Factory created instance but didn't populate - use deserializer
          return this.deserialize(result.instance, item, messageType);
        }
      }

      // No factory method found - fallback
      return this.fallbackDeserialize({}, item);
    });
  }

  /**
   * Convert message type to factory method name
   * "library.v1.Book" -> "newBook"
   */
  protected getFactoryMethodName(messageType: string): string {
    const parts = messageType.split('.');
    const typeName = parts[parts.length - 1]; // Get last part (e.g., "Book")
    return 'new' + typeName;
  }

  /**
   * Fallback deserializer for when no schema is available
   */
  protected fallbackDeserialize<T>(instance: T, data: any): T {
    if (!data || typeof data !== 'object') {
      return instance;
    }

    for (const [key, value] of Object.entries(data)) {
      if (value !== null && value !== undefined) {
        (instance as any)[key] = value;
      }
    }

    return instance;
  }

  /**
   * Create and deserialize a new instance of a message type
   */
  createAndDeserialize<T>(messageType: string, data: any): T {
    // Try to get factory method using cross-package delegation
    let factoryMethod;
    
    if (this.factory.getFactoryMethod) {
      factoryMethod = this.factory.getFactoryMethod(messageType);
    } else {
      // Fallback to simple method name lookup
      const factoryMethodName = this.getFactoryMethodName(messageType);
      factoryMethod = (this.factory as any)[factoryMethodName];
    }

    if (!factoryMethod) {
      throw new Error(`Could not find factory method to deserialize: ${messageType}`)
    }

    const result = factoryMethod(undefined, undefined, undefined, data);
    if (result.fullyLoaded) {
      return result.instance;
    } else {
      return this.deserialize(result.instance, data, messageType);
    }
  }
}
