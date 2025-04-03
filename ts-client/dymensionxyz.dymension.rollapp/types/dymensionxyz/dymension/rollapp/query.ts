/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../../cosmos/base/query/v1beta1/pagination";
import { Params } from "./params";
import { Rollapp, RollappSummary } from "./rollapp";
import { StateInfo, StateInfoIndex } from "./state_info";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetRollappRequest {
  rollappId: string;
}

export interface QueryGetRollappByEIP155Request {
  eip155: number;
}

export interface QueryGetLatestHeightRequest {
  rollappId: string;
  finalized: boolean;
}

export interface QueryGetLatestHeightResponse {
  height: number;
}

export interface QueryGetLatestStateIndexRequest {
  rollappId: string;
  finalized: boolean;
}

export interface QueryGetLatestStateIndexResponse {
  stateIndex: StateInfoIndex | undefined;
}

export interface QueryGetRollappResponse {
  rollapp:
    | Rollapp
    | undefined;
  /** Defines the index of the last rollapp UpdateState. */
  latestStateIndex:
    | StateInfoIndex
    | undefined;
  /** Defines the index of the last rollapp UpdateState that was finalized. */
  latestFinalizedStateIndex: StateInfoIndex | undefined;
  latestHeight: number;
  latestFinalizedHeight: number;
}

export interface QueryAllRollappRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllRollappResponse {
  rollapp: RollappSummary[];
  pagination: PageResponse | undefined;
}

export interface QueryGetStateInfoRequest {
  rollappId: string;
  index: number;
  height: number;
  finalized: boolean;
}

export interface QueryGetStateInfoResponse {
  stateInfo: StateInfo | undefined;
}

export interface QuerySkipDelayRollappRequest {
}

