/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { GenesisState as GenesisState1 } from "../../cosmos/distribution/v1beta1/genesis";
import { Params } from "./params";

export const protobufPackage = "metaearth.wdistri";

/** GenesisState defines the wdistri module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  cosmosDistributionState: GenesisState1 | undefined;
}

function createBaseGenesisState(): GenesisState {
  return { params: undefined, cosmosDistributionState: undefined };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    if (message.cosmosDistributionState !== undefined) {
      GenesisState1.encode(message.cosmosDistributionState, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.cosmosDistributionState = GenesisState1.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      cosmosDistributionState: isSet(object.cosmosDistributionState)
        ? GenesisState1.fromJSON(object.cosmosDistributionState)
        : undefined,
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    message.cosmosDistributionState !== undefined && (obj.cosmosDistributionState = message.cosmosDistributionState
      ? GenesisState1.toJSON(message.cosmosDistributionState)
      : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    message.cosmosDistributionState =
      (object.cosmosDistributionState !== undefined && object.cosmosDistributionState !== null)
        ? GenesisState1.fromPartial(object.cosmosDistributionState)
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
