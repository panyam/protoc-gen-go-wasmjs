import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


import { ModelWithOptionalAndDefaults as ModelWithOptionalAndDefaultsInterface, HelperUtilType as HelperUtilTypeInterface, NestedUtilType as NestedUtilTypeInterface, ParentUtilMessage as ParentUtilMessageInterface, ParentUtilMessage_NestedUtilType as ParentUtilMessage_NestedUtilTypeInterface } from "./interfaces";





export class ModelWithOptionalAndDefaults implements ModelWithOptionalAndDefaultsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.ModelWithOptionalAndDefaults";
  readonly __MESSAGE_TYPE = ModelWithOptionalAndDefaults.MESSAGE_TYPE;

  neededValue: number = 0;
  optionalString?: string | undefined;

  
}



export class HelperUtilType implements HelperUtilTypeInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.HelperUtilType";
  readonly __MESSAGE_TYPE = HelperUtilType.MESSAGE_TYPE;

  /** Just a test me */
  value1: number = 0;
  updateMask?: FieldMask;
  createdAt?: Timestamp;

  
}


/**
 * A top levle NestedUtilType
 */
export class NestedUtilType implements NestedUtilTypeInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.NestedUtilType";
  readonly __MESSAGE_TYPE = NestedUtilType.MESSAGE_TYPE;

  topLevelCount: number = 0;
  topLevelValue: string = "";

  
}


/**
 * Parent message containing a nested type
 */
export class ParentUtilMessage implements ParentUtilMessageInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.ParentUtilMessage";
  readonly __MESSAGE_TYPE = ParentUtilMessage.MESSAGE_TYPE;

  parentValue: number = 0;
  nested?: ParentUtilMessage_NestedUtilType;

  
}


/**
 * Nested type to test cross-package nested imports
 */
export class ParentUtilMessage_NestedUtilType implements ParentUtilMessage_NestedUtilTypeInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.ParentUtilMessage.NestedUtilType";
  readonly __MESSAGE_TYPE = ParentUtilMessage_NestedUtilType.MESSAGE_TYPE;

  nestedValue: string = "";
  nestedCount: number = 0;

  
}


