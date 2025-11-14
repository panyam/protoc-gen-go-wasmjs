/**
 * Generation Stability Tests
 *
 * This file imports all generated TypeScript artifacts to ensure they remain
 * stable and don't change unexpectedly across generator improvements.
 *
 * If this file fails to compile after regeneration, it indicates a breaking change
 * in the generated code structure that needs to be addressed.
 */

// ============================================================================
// Service Clients
// ============================================================================

import { PresenterServiceClient } from './generated/presenter/v1/presenterServiceClient';
import { BrowserAPIClient } from './generated/browser/v1/browserAPIClient';
import { TestServiceClient as TestOnePackageServiceClient } from './generated/test_one_package/v1/services/testServiceClient';
import { TestServiceClient as TestMultiPackagesServiceClient } from './generated/test_multi_packages/v1/services/testServiceClient';

// ============================================================================
// Interfaces (TypeScript type definitions)
// ============================================================================

// Presenter package interfaces
import type {
  LoadUserDataRequest,
  LoadUserDataResponse,
  StateUpdateRequest,
  UIUpdate,
  TestRecord,
  PreferencesRequest,
  PreferencesResponse,
  CallbackDemoRequest,
  CallbackDemoResponse,
} from './generated/presenter/v1/interfaces';

// Browser API interfaces
import type {
  FetchRequest,
  FetchResponse,
  StorageKeyRequest,
  StorageValueResponse,
  StorageSetRequest,
  StorageSetResponse,
  CookieRequest,
  CookieResponse,
  AlertRequest,
  AlertResponse,
  PromptRequest,
  PromptResponse,
  LogRequest,
  LogResponse,
} from './generated/browser/v1/interfaces';

// Test packages interfaces
import type {
  SampleRequest as TestOneSampleRequest,
  SampleResponse as TestOneSampleResponse,
} from './generated/test_one_package/v1/models/interfaces';

import type {
  SecondRequest,
  SecondResponse,
} from './generated/test_one_package/v1/models2/interfaces';

import type {
  SampleRequest as TestMultiSampleRequest,
  SampleResponse as TestMultiSampleResponse,
} from './generated/test_multi_packages/v1/models/interfaces';

// Utils package interfaces
import type {
  HelperUtilType,
  ParentUtilMessage,
  ParentUtilMessage_NestedUtilType,
} from './generated/utils/v1/interfaces';

// WASMJS annotations interfaces
import type {
  ConflictResolution,
  PatchOperation,
  PatchSource,
  StatefulOptions,
  StatefulMethodOptions,
  AsyncMethodOptions,
  MessagePatch,
  PatchBatch,
  PatchResponse,
} from './generated/wasmjs/v1/interfaces';

// ============================================================================
// Models (Concrete class implementations)
// ============================================================================

// Presenter package models
import {
  LoadUserDataRequest as ConcreteLoadUserDataRequest,
  LoadUserDataResponse as ConcreteLoadUserDataResponse,
  StateUpdateRequest as ConcreteStateUpdateRequest,
  UIUpdate as ConcreteUIUpdate,
  TestRecord as ConcreteTestRecord,
  PreferencesRequest as ConcretePreferencesRequest,
  PreferencesResponse as ConcretePreferencesResponse,
  CallbackDemoRequest as ConcreteCallbackDemoRequest,
  CallbackDemoResponse as ConcreteCallbackDemoResponse,
} from './generated/presenter/v1/models';

// Browser API models
import {
  FetchRequest as ConcreteFetchRequest,
  FetchResponse as ConcreteFetchResponse,
  StorageKeyRequest as ConcreteStorageKeyRequest,
  StorageValueResponse as ConcreteStorageValueResponse,
  StorageSetRequest as ConcreteStorageSetRequest,
  StorageSetResponse as ConcreteStorageSetResponse,
  CookieRequest as ConcreteCookieRequest,
  CookieResponse as ConcreteCookieResponse,
  AlertRequest as ConcreteAlertRequest,
  AlertResponse as ConcreteAlertResponse,
  PromptRequest as ConcretePromptRequest,
  PromptResponse as ConcretePromptResponse,
  LogRequest as ConcreteLogRequest,
  LogResponse as ConcreteLogResponse,
} from './generated/browser/v1/models';

// Test packages models
import {
  SampleRequest as ConcreteTestOneSampleRequest,
  SampleResponse as ConcreteTestOneSampleResponse,
} from './generated/test_one_package/v1/models/models';

import {
  SecondRequest as ConcreteSecondRequest,
  SecondResponse as ConcreteSecondResponse,
} from './generated/test_one_package/v1/models2/models';

import {
  SampleRequest as ConcreteTestMultiSampleRequest,
  SampleResponse as ConcreteTestMultiSampleResponse,
} from './generated/test_multi_packages/v1/models/models';

