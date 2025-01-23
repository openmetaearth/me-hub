/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { GroupInfo } from "./group_info";

export const protobufPackage = "metaearth.megroup";

export interface MsgCreateGroup {
  creator: string;
  groupInfo: GroupInfo | undefined;
}

export interface MsgCreateGroupResponse {
  id: number;
}

export interface MsgUpdateGroup {
  creator: string;
  id: number;
  groupInfo: GroupInfo | undefined;
}

export interface MsgUpdateGroupResponse {
}

export interface MsgDeleteGroup {
  creator: string;
  id: number;
}

export interface MsgDeleteGroupResponse {
}

export interface MsgJoinGroup {
  creator: string;
  groupId: number;
  applicantAddress: string;
}

export interface MsgJoinGroupResponse {
}

export interface MsgLeaveGroupRequest {
  creator: string;
  groupId: number;
}

export interface MsgLeaveGroupResponse {
}

function createBaseMsgCreateGroup(): MsgCreateGroup {
  return { creator: "", groupInfo: undefined };
}

export const MsgCreateGroup = {
  encode(message: MsgCreateGroup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.groupInfo !== undefined) {
      GroupInfo.encode(message.groupInfo, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateGroup {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateGroup();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.groupInfo = GroupInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateGroup {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      groupInfo: isSet(object.groupInfo) ? GroupInfo.fromJSON(object.groupInfo) : undefined,
    };
  },

  toJSON(message: MsgCreateGroup): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.groupInfo !== undefined
      && (obj.groupInfo = message.groupInfo ? GroupInfo.toJSON(message.groupInfo) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateGroup>, I>>(object: I): MsgCreateGroup {
    const message = createBaseMsgCreateGroup();
    message.creator = object.creator ?? "";
    message.groupInfo = (object.groupInfo !== undefined && object.groupInfo !== null)
      ? GroupInfo.fromPartial(object.groupInfo)
      : undefined;
    return message;
  },
};

function createBaseMsgCreateGroupResponse(): MsgCreateGroupResponse {
  return { id: 0 };
}

export const MsgCreateGroupResponse = {
  encode(message: MsgCreateGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateGroupResponse {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: MsgCreateGroupResponse): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateGroupResponse>, I>>(object: I): MsgCreateGroupResponse {
    const message = createBaseMsgCreateGroupResponse();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseMsgUpdateGroup(): MsgUpdateGroup {
  return { creator: "", id: 0, groupInfo: undefined };
}

export const MsgUpdateGroup = {
  encode(message: MsgUpdateGroup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    if (message.groupInfo !== undefined) {
      GroupInfo.encode(message.groupInfo, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateGroup {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateGroup();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.groupInfo = GroupInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateGroup {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      id: isSet(object.id) ? Number(object.id) : 0,
      groupInfo: isSet(object.groupInfo) ? GroupInfo.fromJSON(object.groupInfo) : undefined,
    };
  },

  toJSON(message: MsgUpdateGroup): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.id !== undefined && (obj.id = Math.round(message.id));
    message.groupInfo !== undefined
      && (obj.groupInfo = message.groupInfo ? GroupInfo.toJSON(message.groupInfo) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateGroup>, I>>(object: I): MsgUpdateGroup {
    const message = createBaseMsgUpdateGroup();
    message.creator = object.creator ?? "";
    message.id = object.id ?? 0;
    message.groupInfo = (object.groupInfo !== undefined && object.groupInfo !== null)
      ? GroupInfo.fromPartial(object.groupInfo)
      : undefined;
    return message;
  },
};

function createBaseMsgUpdateGroupResponse(): MsgUpdateGroupResponse {
  return {};
}

export const MsgUpdateGroupResponse = {
  encode(_: MsgUpdateGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgUpdateGroupResponse {
    return {};
  },

  toJSON(_: MsgUpdateGroupResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateGroupResponse>, I>>(_: I): MsgUpdateGroupResponse {
    const message = createBaseMsgUpdateGroupResponse();
    return message;
  },
};

function createBaseMsgDeleteGroup(): MsgDeleteGroup {
  return { creator: "", id: 0 };
}

export const MsgDeleteGroup = {
  encode(message: MsgDeleteGroup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.id !== 0) {
      writer.uint32(16).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteGroup {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteGroup();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgDeleteGroup {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      id: isSet(object.id) ? Number(object.id) : 0,
    };
  },

  toJSON(message: MsgDeleteGroup): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDeleteGroup>, I>>(object: I): MsgDeleteGroup {
    const message = createBaseMsgDeleteGroup();
    message.creator = object.creator ?? "";
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseMsgDeleteGroupResponse(): MsgDeleteGroupResponse {
  return {};
}

export const MsgDeleteGroupResponse = {
  encode(_: MsgDeleteGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgDeleteGroupResponse {
    return {};
  },

  toJSON(_: MsgDeleteGroupResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgDeleteGroupResponse>, I>>(_: I): MsgDeleteGroupResponse {
    const message = createBaseMsgDeleteGroupResponse();
    return message;
  },
};

function createBaseMsgJoinGroup(): MsgJoinGroup {
  return { creator: "", groupId: 0, applicantAddress: "" };
}

export const MsgJoinGroup = {
  encode(message: MsgJoinGroup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.groupId !== 0) {
      writer.uint32(16).uint64(message.groupId);
    }
    if (message.applicantAddress !== "") {
      writer.uint32(26).string(message.applicantAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgJoinGroup {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgJoinGroup();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.groupId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.applicantAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgJoinGroup {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      groupId: isSet(object.groupId) ? Number(object.groupId) : 0,
      applicantAddress: isSet(object.applicantAddress) ? String(object.applicantAddress) : "",
    };
  },

  toJSON(message: MsgJoinGroup): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.groupId !== undefined && (obj.groupId = Math.round(message.groupId));
    message.applicantAddress !== undefined && (obj.applicantAddress = message.applicantAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgJoinGroup>, I>>(object: I): MsgJoinGroup {
    const message = createBaseMsgJoinGroup();
    message.creator = object.creator ?? "";
    message.groupId = object.groupId ?? 0;
    message.applicantAddress = object.applicantAddress ?? "";
    return message;
  },
};

function createBaseMsgJoinGroupResponse(): MsgJoinGroupResponse {
  return {};
}

export const MsgJoinGroupResponse = {
  encode(_: MsgJoinGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgJoinGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgJoinGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgJoinGroupResponse {
    return {};
  },

  toJSON(_: MsgJoinGroupResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgJoinGroupResponse>, I>>(_: I): MsgJoinGroupResponse {
    const message = createBaseMsgJoinGroupResponse();
    return message;
  },
};

function createBaseMsgLeaveGroupRequest(): MsgLeaveGroupRequest {
  return { creator: "", groupId: 0 };
}

export const MsgLeaveGroupRequest = {
  encode(message: MsgLeaveGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.groupId !== 0) {
      writer.uint32(16).uint64(message.groupId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgLeaveGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgLeaveGroupRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.groupId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgLeaveGroupRequest {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      groupId: isSet(object.groupId) ? Number(object.groupId) : 0,
    };
  },

  toJSON(message: MsgLeaveGroupRequest): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.groupId !== undefined && (obj.groupId = Math.round(message.groupId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgLeaveGroupRequest>, I>>(object: I): MsgLeaveGroupRequest {
    const message = createBaseMsgLeaveGroupRequest();
    message.creator = object.creator ?? "";
    message.groupId = object.groupId ?? 0;
    return message;
  },
};

function createBaseMsgLeaveGroupResponse(): MsgLeaveGroupResponse {
  return {};
}

export const MsgLeaveGroupResponse = {
  encode(_: MsgLeaveGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgLeaveGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgLeaveGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): MsgLeaveGroupResponse {
    return {};
  },

  toJSON(_: MsgLeaveGroupResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgLeaveGroupResponse>, I>>(_: I): MsgLeaveGroupResponse {
    const message = createBaseMsgLeaveGroupResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  CreateGroup(request: MsgCreateGroup): Promise<MsgCreateGroupResponse>;
  UpdateGroup(request: MsgUpdateGroup): Promise<MsgUpdateGroupResponse>;
  DeleteGroup(request: MsgDeleteGroup): Promise<MsgDeleteGroupResponse>;
  JoinGroup(request: MsgJoinGroup): Promise<MsgJoinGroupResponse>;
  LeaveGroup(request: MsgLeaveGroupRequest): Promise<MsgLeaveGroupResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.CreateGroup = this.CreateGroup.bind(this);
    this.UpdateGroup = this.UpdateGroup.bind(this);
    this.DeleteGroup = this.DeleteGroup.bind(this);
    this.JoinGroup = this.JoinGroup.bind(this);
    this.LeaveGroup = this.LeaveGroup.bind(this);
  }
  CreateGroup(request: MsgCreateGroup): Promise<MsgCreateGroupResponse> {
    const data = MsgCreateGroup.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Msg", "CreateGroup", data);
    return promise.then((data) => MsgCreateGroupResponse.decode(new _m0.Reader(data)));
  }

  UpdateGroup(request: MsgUpdateGroup): Promise<MsgUpdateGroupResponse> {
    const data = MsgUpdateGroup.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Msg", "UpdateGroup", data);
    return promise.then((data) => MsgUpdateGroupResponse.decode(new _m0.Reader(data)));
  }

  DeleteGroup(request: MsgDeleteGroup): Promise<MsgDeleteGroupResponse> {
    const data = MsgDeleteGroup.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Msg", "DeleteGroup", data);
    return promise.then((data) => MsgDeleteGroupResponse.decode(new _m0.Reader(data)));
  }

  JoinGroup(request: MsgJoinGroup): Promise<MsgJoinGroupResponse> {
    const data = MsgJoinGroup.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Msg", "JoinGroup", data);
    return promise.then((data) => MsgJoinGroupResponse.decode(new _m0.Reader(data)));
  }

  LeaveGroup(request: MsgLeaveGroupRequest): Promise<MsgLeaveGroupResponse> {
    const data = MsgLeaveGroupRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Msg", "LeaveGroup", data);
    return promise.then((data) => MsgLeaveGroupResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

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
