import { Reader, Writer } from "protobufjs/minimal";

interface WireCodec<T> {
  encode(message: T): { finish(): Uint8Array };
  decode(input: Uint8Array, length?: number): T;
}

/**
 * Bridge ts-proto 2.x wire codecs to CosmJS 0.27's protobufjs contract.
 * Keep the generated decoder (including its recursion limit) and all generated
 * JSON/fromPartial helpers. Never prefix a complete message with bytes().
 */
export function toCosmjsGeneratedType<T, C>(codec: C & WireCodec<T>): Omit<C, "encode" | "decode"> & {
  encode(message: T, writer?: Writer): Writer;
  decode(input: Uint8Array | Reader, length?: number): T;
} {
  return {
    ...codec,
    encode(message: T, writer: Writer = Writer.create()): Writer {
      // protobufjs has no public raw-byte append API. Copy each wire field via
      // its public API instead; length-delimited payloads remain opaque, so
      // nested messages, strings, packed fields and byte arrays are unchanged.
      const reader = Reader.create(codec.encode(message).finish());
      while (reader.pos < reader.len) {
        const tag = reader.uint32();
        if (tag >>> 3 === 0) {
          throw new Error("Invalid protobuf field number");
        }
        writer.uint32(tag);
        switch (tag & 7) {
          case 0:
            writer.uint64(reader.uint64());
            break;
          case 1:
            writer.fixed32(reader.fixed32());
            writer.fixed32(reader.fixed32());
            break;
          case 2:
            writer.bytes(reader.bytes());
            break;
          case 3:
          case 4:
            // Deprecated groups have no payload beyond the start/end tags.
            break;
          case 5:
            writer.fixed32(reader.fixed32());
            break;
          default:
            throw new Error("Invalid protobuf wire type");
        }
      }
      return writer;
    },
    decode(input: Uint8Array | Reader, length?: number): T {
      // Do not rely on Reader constructor identity: CosmJS can resolve a
      // different protobufjs patch version than this package's direct one.
      if (input instanceof Uint8Array) {
        return codec.decode(input, length);
      }
      const end = length === undefined ? input.len : input.pos + length;
      if (end < input.pos || end > input.len) {
        throw new RangeError("Protobuf message length exceeds reader bounds");
      }
      const message = codec.decode(input.buf.subarray(input.pos, end));
      input.pos = end;
      return message;
    },
  };
}
