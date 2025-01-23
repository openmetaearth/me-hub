/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { DaoAddresses } from "./dao";

export const protobufPackage = "metaearth.dao";

/** GenesisState defines the sudo module's genesis state. */
export interface GenesisState {
  daoAddresses: DaoAddresses | undefined;
}

function createBaseGenesisState(): GenesisState {
  return { daoAddresses: undefined };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.daoAddresses !== undefined) {
      DaoAddresses.encode(message.daoAddresses, writer.uint32(10).fork()).ldelim();
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
          message.daoAddresses = DaoAddresses.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    return { daoAddresses: isSet(object.daoAddresses) ? DaoAddresses.fromJSON(object.daoAddresses) : undefined };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.daoAddresses !== undefined
      && (obj.daoAddresses = message.daoAddresses ? DaoAddresses.toJSON(message.daoAddresses) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.daoAddresses = (object.daoAddresses !== undefined && object.daoAddresses !== null)
      ? DaoAddresses.fromPartial(object.daoAddresses)
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
