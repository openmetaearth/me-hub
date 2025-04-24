/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { Coin } from "../../cosmos/base/v1beta1/coin";
import { Stream } from "./stream";

export const protobufPackage = "metaearth.streamer";

export interface ModuleToDistributeCoinsRequest {
}

export interface ModuleToDistributeCoinsResponse {
  /** Coins that have yet to be distributed */
  coins: Coin[];
}

export interface StreamByIDRequest {
  /** Gague ID being queried */
  id: number;
}

export interface StreamByIDResponse {
  /** Stream that corresponds to provided gague ID */
  stream: Stream | undefined;
}

export interface StreamsRequest {
  /** Pagination defines pagination for the request */
  pagination: PageRequest | undefined;
}

export interface StreamsResponse {
  /** Upcoming and active streams */
  data: Stream[];
  /** Pagination defines pagination for the response */
  pagination: PageResponse | undefined;
}

export interface ActiveStreamsRequest {
  /** Pagination defines pagination for the request */
  pagination: PageRequest | undefined;
}

export interface ActiveStreamsResponse {
  /** Active gagues only */
  data: Stream[];
  /** Pagination defines pagination for the response */
  pagination: PageResponse | undefined;
}

export interface UpcomingStreamsRequest {
  /** Pagination defines pagination for the request */
  pagination: PageRequest | undefined;
}

export interface UpcomingStreamsResponse {
  /** Streams whose distribution is upcoming */
  data: Stream[];
  /** Pagination defines pagination for the response */
  pagination: PageResponse | undefined;
}

function createBaseModuleToDistributeCoinsRequest(): ModuleToDistributeCoinsRequest {
  return {};
}

