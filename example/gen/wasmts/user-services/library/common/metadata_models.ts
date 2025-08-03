import { Metadata as MetadataInterface, HeadersEntry as HeadersEntryInterface, ErrorInfo as ErrorInfoInterface } from "./metadata_interfaces";


/**
 * Metadata provides additional context for requests/responses
 */
export class Metadata implements MetadataInterface {
  requestId: string = "";
  userAgent: string = "";
  headers?: HeadersEntry;
}



export class HeadersEntry implements HeadersEntryInterface {
  key: string = "";
  value: string = "";
}


/**
 * ErrorInfo provides structured error information
 */
export class ErrorInfo implements ErrorInfoInterface {
  code: string = "";
  message: string = "";
  details: string[] = [];
}

