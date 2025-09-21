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

import { FieldSchema, MessageSchema } from './types.js';

/**
 * Base schema registry with utility methods for schema operations
 */
export abstract class BaseSchemaRegistry {
  constructor(protected schemaRegistry: Record<string, MessageSchema>) {}

  /**
   * Get schema for a message type
   */
  getSchema(messageType: string): MessageSchema | undefined {
    return this.schemaRegistry[messageType];
  }

  /**
   * Get field schema by name
   */
  getFieldSchema(messageType: string, fieldName: string): FieldSchema | undefined {
    const schema = this.getSchema(messageType);
    return schema?.fields.find(field => field.name === fieldName);
  }

  /**
   * Get field schema by proto field ID
   */
  getFieldSchemaById(messageType: string, fieldId: number): FieldSchema | undefined {
    const schema = this.getSchema(messageType);
    return schema?.fields.find(field => field.id === fieldId);
  }

  /**
   * Check if field is part of a oneof group
   */
  isOneofField(messageType: string, fieldName: string): boolean {
    const fieldSchema = this.getFieldSchema(messageType, fieldName);
    return fieldSchema?.oneofGroup !== undefined;
  }

  /**
   * Get all fields in a oneof group
   */
  getOneofFields(messageType: string, oneofGroup: string): FieldSchema[] {
    const schema = this.getSchema(messageType);
    return schema?.fields.filter(field => field.oneofGroup === oneofGroup) || [];
  }
}
