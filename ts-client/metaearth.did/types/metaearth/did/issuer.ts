/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "metaearth.did";

/** credential service status enum */
export enum ServiceStatus {
  SERVICE_STATUS_INACTIVE = 0,
  SERVICE_STATUS_ACTIVE = 1,
  UNRECOGNIZED = -1,
}

export function serviceStatusFromJSON(object: any): ServiceStatus {
  switch (object) {
    case 0:
    case "SERVICE_STATUS_INACTIVE":
      return ServiceStatus.SERVICE_STATUS_INACTIVE;
    case 1:
    case "SERVICE_STATUS_ACTIVE":
      return ServiceStatus.SERVICE_STATUS_ACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ServiceStatus.UNRECOGNIZED;
  }
}

export function serviceStatusToJSON(object: ServiceStatus): string {
  switch (object) {
    case ServiceStatus.SERVICE_STATUS_INACTIVE:
      return "SERVICE_STATUS_INACTIVE";
    case ServiceStatus.SERVICE_STATUS_ACTIVE:
      return "SERVICE_STATUS_ACTIVE";
    case ServiceStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** credential service */
export interface Service {
  sid: string;
  name: string;
  description: string;
  issuers: string[];
  status: ServiceStatus;
}

function createBaseService(): Service {
  return { sid: "", name: "", description: "", issuers: [], status: 0 };
}

export const Service = {
  encode(message: Service, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sid !== "") {
      writer.uint32(10).string(message.sid);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    for (const v of message.issuers) {
      writer.uint32(34).string(v!);
    }
    if (message.status !== 0) {
      writer.uint32(40).int32(message.status);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Service {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseService();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sid = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.description = reader.string();
          break;
        case 4:
          message.issuers.push(reader.string());
          break;
        case 5:
          message.status = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Service {
    return {
      sid: isSet(object.sid) ? String(object.sid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      issuers: Array.isArray(object?.issuers) ? object.issuers.map((e: any) => String(e)) : [],
      status: isSet(object.status) ? serviceStatusFromJSON(object.status) : 0,
    };
  },

  toJSON(message: Service): unknown {
    const obj: any = {};
    message.sid !== undefined && (obj.sid = message.sid);
    message.name !== undefined && (obj.name = message.name);
    message.description !== undefined && (obj.description = message.description);
    if (message.issuers) {
      obj.issuers = message.issuers.map((e) => e);
    } else {
      obj.issuers = [];
    }
    message.status !== undefined && (obj.status = serviceStatusToJSON(message.status));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Service>, I>>(object: I): Service {
    const message = createBaseService();
    message.sid = object.sid ?? "";
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.issuers = object.issuers?.map((e) => e) || [];
    message.status = object.status ?? 0;
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
