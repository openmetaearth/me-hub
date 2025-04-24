/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { GroupInfo } from "./group_info";
import { GroupMember } from "./group_member";
import { GroupMemberCount } from "./group_member_count";
import { MemberJoined } from "./member_joined";
import { Params } from "./params";

export const protobufPackage = "metaearth.megroup";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryGetGroupRequest {
  id: number;
}

export interface QueryGetGroupResponse {
  Group: GroupInfo | undefined;
}

export interface QueryAllGroupRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllGroupResponse {
  Group: GroupInfo[];
  pagination: PageResponse | undefined;
}

export interface QueryGroupByMemberRequest {
  address: string;
}

export interface QueryGroupByMemberResponse {
  Group: GroupInfo | undefined;
}

export interface QueryGetGroupMemberRequest {
  address: string;
}

export interface QueryGetGroupMemberResponse {
  GroupMember: GroupMember | undefined;
}

export interface QueryGroupAllMemberRequest {
  groupID: number;
  pagination: PageRequest | undefined;
}

export interface QueryGroupAllMemberResponse {
  groupID: number;
  GroupMember: GroupMember[];
  pagination: PageResponse | undefined;
}

export interface QueryGetMemberJoinedRequest {
  address: string;
}

export interface QueryGetMemberJoinedResponse {
  memberJoined: MemberJoined | undefined;
}

export interface QueryAllMemberJoinedRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllMemberJoinedResponse {
  memberJoined: MemberJoined[];
  pagination: PageResponse | undefined;
}

export interface QueryGetGroupMemberCountRequest {
  groupId: number;
}

export interface QueryGetGroupMemberCountResponse {
  groupMemberCount: number;
}

export interface QueryAllGroupMemberCountRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllGroupMemberCountResponse {
  groupMemberCount: GroupMemberCount[];
  pagination: PageResponse | undefined;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
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

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryGetGroupRequest(): QueryGetGroupRequest {
  return { id: 0 };
}

export const QueryGetGroupRequest = {
  encode(message: QueryGetGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint64(message.id);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGroupRequest();
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

  fromJSON(object: any): QueryGetGroupRequest {
    return { id: isSet(object.id) ? Number(object.id) : 0 };
  },

  toJSON(message: QueryGetGroupRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = Math.round(message.id));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetGroupRequest>, I>>(object: I): QueryGetGroupRequest {
    const message = createBaseQueryGetGroupRequest();
    message.id = object.id ?? 0;
    return message;
  },
};

function createBaseQueryGetGroupResponse(): QueryGetGroupResponse {
  return { Group: undefined };
}

export const QueryGetGroupResponse = {
  encode(message: QueryGetGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Group !== undefined) {
      GroupInfo.encode(message.Group, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Group = GroupInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetGroupResponse {
    return { Group: isSet(object.Group) ? GroupInfo.fromJSON(object.Group) : undefined };
  },

  toJSON(message: QueryGetGroupResponse): unknown {
    const obj: any = {};
    message.Group !== undefined && (obj.Group = message.Group ? GroupInfo.toJSON(message.Group) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetGroupResponse>, I>>(object: I): QueryGetGroupResponse {
    const message = createBaseQueryGetGroupResponse();
    message.Group = (object.Group !== undefined && object.Group !== null)
      ? GroupInfo.fromPartial(object.Group)
      : undefined;
    return message;
  },
};

function createBaseQueryAllGroupRequest(): QueryAllGroupRequest {
  return { pagination: undefined };
}

export const QueryAllGroupRequest = {
  encode(message: QueryAllGroupRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllGroupRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGroupRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllGroupRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllGroupRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllGroupRequest>, I>>(object: I): QueryAllGroupRequest {
    const message = createBaseQueryAllGroupRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllGroupResponse(): QueryAllGroupResponse {
  return { Group: [], pagination: undefined };
}

export const QueryAllGroupResponse = {
  encode(message: QueryAllGroupResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Group) {
      GroupInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllGroupResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGroupResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Group.push(GroupInfo.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllGroupResponse {
    return {
      Group: Array.isArray(object?.Group) ? object.Group.map((e: any) => GroupInfo.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllGroupResponse): unknown {
    const obj: any = {};
    if (message.Group) {
      obj.Group = message.Group.map((e) => e ? GroupInfo.toJSON(e) : undefined);
    } else {
      obj.Group = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllGroupResponse>, I>>(object: I): QueryAllGroupResponse {
    const message = createBaseQueryAllGroupResponse();
    message.Group = object.Group?.map((e) => GroupInfo.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGroupByMemberRequest(): QueryGroupByMemberRequest {
  return { address: "" };
}

export const QueryGroupByMemberRequest = {
  encode(message: QueryGroupByMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGroupByMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGroupByMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGroupByMemberRequest {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryGroupByMemberRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGroupByMemberRequest>, I>>(object: I): QueryGroupByMemberRequest {
    const message = createBaseQueryGroupByMemberRequest();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryGroupByMemberResponse(): QueryGroupByMemberResponse {
  return { Group: undefined };
}

export const QueryGroupByMemberResponse = {
  encode(message: QueryGroupByMemberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Group !== undefined) {
      GroupInfo.encode(message.Group, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGroupByMemberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGroupByMemberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Group = GroupInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGroupByMemberResponse {
    return { Group: isSet(object.Group) ? GroupInfo.fromJSON(object.Group) : undefined };
  },

  toJSON(message: QueryGroupByMemberResponse): unknown {
    const obj: any = {};
    message.Group !== undefined && (obj.Group = message.Group ? GroupInfo.toJSON(message.Group) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGroupByMemberResponse>, I>>(object: I): QueryGroupByMemberResponse {
    const message = createBaseQueryGroupByMemberResponse();
    message.Group = (object.Group !== undefined && object.Group !== null)
      ? GroupInfo.fromPartial(object.Group)
      : undefined;
    return message;
  },
};

function createBaseQueryGetGroupMemberRequest(): QueryGetGroupMemberRequest {
  return { address: "" };
}

export const QueryGetGroupMemberRequest = {
  encode(message: QueryGetGroupMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetGroupMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGroupMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetGroupMemberRequest {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryGetGroupMemberRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetGroupMemberRequest>, I>>(object: I): QueryGetGroupMemberRequest {
    const message = createBaseQueryGetGroupMemberRequest();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryGetGroupMemberResponse(): QueryGetGroupMemberResponse {
  return { GroupMember: undefined };
}

export const QueryGetGroupMemberResponse = {
  encode(message: QueryGetGroupMemberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.GroupMember !== undefined) {
      GroupMember.encode(message.GroupMember, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetGroupMemberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGroupMemberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.GroupMember = GroupMember.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetGroupMemberResponse {
    return { GroupMember: isSet(object.GroupMember) ? GroupMember.fromJSON(object.GroupMember) : undefined };
  },

  toJSON(message: QueryGetGroupMemberResponse): unknown {
    const obj: any = {};
    message.GroupMember !== undefined
      && (obj.GroupMember = message.GroupMember ? GroupMember.toJSON(message.GroupMember) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetGroupMemberResponse>, I>>(object: I): QueryGetGroupMemberResponse {
    const message = createBaseQueryGetGroupMemberResponse();
    message.GroupMember = (object.GroupMember !== undefined && object.GroupMember !== null)
      ? GroupMember.fromPartial(object.GroupMember)
      : undefined;
    return message;
  },
};

function createBaseQueryGroupAllMemberRequest(): QueryGroupAllMemberRequest {
  return { groupID: 0, pagination: undefined };
}

export const QueryGroupAllMemberRequest = {
  encode(message: QueryGroupAllMemberRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.groupID !== 0) {
      writer.uint32(8).uint64(message.groupID);
    }
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGroupAllMemberRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGroupAllMemberRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.groupID = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGroupAllMemberRequest {
    return {
      groupID: isSet(object.groupID) ? Number(object.groupID) : 0,
      pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryGroupAllMemberRequest): unknown {
    const obj: any = {};
    message.groupID !== undefined && (obj.groupID = Math.round(message.groupID));
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGroupAllMemberRequest>, I>>(object: I): QueryGroupAllMemberRequest {
    const message = createBaseQueryGroupAllMemberRequest();
    message.groupID = object.groupID ?? 0;
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGroupAllMemberResponse(): QueryGroupAllMemberResponse {
  return { groupID: 0, GroupMember: [], pagination: undefined };
}

export const QueryGroupAllMemberResponse = {
  encode(message: QueryGroupAllMemberResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.groupID !== 0) {
      writer.uint32(8).uint64(message.groupID);
    }
    for (const v of message.GroupMember) {
      GroupMember.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGroupAllMemberResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGroupAllMemberResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.groupID = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.GroupMember.push(GroupMember.decode(reader, reader.uint32()));
          break;
        case 3:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGroupAllMemberResponse {
    return {
      groupID: isSet(object.groupID) ? Number(object.groupID) : 0,
      GroupMember: Array.isArray(object?.GroupMember)
        ? object.GroupMember.map((e: any) => GroupMember.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryGroupAllMemberResponse): unknown {
    const obj: any = {};
    message.groupID !== undefined && (obj.groupID = Math.round(message.groupID));
    if (message.GroupMember) {
      obj.GroupMember = message.GroupMember.map((e) => e ? GroupMember.toJSON(e) : undefined);
    } else {
      obj.GroupMember = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGroupAllMemberResponse>, I>>(object: I): QueryGroupAllMemberResponse {
    const message = createBaseQueryGroupAllMemberResponse();
    message.groupID = object.groupID ?? 0;
    message.GroupMember = object.GroupMember?.map((e) => GroupMember.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetMemberJoinedRequest(): QueryGetMemberJoinedRequest {
  return { address: "" };
}

export const QueryGetMemberJoinedRequest = {
  encode(message: QueryGetMemberJoinedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetMemberJoinedRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetMemberJoinedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetMemberJoinedRequest {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: QueryGetMemberJoinedRequest): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetMemberJoinedRequest>, I>>(object: I): QueryGetMemberJoinedRequest {
    const message = createBaseQueryGetMemberJoinedRequest();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseQueryGetMemberJoinedResponse(): QueryGetMemberJoinedResponse {
  return { memberJoined: undefined };
}

export const QueryGetMemberJoinedResponse = {
  encode(message: QueryGetMemberJoinedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.memberJoined !== undefined) {
      MemberJoined.encode(message.memberJoined, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetMemberJoinedResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetMemberJoinedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.memberJoined = MemberJoined.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetMemberJoinedResponse {
    return { memberJoined: isSet(object.memberJoined) ? MemberJoined.fromJSON(object.memberJoined) : undefined };
  },

  toJSON(message: QueryGetMemberJoinedResponse): unknown {
    const obj: any = {};
    message.memberJoined !== undefined
      && (obj.memberJoined = message.memberJoined ? MemberJoined.toJSON(message.memberJoined) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetMemberJoinedResponse>, I>>(object: I): QueryGetMemberJoinedResponse {
    const message = createBaseQueryGetMemberJoinedResponse();
    message.memberJoined = (object.memberJoined !== undefined && object.memberJoined !== null)
      ? MemberJoined.fromPartial(object.memberJoined)
      : undefined;
    return message;
  },
};

function createBaseQueryAllMemberJoinedRequest(): QueryAllMemberJoinedRequest {
  return { pagination: undefined };
}

export const QueryAllMemberJoinedRequest = {
  encode(message: QueryAllMemberJoinedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMemberJoinedRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMemberJoinedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllMemberJoinedRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllMemberJoinedRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllMemberJoinedRequest>, I>>(object: I): QueryAllMemberJoinedRequest {
    const message = createBaseQueryAllMemberJoinedRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllMemberJoinedResponse(): QueryAllMemberJoinedResponse {
  return { memberJoined: [], pagination: undefined };
}

export const QueryAllMemberJoinedResponse = {
  encode(message: QueryAllMemberJoinedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.memberJoined) {
      MemberJoined.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMemberJoinedResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMemberJoinedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.memberJoined.push(MemberJoined.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllMemberJoinedResponse {
    return {
      memberJoined: Array.isArray(object?.memberJoined)
        ? object.memberJoined.map((e: any) => MemberJoined.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllMemberJoinedResponse): unknown {
    const obj: any = {};
    if (message.memberJoined) {
      obj.memberJoined = message.memberJoined.map((e) => e ? MemberJoined.toJSON(e) : undefined);
    } else {
      obj.memberJoined = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllMemberJoinedResponse>, I>>(object: I): QueryAllMemberJoinedResponse {
    const message = createBaseQueryAllMemberJoinedResponse();
    message.memberJoined = object.memberJoined?.map((e) => MemberJoined.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryGetGroupMemberCountRequest(): QueryGetGroupMemberCountRequest {
  return { groupId: 0 };
}

export const QueryGetGroupMemberCountRequest = {
  encode(message: QueryGetGroupMemberCountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.groupId !== 0) {
      writer.uint32(8).uint64(message.groupId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetGroupMemberCountRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGroupMemberCountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.groupId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetGroupMemberCountRequest {
    return { groupId: isSet(object.groupId) ? Number(object.groupId) : 0 };
  },

  toJSON(message: QueryGetGroupMemberCountRequest): unknown {
    const obj: any = {};
    message.groupId !== undefined && (obj.groupId = Math.round(message.groupId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetGroupMemberCountRequest>, I>>(
    object: I,
  ): QueryGetGroupMemberCountRequest {
    const message = createBaseQueryGetGroupMemberCountRequest();
    message.groupId = object.groupId ?? 0;
    return message;
  },
};

function createBaseQueryGetGroupMemberCountResponse(): QueryGetGroupMemberCountResponse {
  return { groupMemberCount: 0 };
}

export const QueryGetGroupMemberCountResponse = {
  encode(message: QueryGetGroupMemberCountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.groupMemberCount !== 0) {
      writer.uint32(8).uint64(message.groupMemberCount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetGroupMemberCountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetGroupMemberCountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.groupMemberCount = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryGetGroupMemberCountResponse {
    return { groupMemberCount: isSet(object.groupMemberCount) ? Number(object.groupMemberCount) : 0 };
  },

  toJSON(message: QueryGetGroupMemberCountResponse): unknown {
    const obj: any = {};
    message.groupMemberCount !== undefined && (obj.groupMemberCount = Math.round(message.groupMemberCount));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryGetGroupMemberCountResponse>, I>>(
    object: I,
  ): QueryGetGroupMemberCountResponse {
    const message = createBaseQueryGetGroupMemberCountResponse();
    message.groupMemberCount = object.groupMemberCount ?? 0;
    return message;
  },
};

function createBaseQueryAllGroupMemberCountRequest(): QueryAllGroupMemberCountRequest {
  return { pagination: undefined };
}

export const QueryAllGroupMemberCountRequest = {
  encode(message: QueryAllGroupMemberCountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllGroupMemberCountRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGroupMemberCountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllGroupMemberCountRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllGroupMemberCountRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllGroupMemberCountRequest>, I>>(
    object: I,
  ): QueryAllGroupMemberCountRequest {
    const message = createBaseQueryAllGroupMemberCountRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllGroupMemberCountResponse(): QueryAllGroupMemberCountResponse {
  return { groupMemberCount: [], pagination: undefined };
}

export const QueryAllGroupMemberCountResponse = {
  encode(message: QueryAllGroupMemberCountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.groupMemberCount) {
      GroupMemberCount.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllGroupMemberCountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllGroupMemberCountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.groupMemberCount.push(GroupMemberCount.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllGroupMemberCountResponse {
    return {
      groupMemberCount: Array.isArray(object?.groupMemberCount)
        ? object.groupMemberCount.map((e: any) => GroupMemberCount.fromJSON(e))
        : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllGroupMemberCountResponse): unknown {
    const obj: any = {};
    if (message.groupMemberCount) {
      obj.groupMemberCount = message.groupMemberCount.map((e) => e ? GroupMemberCount.toJSON(e) : undefined);
    } else {
      obj.groupMemberCount = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllGroupMemberCountResponse>, I>>(
    object: I,
  ): QueryAllGroupMemberCountResponse {
    const message = createBaseQueryAllGroupMemberCountResponse();
    message.groupMemberCount = object.groupMemberCount?.map((e) => GroupMemberCount.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a list of Group items. */
  Group(request: QueryGetGroupRequest): Promise<QueryGetGroupResponse>;
  GroupAll(request: QueryAllGroupRequest): Promise<QueryAllGroupResponse>;
  GroupByMember(request: QueryGroupByMemberRequest): Promise<QueryGetGroupResponse>;
  /** Queries a list of GroupMember items. */
  GroupMember(request: QueryGetGroupMemberRequest): Promise<QueryGetGroupMemberResponse>;
  GroupMemberAll(request: QueryGroupAllMemberRequest): Promise<QueryGroupAllMemberResponse>;
  /** Queries a list of GroupMemberCount items. */
  GroupMemberCount(request: QueryGetGroupMemberCountRequest): Promise<QueryGetGroupMemberCountResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.Group = this.Group.bind(this);
    this.GroupAll = this.GroupAll.bind(this);
    this.GroupByMember = this.GroupByMember.bind(this);
    this.GroupMember = this.GroupMember.bind(this);
    this.GroupMemberAll = this.GroupMemberAll.bind(this);
    this.GroupMemberCount = this.GroupMemberCount.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  Group(request: QueryGetGroupRequest): Promise<QueryGetGroupResponse> {
    const data = QueryGetGroupRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "Group", data);
    return promise.then((data) => QueryGetGroupResponse.decode(new _m0.Reader(data)));
  }

  GroupAll(request: QueryAllGroupRequest): Promise<QueryAllGroupResponse> {
    const data = QueryAllGroupRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "GroupAll", data);
    return promise.then((data) => QueryAllGroupResponse.decode(new _m0.Reader(data)));
  }

  GroupByMember(request: QueryGroupByMemberRequest): Promise<QueryGetGroupResponse> {
    const data = QueryGroupByMemberRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "GroupByMember", data);
    return promise.then((data) => QueryGetGroupResponse.decode(new _m0.Reader(data)));
  }

  GroupMember(request: QueryGetGroupMemberRequest): Promise<QueryGetGroupMemberResponse> {
    const data = QueryGetGroupMemberRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "GroupMember", data);
    return promise.then((data) => QueryGetGroupMemberResponse.decode(new _m0.Reader(data)));
  }

  GroupMemberAll(request: QueryGroupAllMemberRequest): Promise<QueryGroupAllMemberResponse> {
    const data = QueryGroupAllMemberRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "GroupMemberAll", data);
    return promise.then((data) => QueryGroupAllMemberResponse.decode(new _m0.Reader(data)));
  }

  GroupMemberCount(request: QueryGetGroupMemberCountRequest): Promise<QueryGetGroupMemberCountResponse> {
    const data = QueryGetGroupMemberCountRequest.encode(request).finish();
    const promise = this.rpc.request("metaearth.megroup.Query", "GroupMemberCount", data);
    return promise.then((data) => QueryGetGroupMemberCountResponse.decode(new _m0.Reader(data)));
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
