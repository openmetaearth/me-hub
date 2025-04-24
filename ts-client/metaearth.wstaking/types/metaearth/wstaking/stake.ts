/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Timestamp } from "../../google/protobuf/timestamp";

export const protobufPackage = "metaearth.wstaking";

/**
 * Stake represents the bond with tokens held by an account. It is
 * owned by one staker, and is associated with the voting power of one
 * validator.
 */
export interface Stake {
  /** staker_address is the bech32-encoded address of the staker. */
  stakerAddress: string;
  /** validator_address is the bech32-encoded address of the validator. */
  validatorAddress: string;
  /** shares define the stake shares received. */
  shares: string;
  startHeight: number;
  rewards: string;
  amount: string;
  unmovable: string;
}

/**
 * UnbondingStake stores all of a single staker's unbonding bonds
 * for a single validator in an time-ordered list.
 */
export interface UnbondingStake {
  /** staker_address is the bech32-encoded address of the staker. */
  stakerAddress: string;
  /** validator_address is the bech32-encoded address of the validator. */
  validatorAddress: string;
  /** entries are the unbonding delegation entries. */
  entries: UnbondingStakeEntry[];
}

/** UnbondingStakeEntry defines an unbonding object with relevant metadata. */
export interface UnbondingStakeEntry {
  /** creation_height is the height which the unbonding took place. */
  creationHeight: number;
  /** completion_time is the unix time for unbonding completion. */
  completionTime:
    | Date
    | undefined;
  /**
   * initial_balance defines the tokens initially scheduled to receive at
   * completion.
   */
  initialBalance: string;
  /** balance defines the tokens to receive at completion. */
  balance: string;
}

/**
 * SVPair is struct that just has a staker-validator pair with no other data.
 * It is intended to be used as a marshalable pointer. For example, a SVPair can
 * be used to construct the key to getting an UnbondingStake from state.
 */
export interface SVPair {
  stakerAddress: string;
  validatorAddress: string;
}

/** SVPairs defines an array of SVPair objects. */
export interface SVPairs {
  pairs: SVPair[];
}

function createBaseStake(): Stake {
  return {
    stakerAddress: "",
    validatorAddress: "",
    shares: "",
    startHeight: 0,
    rewards: "",
    amount: "",
    unmovable: "",
  };
}

