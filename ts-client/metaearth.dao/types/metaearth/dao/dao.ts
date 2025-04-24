/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.dao";

/** GenesisState defines the sudo module's genesis state. */
export interface DaoAddresses {
  globalDao: string;
  meidDao: string;
  devOperator: string;
  airdropAddress: string;
}

function createBaseDaoAddresses(): DaoAddresses {
  return { globalDao: "", meidDao: "", devOperator: "", airdropAddress: "" };
}

export const DaoAddresses = {
  encode(message: DaoAddresses, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.globalDao !== "") {
      writer.uint32(10).string(message.globalDao);
    }
    if (message.meidDao !== "") {
      writer.uint32(18).string(message.meidDao);
    }
    if (message.devOperator !== "") {
      writer.uint32(26).string(message.devOperator);
    }
    if (message.airdropAddress !== "") {
      writer.uint32(34).string(message.airdropAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DaoAddresses {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDaoAddresses();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.globalDao = reader.string();
          break;
        case 2:
          message.meidDao = reader.string();
          break;
        case 3:
          message.devOperator = reader.string();
          break;
        case 4:
          message.airdropAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DaoAddresses {
    return {
      globalDao: isSet(object.globalDao) ? String(object.globalDao) : "",
      meidDao: isSet(object.meidDao) ? String(object.meidDao) : "",
      devOperator: isSet(object.devOperator) ? String(object.devOperator) : "",
      airdropAddress: isSet(object.airdropAddress) ? String(object.airdropAddress) : "",
    };
  },

  toJSON(message: DaoAddresses): unknown {
    const obj: any = {};
    message.globalDao !== undefined && (obj.globalDao = message.globalDao);
    message.meidDao !== undefined && (obj.meidDao = message.meidDao);
    message.devOperator !== undefined && (obj.devOperator = message.devOperator);
    message.airdropAddress !== undefined && (obj.airdropAddress = message.airdropAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DaoAddresses>, I>>(object: I): DaoAddresses {
    const message = createBaseDaoAddresses();
    message.globalDao = object.globalDao ?? "";
    message.meidDao = object.meidDao ?? "";
    message.devOperator = object.devOperator ?? "";
    message.airdropAddress = object.airdropAddress ?? "";
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
