/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Params as Params1 } from "../../cosmos/distribution/v1beta1/distribution";

export const protobufPackage = "metaearth.wdistri";

/** Params defines the parameters for the module. */
export interface Params {
  distributionParams: Params1 | undefined;
}

function createBaseParams(): Params {
  return { distributionParams: undefined };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.distributionParams !== undefined) {
      Params1.encode(message.distributionParams, writer.uint32(10).fork()).ldelim();
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
          message.distributionParams = Params1.decode(reader, reader.uint32());
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
      distributionParams: isSet(object.distributionParams) ? Params1.fromJSON(object.distributionParams) : undefined,
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.distributionParams !== undefined
      && (obj.distributionParams = message.distributionParams ? Params1.toJSON(message.distributionParams) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.distributionParams = (object.distributionParams !== undefined && object.distributionParams !== null)
      ? Params1.fromPartial(object.distributionParams)
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
