/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { TallyResult } from "../../cosmos/gov/v1/gov";

export const protobufPackage = "metaearth.wgov";

export interface QueryMeTallyResultRequest {
  proposalId: number;
}

export interface QueryMeTallyResultResponse {
  tally: TallyResult | undefined;
}

function createBaseQueryMeTallyResultRequest(): QueryMeTallyResultRequest {
  return { proposalId: 0 };
}

export const QueryMeTallyResultRequest = {
  encode(message: QueryMeTallyResultRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.proposalId !== 0) {
      writer.uint32(8).uint64(message.proposalId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMeTallyResultRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMeTallyResultRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.proposalId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryMeTallyResultRequest {
    return { proposalId: isSet(object.proposalId) ? Number(object.proposalId) : 0 };
  },

  toJSON(message: QueryMeTallyResultRequest): unknown {
    const obj: any = {};
    message.proposalId !== undefined && (obj.proposalId = Math.round(message.proposalId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryMeTallyResultRequest>, I>>(object: I): QueryMeTallyResultRequest {
    const message = createBaseQueryMeTallyResultRequest();
    message.proposalId = object.proposalId ?? 0;
    return message;
  },
};

function createBaseQueryMeTallyResultResponse(): QueryMeTallyResultResponse {
  return { tally: undefined };
}

export const QueryMeTallyResultResponse = {
  encode(message: QueryMeTallyResultResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tally !== undefined) {
      TallyResult.encode(message.tally, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMeTallyResultResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMeTallyResultResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tally = TallyResult.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryMeTallyResultResponse {
    return { tally: isSet(object.tally) ? TallyResult.fromJSON(object.tally) : undefined };
  },

  toJSON(message: QueryMeTallyResultResponse): unknown {
    const obj: any = {};
    message.tally !== undefined && (obj.tally = message.tally ? TallyResult.toJSON(message.tally) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryMeTallyResultResponse>, I>>(object: I): QueryMeTallyResultResponse {
    const message = createBaseQueryMeTallyResultResponse();
    message.tally = (object.tally !== undefined && object.tally !== null)
      ? TallyResult.fromPartial(object.tally)
      : undefined;
    return message;
  },
};

export interface Query {
  MeTallyResult(request: QueryMeTallyResultRequest): Promise<QueryMeTallyResultResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.MeTallyResult = this.MeTallyResult.bind(this);
  }
  MeTallyResult(request: QueryMeTallyResultRequest): Promise<QueryMeTallyResultResponse> {
    const data = QueryMeTallyResultRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.wgov.Query", "MeTallyResult", data);
    return promise.then((data) => QueryMeTallyResultResponse.decode(new _m0.Reader(data)));
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
