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

/**
 * Factory result interface
 */
export interface FactoryResult<T> {
  instance: T;
  fullyLoaded: boolean;
}

/**
 * Factory method type for creating instances
 */
export type FactoryMethod<T = any> = (
  parent?: any,
  attributeName?: string,
  attributeKey?: string | number,
  data?: any
) => FactoryResult<T>;

/**
 * Factory interface that deserializer expects
 */
export interface FactoryInterface {
  /**
   * Get factory method for a fully qualified message type
   * This enables cross-package factory delegation
   */
  getFactoryMethod?(messageType: string): FactoryMethod | undefined;
}