export const ModuleToDistributeCoinsRequest = {
  encode(_: ModuleToDistributeCoinsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ModuleToDistributeCoinsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseModuleToDistributeCoinsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): ModuleToDistributeCoinsRequest {
    return {};
  },

  toJSON(_: ModuleToDistributeCoinsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ModuleToDistributeCoinsRequest>, I>>(_: I): ModuleToDistributeCoinsRequest {
    const message = createBaseModuleToDistributeCoinsRequest();
    return message;
  },
};

function createBaseModuleToDistributeCoinsResponse(): ModuleToDistributeCoinsResponse {
  return { coins: [] };
}

export const ModuleToDistributeCoinsResponse = {
  encode(message: ModuleToDistributeCoinsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.coins) {
      Coin.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ModuleToDistributeCoinsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseModuleToDistributeCoinsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.coins.push(Coin.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ModuleToDistributeCoinsResponse {
    return { coins: Array.isArray(object?.coins) ? object.coins.map((e: any) => Coin.fromJSON(e)) : [] };
  },

  toJSON(message: ModuleToDistributeCoinsResponse): unknown {
    const obj: any = {};
    if (message.coins) {
      obj.coins = message.coins.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.coins = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ModuleToDistributeCoinsResponse>, I>>(
    object: I,
  ): ModuleToDistributeCoinsResponse {
    const message = createBaseModuleToDistributeCoinsResponse();
    message.coins = object.coins?.map((e) => Coin.fromPartial(e)) || [];
    return message;
  },
};

function createBaseStreamByIDRequest(): StreamByIDRequest {
  return { id: 0 };
}

export const StreamByIDRequest = {
  encode(message: StreamByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamByIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamByIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StreamByIDRequest {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: StreamByIDRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StreamByIDRequest>, I>>(object: I): StreamByIDRequest {
    const message = createBaseStreamByIDRequest();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseStreamByIDResponse(): StreamByIDResponse {
  return { stream: undefined };
}

export const StreamByIDResponse = {
  encode(message: StreamByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stream !== undefined) {
      Stream.encode(message.stream, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamByIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stream = Stream.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StreamByIDResponse {
    return { stream: isSet(object.stream) ? Stream.fromJSON(object.stream) : undefined };
  },

  toJSON(message: StreamByIDResponse): unknown {
    const obj: any = {};
    message.stream !== undefined && (obj.stream = message.stream ? Stream.toJSON(message.stream) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StreamByIDResponse>, I>>(object: I): StreamByIDResponse {
    const message = createBaseStreamByIDResponse();
    message.stream = (object.stream !== undefined && object.stream !== null)
      ? Stream.fromPartial(object.stream)
      : undefined;
    return message;
  },
};

function createBaseStreamsRequest(): StreamsRequest {
  return { pagination: undefined };
}

export const StreamsRequest = {
  encode(message: StreamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamsRequest();
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

  fromJSON(object: any): StreamsRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: StreamsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StreamsRequest>, I>>(object: I): StreamsRequest {
    const message = createBaseStreamsRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseStreamsResponse(): StreamsResponse {
  return { data: [], pagination: undefined };
}

export const StreamsResponse = {
  encode(message: StreamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.data) {
      Stream.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data.push(Stream.decode(reader, reader.uint32()));
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

  fromJSON(object: any): StreamsResponse {
    return {
      data: Array.isArray(object?.data) ? object.data.map((e: any) => Stream.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: StreamsResponse): unknown {
    const obj: any = {};
    if (message.data) {
      obj.data = message.data.map((e) => e ? Stream.toJSON(e) : undefined);
    } else {
      obj.data = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StreamsResponse>, I>>(object: I): StreamsResponse {
    const message = createBaseStreamsResponse();
    message.data = object.data?.map((e) => Stream.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseActiveStreamsRequest(): ActiveStreamsRequest {
  return { pagination: undefined };
}

export const ActiveStreamsRequest = {
  encode(message: ActiveStreamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActiveStreamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActiveStreamsRequest();
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

  fromJSON(object: any): ActiveStreamsRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: ActiveStreamsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ActiveStreamsRequest>, I>>(object: I): ActiveStreamsRequest {
    const message = createBaseActiveStreamsRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseActiveStreamsResponse(): ActiveStreamsResponse {
  return { data: [], pagination: undefined };
}

export const ActiveStreamsResponse = {
  encode(message: ActiveStreamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.data) {
      Stream.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ActiveStreamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseActiveStreamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data.push(Stream.decode(reader, reader.uint32()));
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

  fromJSON(object: any): ActiveStreamsResponse {
    return {
      data: Array.isArray(object?.data) ? object.data.map((e: any) => Stream.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: ActiveStreamsResponse): unknown {
    const obj: any = {};
    if (message.data) {
      obj.data = message.data.map((e) => e ? Stream.toJSON(e) : undefined);
    } else {
      obj.data = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ActiveStreamsResponse>, I>>(object: I): ActiveStreamsResponse {
    const message = createBaseActiveStreamsResponse();
    message.data = object.data?.map((e) => Stream.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseUpcomingStreamsRequest(): UpcomingStreamsRequest {
  return { pagination: undefined };
}

export const UpcomingStreamsRequest = {
  encode(message: UpcomingStreamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpcomingStreamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpcomingStreamsRequest();
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

  fromJSON(object: any): UpcomingStreamsRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: UpcomingStreamsRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpcomingStreamsRequest>, I>>(object: I): UpcomingStreamsRequest {
    const message = createBaseUpcomingStreamsRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseUpcomingStreamsResponse(): UpcomingStreamsResponse {
  return { data: [], pagination: undefined };
}

export const UpcomingStreamsResponse = {
  encode(message: UpcomingStreamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.data) {
      Stream.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpcomingStreamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpcomingStreamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data.push(Stream.decode(reader, reader.uint32()));
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

  fromJSON(object: any): UpcomingStreamsResponse {
    return {
      data: Array.isArray(object?.data) ? object.data.map((e: any) => Stream.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: UpcomingStreamsResponse): unknown {
    const obj: any = {};
    if (message.data) {
      obj.data = message.data.map((e) => e ? Stream.toJSON(e) : undefined);
    } else {
      obj.data = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpcomingStreamsResponse>, I>>(object: I): UpcomingStreamsResponse {
    const message = createBaseUpcomingStreamsResponse();
    message.data = object.data?.map((e) => Stream.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service */
export interface Query {
  /** ModuleToDistributeCoins returns coins that are going to be distributed */
  ModuleToDistributeCoins(request: ModuleToDistributeCoinsRequest): Promise<ModuleToDistributeCoinsResponse>;
  /** StreamByID returns streams by their respective ID */
  StreamByID(request: StreamByIDRequest): Promise<StreamByIDResponse>;
  /** Streams returns both upcoming and active streams */
  Streams(request: StreamsRequest): Promise<StreamsResponse>;
  /** ActiveStreams returns active streams */
  ActiveStreams(request: ActiveStreamsRequest): Promise<ActiveStreamsResponse>;
  /** Returns scheduled streams that have not yet occurred */
  UpcomingStreams(request: UpcomingStreamsRequest): Promise<UpcomingStreamsResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ModuleToDistributeCoins = this.ModuleToDistributeCoins.bind(this);
    this.StreamByID = this.StreamByID.bind(this);
    this.Streams = this.Streams.bind(this);
    this.ActiveStreams = this.ActiveStreams.bind(this);
    this.UpcomingStreams = this.UpcomingStreams.bind(this);
  }
  ModuleToDistributeCoins(request: ModuleToDistributeCoinsRequest): Promise<ModuleToDistributeCoinsResponse> {
    const data = ModuleToDistributeCoinsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.streamer.Query", "ModuleToDistributeCoins", data);
    return promise.then((data) => ModuleToDistributeCoinsResponse.decode(new _m0.Reader(data)));
  }

  StreamByID(request: StreamByIDRequest): Promise<StreamByIDResponse> {
    const data = StreamByIDRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.streamer.Query", "StreamByID", data);
    return promise.then((data) => StreamByIDResponse.decode(new _m0.Reader(data)));
  }

  Streams(request: StreamsRequest): Promise<StreamsResponse> {
    const data = StreamsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.streamer.Query", "Streams", data);
    return promise.then((data) => StreamsResponse.decode(new _m0.Reader(data)));
  }

  ActiveStreams(request: ActiveStreamsRequest): Promise<ActiveStreamsResponse> {
    const data = ActiveStreamsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.streamer.Query", "ActiveStreams", data);
    return promise.then((data) => ActiveStreamsResponse.decode(new _m0.Reader(data)));
  }

  UpcomingStreams(request: UpcomingStreamsRequest): Promise<UpcomingStreamsResponse> {
    const data = UpcomingStreamsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.streamer.Query", "UpcomingStreams", data);
    return promise.then((data) => UpcomingStreamsResponse.decode(new _m0.Reader(data)));
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
