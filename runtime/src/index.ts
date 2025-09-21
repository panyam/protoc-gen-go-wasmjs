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

// Browser utilities
export { BrowserServiceManager } from './browser/index.js';

// Schema types
export {
  FieldType,
  type FieldSchema,
  type MessageSchema,
  BaseDeserializer,
  BaseSchemaRegistry,
} from './schema/index.js';

// Client types
export {
  type WASMResponse,
  WasmError,
  WASMServiceClient,
  WASMBundle,
  type WASMBundleConfig,
  ServiceClient,
} from './client/index.js';

// Factory and patch types
export {
  type FactoryResult,
  type FactoryMethod,
  type FactoryInterface,
  PatchOperation,
  type MessagePatch,
  type PatchBatch,
  PatchSource,
  type PatchResponse,
  type ChangeTransport,
} from './types/index.js';
