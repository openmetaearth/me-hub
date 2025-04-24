/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../../../cosmos/base/v1beta1/coin";
import { Duration } from "../../../google/protobuf/duration";

export const protobufPackage = "dymensionxyz.dymension.sequencer";

/** Params defines the parameters for the module. */
export interface Params {
  minBond:
    | Coin
    | undefined;
  /** unbonding_time is the time duration of unbonding. */
  unbondingTime: Duration | undefined;
}

function createBaseParams(): Params {
  return { minBond: undefined, unbondingTime: undefined };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.minBond !== undefined) {
      Coin.encode(message.minBond, writer.uint32(10).fork()).ldelim();
    }
    if (message.unbondingTime !== undefined) {
      Duration.encode(message.unbondingTime, writer.uint32(18).fork()).ldelim();
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
          message.minBond = Coin.decode(reader, reader.uint32());
          break;
        case 2:
          message.unbondingTime = Duration.decode(reader, reader.uint32());
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
      minBond: isSet(object.minBond) ? Coin.fromJSON(object.minBond) : undefined,
      unbondingTime: isSet(object.unbondingTime) ? Duration.fromJSON(object.unbondingTime) : undefined,
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.minBond !== undefined && (obj.minBond = message.minBond ? Coin.toJSON(message.minBond) : undefined);
    message.unbondingTime !== undefined
      && (obj.unbondingTime = message.unbondingTime ? Duration.toJSON(message.unbondingTime) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.minBond = (object.minBond !== undefined && object.minBond !== null)
      ? Coin.fromPartial(object.minBond)
      : undefined;
    message.unbondingTime = (object.unbondingTime !== undefined && object.unbondingTime !== null)
      ? Duration.fromPartial(object.unbondingTime)
      : undefined;
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
