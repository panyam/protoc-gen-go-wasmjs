import { LoadUserRequest as LoadUserRequestInterface, LoadUserResponse as LoadUserResponseInterface, StateUpdateRequest as StateUpdateRequestInterface, UIUpdate as UIUpdateInterface, PreferencesRequest as PreferencesRequestInterface, PreferencesResponse as PreferencesResponseInterface } from "./interfaces";
import { PresenterV1Deserializer } from "./deserializer";


/**
 * Request to load user data
 */
export class LoadUserRequest implements LoadUserRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.LoadUserRequest";

  userId: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LoadUserRequest instance or null if creation failed
   */
  static from(data: any) {
    return PresenterV1Deserializer.from<LoadUserRequest>(LoadUserRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response with user data
 */
export class LoadUserResponse implements LoadUserResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.LoadUserResponse";

  username: string = "";
  email: string = "";
  permissions: string[] = [];
  fromCache: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LoadUserResponse instance or null if creation failed
   */
  static from(data: any) {
    return PresenterV1Deserializer.from<LoadUserResponse>(LoadUserResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request to update state
 */
export class StateUpdateRequest implements StateUpdateRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.StateUpdateRequest";

  action: string = "";
  params?: Map<string, string>;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized StateUpdateRequest instance or null if creation failed
   */
  static from(data: any) {
    return PresenterV1Deserializer.from<StateUpdateRequest>(StateUpdateRequest.MESSAGE_TYPE, data);
  }
}


/**
 * UI update message (streamed)
 */
export class UIUpdate implements UIUpdateInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.UIUpdate";

  component: string = "";
  action: string = "";
  data?: Map<string, string>;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized UIUpdate instance or null if creation failed
   */
  static from(data: any) {
    return PresenterV1Deserializer.from<UIUpdate>(UIUpdate.MESSAGE_TYPE, data);
  }
}


/**
 * Request to save preferences
 */
export class PreferencesRequest implements PreferencesRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.PreferencesRequest";

  preferences?: Map<string, string>;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PreferencesRequest instance or null if creation failed
   */
  static from(data: any) {
    return PresenterV1Deserializer.from<PreferencesRequest>(PreferencesRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from preferences save
 */
export class PreferencesResponse implements PreferencesResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "presenter.v1.PreferencesResponse";

  saved: boolean = false;
  itemsSaved: number = 0;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PreferencesResponse instance or null if creation failed
   */
  static from(data: any) {
    return PresenterV1Deserializer.from<PreferencesResponse>(PreferencesResponse.MESSAGE_TYPE, data);
  }
}


