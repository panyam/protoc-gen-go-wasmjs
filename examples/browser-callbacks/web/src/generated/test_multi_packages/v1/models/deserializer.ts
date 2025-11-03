
// Generated TypeScript schema-aware deserializer
// DO NOT EDIT - This file is auto-generated

import { BaseDeserializer, FactoryInterface } from "@protoc-gen-go-wasmjs/runtime";
import { Test_multi_packages_v1_modelsFactory } from "./factory";
import { test_multi_packages_v1_modelsSchemaRegistry } from "./schemas";

// Shared factory instance to avoid creating new instances on every deserializer construction
const DEFAULT_FACTORY = new Test_multi_packages_v1_modelsFactory();

/**
 * Schema-aware deserializer for test_multi_packages.v1.models package
 * Extends BaseDeserializer with package-specific configuration
 */
export class Test_multi_packages_v1_modelsDeserializer extends BaseDeserializer {
  constructor(
    schemaRegistry = test_multi_packages_v1_modelsSchemaRegistry,
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
    const deserializer = new Test_multi_packages_v1_modelsDeserializer(); // Uses default factory and schema registry
    return deserializer.createAndDeserialize<T>(messageType, data);
  }
}
