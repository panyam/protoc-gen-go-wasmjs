// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Request to fetch data from a URL
 */
export interface FetchRequest {
  url: string;
  method: string;
  headers?: Map<string, string>;
  body: string;
}


/**
 * Response from fetch
 */
export interface FetchResponse {
  status: number;
  statusText: string;
  headers?: Map<string, string>;
  body: string;
}


/**
 * Request for localStorage key
 */
export interface StorageKeyRequest {
  key: string;
}


/**
 * Response with localStorage value
 */
export interface StorageValueResponse {
  value: string;
  exists: boolean;
}


/**
 * Request to set localStorage
 */
export interface StorageSetRequest {
  key: string;
  value: string;
}


/**
 * Response from localStorage set
 */
export interface StorageSetResponse {
  success: boolean;
}


/**
 * Request for cookie
 */
export interface CookieRequest {
  name: string;
}


/**
 * Response with cookie value
 */
export interface CookieResponse {
  value: string;
  exists: boolean;
}


/**
 * Request to show alert
 */
export interface AlertRequest {
  message: string;
}


/**
 * Response from alert
 */
export interface AlertResponse {
  shown: boolean;
}


/**
 * Request for user prompt
 */
export interface PromptRequest {
  message: string;
  defaultValue: string;
}


/**
 * Response from user prompt
 */
export interface PromptResponse {
  value: string;
  cancelled: boolean;
}


/**
 * Request to log to window
 */
export interface LogRequest {
  message: string;
  level: string;
}


/**
 * Response from log to window
 */
export interface LogResponse {
  logged: boolean;
}