// Utils package models
import {
  HelperUtilType as ConcreteHelperUtilType,
  ParentUtilMessage as ConcreteParentUtilMessage,
  ParentUtilMessage_NestedUtilType as ConcreteParentUtilMessage_NestedUtilType,
} from './generated/utils/v1/models';

// WASMJS annotations models
import {
  StatefulOptions as ConcreteStatefulOptions,
  StatefulMethodOptions as ConcreteStatefulMethodOptions,
  AsyncMethodOptions as ConcreteAsyncMethodOptions,
  MessagePatch as ConcreteMessagePatch,
  PatchBatch as ConcretePatchBatch,
  PatchResponse as ConcretePatchResponse,
} from './generated/wasmjs/v1/models';

// ============================================================================
// Schemas (Field metadata for runtime introspection)
// ============================================================================

import { browser_v1SchemaRegistry } from './generated/browser/v1/schemas';
import { presenter_v1SchemaRegistry } from './generated/presenter/v1/schemas';
import { test_one_package_v1SchemaRegistry as testOnePackageSchemas } from './generated/test_one_package/v1/models/schemas';
import { test_one_package_v1SchemaRegistry as testOnePackageModels2Schemas } from './generated/test_one_package/v1/models2/schemas';
import { test_multi_packages_v1_modelsSchemaRegistry as testMultiPackageSchemas } from './generated/test_multi_packages/v1/models/schemas';
import { utils_v1SchemaRegistry } from './generated/utils/v1/schemas';
import { wasmjs_v1SchemaRegistry } from './generated/wasmjs/v1/schemas';

// ============================================================================
// Factory + Deserializer (Annotation-based)
// ============================================================================

import {
  Test_one_package_v1Factory,
  Test_one_package_v1Deserializer,
  FactoryResult,
} from './generated/test_one_package/v1/factory';

// Package-level aggregated schemas (used by factory)
import { test_one_package_v1SchemaRegistry as aggregatedTestOnePackageSchemas } from './generated/test_one_package/v1/schemas';

// ============================================================================
// Bundle (Base WASM loader and service manager)
// ============================================================================

import { ExampleBundle } from './generated';

// ============================================================================
// Runtime Types (from @protoc-gen-go-wasmjs/runtime package)
// ============================================================================

import type {
  WASMServiceClient,
  BrowserServiceManager,
  MessageSchema,
  FieldSchema,
  BaseDeserializer,
  FactoryInterface,
} from '@protoc-gen-go-wasmjs/runtime';

// ============================================================================
// Test Functions - Verify all imports work correctly
// ============================================================================

/**
 * Test: Service client instantiation
 */
export function testServiceClients(): void {
  const bundle = new ExampleBundle();

  const presenterClient = new PresenterServiceClient(bundle);
  const browserClient = new BrowserAPIClient(bundle);
  const testOneClient = new TestOnePackageServiceClient(bundle);
  const testMultiClient = new TestMultiPackagesServiceClient(bundle);

  console.log('✅ All service clients instantiated successfully');
}

/**
 * Test: Model instantiation with defaults
 */
export function testModelInstantiation(): void {
  // Presenter models
  const loadUserReq = new ConcreteLoadUserDataRequest();
  const loadUserResp = new ConcreteLoadUserDataResponse();
  const stateUpdateReq = new ConcreteStateUpdateRequest();
  const uiUpdate = new ConcreteUIUpdate();
  const testRecord = new ConcreteTestRecord();
  const prefsReq = new ConcretePreferencesRequest();
  const prefsResp = new ConcretePreferencesResponse();
  const callbackReq = new ConcreteCallbackDemoRequest();
  const callbackResp = new ConcreteCallbackDemoResponse();

  // Browser API models
  const fetchReq = new ConcreteFetchRequest();
  const fetchResp = new ConcreteFetchResponse();
  const storageKeyReq = new ConcreteStorageKeyRequest();
  const storageValueResp = new ConcreteStorageValueResponse();
  const storageSetReq = new ConcreteStorageSetRequest();
  const storageSetResp = new ConcreteStorageSetResponse();
  const cookieReq = new ConcreteCookieRequest();
  const cookieResp = new ConcreteCookieResponse();
  const alertReq = new ConcreteAlertRequest();
  const alertResp = new ConcreteAlertResponse();
  const promptReq = new ConcretePromptRequest();
  const promptResp = new ConcretePromptResponse();
  const logReq = new ConcreteLogRequest();
  const logResp = new ConcreteLogResponse();

  // Test package models
  const testOneSampleReq = new ConcreteTestOneSampleRequest();
  const testOneSampleResp = new ConcreteTestOneSampleResponse();
  const secondReq = new ConcreteSecondRequest();
  const secondResp = new ConcreteSecondResponse();
  const testMultiSampleReq = new ConcreteTestMultiSampleRequest();
  const testMultiSampleResp = new ConcreteTestMultiSampleResponse();

  // Utils models
  const helperUtil = new ConcreteHelperUtilType();
  const parentUtil = new ConcreteParentUtilMessage();
  const nestedUtil = new ConcreteParentUtilMessage_NestedUtilType();

  // WASMJS models
  const statefulOptions = new ConcreteStatefulOptions();
  const statefulMethodOptions = new ConcreteStatefulMethodOptions();
  const asyncMethodOptions = new ConcreteAsyncMethodOptions();
  const messagePatch = new ConcreteMessagePatch();
  const patchBatch = new ConcretePatchBatch();
  const patchResponse = new ConcretePatchResponse();

  console.log('✅ All models instantiated successfully with defaults');
}

