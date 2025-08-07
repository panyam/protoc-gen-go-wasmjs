// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * BaseMessage provides common fields for all library messages
 */
export interface BaseMessage {
  id: string;
  timestamp: number;
  version: string;
}


/**
 * Metadata provides additional context for requests/responses
 */
export interface Metadata {
  requestId: string;
  userAgent: string;
  headers?: Map<string, string>;
}


/**
 * ErrorInfo provides structured error information
 */
export interface ErrorInfo {
  code: string;
  message: string;
  details: string[];
}

