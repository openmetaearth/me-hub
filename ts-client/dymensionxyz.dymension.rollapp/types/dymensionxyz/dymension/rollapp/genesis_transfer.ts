/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

/**
 * Bookkeeping for the genesis transfer bridge protocol.
 * Each rollapp will have one of these items corresponding to it.
 */
export interface GenesisTransfers {
  rollappID: string;
  /**
   * The total number of incoming ibc transfers to be fast tracked in the
   * genesis transfer period
   */
  numTotal: number;
  /**
   * The number of transfers already processed, when this number reaches
   * numTotal the genesis transfer window closes.
   */
  numReceived: number;
}

function createBaseGenesisTransfers(): GenesisTransfers {
  return { rollappID: "", numTotal: 0, numReceived: 0 };
}

export const GenesisTransfers = {
  encode(message: GenesisTransfers, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rollappID !== "") {
      writer.uint32(10).string(message.rollappID);
    }
    if (message.numTotal !== 0) {
      writer.uint32(16).uint64(message.numTotal);
    }
    if (message.numReceived !== 0) {
      writer.uint32(24).uint64(message.numReceived);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisTransfers {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisTransfers();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.rollappID = reader.string();
          break;
        case 2:
          message.numTotal = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.numReceived = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisTransfers {
    return {
      rollappID: isSet(object.rollappID) ? String(object.rollappID) : "",
      numTotal: isSet(object.numTotal) ? Number(object.numTotal) : 0,
      numReceived: isSet(object.numReceived) ? Number(object.numReceived) : 0,
    };
  },

  toJSON(message: GenesisTransfers): unknown {
    const obj: any = {};
    message.rollappID !== undefined && (obj.rollappID = message.rollappID);
    message.numTotal !== undefined && (obj.numTotal = Math.round(message.numTotal));
    message.numReceived !== undefined && (obj.numReceived = Math.round(message.numReceived));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisTransfers>, I>>(object: I): GenesisTransfers {
    const message = createBaseGenesisTransfers();
    message.rollappID = object.rollappID ?? "";
    message.numTotal = object.numTotal ?? 0;
    message.numReceived = object.numReceived ?? 0;
    return message;
  },
};

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
