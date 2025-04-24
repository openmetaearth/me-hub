/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../../../cosmos/base/v1beta1/coin";
import { Any } from "../../../google/protobuf/any";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { Description } from "./description";

export const protobufPackage = "dymensionxyz.dymension.sequencer";

/** MsgCreateSequencer defines a SDK message for creating a new sequencer. */
export interface MsgCreateSequencer {
  /**
   * creator is the bech32-encoded address of the sequencer account which is the
   * account that the message was sent from.
   */
  creator: string;
  /**
   * pubkey is the public key of the sequencers' dymint client, as a Protobuf
   * Any.
   */
  dymintPubKey:
    | Any
    | undefined;
  /** rollappId defines the rollapp to which the sequencer belongs. */
  rollappId: string;
  /** description defines the descriptive terms for the sequencer. */
  description: Description | undefined;
  bond: Coin | undefined;
}

export interface MsgCreateSequencerResponse {
}

/**
 * MsgUnbond defines a SDK message for performing an undelegation from a
 * bond and a sequencer.
 */
export interface MsgUnbond {
  creator: string;
}

/** MsgUnbondResponse defines the Msg/Unbond response type. */
export interface MsgUnbondResponse {
  completionTime: Date | undefined;
}

function createBaseMsgCreateSequencer(): MsgCreateSequencer {
  return { creator: "", dymintPubKey: undefined, rollappId: "", description: undefined, bond: undefined };
}

export const MsgCreateSequencer = {
  encode(message: MsgCreateSequencer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.dymintPubKey !== undefined) {
      Any.encode(message.dymintPubKey, writer.uint32(18).fork()).ldelim();
    }
    if (message.rollappId !== "") {
      writer.uint32(26).string(message.rollappId);
    }
    if (message.description !== undefined) {
      Description.encode(message.description, writer.uint32(34).fork()).ldelim();
    }
    if (message.bond !== undefined) {
      Coin.encode(message.bond, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateSequencer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateSequencer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.dymintPubKey = Any.decode(reader, reader.uint32());
          break;
        case 3:
          message.rollappId = reader.string();
          break;
        case 4:
          message.description = Description.decode(reader, reader.uint32());
          break;
        case 5:
          message.bond = Coin.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCreateSequencer {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      dymintPubKey: isSet(object.dymintPubKey) ? Any.fromJSON(object.dymintPubKey) : undefined,
      rollappId: isSet(object.rollappId) ? String(object.rollappId) : "",
      description: isSet(object.description) ? Description.fromJSON(object.description) : undefined,
      bond: isSet(object.bond) ? Coin.fromJSON(object.bond) : undefined,
    };
  },

  toJSON(message: MsgCreateSequencer): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.dymintPubKey !== undefined
      && (obj.dymintPubKey = message.dymintPubKey ? Any.toJSON(message.dymintPubKey) : undefined);
    message.rollappId !== undefined && (obj.rollappId = message.rollappId);
    message.description !== undefined
      && (obj.description = message.description ? Description.toJSON(message.description) : undefined);
    message.bond !== undefined && (obj.bond = message.bond ? Coin.toJSON(message.bond) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateSequencer>, I>>(object: I): MsgCreateSequencer {
    const message = createBaseMsgCreateSequencer();
    message.creator = object.creator ?? "";
    message.dymintPubKey = (object.dymintPubKey !== undefined && object.dymintPubKey !== null)
      ? Any.fromPartial(object.dymintPubKey)
      : undefined;
    message.rollappId = object.rollappId ?? "";
    message.description = (object.description !== undefined && object.description !== null)
      ? Description.fromPartial(object.description)
      : undefined;
    message.bond = (object.bond !== undefined && object.bond !== null) ? Coin.fromPartial(object.bond) : undefined;
    return message;
  },
};

function createBaseMsgCreateSequencerResponse(): MsgCreateSequencerResponse {
  return {};
}

export const MsgCreateSequencerResponse = {
  encode(_: MsgCreateSequencerResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateSequencerResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateSequencerResponse();
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

  fromJSON(_: any): MsgCreateSequencerResponse {
    return {};
  },

  toJSON(_: MsgCreateSequencerResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCreateSequencerResponse>, I>>(_: I): MsgCreateSequencerResponse {
    const message = createBaseMsgCreateSequencerResponse();
    return message;
  },
};

function createBaseMsgUnbond(): MsgUnbond {
  return { creator: "" };
}

export const MsgUnbond = {
  encode(message: MsgUnbond, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnbond {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnbond();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUnbond {
    return { creator: isSet(object.creator) ? String(object.creator) : "" };
  },

  toJSON(message: MsgUnbond): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnbond>, I>>(object: I): MsgUnbond {
    const message = createBaseMsgUnbond();
    message.creator = object.creator ?? "";
    return message;
  },
};

function createBaseMsgUnbondResponse(): MsgUnbondResponse {
  return { completionTime: undefined };
}

export const MsgUnbondResponse = {
  encode(message: MsgUnbondResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.completionTime !== undefined) {
      Timestamp.encode(toTimestamp(message.completionTime), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUnbondResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUnbondResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.completionTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgUnbondResponse {
    return { completionTime: isSet(object.completionTime) ? fromJsonTimestamp(object.completionTime) : undefined };
  },

  toJSON(message: MsgUnbondResponse): unknown {
    const obj: any = {};
    message.completionTime !== undefined && (obj.completionTime = message.completionTime.toISOString());
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgUnbondResponse>, I>>(object: I): MsgUnbondResponse {
    const message = createBaseMsgUnbondResponse();
    message.completionTime = object.completionTime ?? undefined;
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  /** CreateSequencer defines a method for creating a new sequencer. */
  CreateSequencer(request: MsgCreateSequencer): Promise<MsgCreateSequencerResponse>;
  /** Unbond defines a method for removing coins from sequencer's bond */
  Unbond(request: MsgUnbond): Promise<MsgUnbondResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.CreateSequencer = this.CreateSequencer.bind(this);
    this.Unbond = this.Unbond.bind(this);
  }
  CreateSequencer(request: MsgCreateSequencer): Promise<MsgCreateSequencerResponse> {
    const data = MsgCreateSequencer.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.sequencer.Msg", "CreateSequencer", data);
    return promise.then((data) => MsgCreateSequencerResponse.decode(new _m0.Reader(data)));
  }

  Unbond(request: MsgUnbond): Promise<MsgUnbondResponse> {
    const data = MsgUnbond.encode(request).finish();
    const promise = this.rpc.request("dymensionxyz.dymension.sequencer.Msg", "Unbond", data);
    return promise.then((data) => MsgUnbondResponse.decode(new _m0.Reader(data)));
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

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
