import { FetchRequest as FetchRequestInterface, FetchResponse as FetchResponseInterface, StorageKeyRequest as StorageKeyRequestInterface, StorageValueResponse as StorageValueResponseInterface, StorageSetRequest as StorageSetRequestInterface, StorageSetResponse as StorageSetResponseInterface, CookieRequest as CookieRequestInterface, CookieResponse as CookieResponseInterface, AlertRequest as AlertRequestInterface, AlertResponse as AlertResponseInterface, PromptRequest as PromptRequestInterface, PromptResponse as PromptResponseInterface, LogRequest as LogRequestInterface, LogResponse as LogResponseInterface } from "./interfaces";
import { Browser_v1Deserializer } from "./deserializer";


/**
 * Request to fetch data from a URL
 */
export class FetchRequest implements FetchRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.FetchRequest";

  url: string = "";
  method: string = "";
  headers: Record<string, string> = {};
  body: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized FetchRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<FetchRequest>(FetchRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from fetch
 */
export class FetchResponse implements FetchResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.FetchResponse";

  status: number = 0;
  statusText: string = "";
  headers: Record<string, string> = {};
  body: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized FetchResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<FetchResponse>(FetchResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request for localStorage key
 */
export class StorageKeyRequest implements StorageKeyRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageKeyRequest";

  key: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized StorageKeyRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<StorageKeyRequest>(StorageKeyRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response with localStorage value
 */
export class StorageValueResponse implements StorageValueResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageValueResponse";

  value: string = "";
  exists: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized StorageValueResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<StorageValueResponse>(StorageValueResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request to set localStorage
 */
export class StorageSetRequest implements StorageSetRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageSetRequest";

  key: string = "";
  value: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized StorageSetRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<StorageSetRequest>(StorageSetRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from localStorage set
 */
export class StorageSetResponse implements StorageSetResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.StorageSetResponse";

  success: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized StorageSetResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<StorageSetResponse>(StorageSetResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request for cookie
 */
export class CookieRequest implements CookieRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.CookieRequest";

  name: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CookieRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<CookieRequest>(CookieRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response with cookie value
 */
export class CookieResponse implements CookieResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.CookieResponse";

  value: string = "";
  exists: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized CookieResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<CookieResponse>(CookieResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request to show alert
 */
export class AlertRequest implements AlertRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.AlertRequest";

  message: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized AlertRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<AlertRequest>(AlertRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from alert
 */
export class AlertResponse implements AlertResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.AlertResponse";

  shown: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized AlertResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<AlertResponse>(AlertResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request for user prompt
 */
export class PromptRequest implements PromptRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.PromptRequest";

  message: string = "";
  defaultValue: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PromptRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<PromptRequest>(PromptRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from user prompt
 */
export class PromptResponse implements PromptResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.PromptResponse";

  value: string = "";
  cancelled: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized PromptResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<PromptResponse>(PromptResponse.MESSAGE_TYPE, data);
  }
}


/**
 * Request to log to window
 */
export class LogRequest implements LogRequestInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.LogRequest";

  message: string = "";
  level: string = "";

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LogRequest instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<LogRequest>(LogRequest.MESSAGE_TYPE, data);
  }
}


/**
 * Response from log to window
 */
export class LogResponse implements LogResponseInterface {
  /**
   * Fully qualified message type for schema resolution
   */
  static readonly MESSAGE_TYPE = "browser.v1.LogResponse";

  logged: boolean = false;

  /**
   * Create and deserialize an instance from raw data
   * @param data Raw data to deserialize
   * @returns Deserialized LogResponse instance or null if creation failed
   */
  static from(data: any) {
    return Browser_v1Deserializer.from<LogResponse>(LogResponse.MESSAGE_TYPE, data);
  }
}


