/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.eibc";

/** Params defines the parameters for the module. */
export interface Params {
  epochIdentifier: string;
  timeoutFee: string;
  errackFee: string;
}

function createBaseParams(): Params {
  return { epochIdentifier: "", timeoutFee: "", errackFee: "" };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.epochIdentifier !== "") {
      writer.uint32(10).string(message.epochIdentifier);
    }
    if (message.timeoutFee !== "") {
      writer.uint32(18).string(message.timeoutFee);
    }
    if (message.errackFee !== "") {
      writer.uint32(26).string(message.errackFee);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.epochIdentifier = reader.string();
          break;
        case 2:
          message.timeoutFee = reader.string();
          break;
        case 3:
          message.errackFee = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Params {
    return {
      epochIdentifier: isSet(object.epochIdentifier) ? String(object.epochIdentifier) : "",
      timeoutFee: isSet(object.timeoutFee) ? String(object.timeoutFee) : "",
      errackFee: isSet(object.errackFee) ? String(object.errackFee) : "",
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.epochIdentifier !== undefined && (obj.epochIdentifier = message.epochIdentifier);
    message.timeoutFee !== undefined && (obj.timeoutFee = message.timeoutFee);
    message.errackFee !== undefined && (obj.errackFee = message.errackFee);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.epochIdentifier = object.epochIdentifier ?? "";
    message.timeoutFee = object.timeoutFee ?? "";
    message.errackFee = object.errackFee ?? "";
    return message;
  },
};

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
