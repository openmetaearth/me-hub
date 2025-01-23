/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { GroupInfo } from "./group_info";
import { GroupMember } from "./group_member";
import { GroupMemberCount } from "./group_member_count";
import { MemberJoined } from "./member_joined";
import { Params } from "./params";

export const protobufPackage = "metaearth.megroup";

/** GenesisState defines the megroup module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  groups: GroupInfo[];
  groupCount: number;
  groupMembers: GroupMember[];
  memberJoinedList: MemberJoined[];
  groupMemberCountList: GroupMemberCount[];
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    groups: [],
    groupCount: 0,
    groupMembers: [],
    memberJoinedList: [],
    groupMemberCountList: [],
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.groups) {
      GroupInfo.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.groupCount !== 0) {
      writer.uint32(24).uint64(message.groupCount);
    }
    for (const v of message.groupMembers) {
      GroupMember.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.memberJoinedList) {
      MemberJoined.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.groupMemberCountList) {
      GroupMemberCount.encode(v!, writer.uint32(50).fork()).ldelim();
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
          message.groups.push(GroupInfo.decode(reader, reader.uint32()));
          break;
        case 3:
          message.groupCount = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.groupMembers.push(GroupMember.decode(reader, reader.uint32()));
          break;
        case 5:
          message.memberJoinedList.push(MemberJoined.decode(reader, reader.uint32()));
          break;
        case 6:
          message.groupMemberCountList.push(GroupMemberCount.decode(reader, reader.uint32()));
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
      groups: Array.isArray(object?.groups) ? object.groups.map((e: any) => GroupInfo.fromJSON(e)) : [],
      groupCount: isSet(object.groupCount) ? Number(object.groupCount) : 0,
      groupMembers: Array.isArray(object?.groupMembers)
        ? object.groupMembers.map((e: any) => GroupMember.fromJSON(e))
        : [],
      memberJoinedList: Array.isArray(object?.memberJoinedList)
        ? object.memberJoinedList.map((e: any) => MemberJoined.fromJSON(e))
        : [],
      groupMemberCountList: Array.isArray(object?.groupMemberCountList)
        ? object.groupMemberCountList.map((e: any) => GroupMemberCount.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.groups) {
      obj.groups = message.groups.map((e) => e ? GroupInfo.toJSON(e) : undefined);
    } else {
      obj.groups = [];
    }
    message.groupCount !== undefined && (obj.groupCount = Math.round(message.groupCount));
    if (message.groupMembers) {
      obj.groupMembers = message.groupMembers.map((e) => e ? GroupMember.toJSON(e) : undefined);
    } else {
      obj.groupMembers = [];
    }
    if (message.memberJoinedList) {
      obj.memberJoinedList = message.memberJoinedList.map((e) => e ? MemberJoined.toJSON(e) : undefined);
    } else {
      obj.memberJoinedList = [];
    }
    if (message.groupMemberCountList) {
      obj.groupMemberCountList = message.groupMemberCountList.map((e) => e ? GroupMemberCount.toJSON(e) : undefined);
    } else {
      obj.groupMemberCountList = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    message.groups = object.groups?.map((e) => GroupInfo.fromPartial(e)) || [];
    message.groupCount = object.groupCount ?? 0;
    message.groupMembers = object.groupMembers?.map((e) => GroupMember.fromPartial(e)) || [];
    message.memberJoinedList = object.memberJoinedList?.map((e) => MemberJoined.fromPartial(e)) || [];
    message.groupMemberCountList = object.groupMemberCountList?.map((e) => GroupMemberCount.fromPartial(e)) || [];
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
