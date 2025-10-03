// Generated TypeScript interfaces from proto file
// DO NOT EDIT - This file is auto-generated

import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";




export interface ModelWithOptionalAndDefaults {
  neededValue: number;
  optionalString?: string | undefined;
}



export interface HelperUtilType {
  /** Just a test me */
  value1: number;
  updateMask?: FieldMask;
  createdAt?: Timestamp;
}


/**
 * A top levle NestedUtilType
 */
export interface NestedUtilType {
  topLevelCount: number;
  topLevelValue: string;
}


/**
 * Parent message containing a nested type
 */
export interface ParentUtilMessage {
  parentValue: number;
  nested?: ParentUtilMessage_NestedUtilType;
}


/**
 * Nested type to test cross-package nested imports
 */
export interface ParentUtilMessage_NestedUtilType {
  nestedValue: string;
  nestedCount: number;
}

