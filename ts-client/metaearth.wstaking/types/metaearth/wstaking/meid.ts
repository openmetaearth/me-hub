/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.wstaking";

export interface Meid {
  account: string;
  creator: string;
  regionId: string;
  regionName: string;
  RewardType: number;
}

function createBaseMeid(): Meid {
  return { account: "", creator: "", regionId: "", regionName: "", RewardType: 0 };
}

export const Meid = {
  encode(message: Meid, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.account !== "") {
      writer.uint32(10).string(message.account);
    }
    if (message.creator !== "") {
      writer.uint32(18).string(message.creator);
    }
    if (message.regionId !== "") {
      writer.uint32(26).string(message.regionId);
    }
    if (message.regionName !== "") {
      writer.uint32(34).string(message.regionName);
    }
    if (message.RewardType !== 0) {
      writer.uint32(40).int32(message.RewardType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Meid {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMeid();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.account = reader.string();
          break;
        case 2:
          message.creator = reader.string();
          break;
        case 3:
          message.regionId = reader.string();
          break;
        case 4:
          message.regionName = reader.string();
          break;
        case 5:
          message.RewardType = reader.int32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Meid {
    return {
      account: isSet(object.account) ? String(object.account) : "",
      creator: isSet(object.creator) ? String(object.creator) : "",
      regionId: isSet(object.regionId) ? String(object.regionId) : "",
      regionName: isSet(object.regionName) ? String(object.regionName) : "",
      RewardType: isSet(object.RewardType) ? Number(object.RewardType) : 0,
    };
  },

  toJSON(message: Meid): unknown {
    const obj: any = {};
    message.account !== undefined && (obj.account = message.account);
    message.creator !== undefined && (obj.creator = message.creator);
    message.regionId !== undefined && (obj.regionId = message.regionId);
    message.regionName !== undefined && (obj.regionName = message.regionName);
    message.RewardType !== undefined && (obj.RewardType = Math.round(message.RewardType));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Meid>, I>>(object: I): Meid {
    const message = createBaseMeid();
    message.account = object.account ?? "";
    message.creator = object.creator ?? "";
    message.regionId = object.regionId ?? "";
    message.regionName = object.regionName ?? "";
    message.RewardType = object.RewardType ?? 0;
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
