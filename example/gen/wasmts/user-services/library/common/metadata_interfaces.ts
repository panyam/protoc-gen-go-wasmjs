// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Metadata provides additional context for requests/responses
 */
export interface Metadata {
  requestId: string;
  userAgent: string;
  headers?: HeadersEntry;
}



export interface HeadersEntry {
  key: string;
  value: string;
}


/**
 * ErrorInfo provides structured error information
 */
export interface ErrorInfo {
  code: string;
  message: string;
  details: string[];
}