export interface QuerySkipDelayRollappResponse {
  rollapps: string[];
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
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

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryGetRollappRequest(): QueryGetRollappRequest {
  return { rollappId: "" };
}

export const QueryGetRollappRequest = {
  encode(message: QueryGetRollappRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRollappRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRollappRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetRollappRequest {
    return { rollappId: isSet(object.rollappId) ? String(object.rollappId) : "" };
  },

  toJSON(message: QueryGetRollappRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRollappRequest>, I>>(object: I): QueryGetRollappRequest {
    const message = createBaseQueryGetRollappRequest();
    message.rollappId = object.rollappId ?? "";
    return message;
  },
};

function createBaseQueryGetRollappByEIP155Request(): QueryGetRollappByEIP155Request {
  return { eip155: 0 };
}

export const QueryGetRollappByEIP155Request = {
  encode(message: QueryGetRollappByEIP155Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.eip155 !== 0) {
      writer.uint32(8).uint64(message.eip155);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRollappByEIP155Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRollappByEIP155Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.eip155 = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetRollappByEIP155Request {
    return { eip155: isSet(object.eip155) ? Number(object.eip155) : 0 };
  },

  toJSON(message: QueryGetRollappByEIP155Request): unknown {
    const obj: any = {};
    message.eip155 !== undefined && (obj.eip155 = Math.round(message.eip155));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRollappByEIP155Request>, I>>(
    object: I,
  ): QueryGetRollappByEIP155Request {
    const message = createBaseQueryGetRollappByEIP155Request();
    message.eip155 = object.eip155 ?? 0;
    return message;
  },
};

function createBaseQueryGetLatestHeightRequest(): QueryGetLatestHeightRequest {
  return { rollappId: "", finalized: false };
}

export const QueryGetLatestHeightRequest = {
  encode(message: QueryGetLatestHeightRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.finalized === true) {
      writer.uint32(16).bool(message.finalized);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetLatestHeightRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetLatestHeightRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.finalized = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetLatestHeightRequest {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      finalized: isSet(object.finalized) ? Boolean(object.finalized) : false,
    };
  },

  toJSON(message: QueryGetLatestHeightRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.finalized !== undefined && (obj.finalized = message.finalized);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetLatestHeightRequest>, I>>(object: I): QueryGetLatestHeightRequest {
    const message = createBaseQueryGetLatestHeightRequest();
    message.rollappId = object.rollappId ?? "";
    message.finalized = object.finalized ?? false;
    return message;
  },
};

function createBaseQueryGetLatestHeightResponse(): QueryGetLatestHeightResponse {
  return { height: 0 };
}

export const QueryGetLatestHeightResponse = {
  encode(message: QueryGetLatestHeightResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.height !== 0) {
      writer.uint32(8).uint64(message.height);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetLatestHeightResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetLatestHeightResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetLatestHeightResponse {
    return { height: isSet(object.height) ? Number(object.height) : 0 };
  },

  toJSON(message: QueryGetLatestHeightResponse): unknown {
    const obj: any = {};
    message.height !== undefined && (obj.height = Math.round(message.height));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetLatestHeightResponse>, I>>(object: I): QueryGetLatestHeightResponse {
    const message = createBaseQueryGetLatestHeightResponse();
    message.height = object.height ?? 0;
    return message;
  },
};

function createBaseQueryGetLatestStateIndexRequest(): QueryGetLatestStateIndexRequest {
  return { rollappId: "", finalized: false };
}

export const QueryGetLatestStateIndexRequest = {
  encode(message: QueryGetLatestStateIndexRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.finalized === true) {
      writer.uint32(16).bool(message.finalized);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetLatestStateIndexRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetLatestStateIndexRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.finalized = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetLatestStateIndexRequest {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      finalized: isSet(object.finalized) ? Boolean(object.finalized) : false,
    };
  },

  toJSON(message: QueryGetLatestStateIndexRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.finalized !== undefined && (obj.finalized = message.finalized);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetLatestStateIndexRequest>, I>>(
    object: I,
  ): QueryGetLatestStateIndexRequest {
    const message = createBaseQueryGetLatestStateIndexRequest();
    message.rollappId = object.rollappId ?? "";
    message.finalized = object.finalized ?? false;
    return message;
  },
};

function createBaseQueryGetLatestStateIndexResponse(): QueryGetLatestStateIndexResponse {
  return { stateIndex: undefined };
}

export const QueryGetLatestStateIndexResponse = {
  encode(message: QueryGetLatestStateIndexResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stateIndex !== undefined) {
      StateInfoIndex.encode(message.stateIndex, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetLatestStateIndexResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetLatestStateIndexResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stateIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetLatestStateIndexResponse {
    return { stateIndex: isSet(object.stateIndex) ? StateInfoIndex.fromJSON(object.stateIndex) : undefined };
  },

  toJSON(message: QueryGetLatestStateIndexResponse): unknown {
    const obj: any = {};
    message.stateIndex !== undefined
      && (obj.stateIndex = message.stateIndex ? StateInfoIndex.toJSON(message.stateIndex) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetLatestStateIndexResponse>, I>>(
    object: I,
  ): QueryGetLatestStateIndexResponse {
    const message = createBaseQueryGetLatestStateIndexResponse();
    message.stateIndex = (object.stateIndex !== undefined && object.stateIndex !== null)
      ? StateInfoIndex.fromPartial(object.stateIndex)
      : undefined;
    return message;
  },
};

function createBaseQueryGetRollappResponse(): QueryGetRollappResponse {
  return {
    rollapp: undefined,
    latestStateIndex: undefined,
    latestFinalizedStateIndex: undefined,
    latestHeight: 0,
    latestFinalizedHeight: 0,
  };
}

export const QueryGetRollappResponse = {
  encode(message: QueryGetRollappResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollapp !== undefined) {
      Rollapp.encode(message.rollapp, writer.uint32(10).fork()).ldelim();
    }
    if (message.latestStateIndex !== undefined) {
      StateInfoIndex.encode(message.latestStateIndex, writer.uint32(18).fork()).ldelim();
    }
    if (message.latestFinalizedStateIndex !== undefined) {
      StateInfoIndex.encode(message.latestFinalizedStateIndex, writer.uint32(26).fork()).ldelim();
    }
    if (message.latestHeight !== 0) {
      writer.uint32(32).uint64(message.latestHeight);
    }
    if (message.latestFinalizedHeight !== 0) {
      writer.uint32(40).uint64(message.latestFinalizedHeight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetRollappResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollapp = Rollapp.decode(reader, reader.uint32());
          break;
        case 2:
          message.latestStateIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        case 3:
          message.latestFinalizedStateIndex = StateInfoIndex.decode(reader, reader.uint32());
          break;
        case 4:
          message.latestHeight = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.latestFinalizedHeight = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetRollappResponse {
    return {
      rollapp: isSet(object.rollapp) ? Rollapp.fromJSON(object.rollapp) : undefined,
      latestStateIndex: isSet(object.latestStateIndex) ? StateInfoIndex.fromJSON(object.latestStateIndex) : undefined,
      latestFinalizedStateIndex: isSet(object.latestFinalizedStateIndex)
        ? StateInfoIndex.fromJSON(object.latestFinalizedStateIndex)
        : undefined,
      latestHeight: isSet(object.latestHeight) ? Number(object.latestHeight) : 0,
      latestFinalizedHeight: isSet(object.latestFinalizedHeight) ? Number(object.latestFinalizedHeight) : 0,
    };
  },

  toJSON(message: QueryGetRollappResponse): unknown {
    const obj: any = {};
    message.rollapp !== undefined && (obj.rollapp = message.rollapp ? Rollapp.toJSON(message.rollapp) : undefined);
    message.latestStateIndex !== undefined
      && (obj.latestStateIndex = message.latestStateIndex
        ? StateInfoIndex.toJSON(message.latestStateIndex)
        : undefined);
    message.latestFinalizedStateIndex !== undefined
      && (obj.latestFinalizedStateIndex = message.latestFinalizedStateIndex
        ? StateInfoIndex.toJSON(message.latestFinalizedStateIndex)
        : undefined);
    message.latestHeight !== undefined && (obj.latestHeight = Math.round(message.latestHeight));
    message.latestFinalizedHeight !== undefined
      && (obj.latestFinalizedHeight = Math.round(message.latestFinalizedHeight));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetRollappResponse>, I>>(object: I): QueryGetRollappResponse {
    const message = createBaseQueryGetRollappResponse();
    message.rollapp = (object.rollapp !== undefined && object.rollapp !== null)
      ? Rollapp.fromPartial(object.rollapp)
      : undefined;
    message.latestStateIndex = (object.latestStateIndex !== undefined && object.latestStateIndex !== null)
      ? StateInfoIndex.fromPartial(object.latestStateIndex)
      : undefined;
    message.latestFinalizedStateIndex =
      (object.latestFinalizedStateIndex !== undefined && object.latestFinalizedStateIndex !== null)
        ? StateInfoIndex.fromPartial(object.latestFinalizedStateIndex)
        : undefined;
    message.latestHeight = object.latestHeight ?? 0;
    message.latestFinalizedHeight = object.latestFinalizedHeight ?? 0;
    return message;
  },
};

function createBaseQueryAllRollappRequest(): QueryAllRollappRequest {
  return { pagination: undefined };
}

export const QueryAllRollappRequest = {
  encode(message: QueryAllRollappRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRollappRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRollappRequest();
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

  fromJSON(object: any): QueryAllRollappRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllRollappRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRollappRequest>, I>>(object: I): QueryAllRollappRequest {
    const message = createBaseQueryAllRollappRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllRollappResponse(): QueryAllRollappResponse {
  return { rollapp: [], pagination: undefined };
}

export const QueryAllRollappResponse = {
  encode(message: QueryAllRollappResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.rollapp) {
      RollappSummary.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllRollappResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollapp.push(RollappSummary.decode(reader, reader.uint32()));
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

  fromJSON(object: any): QueryAllRollappResponse {
    return {
      rollapp: Array.isArray(object?.rollapp) ? object.rollapp.map((e: any) => RollappSummary.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllRollappResponse): unknown {
    const obj: any = {};
    if (message.rollapp) {
      obj.rollapp = message.rollapp.map((e) => e ? RollappSummary.toJSON(e) : undefined);
    } else {
      obj.rollapp = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllRollappResponse>, I>>(object: I): QueryAllRollappResponse {
    const message = createBaseQueryAllRollappResponse();
    message.rollapp = object.rollapp?.map((e) => RollappSummary.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetStateInfoRequest(): QueryGetStateInfoRequest {
  return { rollappId: "", index: 0, height: 0, finalized: false };
}

export const QueryGetStateInfoRequest = {
  encode(message: QueryGetStateInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.index !== 0) {
      writer.uint32(16).uint64(message.index);
    }
    if (message.height !== 0) {
      writer.uint32(24).uint64(message.height);
    }
    if (message.finalized === true) {
      writer.uint32(32).bool(message.finalized);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetStateInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetStateInfoRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.index = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.height = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.finalized = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetStateInfoRequest {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      index: isSet(object.index) ? Number(object.index) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
      finalized: isSet(object.finalized) ? Boolean(object.finalized) : false,
    };
  },

  toJSON(message: QueryGetStateInfoRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.index !== undefined && (obj.index = Math.round(message.index));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.finalized !== undefined && (obj.finalized = message.finalized);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetStateInfoRequest>, I>>(object: I): QueryGetStateInfoRequest {
    const message = createBaseQueryGetStateInfoRequest();
    message.rollappId = object.rollappId ?? "";
    message.index = object.index ?? 0;
    message.height = object.height ?? 0;
    message.finalized = object.finalized ?? false;
    return message;
  },
};

function createBaseQueryGetStateInfoResponse(): QueryGetStateInfoResponse {
  return { stateInfo: undefined };
}

export const QueryGetStateInfoResponse = {
  encode(message: QueryGetStateInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stateInfo !== undefined) {
      StateInfo.encode(message.stateInfo, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetStateInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetStateInfoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stateInfo = StateInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetStateInfoResponse {
    return { stateInfo: isSet(object.stateInfo) ? StateInfo.fromJSON(object.stateInfo) : undefined };
  },

  toJSON(message: QueryGetStateInfoResponse): unknown {
    const obj: any = {};
    message.stateInfo !== undefined
      && (obj.stateInfo = message.stateInfo ? StateInfo.toJSON(message.stateInfo) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetStateInfoResponse>, I>>(object: I): QueryGetStateInfoResponse {
    const message = createBaseQueryGetStateInfoResponse();
    message.stateInfo = (object.stateInfo !== undefined && object.stateInfo !== null)
      ? StateInfo.fromPartial(object.stateInfo)
      : undefined;
    return message;
  },
};

function createBaseQuerySkipDelayRollappRequest(): QuerySkipDelayRollappRequest {
  return {};
}

export const QuerySkipDelayRollappRequest = {
  encode(_: QuerySkipDelayRollappRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySkipDelayRollappRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySkipDelayRollappRequest();
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

  fromJSON(_: any): QuerySkipDelayRollappRequest {
    return {};
  },

  toJSON(_: QuerySkipDelayRollappRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QuerySkipDelayRollappRequest>, I>>(_: I): QuerySkipDelayRollappRequest {
    const message = createBaseQuerySkipDelayRollappRequest();
    return message;
  },
};

function createBaseQuerySkipDelayRollappResponse(): QuerySkipDelayRollappResponse {
  return { rollapps: [] };
}

export const QuerySkipDelayRollappResponse = {
  encode(message: QuerySkipDelayRollappResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.rollapps) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySkipDelayRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySkipDelayRollappResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollapps.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QuerySkipDelayRollappResponse {
    return { rollapps: Array.isArray(object?.rollapps) ? object.rollapps.map((e: any) => String(e)) : [] };
  },

  toJSON(message: QuerySkipDelayRollappResponse): unknown {
    const obj: any = {};
    if (message.rollapps) {
      obj.rollapps = message.rollapps.map((e) => e);
    } else {
      obj.rollapps = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QuerySkipDelayRollappResponse>, I>>(
    object: I,
  ): QuerySkipDelayRollappResponse {
    const message = createBaseQuerySkipDelayRollappResponse();
    message.rollapps = object.rollapps?.map((e) => e) || [];
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a Rollapp by index. */
  Rollapp(request: QueryGetRollappRequest): Promise<QueryGetRollappResponse>;
  /** Queries a Rollapp by index. */
  RollappByEIP155(request: QueryGetRollappByEIP155Request): Promise<QueryGetRollappResponse>;
  /** Queries a list of Rollapp items. */
  RollappAll(request: QueryAllRollappRequest): Promise<QueryAllRollappResponse>;
  /** Queries a LatestHeight by rollapp-id. */
  LatestHeight(request: QueryGetLatestHeightRequest): Promise<QueryGetLatestHeightResponse>;
  /** Queries a LatestStateIndex by rollapp-id. */
  LatestStateIndex(request: QueryGetLatestStateIndexRequest): Promise<QueryGetLatestStateIndexResponse>;
  /** Queries a StateInfo by index. */
  StateInfo(request: QueryGetStateInfoRequest): Promise<QueryGetStateInfoResponse>;
  SkipDelayRollapp(request: QuerySkipDelayRollappRequest): Promise<QuerySkipDelayRollappResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.Rollapp = this.Rollapp.bind(this);
    this.RollappByEIP155 = this.RollappByEIP155.bind(this);
    this.RollappAll = this.RollappAll.bind(this);
    this.LatestHeight = this.LatestHeight.bind(this);
    this.LatestStateIndex = this.LatestStateIndex.bind(this);
    this.StateInfo = this.StateInfo.bind(this);
    this.SkipDelayRollapp = this.SkipDelayRollapp.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  Rollapp(request: QueryGetRollappRequest): Promise<QueryGetRollappResponse> {
    const data = QueryGetRollappRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "Rollapp", data);
    return promise.then((data) => QueryGetRollappResponse.decode(new _m0.Reader(data)));
  }

  RollappByEIP155(request: QueryGetRollappByEIP155Request): Promise<QueryGetRollappResponse> {
    const data = QueryGetRollappByEIP155Request.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "RollappByEIP155", data);
    return promise.then((data) => QueryGetRollappResponse.decode(new _m0.Reader(data)));
  }

  RollappAll(request: QueryAllRollappRequest): Promise<QueryAllRollappResponse> {
    const data = QueryAllRollappRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "RollappAll", data);
    return promise.then((data) => QueryAllRollappResponse.decode(new _m0.Reader(data)));
  }

  LatestHeight(request: QueryGetLatestHeightRequest): Promise<QueryGetLatestHeightResponse> {
    const data = QueryGetLatestHeightRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "LatestHeight", data);
    return promise.then((data) => QueryGetLatestHeightResponse.decode(new _m0.Reader(data)));
  }

  LatestStateIndex(request: QueryGetLatestStateIndexRequest): Promise<QueryGetLatestStateIndexResponse> {
    const data = QueryGetLatestStateIndexRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "LatestStateIndex", data);
    return promise.then((data) => QueryGetLatestStateIndexResponse.decode(new _m0.Reader(data)));
  }

  StateInfo(request: QueryGetStateInfoRequest): Promise<QueryGetStateInfoResponse> {
    const data = QueryGetStateInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "StateInfo", data);
    return promise.then((data) => QueryGetStateInfoResponse.decode(new _m0.Reader(data)));
  }

  SkipDelayRollapp(request: QuerySkipDelayRollappRequest): Promise<QuerySkipDelayRollappResponse> {
    const data = QuerySkipDelayRollappRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.rollapp.Query", "SkipDelayRollapp", data);
    return promise.then((data) => QuerySkipDelayRollappResponse.decode(new _m0.Reader(data)));
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