/**
 * Test: Schema registry access
 */
export function testSchemaRegistries(): void {
  // Verify schema registries exist and are objects
  const registries = [
    browser_v1SchemaRegistry,
    presenter_v1SchemaRegistry,
    testOnePackageSchemas,
    testOnePackageModels2Schemas,
    testMultiPackageSchemas,
    utils_v1SchemaRegistry,
    wasmjs_v1SchemaRegistry,
    aggregatedTestOnePackageSchemas,
  ];

  registries.forEach((registry, index) => {
    if (typeof registry !== 'object') {
      throw new Error(`Schema registry ${index} is not an object`);
    }
  });

  console.log('✅ All schema registries accessible');
}

/**
 * Test: Factory and deserializer functionality
 */
export function testFactoryAndDeserializer(): void {
  const factory = new Test_one_package_v1Factory();
  const deserializer = new Test_one_package_v1Deserializer();

  // Test factory methods exist
  const sampleReqResult: FactoryResult<TestOneSampleRequest> = factory.newSampleRequest();
  const sampleRespResult: FactoryResult<TestOneSampleResponse> = factory.newSampleResponse();
  const secondReqResult: FactoryResult<SecondRequest> = factory.newSecondRequest();
  const secondRespResult: FactoryResult<SecondResponse> = factory.newSecondResponse();

  // Verify instances were created
  if (!sampleReqResult.instance || !sampleRespResult.instance) {
    throw new Error('Factory failed to create instances');
  }

  // Test static deserializer utility
  const bundle = new ExampleBundle();
  const testData = { name: 'test', value: 42 };
  const deserialized = Test_one_package_v1Deserializer.from(ConcreteTestOneSampleRequest, testData);
  const client = new TestOnePackageServiceClient(bundle)
  client.sample(deserialized!)

  console.log('✅ Factory and deserializer working correctly');
}

/**
 * Test: Type checking - interfaces are assignable to concrete models
 */
export function testTypeCompatibility(): void {
  // Interface should be assignable from concrete model
  const concreteRequest = new ConcreteLoadUserDataRequest();
  const interfaceRequest: LoadUserDataRequest = concreteRequest;

  // Verify type compatibility
  const userId: string = interfaceRequest.userId || '';

  console.log('✅ Interface and model types are compatible');
}

/**
 * Test: Cross-package imports work correctly
 */
export function testCrossPackageImports(): void {
  // HelperUtilType is from utils.v1 package
  const helper = new ConcreteHelperUtilType();

  // ParentUtilMessage and nested type
  const parent = new ConcreteParentUtilMessage();
  const nested = new ConcreteParentUtilMessage_NestedUtilType();

  // These should be importable and instantiable without errors
  console.log('✅ Cross-package imports working correctly');
}

/**
 * Run all stability tests
 */
export function runAllStabilityTests(): void {
  console.log('🧪 Running generation stability tests...\n');

  try {
    testServiceClients();
    testModelInstantiation();
    testSchemaRegistries();
    testFactoryAndDeserializer();
    testTypeCompatibility();
    testCrossPackageImports();

    console.log('\n✅ All generation stability tests passed!');
  } catch (error) {
    console.error('\n❌ Stability tests failed:', error);
    throw error;
  }
}

// Export everything for use in other test files
export {
  // Service clients
  PresenterServiceClient,
  BrowserAPIClient,
  TestOnePackageServiceClient,
  TestMultiPackagesServiceClient,

  // Bundle
  ExampleBundle,

  // Factory/Deserializer
  Test_one_package_v1Factory,
  Test_one_package_v1Deserializer,

  // Schema registries
  browser_v1SchemaRegistry,
  presenter_v1SchemaRegistry,
  testOnePackageSchemas,
  testOnePackageModels2Schemas,
  testMultiPackageSchemas,
  utils_v1SchemaRegistry,
  wasmjs_v1SchemaRegistry,
  aggregatedTestOnePackageSchemas,
};
