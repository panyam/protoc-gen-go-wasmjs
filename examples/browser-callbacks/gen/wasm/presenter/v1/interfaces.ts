// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated


/**
 * Request to load user data
 */
export interface LoadUserRequest {
  userId: string;
}


/**
 * Response with user data
 */
export interface LoadUserResponse {
  username: string;
  email: string;
  permissions: string[];
  fromCache: boolean;
}


/**
 * Request to update state
 */
export interface StateUpdateRequest {
  action: string;
  params?: Map<string, string>;
}


/**
 * UI update message (streamed)
 */
export interface UIUpdate {
  component: string;
  action: string;
  data?: Map<string, string>;
}


/**
 * Request to save preferences
 */
export interface PreferencesRequest {
  preferences?: Map<string, string>;
}


/**
 * Response from preferences save
 */
export interface PreferencesResponse {
  saved: boolean;
  itemsSaved: number;
}

