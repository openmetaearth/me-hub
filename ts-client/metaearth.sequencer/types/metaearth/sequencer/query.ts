/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { OperatingStatus, operatingStatusFromJSON, operatingStatusToJSON } from "./operating_status";
import { Params } from "./params";
import { MsgStoreReplaceProposer, Sequencer } from "./sequencer";

export const protobufPackage = "metaearth.sequencer";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetSequencerRequest {
  sequencerAddress: string;
}

export interface QueryGetSequencerResponse {
  sequencer: Sequencer | undefined;
}

export interface QuerySequencersRequest {
  pagination: PageRequest | undefined;
}

export interface QuerySequencersResponse {
  sequencers: Sequencer[];
  pagination: PageResponse | undefined;
}

export interface QueryGetSequencersByRollappRequest {
  rollappId: string;
}

export interface QueryGetSequencersByRollappResponse {
  sequencers: Sequencer[];
}

export interface QueryGetSequencersByRollappByStatusRequest {
  rollappId: string;
  status: OperatingStatus;
}

export interface QueryGetSequencersByRollappByStatusResponse {
  sequencers: Sequencer[];
}

export interface QueryGetUnConfirmSequencersAddrByRollappRequest {
  rollappId: string;
  blockHeight: number;
}

export interface QueryGetUnConfirmSequencersAddrByRollappResponse {
  newSequencer: Sequencer | undefined;
  startHeight: number;
  unconfirmCacheHeight: number;
}

export interface QueryReplaceProposerInfoRequest {
  rollappId: string;
}

