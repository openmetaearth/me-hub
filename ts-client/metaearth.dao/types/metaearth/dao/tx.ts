/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { DaoAddresses } from "./dao";

export const protobufPackage = "metaearth.dao";

export interface MsgUpdateGlobalDao {
  creator: string;
  daoAddresses: DaoAddresses | undefined;
}

export interface MsgUpdateGlobalDaoResponse {
}

function createBaseMsgUpdateGlobalDao(): MsgUpdateGlobalDao {
  return { creator: "", daoAddresses: undefined };
}

export const MsgUpdateGlobalDao = {
  encode(message: MsgUpdateGlobalDao, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.daoAddresses !== undefined) {
      DaoAddresses.encode(message.daoAddresses, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateGlobalDao {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateGlobalDao();
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

  fromJSON(object: any): MsgUpdateGlobalDao {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      daoAddresses: isSet(object.daoAddresses) ? DaoAddresses.fromJSON(object.daoAddresses) : undefined,
    };
  },

  toJSON(message: MsgUpdateGlobalDao): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.daoAddresses !== undefined
      && (obj.daoAddresses = message.daoAddresses ? DaoAddresses.toJSON(message.daoAddresses) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateGlobalDao>, I>>(object: I): MsgUpdateGlobalDao {
    const message = createBaseMsgUpdateGlobalDao();
    message.creator = object.creator ?? "";
    message.daoAddresses = (object.daoAddresses !== undefined && object.daoAddresses !== null)
      ? DaoAddresses.fromPartial(object.daoAddresses)
      : undefined;
    return message;
  },
};

function createBaseMsgUpdateGlobalDaoResponse(): MsgUpdateGlobalDaoResponse {
  return {};
}

export const MsgUpdateGlobalDaoResponse = {
  encode(_: MsgUpdateGlobalDaoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateGlobalDaoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateGlobalDaoResponse();
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

  fromJSON(_: any): MsgUpdateGlobalDaoResponse {
    return {};
  },

  toJSON(_: MsgUpdateGlobalDaoResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUpdateGlobalDaoResponse>, I>>(_: I): MsgUpdateGlobalDaoResponse {
    const message = createBaseMsgUpdateGlobalDaoResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  UpdateGlobalDao(request: MsgUpdateGlobalDao): Promise<MsgUpdateGlobalDaoResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.UpdateGlobalDao = this.UpdateGlobalDao.bind(this);
  }
  UpdateGlobalDao(request: MsgUpdateGlobalDao): Promise<MsgUpdateGlobalDaoResponse> {
    const data = MsgUpdateGlobalDao.encode(request).finish();
    const promise = this.rpc.request("metaearth.dao.Msg", "UpdateGlobalDao", data);
    return promise.then((data) => MsgUpdateGlobalDaoResponse.decode(new _m0.Reader(data)));
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
