import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


import { ModelWithOptionalAndDefaults as ModelWithOptionalAndDefaultsInterface, HelperUtilType as HelperUtilTypeInterface, NestedUtilType as NestedUtilTypeInterface, ParentUtilMessage as ParentUtilMessageInterface, ParentUtilMessage_NestedUtilType as ParentUtilMessage_NestedUtilTypeInterface } from "./interfaces";





export class ModelWithOptionalAndDefaults implements ModelWithOptionalAndDefaultsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.ModelWithOptionalAndDefaults";

  neededValue: number = 0;
  optionalString?: string | undefined;

  
}



export class HelperUtilType implements HelperUtilTypeInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.HelperUtilType";

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

  nestedValue: string = "";
  nestedCount: number = 0;

  
}


