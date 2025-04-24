/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

export interface SubmitFraudProposal {
  title: string;
  description: string;
  /** The rollapp id */
  rollappId: string;
  /** The ibc client id of the rollapp */
  ibcClientId: string;
  /** The height of the fraudelent block */
  fraudelentHeight: number;
  /** The address of the fraudelent sequencer */
  fraudelentSequencerAddress: string;
}

function createBaseSubmitFraudProposal(): SubmitFraudProposal {
  return {
    title: "",
    description: "",
    rollappId: "",
    ibcClientId: "",
    fraudelentHeight: 0,
    fraudelentSequencerAddress: "",
  };
}

export const SubmitFraudProposal = {
  encode(message: SubmitFraudProposal, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.rollappId !== "") {
      writer.uint32(26).string(message.rollappId);
    }
    if (message.ibcClientId !== "") {
      writer.uint32(34).string(message.ibcClientId);
    }
    if (message.fraudelentHeight !== 0) {
      writer.uint32(40).uint64(message.fraudelentHeight);
    }
    if (message.fraudelentSequencerAddress !== "") {
      writer.uint32(50).string(message.fraudelentSequencerAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubmitFraudProposal {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubmitFraudProposal();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.rollappId = reader.string();
          break;
        case 4:
          message.ibcClientId = reader.string();
          break;
        case 5:
          message.fraudelentHeight = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.fraudelentSequencerAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SubmitFraudProposal {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      description: isSet(object.description) ? String(object.description) : "",
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      ibcClientId: isSet(object.ibcClientId) ? String(object.ibcClientId) : "",
      fraudelentHeight: isSet(object.fraudelentHeight) ? Number(object.fraudelentHeight) : 0,
      fraudelentSequencerAddress: isSet(object.fraudelentSequencerAddress)
        ? String(object.fraudelentSequencerAddress)
        : "",
    };
  },

  toJSON(message: SubmitFraudProposal): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.description !== undefined && (obj.description = message.description);
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.ibcClientId !== undefined && (obj.ibcClientId = message.ibcClientId);
    message.fraudelentHeight !== undefined && (obj.fraudelentHeight = Math.round(message.fraudelentHeight));
    message.fraudelentSequencerAddress !== undefined
      && (obj.fraudelentSequencerAddress = message.fraudelentSequencerAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SubmitFraudProposal>, I>>(object: I): SubmitFraudProposal {
    const message = createBaseSubmitFraudProposal();
    message.title = object.title ?? "";
    message.description = object.description ?? "";
    message.rollappId = object.rollappId ?? "";
    message.ibcClientId = object.ibcClientId ?? "";
    message.fraudelentHeight = object.fraudelentHeight ?? 0;
    message.fraudelentSequencerAddress = object.fraudelentSequencerAddress ?? "";
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
