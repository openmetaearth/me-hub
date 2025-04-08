/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dymensionxyz.dymension.delayedack";

/** Params defines the parameters for the module. */
export interface Params {
  epochIdentifier: string;
  bridgingFee: string;
  /**
   * `delete_packets_epoch_limit` is the hard limit of the number of finalized
   * rollapp packets that will be deleted from the store on every epoch end. As
   * deleting finalized rollapp packets is meant to keep the store from growing,
   * it is more of a "nice to have" rather than a "must have" feature,
   * this is a way to limit the time it takes to do so,
   * even if it means potentially causing the store to temporarily grow by
   * piling up packets that weren't deleted but rather "postponed", to
   * subsequent epochs.
   */
  deletePacketsEpochLimit: number;
}

function createBaseParams(): Params {
  return { epochIdentifier: "", bridgingFee: "", deletePacketsEpochLimit: 0 };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.epochIdentifier !== "") {
      writer.uint32(10).string(message.epochIdentifier);
    }
    if (message.bridgingFee !== "") {
      writer.uint32(18).string(message.bridgingFee);
    }
    if (message.deletePacketsEpochLimit !== 0) {
      writer.uint32(24).int32(message.deletePacketsEpochLimit);
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
          message.epochIdentifier = reader.string();
          break;
        case 2:
          message.bridgingFee = reader.string();
          break;
        case 3:
          message.deletePacketsEpochLimit = reader.int32();
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
      epochIdentifier: isSet(object.epochIdentifier) ? String(object.epochIdentifier) : "",
      bridgingFee: isSet(object.bridgingFee) ? String(object.bridgingFee) : "",
      deletePacketsEpochLimit: isSet(object.deletePacketsEpochLimit) ? Number(object.deletePacketsEpochLimit) : 0,
    };
  },

  toJSON(message: Params): unknown {
    const obj: any = {};
    message.epochIdentifier !== undefined && (obj.epochIdentifier = message.epochIdentifier);
    message.bridgingFee !== undefined && (obj.bridgingFee = message.bridgingFee);
    message.deletePacketsEpochLimit !== undefined
      && (obj.deletePacketsEpochLimit = Math.round(message.deletePacketsEpochLimit));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Params>, I>>(object: I): Params {
    const message = createBaseParams();
    message.epochIdentifier = object.epochIdentifier ?? "";
    message.bridgingFee = object.bridgingFee ?? "";
    message.deletePacketsEpochLimit = object.deletePacketsEpochLimit ?? 0;
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
