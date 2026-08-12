/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { DaoAddresses } from "./dao";

export const protobufPackage = "metaearth.dao";

export interface MsgUpdateDao {
  creator: string;
  daoAddresses: DaoAddresses | undefined;
}

export interface MsgUpdateDaoResponse {
}

export interface MsgFreeGasAccount {
  creator: string;
  accounts: FreeGasAccount[];
}

export interface FreeGasAccount {
  address: string;
  isFree: boolean;
}

export interface MsgFreeGasAccountResponse {
}

function createBaseMsgUpdateDao(): MsgUpdateDao {
  return { creator: "", daoAddresses: undefined };
}

export const MsgUpdateDao = {
  encode(message: MsgUpdateDao, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.daoAddresses !== undefined) {
      DaoAddresses.encode(message.daoAddresses, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDao {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDao();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.daoAddresses = DaoAddresses.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUpdateDao {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      daoAddresses: isSet(object.daoAddresses) ? DaoAddresses.fromJSON(object.daoAddresses) : undefined,
    };
  },

  toJSON(message: MsgUpdateDao): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.daoAddresses !== undefined
      && (obj.daoAddresses = message.daoAddresses ? DaoAddresses.toJSON(message.daoAddresses) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateDao>, I>>(object: I): MsgUpdateDao {
    const message = createBaseMsgUpdateDao();
    message.creator = object.creator ?? "";
    message.daoAddresses = (object.daoAddresses !== undefined && object.daoAddresses !== null)
      ? DaoAddresses.fromPartial(object.daoAddresses)
      : undefined;
    return message;
  },
};

function createBaseMsgUpdateDaoResponse(): MsgUpdateDaoResponse {
  return {};
}

export const MsgUpdateDaoResponse = {
  encode(_: MsgUpdateDaoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDaoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDaoResponse();
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

  fromJSON(_: any): MsgUpdateDaoResponse {
    return {};
  },

  toJSON(_: MsgUpdateDaoResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateDaoResponse>, I>>(_: I): MsgUpdateDaoResponse {
    const message = createBaseMsgUpdateDaoResponse();
    return message;
  },
};

function createBaseMsgFreeGasAccount(): MsgFreeGasAccount {
  return { creator: "", accounts: [] };
}

export const MsgFreeGasAccount = {
  encode(message: MsgFreeGasAccount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    for (const v of message.accounts) {
      FreeGasAccount.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFreeGasAccount {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFreeGasAccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.accounts.push(FreeGasAccount.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgFreeGasAccount {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      accounts: Array.isArray(object?.accounts) ? object.accounts.map((e: any) => FreeGasAccount.fromJSON(e)) : [],
    };
  },

  toJSON(message: MsgFreeGasAccount): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    if (message.accounts) {
      obj.accounts = message.accounts.map((e) => e ? FreeGasAccount.toJSON(e) : undefined);
    } else {
      obj.accounts = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFreeGasAccount>, I>>(object: I): MsgFreeGasAccount {
    const message = createBaseMsgFreeGasAccount();
    message.creator = object.creator ?? "";
    message.accounts = object.accounts?.map((e) => FreeGasAccount.fromPartial(e)) || [];
    return message;
  },
};

function createBaseFreeGasAccount(): FreeGasAccount {
  return { address: "", isFree: false };
}

export const FreeGasAccount = {
  encode(message: FreeGasAccount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.isFree === true) {
      writer.uint32(16).bool(message.isFree);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FreeGasAccount {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFreeGasAccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;
        case 2:
          message.isFree = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FreeGasAccount {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      isFree: isSet(object.isFree) ? Boolean(object.isFree) : false,
    };
  },

  toJSON(message: FreeGasAccount): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.isFree !== undefined && (obj.isFree = message.isFree);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FreeGasAccount>, I>>(object: I): FreeGasAccount {
    const message = createBaseFreeGasAccount();
    message.address = object.address ?? "";
    message.isFree = object.isFree ?? false;
    return message;
  },
};

function createBaseMsgFreeGasAccountResponse(): MsgFreeGasAccountResponse {
  return {};
}

export const MsgFreeGasAccountResponse = {
  encode(_: MsgFreeGasAccountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgFreeGasAccountResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgFreeGasAccountResponse();
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

  fromJSON(_: any): MsgFreeGasAccountResponse {
    return {};
  },

  toJSON(_: MsgFreeGasAccountResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgFreeGasAccountResponse>, I>>(_: I): MsgFreeGasAccountResponse {
    const message = createBaseMsgFreeGasAccountResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  UpdateDao(request: MsgUpdateDao): Promise<MsgUpdateDaoResponse>;
  FreeGasAccount(request: MsgFreeGasAccount): Promise<MsgFreeGasAccountResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.UpdateDao = this.UpdateDao.bind(this);
    this.FreeGasAccount = this.FreeGasAccount.bind(this);
  }
  UpdateDao(request: MsgUpdateDao): Promise<MsgUpdateDaoResponse> {
    const data = MsgUpdateDao.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Msg", "UpdateDao", data);
    return promise.then((data) => MsgUpdateDaoResponse.decode(new _m0.Reader(data)));
  }

  FreeGasAccount(request: MsgFreeGasAccount): Promise<MsgFreeGasAccountResponse> {
    const data = MsgFreeGasAccount.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Msg", "FreeGasAccount", data);
    return promise.then((data) => MsgFreeGasAccountResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

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
