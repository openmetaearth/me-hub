/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.rollapp";

export interface DeployerParams {
  /**
   * address is a bech32-encoded address of the
   * accounts that are allowed to create a rollapp.
   */
  address: string;
}

/** Params defines the parameters for the module. */
export interface Params {
  /**
   * dispute_period_in_blocks the number of blocks it takes
   * to change a status of a state from received to finalized.
   * during that period, any user could submit fraud proof
   */
  disputePeriodInBlocks: number;
  /**
   * deployer_whitelist is a list of the
   * accounts that are allowed to create a rollapp and maximum number of
   * rollapps. In the case of an empty list, there are no restrictions
   */
  deployerWhitelist: DeployerParams[];
  rollappsEnabled: boolean;
}

function createBaseDeployerParams(): DeployerParams {
  return { address: "" };
}

export const DeployerParams = {
  encode(message: DeployerParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeployerParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeployerParams();
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

  fromJSON(object: any): DeployerParams {
    return { address: isSet(object.address) ? String(object.address) : "" };
  },

  toJSON(message: DeployerParams): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeployerParams>, I>>(object: I): DeployerParams {
    const message = createBaseDeployerParams();
    message.address = object.address ?? "";
    return message;
  },
};

function createBaseParams(): Params {
  return { disputePeriodInBlocks: 0, deployerWhitelist: [], rollappsEnabled: false };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.disputePeriodInBlocks !== 0) {
      writer.uint32(8).uint64(message.disputePeriodInBlocks);
    }
    for (const v of message.deployerWhitelist) {
      DeployerParams.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.rollappsEnabled === true) {
      writer.uint32(24).bool(message.rollappsEnabled);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.disputePeriodInBlocks = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.deployerWhitelist.push(DeployerParams.decode(reader, reader.uint32()));
          break;
        case 3:
          message.rollappsEnabled = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Params {
    return {
      disputePeriodInBlocks: isSet(object.disputePeriodInBlocks) ? Number(object.disputePeriodInBlocks) : 0,
      deployerWhitelist: Array.isArray(object?.deployerWhitelist)
        ? object.deployerWhitelist.map((e: any) => DeployerParams.fromJSON(e))
        : [],
      rollappsEnabled: isSet(object.rollappsEnabled) ? Boolean(object.rollappsEnabled) : false,
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.disputePeriodInBlocks !== undefined
      && (obj.disputePeriodInBlocks = Math.round(message.disputePeriodInBlocks));
    if (message.deployerWhitelist) {
      obj.deployerWhitelist = message.deployerWhitelist.map((e) => e ? DeployerParams.toJSON(e) : undefined);
    } else {
      obj.deployerWhitelist = [];
    }
    message.rollappsEnabled !== undefined && (obj.rollappsEnabled = message.rollappsEnabled);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.disputePeriodInBlocks = object.disputePeriodInBlocks ?? 0;
    message.deployerWhitelist = object.deployerWhitelist?.map((e) => DeployerParams.fromPartial(e)) || [];
    message.rollappsEnabled = object.rollappsEnabled ?? false;
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
