/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Credential } from "./credential";
import { DidInfo } from "./did";
import { FilterLogger } from "./filter";
import { Service } from "./issuer";

export const protobufPackage = "metaearth.did";

/** GenesisState defines the kyc module's genesis state. */
export interface GenesisState {
  infos: DidInfo[];
  svcs: Service[];
  vcs: Credential[];
  flogs: FilterLogger[];
}

function createBaseGenesisState(): GenesisState {
  return { infos: [], svcs: [], vcs: [], flogs: [] };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.infos) {
      DidInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.svcs) {
      Service.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.vcs) {
      Credential.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.flogs) {
      FilterLogger.encode(v!, writer.uint32(34).fork()).ldelim();
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
          message.infos.push(DidInfo.decode(reader, reader.uint32()));
          break;
        case 2:
          message.svcs.push(Service.decode(reader, reader.uint32()));
          break;
        case 3:
          message.vcs.push(Credential.decode(reader, reader.uint32()));
          break;
        case 4:
          message.flogs.push(FilterLogger.decode(reader, reader.uint32()));
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
      infos: Array.isArray(object?.infos) ? object.infos.map((e: any) => DidInfo.fromJSON(e)) : [],
      svcs: Array.isArray(object?.svcs) ? object.svcs.map((e: any) => Service.fromJSON(e)) : [],
      vcs: Array.isArray(object?.vcs) ? object.vcs.map((e: any) => Credential.fromJSON(e)) : [],
      flogs: Array.isArray(object?.flogs) ? object.flogs.map((e: any) => FilterLogger.fromJSON(e)) : [],
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    if (message.infos) {
      obj.infos = message.infos.map((e) => e ? DidInfo.toJSON(e) : undefined);
    } else {
      obj.infos = [];
    }
    if (message.svcs) {
      obj.svcs = message.svcs.map((e) => e ? Service.toJSON(e) : undefined);
    } else {
      obj.svcs = [];
    }
    if (message.vcs) {
      obj.vcs = message.vcs.map((e) => e ? Credential.toJSON(e) : undefined);
    } else {
      obj.vcs = [];
    }
    if (message.flogs) {
      obj.flogs = message.flogs.map((e) => e ? FilterLogger.toJSON(e) : undefined);
    } else {
      obj.flogs = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.infos = object.infos?.map((e) => DidInfo.fromPartial(e)) || [];
    message.svcs = object.svcs?.map((e) => Service.fromPartial(e)) || [];
    message.vcs = object.vcs?.map((e) => Credential.fromPartial(e)) || [];
    message.flogs = object.flogs?.map((e) => FilterLogger.fromPartial(e)) || [];
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
