import { FieldMask, Timestamp } from "@bufbuild/protobuf/wkt";


import { ModelWithOptionalAndDefaults as ModelWithOptionalAndDefaultsInterface, HelperUtilType as HelperUtilTypeInterface, NestedUtilType as NestedUtilTypeInterface, ParentUtilMessage as ParentUtilMessageInterface, ParentUtilMessage_NestedUtilType as ParentUtilMessage_NestedUtilTypeInterface } from "./interfaces";
import { Utils_v1Deserializer } from "./deserializer";



export class ModelWithOptionalAndDefaults implements ModelWithOptionalAndDefaultsInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "utils.v1.ModelWithOptionalAndDefaults";

  neededValue: number = 0;
  optionalString?: string | undefined;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized ModelWithOptionalAndDefaults instance or null if creation failed
   */
  static from(data: any) {
    return Utils_v1Deserializer.from<ModelWithOptionalAndDefaults>(ModelWithOptionalAndDefaults.MESSAGE_TYPE, data);
  }
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

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized HelperUtilType instance or null if creation failed
   */
  static from(data: any) {
    return Utils_v1Deserializer.from<HelperUtilType>(HelperUtilType.MESSAGE_TYPE, data);
  }
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

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized NestedUtilType instance or null if creation failed
   */
  static from(data: any) {
    return Utils_v1Deserializer.from<NestedUtilType>(NestedUtilType.MESSAGE_TYPE, data);
  }
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

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized ParentUtilMessage instance or null if creation failed
   */
  static from(data: any) {
    return Utils_v1Deserializer.from<ParentUtilMessage>(ParentUtilMessage.MESSAGE_TYPE, data);
  }
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

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized ParentUtilMessage_NestedUtilType instance or null if creation failed
   */
  static from(data: any) {
    return Utils_v1Deserializer.from<ParentUtilMessage_NestedUtilType>(ParentUtilMessage_NestedUtilType.MESSAGE_TYPE, data);
  }
}


