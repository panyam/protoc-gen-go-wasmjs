import { BaseMessage as BaseMessageInterface } from "./base_interfaces";


/**
 * BaseMessage provides common fields for all library messages
 */
export class BaseMessage implements BaseMessageInterface {
  id: string = "";
  timestamp: number = 0;
  version: string = "";
}

