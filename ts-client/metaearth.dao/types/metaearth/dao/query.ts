/* eslint-disable */
import _m0 from "protobufjs/minimal";
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

/** Query defines the gRPC querier service. */
export interface Query {
  /** Queries a list of admin items. */
  GlobalDao(request: QueryGlobalDaoRequest): Promise<QueryGlobalDaoResponse>;
  GlobalDaoFeePool(request: QueryGlobalDaoFeePoolReq): Promise<QueryGlobalDaoFeePoolResp>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GlobalDao = this.GlobalDao.bind(this);
    this.GlobalDaoFeePool = this.GlobalDaoFeePool.bind(this);
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
