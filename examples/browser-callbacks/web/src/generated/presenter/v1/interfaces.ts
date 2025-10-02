// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";



/**
 * Request to load user data
 */
export interface LoadUserDataRequest {
  userId: string;
}


/**
 * Response with user data
 */
export interface LoadUserDataResponse {
  username: string;
  email: string;
  permissions: string[];
  fromCache: boolean;
  createdAt?: Timestamp;
}


/**
 * Request to update state
 */
export interface StateUpdateRequest {
  action: string;
  params: Record<string, string>;
  updateMask?: FieldMask;
}


/**
 * UI update message (streamed)
 */
export interface UIUpdate {
  component: string;
  action: string;
  data: Record<string, string>;
}



export interface TestRecord {
  helperRecord?: HelperUtilType;
}


/**
 * Request to save preferences
 */
export interface PreferencesRequest {
  preferences: Record<string, string>;
}


/**
 * Response from preferences save
 */
export interface PreferencesResponse {
  saved: boolean;
  itemsSaved: number;
}


/**
 * Request to run callback demo
 */
export interface CallbackDemoRequest {
  demoName: string;
}


/**
 * Response from callback demo
 */
export interface CallbackDemoResponse {
  collectedInputs: string[];
  completed: boolean;
}