export const Stake = {
  encode(message: Stake, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakerAddress !== "") {
      writer.uint32(10).string(message.stakerAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    if (message.shares !== "") {
      writer.uint32(26).string(message.shares);
    }
    if (message.startHeight !== 0) {
      writer.uint32(32).int64(message.startHeight);
    }
    if (message.rewards !== "") {
      writer.uint32(42).string(message.rewards);
    }
    if (message.amount !== "") {
      writer.uint32(50).string(message.amount);
    }
    if (message.unmovable !== "") {
      writer.uint32(58).string(message.unmovable);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Stake {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStake();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakerAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.shares = reader.string();
          break;
        case 4:
          message.startHeight = longToNumber(reader.int64() as Long);
          break;
        case 5:
          message.rewards = reader.string();
          break;
        case 6:
          message.amount = reader.string();
          break;
        case 7:
          message.unmovable = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Stake {
    return {
      stakerAddress: isSet(object.stakerAddress) ? String(object.stakerAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
      shares: isSet(object.shares) ? String(object.shares) : "",
      startHeight: isSet(object.startHeight) ? Number(object.startHeight) : 0,
      rewards: isSet(object.rewards) ? String(object.rewards) : "",
      amount: isSet(object.amount) ? String(object.amount) : "",
      unmovable: isSet(object.unmovable) ? String(object.unmovable) : "",
    };
  },

  toJSON(message: Stake): unknown {
    const obj: any = {};
    message.stakerAddress !== undefined && (obj.stakerAddress = message.stakerAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    message.shares !== undefined && (obj.shares = message.shares);
    message.startHeight !== undefined && (obj.startHeight = Math.round(message.startHeight));
    message.rewards !== undefined && (obj.rewards = message.rewards);
    message.amount !== undefined && (obj.amount = message.amount);
    message.unmovable !== undefined && (obj.unmovable = message.unmovable);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Stake>, I>>(object: I): Stake {
    const message = createBaseStake();
    message.stakerAddress = object.stakerAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    message.shares = object.shares ?? "";
    message.startHeight = object.startHeight ?? 0;
    message.rewards = object.rewards ?? "";
    message.amount = object.amount ?? "";
    message.unmovable = object.unmovable ?? "";
    return message;
  },
};

function createBaseUnbondingStake(): UnbondingStake {
  return { stakerAddress: "", validatorAddress: "", entries: [] };
}

export const UnbondingStake = {
  encode(message: UnbondingStake, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakerAddress !== "") {
      writer.uint32(10).string(message.stakerAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    for (const v of message.entries) {
      UnbondingStakeEntry.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnbondingStake {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnbondingStake();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakerAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        case 3:
          message.entries.push(UnbondingStakeEntry.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbondingStake {
    return {
      stakerAddress: isSet(object.stakerAddress) ? String(object.stakerAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
      entries: Array.isArray(object?.entries) ? object.entries.map((e: any) => UnbondingStakeEntry.fromJSON(e)) : [],
    };
  },

  toJSON(message: UnbondingStake): unknown {
    const obj: any = {};
    message.stakerAddress !== undefined && (obj.stakerAddress = message.stakerAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    if (message.entries) {
      obj.entries = message.entries.map((e) => e ? UnbondingStakeEntry.toJSON(e) : undefined);
    } else {
      obj.entries = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UnbondingStake>, I>>(object: I): UnbondingStake {
    const message = createBaseUnbondingStake();
    message.stakerAddress = object.stakerAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    message.entries = object.entries?.map((e) => UnbondingStakeEntry.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUnbondingStakeEntry(): UnbondingStakeEntry {
  return { creationHeight: 0, completionTime: undefined, initialBalance: "", balance: "" };
}

export const UnbondingStakeEntry = {
  encode(message: UnbondingStakeEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creationHeight !== 0) {
      writer.uint32(8).int64(message.creationHeight);
    }
    if (message.completionTime !== undefined) {
      Timestamp.encode(toTimestamp(message.completionTime), writer.uint32(18).fork()).ldelim();
    }
    if (message.initialBalance !== "") {
      writer.uint32(26).string(message.initialBalance);
    }
    if (message.balance !== "") {
      writer.uint32(34).string(message.balance);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnbondingStakeEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnbondingStakeEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creationHeight = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.completionTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 3:
          message.initialBalance = reader.string();
          break;
        case 4:
          message.balance = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UnbondingStakeEntry {
    return {
      creationHeight: isSet(object.creationHeight) ? Number(object.creationHeight) : 0,
      completionTime: isSet(object.completionTime) ? fromJsonTimestamp(object.completionTime) : undefined,
      initialBalance: isSet(object.initialBalance) ? String(object.initialBalance) : "",
      balance: isSet(object.balance) ? String(object.balance) : "",
    };
  },

  toJSON(message: UnbondingStakeEntry): unknown {
    const obj: any = {};
    message.creationHeight !== undefined && (obj.creationHeight = Math.round(message.creationHeight));
    message.completionTime !== undefined && (obj.completionTime = message.completionTime.toISOString());
    message.initialBalance !== undefined && (obj.initialBalance = message.initialBalance);
    message.balance !== undefined && (obj.balance = message.balance);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UnbondingStakeEntry>, I>>(object: I): UnbondingStakeEntry {
    const message = createBaseUnbondingStakeEntry();
    message.creationHeight = object.creationHeight ?? 0;
    message.completionTime = object.completionTime ?? undefined;
    message.initialBalance = object.initialBalance ?? "";
    message.balance = object.balance ?? "";
    return message;
  },
};

function createBaseSVPair(): SVPair {
  return { stakerAddress: "", validatorAddress: "" };
}

export const SVPair = {
  encode(message: SVPair, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakerAddress !== "") {
      writer.uint32(10).string(message.stakerAddress);
    }
    if (message.validatorAddress !== "") {
      writer.uint32(18).string(message.validatorAddress);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SVPair {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSVPair();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.stakerAddress = reader.string();
          break;
        case 2:
          message.validatorAddress = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SVPair {
    return {
      stakerAddress: isSet(object.stakerAddress) ? String(object.stakerAddress) : "",
      validatorAddress: isSet(object.validatorAddress) ? String(object.validatorAddress) : "",
    };
  },

  toJSON(message: SVPair): unknown {
    const obj: any = {};
    message.stakerAddress !== undefined && (obj.stakerAddress = message.stakerAddress);
    message.validatorAddress !== undefined && (obj.validatorAddress = message.validatorAddress);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SVPair>, I>>(object: I): SVPair {
    const message = createBaseSVPair();
    message.stakerAddress = object.stakerAddress ?? "";
    message.validatorAddress = object.validatorAddress ?? "";
    return message;
  },
};

function createBaseSVPairs(): SVPairs {
  return { pairs: [] };
}

export const SVPairs = {
  encode(message: SVPairs, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.pairs) {
      SVPair.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SVPairs {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSVPairs();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pairs.push(SVPair.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SVPairs {
    return { pairs: Array.isArray(object?.pairs) ? object.pairs.map((e: any) => SVPair.fromJSON(e)) : [] };
  },

  toJSON(message: SVPairs): unknown {
    const obj: any = {};
    if (message.pairs) {
      obj.pairs = message.pairs.map((e) => e ? SVPair.toJSON(e) : undefined);
    } else {
      obj.pairs = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SVPairs>, I>>(object: I): SVPairs {
    const message = createBaseSVPairs();
    message.pairs = object.pairs?.map((e) => SVPair.fromPartial(e)) || [];
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
