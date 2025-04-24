/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.wnft";

/** QueryClassRequest is the request type for the Query/Class RPC method */
export interface QueryClassAddressRequest {
  /** class_id associated with the nft */
  classId: string;
  address: string;
}

/** QueryClassesResponse is the response type for the Query/Classes RPC method */
export interface QueryClassAddressResponse {
  exists: boolean;
  totalSupply: number;
  nfts: string[];
}

/** QueryNftFilterResponse is the request type for the Query/NftFilter RPC method */
export interface QueryNftFilterRequest {
  owner: string;
  classId: string;
  tokenId: string;
}

/**
 * QueryNftFilterResponse is the response type for the Query/NftFilter RPC
 * method
 */
export interface QueryNftFilterResponse {
  nfts: NftList[];
}

export interface NftList {
  classId: string;
  tokenId: string;
  uri: string;
  owner: string;
}

function createBaseQueryClassAddressRequest(): QueryClassAddressRequest {
  return { classId: "", address: "" };
}

export const QueryClassAddressRequest = {
  encode(message: QueryClassAddressRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.classId !== "") {
      writer.uint32(10).string(message.classId);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClassAddressRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClassAddressRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.classId = reader.string();
          break;
        case 2:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryClassAddressRequest {
    return {
      classId: isSet(object.classId) ? String(object.classId) : "",
      address: isSet(object.address) ? String(object.address) : "",
    };
  },

  toJSON(message: QueryClassAddressRequest): unknown {
    const obj: any = {};
    message.classId !== undefined && (obj.classId = message.classId);
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryClassAddressRequest>, I>>(object: I): QueryClassAddressRequest {
    const message = createBaseQueryClassAddressRequest();
    message.classId = object.classId ?? "";
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryClassAddressResponse(): QueryClassAddressResponse {
  return { exists: false, totalSupply: 0, nfts: [] };
}

export const QueryClassAddressResponse = {
  encode(message: QueryClassAddressResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exists === true) {
      writer.uint32(8).bool(message.exists);
    }
    if (message.totalSupply !== 0) {
      writer.uint32(16).uint64(message.totalSupply);
    }
    for (const v of message.nfts) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClassAddressResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClassAddressResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.exists = reader.bool();
          break;
        case 2:
          message.totalSupply = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.nfts.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryClassAddressResponse {
    return {
      exists: isSet(object.exists) ? Boolean(object.exists) : false,
      totalSupply: isSet(object.totalSupply) ? Number(object.totalSupply) : 0,
      nfts: Array.isArray(object?.nfts) ? object.nfts.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: QueryClassAddressResponse): unknown {
    const obj: any = {};
    message.exists !== undefined && (obj.exists = message.exists);
    message.totalSupply !== undefined && (obj.totalSupply = Math.round(message.totalSupply));
    if (message.nfts) {
      obj.nfts = message.nfts.map((e) => e);
    } else {
      obj.nfts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryClassAddressResponse>, I>>(object: I): QueryClassAddressResponse {
    const message = createBaseQueryClassAddressResponse();
    message.exists = object.exists ?? false;
    message.totalSupply = object.totalSupply ?? 0;
    message.nfts = object.nfts?.map((e) => e) || [];
    return message;
  },
};

function createBaseQueryNftFilterRequest(): QueryNftFilterRequest {
  return { owner: "", classId: "", tokenId: "" };
}

export const QueryNftFilterRequest = {
  encode(message: QueryNftFilterRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }
    if (message.classId !== "") {
      writer.uint32(18).string(message.classId);
    }
    if (message.tokenId !== "") {
      writer.uint32(26).string(message.tokenId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNftFilterRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNftFilterRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;
        case 2:
          message.classId = reader.string();
          break;
        case 3:
          message.tokenId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryNftFilterRequest {
    return {
      owner: isSet(object.owner) ? String(object.owner) : "",
      classId: isSet(object.classId) ? String(object.classId) : "",
      tokenId: isSet(object.tokenId) ? String(object.tokenId) : "",
    };
  },

  toJSON(message: QueryNftFilterRequest): unknown {
    const obj: any = {};
    message.owner !== undefined && (obj.owner = message.owner);
    message.classId !== undefined && (obj.classId = message.classId);
    message.tokenId !== undefined && (obj.tokenId = message.tokenId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryNftFilterRequest>, I>>(object: I): QueryNftFilterRequest {
    const message = createBaseQueryNftFilterRequest();
    message.owner = object.owner ?? "";
    message.classId = object.classId ?? "";
    message.tokenId = object.tokenId ?? "";
    return message;
  },
};

function createBaseQueryNftFilterResponse(): QueryNftFilterResponse {
  return { nfts: [] };
}

export const QueryNftFilterResponse = {
  encode(message: QueryNftFilterResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.nfts) {
      NftList.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNftFilterResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNftFilterResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.nfts.push(NftList.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryNftFilterResponse {
    return { nfts: Array.isArray(object?.nfts) ? object.nfts.map((e: any) => NftList.fromJSON(e)) : [] };
  },

  toJSON(message: QueryNftFilterResponse): unknown {
    const obj: any = {};
    if (message.nfts) {
      obj.nfts = message.nfts.map((e) => e ? NftList.toJSON(e) : undefined);
    } else {
      obj.nfts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryNftFilterResponse>, I>>(object: I): QueryNftFilterResponse {
    const message = createBaseQueryNftFilterResponse();
    message.nfts = object.nfts?.map((e) => NftList.fromPartial(e)) || [];
    return message;
  },
};

function createBaseNftList(): NftList {
  return { classId: "", tokenId: "", uri: "", owner: "" };
}

export const NftList = {
  encode(message: NftList, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.classId !== "") {
      writer.uint32(10).string(message.classId);
    }
    if (message.tokenId !== "") {
      writer.uint32(18).string(message.tokenId);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.owner !== "") {
      writer.uint32(34).string(message.owner);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NftList {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNftList();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.classId = reader.string();
          break;
        case 2:
          message.tokenId = reader.string();
          break;
        case 3:
          message.uri = reader.string();
          break;
        case 4:
          message.owner = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): NftList {
    return {
      classId: isSet(object.classId) ? String(object.classId) : "",
      tokenId: isSet(object.tokenId) ? String(object.tokenId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      owner: isSet(object.owner) ? String(object.owner) : "",
    };
  },

  toJSON(message: NftList): unknown {
    const obj: any = {};
    message.classId !== undefined && (obj.classId = message.classId);
    message.tokenId !== undefined && (obj.tokenId = message.tokenId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.owner !== undefined && (obj.owner = message.owner);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<NftList>, I>>(object: I): NftList {
    const message = createBaseNftList();
    message.classId = object.classId ?? "";
    message.tokenId = object.tokenId ?? "";
    message.uri = object.uri ?? "";
    message.owner = object.owner ?? "";
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Classes queries Address */
  ClassAddress(request: QueryClassAddressRequest): Promise<QueryClassAddressResponse>;
  /** queries nft filter */
  NftFilter(request: QueryNftFilterRequest): Promise<QueryNftFilterResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ClassAddress = this.ClassAddress.bind(this);
    this.NftFilter = this.NftFilter.bind(this);
  }
  ClassAddress(request: QueryClassAddressRequest): Promise<QueryClassAddressResponse> {
    const data = QueryClassAddressRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wnft.Query", "ClassAddress", data);
    return promise.then((data) => QueryClassAddressResponse.decode(new _m0.Reader(data)));
  }

  NftFilter(request: QueryNftFilterRequest): Promise<QueryNftFilterResponse> {
    const data = QueryNftFilterRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wnft.Query", "NftFilter", data);
    return promise.then((data) => QueryNftFilterResponse.decode(new _m0.Reader(data)));
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
