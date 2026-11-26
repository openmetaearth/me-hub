// Run with `npm run typecheck && npm run test:registry` after installing this
// client in an isolated copy. Transpile only for runtime execution; the separate
// noEmit command checks the entire client using its unchanged tsconfig.
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const ts = require("typescript");
const { Registry } = require("@cosmjs/proto-signing");
const { Reader, Writer } = require("protobufjs/minimal");
const { BinaryWriter } = require("@bufbuild/protobuf/wire");

require.extensions[".ts"] = (module, filename) => {
  const { outputText } = ts.transpileModule(fs.readFileSync(filename, "utf8"), {
    compilerOptions: {
      target: ts.ScriptTarget.ES2020,
      module: ts.ModuleKind.CommonJS,
      esModuleInterop: true,
    },
    fileName: filename,
  });
  module._compile(outputText, filename);
};

const { toCosmjsGeneratedType } = require("../registry-codec.ts");
const regiongov = require("../metaearth.regiongov/types/metaearth/regiongov/tx.ts");
const { msgTypes: regiongovTypes, registry: regiongovRegistry } = require("../metaearth.regiongov/registry.ts");
const kyc = require("../metaearth.kyc/types/metaearth/kyc/tx.ts");
const { msgTypes: kycTypes } = require("../metaearth.kyc/registry.ts");
const root = require("../index.ts");
const registry = new Registry([...kycTypes, ...regiongovTypes]);

// Non-default scalars, repeated values, nested messages and decimal coin
// amounts ensure this is not merely an empty-message compatibility check.
const fixtures = {
  MsgSubmitProposal: {
    proposer: "me1proposer",
    kind: 1,
    title: "Treasury proposal",
    summary: "nested coin payload",
    votingPeriod: 172800,
    treasury: {
      withdrawals: [
        { regionId: "cn", receiver: "me1receiver", amount: [{ denom: "ume", amount: "18446744073709551615" }] },
        { regionId: "us", receiver: "me1other", amount: [{ denom: "ucoin", amount: "12345" }] },
      ],
    },
  },
  MsgVote: { voter: "me1voter", proposalId: 4294967303, option: 2 },
  MsgSubmitEquityUpdate: { proposer: "me1proposer", regionId: "cn", addresses: ["me1a", "me1b"] },
  MsgSubmitFreezeEquity: { proposer: "me1proposer", regionId: "cn", addresses: ["me1a", "me1b"] },
  MsgSubmitUnfreezeEquity: { proposer: "me1proposer", regionId: "cn", addresses: ["me1a", "me1b"] },
  MsgApprovePendingAction: { authority: "me1dao", actionId: 4294967305 },
  MsgRejectPendingAction: { authority: "me1dao", actionId: 91 },
  MsgFreezeRegion: { authority: "me1dao", regionIds: ["cn", "us"] },
  MsgUnfreezeRegion: { authority: "me1dao", regionIds: ["cn", "us"] },
  MsgFreezeEquity: { authority: "me1dao", regionId: "cn", address: "me1equity" },
  MsgUnfreezeEquity: { authority: "me1dao", regionId: "cn", address: "me1equity" },
};
assert.equal(regiongovTypes.length, 11);
assert.equal(root.MetaearthRegiongov.msgTypes.length, 11);

let roundtrips = 0;
function checkMessage(typeUrl, codec, partial) {
  const value = codec.fromPartial(partial);
  const expected = codec.encode(value).finish();
  const object = { typeUrl, value };
  const registered = registry.lookupType(typeUrl);
  assert.ok(registered, typeUrl);
  assert.deepEqual(Uint8Array.from(registry.encode(object)), expected, `${typeUrl}: exact wire, no prefix`);
  assert.deepEqual(registry.decode({ typeUrl, value: expected }), value, `${typeUrl}: Registry roundtrip`);
  assert.deepEqual(Uint8Array.from(registry.encodeAsAny(object).value), expected);
  assert.deepEqual(root.registry.decode({ typeUrl, value: root.registry.encode(object) }), value);
  const body = registry.decodeTxBody(registry.encodeTxBody({ messages: [object], memo: "adapter regression" }));
  assert.deepEqual(body.messages[0], value, `${typeUrl}: TxBody Any roundtrip`);
  assert.equal(registered.fromPartial, codec.fromPartial);
  assert.equal(registered.fromJSON, codec.fromJSON);
  assert.deepEqual(registered.toJSON(value), codec.toJSON(value));

  // The supplied writer must be appended in place, including fork/ldelim.
  const supplied = Writer.create().uint32(10).fork();
  assert.equal(registered.encode(value, supplied), supplied);
  const framed = Reader.create(supplied.ldelim().uint32(16).uint32(73).finish());
  assert.equal(framed.uint32(), 10);
  const length = framed.uint32();
  assert.deepEqual(codec.toJSON(registered.decode(framed, length)), codec.toJSON(value));
  assert.equal(framed.uint32(), 16);
  assert.equal(framed.uint32(), 73);
  assert.equal(framed.pos, framed.len);
  const invalid = Reader.create(expected);
  assert.throws(() => registered.decode(invalid, expected.length + 1), RangeError);
  assert.equal(invalid.pos, 0);
  roundtrips++;
}

