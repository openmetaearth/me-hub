/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { RollappPacket_Type, rollappPacket_TypeFromJSON, rollappPacket_TypeToJSON } from "../common/rollapp_packet";
import { Status, statusFromJSON, statusToJSON } from "../common/status";
import { DemandOrder } from "./demand_order";
import { Params } from "./params";

export const protobufPackage = "dymensionxyz.dymension.eibc";

export enum FulfillmentState {
  UNDEFINED = 0,
  FULFILLED = 1,
  UNFULFILLED = 2,
  UNRECOGNIZED = -1,
}

export function fulfillmentStateFromJSON(object: any): FulfillmentState {
  switch (object) {
    case 0:
    case "UNDEFINED":
      return FulfillmentState.UNDEFINED;
    case 1:
    case "FULFILLED":
      return FulfillmentState.FULFILLED;
    case 2:
    case "UNFULFILLED":
      return FulfillmentState.UNFULFILLED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FulfillmentState.UNRECOGNIZED;
  }
}

export function fulfillmentStateToJSON(object: FulfillmentState): string {
  switch (object) {
    case FulfillmentState.UNDEFINED:
      return "UNDEFINED";
    case FulfillmentState.FULFILLED:
      return "FULFILLED";
    case FulfillmentState.UNFULFILLED:
      return "UNFULFILLED";
    case FulfillmentState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

/**
 * QueryGetDemandOrderRequest is the request type for the Query/GetDemandOrder
 * RPC method.
 */
export interface QueryGetDemandOrderRequest {
  /** id of the demand order to get */
  id: string;
}

/**
 * QueryDemandOrdersByStatusRequest is the request type for the
 * Query/GetDemandOrdersByStatus RPC method.
 */
export interface QueryDemandOrdersByStatusRequest {
  /** status of the demand order */
  status: Status;
  /** optional type */
  type: RollappPacket_Type;
  /** optional rollapp_id */
  rollappId: string;
  /** optional limit */
  limit: number;
  /** optional fulfillment state */
  fulfillmentState: FulfillmentState;
  /** optional fulfiller address */
  fulfiller: string;
  /** optional denom */
  denom: string;
  /** optional recipient address */
  recipient: string;
}

/**
 * QueryGetDemandOrderResponse is the response type for the Query/GetDemandOrder
 * RPC method.
 */
export interface QueryGetDemandOrderResponse {
  /** demand order with the given id */
  demandOrder: DemandOrder | undefined;
}

/**
 * QueryDemandOrdersByStatusResponse is the response type for the
 * Query/GetDemandOrdersByStatus RPC method.
 */
export interface QueryDemandOrdersByStatusResponse {
  /** A list of demand orders with the given status */
  demandOrders: DemandOrder[];
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

function createBaseQueryGetDemandOrderRequest(): QueryGetDemandOrderRequest {
  return { id: "" };
}

export const QueryGetDemandOrderRequest = {
  encode(message: QueryGetDemandOrderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetDemandOrderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetDemandOrderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDemandOrderRequest {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: QueryGetDemandOrderRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetDemandOrderRequest>, I>>(object: I): QueryGetDemandOrderRequest {
    const message = createBaseQueryGetDemandOrderRequest();
    message.id = object.id ?? "";
    return message;
  },
};

function createBaseQueryDemandOrdersByStatusRequest(): QueryDemandOrdersByStatusRequest {
  return { status: 0, type: 0, rollappId: "", limit: 0, fulfillmentState: 0, fulfiller: "", denom: "", recipient: "" };
}

export const QueryDemandOrdersByStatusRequest = {
  encode(message: QueryDemandOrdersByStatusRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.type !== 0) {
      writer.uint32(16).int32(message.type);
    }
    if (message.rollappId !== "") {
      writer.uint32(26).string(message.rollappId);
    }
    if (message.limit !== 0) {
      writer.uint32(32).int32(message.limit);
    }
    if (message.fulfillmentState !== 0) {
      writer.uint32(40).int32(message.fulfillmentState);
    }
    if (message.fulfiller !== "") {
      writer.uint32(50).string(message.fulfiller);
    }
    if (message.denom !== "") {
      writer.uint32(58).string(message.denom);
    }
    if (message.recipient !== "") {
      writer.uint32(66).string(message.recipient);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDemandOrdersByStatusRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDemandOrdersByStatusRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.type = reader.int32() as any;
          break;
        case 3:
          message.rollappId = reader.string();
          break;
        case 4:
          message.limit = reader.int32();
          break;
        case 5:
          message.fulfillmentState = reader.int32() as any;
          break;
        case 6:
          message.fulfiller = reader.string();
          break;
        case 7:
          message.denom = reader.string();
          break;
        case 8:
          message.recipient = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDemandOrdersByStatusRequest {
    return {
      status: isSet(object.status) ? statusFromJSON(object.status) : 0,
      type: isSet(object.type) ? rollappPacket_TypeFromJSON(object.type) : 0,
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      limit: isSet(object.limit) ? Number(object.limit) : 0,
      fulfillmentState: isSet(object.fulfillmentState) ? fulfillmentStateFromJSON(object.fulfillmentState) : 0,
      fulfiller: isSet(object.fulfiller) ? String(object.fulfiller) : "",
      denom: isSet(object.denom) ? String(object.denom) : "",
      recipient: isSet(object.recipient) ? String(object.recipient) : "",
    };
  },

  toJSON(message: QueryDemandOrdersByStatusRequest): unknown {
    const obj: any = {};
    message.status !== undefined && (obj.status = statusToJSON(message.status));
    message.type !== undefined && (obj.type = rollappPacket_TypeToJSON(message.type));
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.limit !== undefined && (obj.limit = Math.round(message.limit));
    message.fulfillmentState !== undefined && (obj.fulfillmentState = fulfillmentStateToJSON(message.fulfillmentState));
    message.fulfiller !== undefined && (obj.fulfiller = message.fulfiller);
    message.denom !== undefined && (obj.denom = message.denom);
    message.recipient !== undefined && (obj.recipient = message.recipient);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDemandOrdersByStatusRequest>, I>>(
    object: I,
  ): QueryDemandOrdersByStatusRequest {
    const message = createBaseQueryDemandOrdersByStatusRequest();
    message.status = object.status ?? 0;
    message.type = object.type ?? 0;
    message.rollappId = object.rollappId ?? "";
    message.limit = object.limit ?? 0;
    message.fulfillmentState = object.fulfillmentState ?? 0;
    message.fulfiller = object.fulfiller ?? "";
    message.denom = object.denom ?? "";
    message.recipient = object.recipient ?? "";
    return message;
  },
};

function createBaseQueryGetDemandOrderResponse(): QueryGetDemandOrderResponse {
  return { demandOrder: undefined };
}

export const QueryGetDemandOrderResponse = {
  encode(message: QueryGetDemandOrderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.demandOrder !== undefined) {
      DemandOrder.encode(message.demandOrder, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetDemandOrderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetDemandOrderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.demandOrder = DemandOrder.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetDemandOrderResponse {
    return { demandOrder: isSet(object.demandOrder) ? DemandOrder.fromJSON(object.demandOrder) : undefined };
  },

  toJSON(message: QueryGetDemandOrderResponse): unknown {
    const obj: any = {};
    message.demandOrder !== undefined
      && (obj.demandOrder = message.demandOrder ? DemandOrder.toJSON(message.demandOrder) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetDemandOrderResponse>, I>>(object: I): QueryGetDemandOrderResponse {
    const message = createBaseQueryGetDemandOrderResponse();
    message.demandOrder = (object.demandOrder !== undefined && object.demandOrder !== null)
      ? DemandOrder.fromPartial(object.demandOrder)
      : undefined;
    return message;
  },
};

function createBaseQueryDemandOrdersByStatusResponse(): QueryDemandOrdersByStatusResponse {
  return { demandOrders: [] };
}

export const QueryDemandOrdersByStatusResponse = {
  encode(message: QueryDemandOrdersByStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.demandOrders) {
      DemandOrder.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDemandOrdersByStatusResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDemandOrdersByStatusResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.demandOrders.push(DemandOrder.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryDemandOrdersByStatusResponse {
    return {
      demandOrders: Array.isArray(object?.demandOrders)
        ? object.demandOrders.map((e: any) => DemandOrder.fromJSON(e))
        : [],
    };
  },

  toJSON(message: QueryDemandOrdersByStatusResponse): unknown {
    const obj: any = {};
    if (message.demandOrders) {
      obj.demandOrders = message.demandOrders.map((e) => e ? DemandOrder.toJSON(e) : undefined);
    } else {
      obj.demandOrders = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryDemandOrdersByStatusResponse>, I>>(
    object: I,
  ): QueryDemandOrdersByStatusResponse {
    const message = createBaseQueryDemandOrdersByStatusResponse();
    message.demandOrders = object.demandOrders?.map((e) => DemandOrder.fromPartial(e)) || [];
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a Demand Order by id. */
  DemandOrderById(request: QueryGetDemandOrderRequest): Promise<QueryGetDemandOrderResponse>;
  /** Queries a list of demand orders by status. */
  DemandOrdersByStatus(request: QueryDemandOrdersByStatusRequest): Promise<QueryDemandOrdersByStatusResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.DemandOrderById = this.DemandOrderById.bind(this);
    this.DemandOrdersByStatus = this.DemandOrdersByStatus.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.eibc.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  DemandOrderById(request: QueryGetDemandOrderRequest): Promise<QueryGetDemandOrderResponse> {
    const data = QueryGetDemandOrderRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.eibc.Query", "DemandOrderById", data);
    return promise.then((data) => QueryGetDemandOrderResponse.decode(new _m0.Reader(data)));
  }

  DemandOrdersByStatus(request: QueryDemandOrdersByStatusRequest): Promise<QueryDemandOrdersByStatusResponse> {
    const data = QueryDemandOrdersByStatusRequest.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.eibc.Query", "DemandOrdersByStatus", data);
    return promise.then((data) => QueryDemandOrdersByStatusResponse.decode(new _m0.Reader(data)));
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
