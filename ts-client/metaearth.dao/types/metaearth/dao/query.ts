/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { DaoAddresses } from "./dao";

export const protobufPackage = "metaearth.dao";

export interface QueryGlobalDaoRequest {
}

export interface QueryGlobalDaoResponse {
  daoAddresses: DaoAddresses | undefined;
}

export interface QueryGlobalDaoFeePoolReq {
}

export interface QueryGlobalDaoFeePoolResp {
  globalDaoFeePool: string;
}

export interface QueryFreeGasAccountsReq {
  pagination: PageRequest | undefined;
}

export interface QueryFreeGasAccountsResp {
  addresses: string[];
  pagination: PageResponse | undefined;
}

export interface QueryIsFreeGasAccountReq {
  address: string;
}

export interface QueryIsFreeGasAccountResp {
  isFree: boolean;
}

function createBaseQueryGlobalDaoRequest(): QueryGlobalDaoRequest {
  return {};
}

export const QueryGlobalDaoRequest = {
  encode(_: QueryGlobalDaoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGlobalDaoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGlobalDaoRequest();
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

  fromJSON(_: any): QueryGlobalDaoRequest {
    return {};
  },

  toJSON(_: QueryGlobalDaoRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGlobalDaoRequest>, I>>(_: I): QueryGlobalDaoRequest {
    const message = createBaseQueryGlobalDaoRequest();
    return message;
  },
};

function createBaseQueryGlobalDaoResponse(): QueryGlobalDaoResponse {
  return { daoAddresses: undefined };
}

export const QueryGlobalDaoResponse = {
  encode(message: QueryGlobalDaoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.daoAddresses !== undefined) {
      DaoAddresses.encode(message.daoAddresses, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGlobalDaoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGlobalDaoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.daoAddresses = DaoAddresses.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGlobalDaoResponse {
    return { daoAddresses: isSet(object.daoAddresses) ? DaoAddresses.fromJSON(object.daoAddresses) : undefined };
  },

  toJSON(message: QueryGlobalDaoResponse): unknown {
    const obj: any = {};
    message.daoAddresses !== undefined
      && (obj.daoAddresses = message.daoAddresses ? DaoAddresses.toJSON(message.daoAddresses) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGlobalDaoResponse>, I>>(object: I): QueryGlobalDaoResponse {
    const message = createBaseQueryGlobalDaoResponse();
    message.daoAddresses = (object.daoAddresses !== undefined && object.daoAddresses !== null)
      ? DaoAddresses.fromPartial(object.daoAddresses)
      : undefined;
    return message;
  },
};

function createBaseQueryGlobalDaoFeePoolReq(): QueryGlobalDaoFeePoolReq {
  return {};
}

export const QueryGlobalDaoFeePoolReq = {
  encode(_: QueryGlobalDaoFeePoolReq, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGlobalDaoFeePoolReq {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGlobalDaoFeePoolReq();
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

  fromJSON(_: any): QueryGlobalDaoFeePoolReq {
    return {};
  },

  toJSON(_: QueryGlobalDaoFeePoolReq): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGlobalDaoFeePoolReq>, I>>(_: I): QueryGlobalDaoFeePoolReq {
    const message = createBaseQueryGlobalDaoFeePoolReq();
    return message;
  },
};

function createBaseQueryGlobalDaoFeePoolResp(): QueryGlobalDaoFeePoolResp {
  return { globalDaoFeePool: "" };
}

export const QueryGlobalDaoFeePoolResp = {
  encode(message: QueryGlobalDaoFeePoolResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.globalDaoFeePool !== "") {
      writer.uint32(10).string(message.globalDaoFeePool);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGlobalDaoFeePoolResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGlobalDaoFeePoolResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.globalDaoFeePool = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGlobalDaoFeePoolResp {
    return { globalDaoFeePool: isSet(object.globalDaoFeePool) ? String(object.globalDaoFeePool) : "" };
  },

  toJSON(message: QueryGlobalDaoFeePoolResp): unknown {
    const obj: any = {};
    message.globalDaoFeePool !== undefined && (obj.globalDaoFeePool = message.globalDaoFeePool);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGlobalDaoFeePoolResp>, I>>(object: I): QueryGlobalDaoFeePoolResp {
    const message = createBaseQueryGlobalDaoFeePoolResp();
    message.globalDaoFeePool = object.globalDaoFeePool ?? "";
    return message;
  },
};

function createBaseQueryFreeGasAccountsReq(): QueryFreeGasAccountsReq {
  return { pagination: undefined };
}

export const QueryFreeGasAccountsReq = {
  encode(message: QueryFreeGasAccountsReq, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFreeGasAccountsReq {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFreeGasAccountsReq();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFreeGasAccountsReq {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryFreeGasAccountsReq): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFreeGasAccountsReq>, I>>(object: I): QueryFreeGasAccountsReq {
    const message = createBaseQueryFreeGasAccountsReq();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryFreeGasAccountsResp(): QueryFreeGasAccountsResp {
  return { addresses: [], pagination: undefined };
}

export const QueryFreeGasAccountsResp = {
  encode(message: QueryFreeGasAccountsResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.addresses) {
      writer.uint32(10).string(v!);
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFreeGasAccountsResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFreeGasAccountsResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.addresses.push(reader.string());
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

  fromJSON(object: any): QueryFreeGasAccountsResp {
    return {
      addresses: Array.isArray(object?.addresses) ? object.addresses.map((e: any) => String(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryFreeGasAccountsResp): unknown {
    const obj: any = {};
    if (message.addresses) {
      obj.addresses = message.addresses.map((e) => e);
    } else {
      obj.addresses = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFreeGasAccountsResp>, I>>(object: I): QueryFreeGasAccountsResp {
    const message = createBaseQueryFreeGasAccountsResp();
    message.addresses = object.addresses?.map((e) => e) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryIsFreeGasAccountReq(): QueryIsFreeGasAccountReq {
  return { address: "" };
}

export const QueryIsFreeGasAccountReq = {
  encode(message: QueryIsFreeGasAccountReq, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryIsFreeGasAccountReq {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryIsFreeGasAccountReq();
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

  fromJSON(object: any): QueryIsFreeGasAccountReq {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryIsFreeGasAccountReq): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryIsFreeGasAccountReq>, I>>(object: I): QueryIsFreeGasAccountReq {
    const message = createBaseQueryIsFreeGasAccountReq();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryIsFreeGasAccountResp(): QueryIsFreeGasAccountResp {
  return { isFree: false };
}

export const QueryIsFreeGasAccountResp = {
  encode(message: QueryIsFreeGasAccountResp, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isFree === true) {
      writer.uint32(8).bool(message.isFree);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryIsFreeGasAccountResp {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryIsFreeGasAccountResp();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.isFree = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryIsFreeGasAccountResp {
    return { isFree: isSet(object.isFree) ? Boolean(object.isFree) : false };
  },

  toJSON(message: QueryIsFreeGasAccountResp): unknown {
    const obj: any = {};
    message.isFree !== undefined && (obj.isFree = message.isFree);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryIsFreeGasAccountResp>, I>>(object: I): QueryIsFreeGasAccountResp {
    const message = createBaseQueryIsFreeGasAccountResp();
    message.isFree = object.isFree ?? false;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Queries a list of admin items. */
  GlobalDao(request: QueryGlobalDaoRequest): Promise<QueryGlobalDaoResponse>;
  GlobalDaoFeePool(request: QueryGlobalDaoFeePoolReq): Promise<QueryGlobalDaoFeePoolResp>;
  FreeGasAccounts(request: QueryFreeGasAccountsReq): Promise<QueryFreeGasAccountsResp>;
  IsFreeGasAccount(request: QueryIsFreeGasAccountReq): Promise<QueryIsFreeGasAccountResp>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GlobalDao = this.GlobalDao.bind(this);
    this.GlobalDaoFeePool = this.GlobalDaoFeePool.bind(this);
    this.FreeGasAccounts = this.FreeGasAccounts.bind(this);
    this.IsFreeGasAccount = this.IsFreeGasAccount.bind(this);
  }
  GlobalDao(request: QueryGlobalDaoRequest): Promise<QueryGlobalDaoResponse> {
    const data = QueryGlobalDaoRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Query", "GlobalDao", data);
    return promise.then((data) => QueryGlobalDaoResponse.decode(new _m0.Reader(data)));
  }

  GlobalDaoFeePool(request: QueryGlobalDaoFeePoolReq): Promise<QueryGlobalDaoFeePoolResp> {
    const data = QueryGlobalDaoFeePoolReq.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Query", "GlobalDaoFeePool", data);
    return promise.then((data) => QueryGlobalDaoFeePoolResp.decode(new _m0.Reader(data)));
  }

  FreeGasAccounts(request: QueryFreeGasAccountsReq): Promise<QueryFreeGasAccountsResp> {
    const data = QueryFreeGasAccountsReq.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Query", "FreeGasAccounts", data);
    return promise.then((data) => QueryFreeGasAccountsResp.decode(new _m0.Reader(data)));
  }

  IsFreeGasAccount(request: QueryIsFreeGasAccountReq): Promise<QueryIsFreeGasAccountResp> {
    const data = QueryIsFreeGasAccountReq.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Query", "IsFreeGasAccount", data);
    return promise.then((data) => QueryIsFreeGasAccountResp.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
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
