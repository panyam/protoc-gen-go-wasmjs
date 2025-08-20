// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Test message for streaming
 */
export interface TickRequest {
  count: number;
  intervalMs: number;
  message: string;
}



export interface TickResponse {
  tickNumber: number;
  timestamp: number;
  message: string;
  isFinal: boolean;
}

