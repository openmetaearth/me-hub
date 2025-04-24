/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { Credential } from "./credential";
import { DidDocument, DidInfo } from "./did";
import { Service } from "./issuer";

export const protobufPackage = "metaearth.did";

export interface QueryDid {
  address: string;
}

export interface QueryDidResponse {
  info: DidInfo | undefined;
}

export interface QueryDidDocument {
  did: string;
}

export interface QueryDidDocumentResponse {
  doc: DidDocument | undefined;
}

export interface QueryService {
  sid: string;
}

export interface QueryServiceResponse {
  service: Service | undefined;
}

export interface QueryCredential {
  did: string;
  sid: string;
}

export interface QueryCredentialResponse {
  credential: Credential | undefined;
}

export interface QueryCredentials {
  sid: string;
  filter: Uint8Array;
  pagination: PageRequest | undefined;
}

export interface QueryCredentialsResponse {
  credentials: Credential[];
  pagination: PageResponse | undefined;
}

export interface QueryDidInfosRequest {
  pagination: PageRequest | undefined;
}

export interface QueryDidInfosResponse {
  infos: DidInfo[];
  pagination: PageResponse | undefined;
}

function createBaseQueryDid(): QueryDid {
  return { address: "" };
}

export const QueryDid = {
  encode(message: QueryDid, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDid {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDid();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDid {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryDid): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDid>, I>>(object: I): QueryDid {
    const message = createBaseQueryDid();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryDidResponse(): QueryDidResponse {
  return { info: undefined };
}

export const QueryDidResponse = {
  encode(message: QueryDidResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      DidInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDidResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDidResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.info = DidInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDidResponse {
    return { info: isSet(object.info) ? DidInfo.fromJSON(object.info) : undefined };
  },

  toJSON(message: QueryDidResponse): unknown {
    const obj: any = {};
    message.info !== undefined && (obj.info = message.info ? DidInfo.toJSON(message.info) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDidResponse>, I>>(object: I): QueryDidResponse {
    const message = createBaseQueryDidResponse();
    message.info = (object.info !== undefined && object.info !== null) ? DidInfo.fromPartial(object.info) : undefined;
    return message;
  },
};

function createBaseQueryDidDocument(): QueryDidDocument {
  return { did: "" };
}

export const QueryDidDocument = {
  encode(message: QueryDidDocument, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.did !== "") {
      writer.uint32(10).string(message.did);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDidDocument {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDidDocument();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.did = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDidDocument {
    return { did: isSet(object.did) ? String(object.did) : "" };
  },

  toJSON(message: QueryDidDocument): unknown {
    const obj: any = {};
    message.did !== undefined && (obj.did = message.did);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDidDocument>, I>>(object: I): QueryDidDocument {
    const message = createBaseQueryDidDocument();
    message.did = object.did ?? "";
    return message;
  },
};

function createBaseQueryDidDocumentResponse(): QueryDidDocumentResponse {
  return { doc: undefined };
}

export const QueryDidDocumentResponse = {
  encode(message: QueryDidDocumentResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.doc !== undefined) {
      DidDocument.encode(message.doc, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDidDocumentResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDidDocumentResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.doc = DidDocument.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDidDocumentResponse {
    return { doc: isSet(object.doc) ? DidDocument.fromJSON(object.doc) : undefined };
  },

  toJSON(message: QueryDidDocumentResponse): unknown {
    const obj: any = {};
    message.doc !== undefined && (obj.doc = message.doc ? DidDocument.toJSON(message.doc) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDidDocumentResponse>, I>>(object: I): QueryDidDocumentResponse {
    const message = createBaseQueryDidDocumentResponse();
    message.doc = (object.doc !== undefined && object.doc !== null) ? DidDocument.fromPartial(object.doc) : undefined;
    return message;
  },
};

function createBaseQueryService(): QueryService {
  return { sid: "" };
}

export const QueryService = {
  encode(message: QueryService, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sid !== "") {
      writer.uint32(10).string(message.sid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryService {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryService();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryService {
    return { sid: isSet(object.sid) ? String(object.sid) : "" };
  },

  toJSON(message: QueryService): unknown {
    const obj: any = {};
    message.sid !== undefined && (obj.sid = message.sid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryService>, I>>(object: I): QueryService {
    const message = createBaseQueryService();
    message.sid = object.sid ?? "";
    return message;
  },
};

function createBaseQueryServiceResponse(): QueryServiceResponse {
  return { service: undefined };
}

export const QueryServiceResponse = {
  encode(message: QueryServiceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.service !== undefined) {
      Service.encode(message.service, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryServiceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryServiceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.service = Service.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryServiceResponse {
    return { service: isSet(object.service) ? Service.fromJSON(object.service) : undefined };
  },

  toJSON(message: QueryServiceResponse): unknown {
    const obj: any = {};
    message.service !== undefined && (obj.service = message.service ? Service.toJSON(message.service) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryServiceResponse>, I>>(object: I): QueryServiceResponse {
    const message = createBaseQueryServiceResponse();
    message.service = (object.service !== undefined && object.service !== null)
      ? Service.fromPartial(object.service)
      : undefined;
    return message;
  },
};

function createBaseQueryCredential(): QueryCredential {
  return { did: "", sid: "" };
}

export const QueryCredential = {
  encode(message: QueryCredential, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.did !== "") {
      writer.uint32(10).string(message.did);
    }
    if (message.sid !== "") {
      writer.uint32(18).string(message.sid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCredential {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCredential();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.did = reader.string();
          break;
        case 2:
          message.sid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCredential {
    return { did: isSet(object.did) ? String(object.did) : "", sid: isSet(object.sid) ? String(object.sid) : "" };
  },

  toJSON(message: QueryCredential): unknown {
    const obj: any = {};
    message.did !== undefined && (obj.did = message.did);
    message.sid !== undefined && (obj.sid = message.sid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryCredential>, I>>(object: I): QueryCredential {
    const message = createBaseQueryCredential();
    message.did = object.did ?? "";
    message.sid = object.sid ?? "";
    return message;
  },
};

function createBaseQueryCredentialResponse(): QueryCredentialResponse {
  return { credential: undefined };
}

export const QueryCredentialResponse = {
  encode(message: QueryCredentialResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.credential !== undefined) {
      Credential.encode(message.credential, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCredentialResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCredentialResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.credential = Credential.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCredentialResponse {
    return { credential: isSet(object.credential) ? Credential.fromJSON(object.credential) : undefined };
  },

  toJSON(message: QueryCredentialResponse): unknown {
    const obj: any = {};
    message.credential !== undefined
      && (obj.credential = message.credential ? Credential.toJSON(message.credential) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryCredentialResponse>, I>>(object: I): QueryCredentialResponse {
    const message = createBaseQueryCredentialResponse();
    message.credential = (object.credential !== undefined && object.credential !== null)
      ? Credential.fromPartial(object.credential)
      : undefined;
    return message;
  },
};

function createBaseQueryCredentials(): QueryCredentials {
  return { sid: "", filter: new Uint8Array(), pagination: undefined };
}

export const QueryCredentials = {
  encode(message: QueryCredentials, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sid !== "") {
      writer.uint32(10).string(message.sid);
    }
    if (message.filter.length !== 0) {
      writer.uint32(18).bytes(message.filter);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCredentials {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCredentials();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sid = reader.string();
          break;
        case 2:
          message.filter = reader.bytes();
          break;
        case 3:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCredentials {
    return {
      sid: isSet(object.sid) ? String(object.sid) : "",
      filter: isSet(object.filter) ? bytesFromBase64(object.filter) : new Uint8Array(),
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryCredentials): unknown {
    const obj: any = {};
    message.sid !== undefined && (obj.sid = message.sid);
    message.filter !== undefined
      && (obj.filter = base64FromBytes(message.filter !== undefined ? message.filter : new Uint8Array()));
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryCredentials>, I>>(object: I): QueryCredentials {
    const message = createBaseQueryCredentials();
    message.sid = object.sid ?? "";
    message.filter = object.filter ?? new Uint8Array();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryCredentialsResponse(): QueryCredentialsResponse {
  return { credentials: [], pagination: undefined };
}

export const QueryCredentialsResponse = {
  encode(message: QueryCredentialsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.credentials) {
      Credential.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCredentialsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCredentialsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.credentials.push(Credential.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryCredentialsResponse {
    return {
      credentials: Array.isArray(object?.credentials) ? object.credentials.map((e: any) => Credential.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryCredentialsResponse): unknown {
    const obj: any = {};
    if (message.credentials) {
      obj.credentials = message.credentials.map((e) => e ? Credential.toJSON(e) : undefined);
    } else {
      obj.credentials = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryCredentialsResponse>, I>>(object: I): QueryCredentialsResponse {
    const message = createBaseQueryCredentialsResponse();
    message.credentials = object.credentials?.map((e) => Credential.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryDidInfosRequest(): QueryDidInfosRequest {
  return { pagination: undefined };
}

export const QueryDidInfosRequest = {
  encode(message: QueryDidInfosRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDidInfosRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDidInfosRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDidInfosRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryDidInfosRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDidInfosRequest>, I>>(object: I): QueryDidInfosRequest {
    const message = createBaseQueryDidInfosRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryDidInfosResponse(): QueryDidInfosResponse {
  return { infos: [], pagination: undefined };
}

export const QueryDidInfosResponse = {
  encode(message: QueryDidInfosResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.infos) {
      DidInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDidInfosResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDidInfosResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.infos.push(DidInfo.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDidInfosResponse {
    return {
      infos: Array.isArray(object?.infos) ? object.infos.map((e: any) => DidInfo.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryDidInfosResponse): unknown {
    const obj: any = {};
    if (message.infos) {
      obj.infos = message.infos.map((e) => e ? DidInfo.toJSON(e) : undefined);
    } else {
      obj.infos = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDidInfosResponse>, I>>(object: I): QueryDidInfosResponse {
    const message = createBaseQueryDidInfosResponse();
    message.infos = object.infos?.map((e) => DidInfo.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

export interface Query {
  /** query user did doc */
  Did(request: QueryDid): Promise<QueryDidResponse>;
  DidDocument(request: QueryDidDocument): Promise<QueryDidDocumentResponse>;
  Service(request: QueryService): Promise<QueryServiceResponse>;
  Credential(request: QueryCredential): Promise<QueryCredentialResponse>;
  Credentials(request: QueryCredentials): Promise<QueryCredentialsResponse>;
  /** Queries a list of DidInfos items. */
  DidInfos(request: QueryDidInfosRequest): Promise<QueryDidInfosResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Did = this.Did.bind(this);
    this.DidDocument = this.DidDocument.bind(this);
    this.Service = this.Service.bind(this);
    this.Credential = this.Credential.bind(this);
    this.Credentials = this.Credentials.bind(this);
    this.DidInfos = this.DidInfos.bind(this);
  }
  Did(request: QueryDid): Promise<QueryDidResponse> {
    const data = QueryDid.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Query", "Did", data);
    return promise.then((data) => QueryDidResponse.decode(new _m0.Reader(data)));
  }

  DidDocument(request: QueryDidDocument): Promise<QueryDidDocumentResponse> {
    const data = QueryDidDocument.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Query", "DidDocument", data);
    return promise.then((data) => QueryDidDocumentResponse.decode(new _m0.Reader(data)));
  }

  Service(request: QueryService): Promise<QueryServiceResponse> {
    const data = QueryService.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Query", "Service", data);
    return promise.then((data) => QueryServiceResponse.decode(new _m0.Reader(data)));
  }

  Credential(request: QueryCredential): Promise<QueryCredentialResponse> {
    const data = QueryCredential.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Query", "Credential", data);
    return promise.then((data) => QueryCredentialResponse.decode(new _m0.Reader(data)));
  }

  Credentials(request: QueryCredentials): Promise<QueryCredentialsResponse> {
    const data = QueryCredentials.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Query", "Credentials", data);
    return promise.then((data) => QueryCredentialsResponse.decode(new _m0.Reader(data)));
  }

  DidInfos(request: QueryDidInfosRequest): Promise<QueryDidInfosResponse> {
    const data = QueryDidInfosRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.did.Query", "DidInfos", data);
    return promise.then((data) => QueryDidInfosResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