for (const [typeUrl] of regiongovTypes) {
  const name = typeUrl.substring(typeUrl.lastIndexOf(".") + 1);
  checkMessage(typeUrl, regiongov[name], fixtures[name]);
  assert.deepEqual(regiongovRegistry.decode({ typeUrl, value: regiongovRegistry.encode({ typeUrl, value: fixtures[name] }) }), regiongov[name].fromPartial(fixtures[name]));
}
// Exercise the other proposal variant too.
checkMessage("/metaearth.regiongov.MsgSubmitProposal", regiongov.MsgSubmitProposal, {
  proposer: "me1operator", kind: 2, title: "Replace", summary: "new operator", votingPeriod: 86400,
  replaceOperator: { regionId: "cn", newOperator: "mevaloper1new" },
});
for (const optionType of [0, 1]) {
  checkMessage("/metaearth.kyc.MsgUpdateKycIssuer", kyc.MsgUpdateKycIssuer, {
    creator: "me1dao", addresses: ["me1issuer1", "me1issuer2"], optionType,
  });
}
// Opaque bytes must survive a real KYC registry entry as well.
checkMessage("/metaearth.kyc.MsgCreateSBT", kyc.MsgCreateSBT, {
  issuer: "me1issuer", did: "did:me:123", uri: "https://example.invalid/sbt", uriHash: "hash",
  data: Uint8Array.from([0, 1, 127, 128, 255]),
});

// Cover every supported protobuf wire type and 64-bit boundary bit patterns.
const bytes = new BinaryWriter()
  .uint32(8).uint64("18446744073709551615")
  .uint32(16).int64(-1)
  .uint32(25).fixed64("18446744073709551614")
  .uint32(34).bytes(Uint8Array.from([0, 255, 128, 1]))
  .uint32(45).fixed32(4294967294)
  .uint32(51).uint32(8).uint64("9223372036854775808").uint32(52)
  .finish();
const wireCodec = {
  encode(value) { return { finish: () => value }; },
  decode(input) { return input; },
};
const wireAdapter = toCosmjsGeneratedType(wireCodec);
assert.deepEqual(Uint8Array.from(wireAdapter.encode(bytes).finish()), bytes, "all wire types retain exact bytes");
const append = Writer.create().uint32(8).uint32(7);
assert.equal(wireAdapter.encode(bytes, append), append);
assert.deepEqual(Uint8Array.from(append.finish()), Uint8Array.from([8, 7, ...bytes]));

// Use a different protobufjs installation when CosmJS resolves a nested copy.
const signingDir = path.dirname(require.resolve("@cosmjs/proto-signing/package.json"));
const otherProtobuf = require(require.resolve("protobufjs/minimal", { paths: [signingDir] }));
const foreignReader = otherProtobuf.Reader.create(bytes);
assert.deepEqual(wireAdapter.decode(foreignReader), bytes);
assert.equal(foreignReader.pos, foreignReader.len);

// Generated DescriptorProto is genuinely recursive. Its retained ts-proto
// limit must reject hostile nesting through both Registry and Reader entrypoints.
const { DescriptorProto } = require("../metaearth.regiongov/types/google/protobuf/descriptor.ts");
const recursive = toCosmjsGeneratedType(DescriptorProto);
const recursiveRegistry = new Registry([["/google.protobuf.DescriptorProto", recursive]]);
let nested = new Uint8Array();
for (let i = 0; i < 101; i++) nested = new BinaryWriter().uint32(26).bytes(nested).finish();
assert.throws(() => recursiveRegistry.decode({ typeUrl: "/google.protobuf.DescriptorProto", value: nested }), /protobuf decode recursion limit exceeded/);
assert.throws(() => recursive.decode(Reader.create(nested)), /protobuf decode recursion limit exceeded/);
const shallow = DescriptorProto.fromPartial({ name: "Root", nestedType: [{ name: "Child" }] });
assert.deepEqual(recursive.decode(recursive.encode(shallow).finish()), shallow);

const rollappTypes = require("../metaearth.rollapp/registry.ts").msgTypes;
assert.ok(rollappTypes.some(([url]) => url === "/metaearth.rollapp.MsgUpdateApp"));
assert.ok(!rollappTypes.some(([url]) => url === "/metaearth.rollapp.MsgUpdateRollapp"));
assert.ok(!require("../metaearth.wstaking/registry.ts").msgTypes.some(([url]) => url === "/metaearth.wstaking.MsgWithdrawFromRegion"));
console.log(`PASS: ${roundtrips} non-default KYC/regiongov Registry and supplied-Writer roundtrips; exact 64-bit/bytes wire preservation; recursion limit retained`);