export interface QueryReplaceProposerInfoResponse {
  replaceProposer: MsgStoreReplaceProposer | undefined;
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

function createBaseQueryGetSequencerRequest(): QueryGetSequencerRequest {
  return { sequencerAddress: "" };
}

export const QueryGetSequencerRequest = {
  encode(message: QueryGetSequencerRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sequencerAddress !== "") {
      writer.uint32(10).string(message.sequencerAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSequencerRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSequencerRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequencerAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequencerRequest {
    return { sequencerAddress: isSet(object.sequencerAddress) ? String(object.sequencerAddress) : "" };
  },

  toJSON(message: QueryGetSequencerRequest): unknown {
    const obj: any = {};
    message.sequencerAddress !== undefined && (obj.sequencerAddress = message.sequencerAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSequencerRequest>, I>>(object: I): QueryGetSequencerRequest {
    const message = createBaseQueryGetSequencerRequest();
    message.sequencerAddress = object.sequencerAddress ?? "";
    return message;
  },
};

function createBaseQueryGetSequencerResponse(): QueryGetSequencerResponse {
  return { sequencer: undefined };
}

export const QueryGetSequencerResponse = {
  encode(message: QueryGetSequencerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sequencer !== undefined) {
      Sequencer.encode(message.sequencer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSequencerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSequencerResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequencer = Sequencer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequencerResponse {
    return { sequencer: isSet(object.sequencer) ? Sequencer.fromJSON(object.sequencer) : undefined };
  },

  toJSON(message: QueryGetSequencerResponse): unknown {
    const obj: any = {};
    message.sequencer !== undefined
      && (obj.sequencer = message.sequencer ? Sequencer.toJSON(message.sequencer) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSequencerResponse>, I>>(object: I): QueryGetSequencerResponse {
    const message = createBaseQueryGetSequencerResponse();
    message.sequencer = (object.sequencer !== undefined && object.sequencer !== null)
      ? Sequencer.fromPartial(object.sequencer)
      : undefined;
    return message;
  },
};

function createBaseQuerySequencersRequest(): QuerySequencersRequest {
  return { pagination: undefined };
}

export const QuerySequencersRequest = {
  encode(message: QuerySequencersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySequencersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySequencersRequest();
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

  fromJSON(object: any): QuerySequencersRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QuerySequencersRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QuerySequencersRequest>, I>>(object: I): QuerySequencersRequest {
    const message = createBaseQuerySequencersRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQuerySequencersResponse(): QuerySequencersResponse {
  return { sequencers: [], pagination: undefined };
}

export const QuerySequencersResponse = {
  encode(message: QuerySequencersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.sequencers) {
      Sequencer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySequencersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySequencersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequencers.push(Sequencer.decode(reader, reader.uint32()));
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

  fromJSON(object: any): QuerySequencersResponse {
    return {
      sequencers: Array.isArray(object?.sequencers) ? object.sequencers.map((e: any) => Sequencer.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QuerySequencersResponse): unknown {
    const obj: any = {};
    if (message.sequencers) {
      obj.sequencers = message.sequencers.map((e) => e ? Sequencer.toJSON(e) : undefined);
    } else {
      obj.sequencers = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QuerySequencersResponse>, I>>(object: I): QuerySequencersResponse {
    const message = createBaseQuerySequencersResponse();
    message.sequencers = object.sequencers?.map((e) => Sequencer.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetSequencersByRollappRequest(): QueryGetSequencersByRollappRequest {
  return { rollappId: "" };
}

export const QueryGetSequencersByRollappRequest = {
  encode(message: QueryGetSequencersByRollappRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSequencersByRollappRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSequencersByRollappRequest();
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

  fromJSON(object: any): QueryGetSequencersByRollappRequest {
    return { rollappId: isSet(object.rollappId) ? String(object.rollappId) : "" };
  },

  toJSON(message: QueryGetSequencersByRollappRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSequencersByRollappRequest>, I>>(
    object: I,
  ): QueryGetSequencersByRollappRequest {
    const message = createBaseQueryGetSequencersByRollappRequest();
    message.rollappId = object.rollappId ?? "";
    return message;
  },
};

function createBaseQueryGetSequencersByRollappResponse(): QueryGetSequencersByRollappResponse {
  return { sequencers: [] };
}

export const QueryGetSequencersByRollappResponse = {
  encode(message: QueryGetSequencersByRollappResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.sequencers) {
      Sequencer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSequencersByRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSequencersByRollappResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequencers.push(Sequencer.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequencersByRollappResponse {
    return {
      sequencers: Array.isArray(object?.sequencers) ? object.sequencers.map((e: any) => Sequencer.fromJSON(e)) : [],
    };
  },

  toJSON(message: QueryGetSequencersByRollappResponse): unknown {
    const obj: any = {};
    if (message.sequencers) {
      obj.sequencers = message.sequencers.map((e) => e ? Sequencer.toJSON(e) : undefined);
    } else {
      obj.sequencers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSequencersByRollappResponse>, I>>(
    object: I,
  ): QueryGetSequencersByRollappResponse {
    const message = createBaseQueryGetSequencersByRollappResponse();
    message.sequencers = object.sequencers?.map((e) => Sequencer.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQueryGetSequencersByRollappByStatusRequest(): QueryGetSequencersByRollappByStatusRequest {
  return { rollappId: "", status: 0 };
}

export const QueryGetSequencersByRollappByStatusRequest = {
  encode(message: QueryGetSequencersByRollappByStatusRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.status !== 0) {
      writer.uint32(16).int32(message.status);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSequencersByRollappByStatusRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSequencersByRollappByStatusRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.status = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequencersByRollappByStatusRequest {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      status: isSet(object.status) ? operatingStatusFromJSON(object.status) : 0,
    };
  },

  toJSON(message: QueryGetSequencersByRollappByStatusRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.status !== undefined && (obj.status = operatingStatusToJSON(message.status));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSequencersByRollappByStatusRequest>, I>>(
    object: I,
  ): QueryGetSequencersByRollappByStatusRequest {
    const message = createBaseQueryGetSequencersByRollappByStatusRequest();
    message.rollappId = object.rollappId ?? "";
    message.status = object.status ?? 0;
    return message;
  },
};

function createBaseQueryGetSequencersByRollappByStatusResponse(): QueryGetSequencersByRollappByStatusResponse {
  return { sequencers: [] };
}

export const QueryGetSequencersByRollappByStatusResponse = {
  encode(message: QueryGetSequencersByRollappByStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.sequencers) {
      Sequencer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetSequencersByRollappByStatusResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSequencersByRollappByStatusResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sequencers.push(Sequencer.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetSequencersByRollappByStatusResponse {
    return {
      sequencers: Array.isArray(object?.sequencers) ? object.sequencers.map((e: any) => Sequencer.fromJSON(e)) : [],
    };
  },

  toJSON(message: QueryGetSequencersByRollappByStatusResponse): unknown {
    const obj: any = {};
    if (message.sequencers) {
      obj.sequencers = message.sequencers.map((e) => e ? Sequencer.toJSON(e) : undefined);
    } else {
      obj.sequencers = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetSequencersByRollappByStatusResponse>, I>>(
    object: I,
  ): QueryGetSequencersByRollappByStatusResponse {
    const message = createBaseQueryGetSequencersByRollappByStatusResponse();
    message.sequencers = object.sequencers?.map((e) => Sequencer.fromPartial(e)) || [];
    return message;
  },
};

function createBaseQueryGetUnConfirmSequencersAddrByRollappRequest(): QueryGetUnConfirmSequencersAddrByRollappRequest {
  return { rollappId: "", blockHeight: 0 };
}

export const QueryGetUnConfirmSequencersAddrByRollappRequest = {
  encode(
    message: QueryGetUnConfirmSequencersAddrByRollappRequest,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    if (message.blockHeight !== 0) {
      writer.uint32(16).int64(message.blockHeight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetUnConfirmSequencersAddrByRollappRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetUnConfirmSequencersAddrByRollappRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappId = reader.string();
          break;
        case 2:
          message.blockHeight = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetUnConfirmSequencersAddrByRollappRequest {
    return {
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      blockHeight: isSet(object.blockHeight) ? Number(object.blockHeight) : 0,
    };
  },

  toJSON(message: QueryGetUnConfirmSequencersAddrByRollappRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.blockHeight !== undefined && (obj.blockHeight = Math.round(message.blockHeight));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetUnConfirmSequencersAddrByRollappRequest>, I>>(
    object: I,
  ): QueryGetUnConfirmSequencersAddrByRollappRequest {
    const message = createBaseQueryGetUnConfirmSequencersAddrByRollappRequest();
    message.rollappId = object.rollappId ?? "";
    message.blockHeight = object.blockHeight ?? 0;
    return message;
  },
};

function createBaseQueryGetUnConfirmSequencersAddrByRollappResponse(): QueryGetUnConfirmSequencersAddrByRollappResponse {
  return { newSequencer: undefined, startHeight: 0, unconfirmCacheHeight: 0 };
}

export const QueryGetUnConfirmSequencersAddrByRollappResponse = {
  encode(
    message: QueryGetUnConfirmSequencersAddrByRollappResponse,
    writer: _m0.Writer = _m0.Writer.create(),
  ): _m0.Writer {
    if (message.newSequencer !== undefined) {
      Sequencer.encode(message.newSequencer, writer.uint32(10).fork()).ldelim();
    }
    if (message.startHeight !== 0) {
      writer.uint32(16).int64(message.startHeight);
    }
    if (message.unconfirmCacheHeight !== 0) {
      writer.uint32(24).int64(message.unconfirmCacheHeight);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetUnConfirmSequencersAddrByRollappResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetUnConfirmSequencersAddrByRollappResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.newSequencer = Sequencer.decode(reader, reader.uint32());
          break;
        case 2:
          message.startHeight = longToNumber(reader.int64() as Long);
          break;
        case 3:
          message.unconfirmCacheHeight = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetUnConfirmSequencersAddrByRollappResponse {
    return {
      newSequencer: isSet(object.newSequencer) ? Sequencer.fromJSON(object.newSequencer) : undefined,
      startHeight: isSet(object.startHeight) ? Number(object.startHeight) : 0,
      unconfirmCacheHeight: isSet(object.unconfirmCacheHeight) ? Number(object.unconfirmCacheHeight) : 0,
    };
  },

  toJSON(message: QueryGetUnConfirmSequencersAddrByRollappResponse): unknown {
    const obj: any = {};
    message.newSequencer !== undefined
      && (obj.newSequencer = message.newSequencer ? Sequencer.toJSON(message.newSequencer) : undefined);
    message.startHeight !== undefined && (obj.startHeight = Math.round(message.startHeight));
    message.unconfirmCacheHeight !== undefined && (obj.unconfirmCacheHeight = Math.round(message.unconfirmCacheHeight));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetUnConfirmSequencersAddrByRollappResponse>, I>>(
    object: I,
  ): QueryGetUnConfirmSequencersAddrByRollappResponse {
    const message = createBaseQueryGetUnConfirmSequencersAddrByRollappResponse();
    message.newSequencer = (object.newSequencer !== undefined && object.newSequencer !== null)
      ? Sequencer.fromPartial(object.newSequencer)
      : undefined;
    message.startHeight = object.startHeight ?? 0;
    message.unconfirmCacheHeight = object.unconfirmCacheHeight ?? 0;
    return message;
  },
};

function createBaseQueryReplaceProposerInfoRequest(): QueryReplaceProposerInfoRequest {
  return { rollappId: "" };
}

export const QueryReplaceProposerInfoRequest = {
  encode(message: QueryReplaceProposerInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappId !== "") {
      writer.uint32(10).string(message.rollappId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryReplaceProposerInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryReplaceProposerInfoRequest();
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

  fromJSON(object: any): QueryReplaceProposerInfoRequest {
    return { rollappId: isSet(object.rollappId) ? String(object.rollappId) : "" };
  },

  toJSON(message: QueryReplaceProposerInfoRequest): unknown {
    const obj: any = {};
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryReplaceProposerInfoRequest>, I>>(
    object: I,
  ): QueryReplaceProposerInfoRequest {
    const message = createBaseQueryReplaceProposerInfoRequest();
    message.rollappId = object.rollappId ?? "";
    return message;
  },
};

function createBaseQueryReplaceProposerInfoResponse(): QueryReplaceProposerInfoResponse {
  return { replaceProposer: undefined };
}

export const QueryReplaceProposerInfoResponse = {
  encode(message: QueryReplaceProposerInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.replaceProposer !== undefined) {
      MsgStoreReplaceProposer.encode(message.replaceProposer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryReplaceProposerInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryReplaceProposerInfoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.replaceProposer = MsgStoreReplaceProposer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryReplaceProposerInfoResponse {
    return {
      replaceProposer: isSet(object.replaceProposer)
        ? MsgStoreReplaceProposer.fromJSON(object.replaceProposer)
        : undefined,
    };
  },

  toJSON(message: QueryReplaceProposerInfoResponse): unknown {
    const obj: any = {};
    message.replaceProposer !== undefined && (obj.replaceProposer = message.replaceProposer
      ? MsgStoreReplaceProposer.toJSON(message.replaceProposer)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryReplaceProposerInfoResponse>, I>>(
    object: I,
  ): QueryReplaceProposerInfoResponse {
    const message = createBaseQueryReplaceProposerInfoResponse();
    message.replaceProposer = (object.replaceProposer !== undefined && object.replaceProposer !== null)
      ? MsgStoreReplaceProposer.fromPartial(object.replaceProposer)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a Sequencer by address. */
  Sequencer(request: QueryGetSequencerRequest): Promise<QueryGetSequencerResponse>;
  /** Queries a list of Sequencer items. */
  Sequencers(request: QuerySequencersRequest): Promise<QuerySequencersResponse>;
  /** Queries a SequencersByRollapp by rollappId. */
  SequencersByRollapp(request: QueryGetSequencersByRollappRequest): Promise<QueryGetSequencersByRollappResponse>;
  /** Queries a SequencersByRollappByStatus */
  SequencersByRollappByStatus(
    request: QueryGetSequencersByRollappByStatusRequest,
  ): Promise<QueryGetSequencersByRollappByStatusResponse>;
  UnConfirmSequencerAddressByRollappByStatus(
    request: QueryGetUnConfirmSequencersAddrByRollappRequest,
  ): Promise<QueryGetUnConfirmSequencersAddrByRollappResponse>;
  ReplaceProposerInfo(request: QueryReplaceProposerInfoRequest): Promise<QueryReplaceProposerInfoResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.Sequencer = this.Sequencer.bind(this);
    this.Sequencers = this.Sequencers.bind(this);
    this.SequencersByRollapp = this.SequencersByRollapp.bind(this);
    this.SequencersByRollappByStatus = this.SequencersByRollappByStatus.bind(this);
    this.UnConfirmSequencerAddressByRollappByStatus = this.UnConfirmSequencerAddressByRollappByStatus.bind(this);
    this.ReplaceProposerInfo = this.ReplaceProposerInfo.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  Sequencer(request: QueryGetSequencerRequest): Promise<QueryGetSequencerResponse> {
    const data = QueryGetSequencerRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "Sequencer", data);
    return promise.then((data) => QueryGetSequencerResponse.decode(new _m0.Reader(data)));
  }

  Sequencers(request: QuerySequencersRequest): Promise<QuerySequencersResponse> {
    const data = QuerySequencersRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "Sequencers", data);
    return promise.then((data) => QuerySequencersResponse.decode(new _m0.Reader(data)));
  }

  SequencersByRollapp(request: QueryGetSequencersByRollappRequest): Promise<QueryGetSequencersByRollappResponse> {
    const data = QueryGetSequencersByRollappRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "SequencersByRollapp", data);
    return promise.then((data) => QueryGetSequencersByRollappResponse.decode(new _m0.Reader(data)));
  }

  SequencersByRollappByStatus(
    request: QueryGetSequencersByRollappByStatusRequest,
  ): Promise<QueryGetSequencersByRollappByStatusResponse> {
    const data = QueryGetSequencersByRollappByStatusRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "SequencersByRollappByStatus", data);
    return promise.then((data) => QueryGetSequencersByRollappByStatusResponse.decode(new _m0.Reader(data)));
  }

  UnConfirmSequencerAddressByRollappByStatus(
    request: QueryGetUnConfirmSequencersAddrByRollappRequest,
  ): Promise<QueryGetUnConfirmSequencersAddrByRollappResponse> {
    const data = QueryGetUnConfirmSequencersAddrByRollappRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "UnConfirmSequencerAddressByRollappByStatus", data);
    return promise.then((data) => QueryGetUnConfirmSequencersAddrByRollappResponse.decode(new _m0.Reader(data)));
  }

  ReplaceProposerInfo(request: QueryReplaceProposerInfoRequest): Promise<QueryReplaceProposerInfoResponse> {
    const data = QueryReplaceProposerInfoRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.sequencer.Query", "ReplaceProposerInfo", data);
    return promise.then((data) => QueryReplaceProposerInfoResponse.decode(new _m0.Reader(data)));
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
