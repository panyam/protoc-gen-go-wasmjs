
// Generated TypeScript schema-aware deserializer
// DO NOT EDIT - This file is auto-generated

import { BaseDeserializer, FactoryInterface } from "@protoc-gen-go-wasmjs/runtime";
import { Test_one_package_v1Factory } from "./factory";
import { test_one_package_v1SchemaRegistry } from "./schemas";

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
  static from<T>(messageType: string, data: any) {
    const deserializer = new Test_one_package_v1Deserializer(); // Uses default factory and schema registry
    return deserializer.createAndDeserialize<T>(messageType, data);
  }
}
