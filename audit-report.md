# PR Audit Report

Generated: 2026-06-22 01:18 UTC
Total open PRs fetched: **267**

## Summary

| Category | Count |
|----------|-------|
| Auto-close candidates | **0** |
| Keep open (manual review) | **267** |

**Auto-close criteria:**
- Has merge conflicts
- Bot/AI account (detected via GitHub API type=Bot or app/ prefix) with zero human engagement
- Stale draft PR with zero engagement

**Bot accounts detected:** none

**Keep criteria (any one is enough to preserve):**
- Author is not a known bot account
- PR has at least one comment or review

> **Next steps:**
> 1. Review the Auto-Close list below, remove any lines you want to preserve
> 2. Dry-run: `./scripts/close-prs.sh --dry-run`
> 3. Execute: `./scripts/close-prs.sh`

---

## Auto-Close List (0 PRs)

| # | Title | Author | Reason | Date | Files |
|---|-------|--------|--------|------|-------|

---

## Keep Open List (manual review)

| # | Title | Author | Reason | Date | Files |
|---|-------|--------|--------|------|-------|
| #1361 | fix(ante): require single wasm message for creator fee | jamilahmadzai | human contributor | 2026-06-19 | 2 |
| #1358 | fix(denommetadata): return update proposal type | jamilahmadzai | human contributor | 2026-06-19 | 2 |
| #1357 | fix(gravity): paginate pending outgoing address query | jamilahmadzai | human contributor | 2026-06-18 | 5 |
| #1356 | fix(eibc): make order fulfillment atomic | jamilahmadzai | human contributor | 2026-06-18 | 2 |
| #1286 | fix(gravity): convert added bridge fee to external prec | sureshchouksey8 | human contributor | 2026-06-15 | 2 |
| #1285 | fix(wstaking): block region reassignment if old region  | sureshchouksey8 | human contributor | 2026-06-15 | 2 |
| #1284 | fix(gravity): clamp batch timeout projection height | jamilahmadzai | human contributor | 2026-06-14 | 2 |
| #1278 | fix(gravity): avoid relayer power overflow | jamilahmadzai | human contributor | 2026-06-14 | 2 |
| #1271 | fix(delayedack): bind packet finalization to state root | jamilahmadzai | human contributor | 2026-06-13 | 8 |
| #1265 | fix(gravity): defer batch cancellation until after iter | jamilahmadzai | human contributor | 2026-06-13 | 3 |
| #1234 | fix(wstaking): block region moves with delegations | jamilahmadzai | human contributor | 2026-06-12 | 2 |
| #1213 | fix(sequencer): allow permissioned sequencers past publ | jamilahmadzai | human contributor | 2026-06-11 | 3 |
| #1193 | fix(rollapp): authenticate transfers from packet endpoi | jamilahmadzai | human contributor | 2026-06-10 | 10 |
| #1088 | fix: reject join when group admin is invalid | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1087 | fix(wstaking): reject blocked region withdraw receivers | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1086 | fix(sequencer): preserve replacement unbonding period | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #1059 | fix(kyc): preserve non-dao issuers on dao update | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1058 | fix(wstaking): allow historical fixed deposit migration | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1057 | fix: restrict rollapp skip-delay to GlobalDao | Peter7896 | human contributor | 2026-06-06 | 3 |
| #1055 | fix(sequencer): return not found for missing handoff | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1053 | fix(dao): reject nil free gas query requests | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1050 | fix(wstaking): return unbond pool transfer errors | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 1 |
| #1048 | fix(wmint): avoid begin blocker mint panics | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1045 | fix(wnft): reject invalid signer addresses | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1044 | fix(rollapp): skip missing fraud rollback states | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1043 | fix(gravity): ignore unknown executed batches | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1042 | fix(bridgingfee): skip fees above transfer amount | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1040 | fix(kyc): burn SBT on KYC removal | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1038 | fix(wstaking): restore runtime genesis state | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1037 | fix(kyc): clear stale DID address mappings | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #1036 | fix(ante): account multisend inputs separately | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1035 | fix(did): revoke issuers on deactivation | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1034 | fix(sequencer): queue replaced proposers for unbonding | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1032 | fix(delayedack): keep refund failures pending | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1031 | fix(delayedack): roll back failed OnRecv setup | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #1030 | fix(upgrades): let module migrations run | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #1028 | fix(rollapp): bound fraud rollback by height | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 6 |
| #1026 | fix(wstaking): export pending consensus key replacement | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 4 |
| #1025 | fix(sequencer): clear stale replace proposer | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1023 | fix(app): allow bridge fee pool receipts | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1020 | fix(rollapp): reject frozen state updates | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1019 | fix(ante): reject oversized gas limits | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1013 | fix(wstaking): ignore inactive unrelated fixed terms | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #1011 | fix(gravity): run relayer slashing in endblock | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 1 |
| #1008 | fix(wstaking): dedupe validator updates | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #1007 | fix(gravity): preserve confirmed timed out batches | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 5 |
| #1006 | fix(gravity): ignore offline attestation votes | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1003 | fix(app): validate gentx message shapes | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #1002 | fix(eibc): require positive fulfill fees | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 5 |
| #1001 | fix: prevent RollApp packet prefix collisions | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 4 |
| #999 | fix: reject empty EVM bank contract denom query | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #998 | fix(eibc): filter demand orders by fulfiller | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #997 | fix(gravity): reject offline relayer confirms | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #993 | fix(megroup): skip duplicate self admin reward | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #992 | fix(wnft): reject non-canonical token ids | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 3 |
| #990 | fix(did): length-delimit credential filter keys | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 6 |
| #988 | fix(tron): avoid panics for invalid external addresses | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #987 | fix(wstaking): reject unsupported transfer region tx | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 4 |
| #986 | fix(gravity): reject nil direct query requests | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #982 | fix(gravity): reject nil router query requests | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #981 | fix(wstaking): reject nil non-region query requests | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #980 | fix(wnft): reject nil nft filter queries | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 1 |
| #979 | fix(delayedack): paginate packet queries | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #978 | fix(wstaking): validate fixed deposit config status | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #977 | fix(sequencer): reject nil replace proposer queries | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #976 | fix(megroup): decode group member counts from raw uint6 | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 1 |
| #975 | fix(wgov): avoid tally division by zero | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 1 |
| #973 | fix(dao): reject module accounts in dao updates | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #972 | fix(wstaking): require combined capacity on unstake | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #971 | fix: commit canonical channel only after recv success | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #970 | fix(ante): require fees for join group transactions | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 1 |
| #968 | fix(gravity): use tron confirm material in cli | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #967 | fix(kyc): preserve issuers during dao rotation | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #966 | fix(wstaking): protect region reserves on withdrawal | modelsbridgeaicom-ship-it | human contributor | 2026-06-06 | 2 |
| #962 | fix(gravity): reject offline relayer confirmations | jamilahmadzai | human contributor | 2026-06-06 | 4 |
| #957 | fix(rollapp): validate genesis rollapp invariants | jamilahmadzai | human contributor | 2026-06-05 | 4 |
| #955 | fix(rollapp): validate genesis state indexes | jamilahmadzai | human contributor | 2026-06-05 | 2 |
| #953 | fix: remove rate equality check in transferDeposit | q3515 | human contributor | 2026-06-05 | 1 |
| #949 | fix: guard relayer set observation nonce | jamilahmadzai | human contributor | 2026-06-05 | 4 |
| #948 | [codex] fix gravity batch timeout underflow | yuzhiyang1 | human contributor | 2026-06-05 | 2 |
| #922 | fix(gravity): only commit state changes on successful a | sureshchouksey8 | human contributor | 2026-06-05 | 1 |
| #920 | fix(gravity): reject any relayer claim whose nonce is n | sureshchouksey8 | human contributor | 2026-06-05 | 1 |
| #919 | fix(rollapp): correct swapped RevisionNumber/RevisionHe | sureshchouksey8 | human contributor | 2026-06-05 | 3 |
| #918 | fix(wdistri): prevent silent int64 truncation in Alloca | sureshchouksey8 | human contributor | 2026-06-05 | 2 |
| #917 | fix(sequencer): add rollapp membership checks in Replac | sureshchouksey8 | human contributor | 2026-06-05 | 2 |
| #916 | fix(sequencer): prevent uint64-to-int64 overflow in Rep | sureshchouksey8 | human contributor | 2026-06-05 | 2 |
| #888 | fix(wnft): prevent DoS by implementing pagination in ow | sureshchouksey8 | human contributor | 2026-06-04 | 5 |
| #886 | fix(gravity): restrict relayer proposals to GlobalDao | x0tta6bl4-ai | human contributor | 2026-06-04 | 3 |
| #885 | fix(wstaking): preserve downtime state on pubkey rotati | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #884 | fix(wstaking): validate global fee pool withdrawal amou | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #882 | fix(vfc): reject stale VFBC metadata updates | x0tta6bl4-ai | human contributor | 2026-06-04 | 4 |
| #879 | fix(gravity): reject duplicate outgoing txs in genesis | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #878 | fix(megroup): allow higher KYC levels to join groups | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #877 | fix(rollapp): preserve skip-delay policy in genesis | x0tta6bl4-ai | human contributor | 2026-06-04 | 5 |
| #876 | fix(dao): preserve free gas accounts in genesis | x0tta6bl4-ai | human contributor | 2026-06-04 | 6 |
| #873 | fix(sequencer): rotate proposer by bond denom determini | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #869 | fix(wstaking): preserve validator updates in genesis | x0tta6bl4-ai | human contributor | 2026-06-04 | 3 |
| #867 | fix(sequencer): block unbond during proposer replacemen | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #866 | fix(rollapp): make fraud handling atomic | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #864 | fix(gravity): preserve higher fee live batches | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #863 | fix(wstaking): price halving boundary rewards correctly | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #855 | fix(rollapp): reject first state height gaps | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #853 | fix(rollapp): restrict frozen eip155 recovery creator | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #852 | fix(dao): reject empty GlobalDao/MeidDao in MsgUpdateDa | AnwarSup | human contributor | 2026-06-04 | 1 |
| #848 | fix(kyc): avoid duplicate validator rewards on re-appro | x0tta6bl4-ai | human contributor | 2026-06-04 | 4 |
| #846 | fix(kyc): clear stale DID address mappings on removal | x0tta6bl4-ai | human contributor | 2026-06-04 | 4 |
| #845 | fix(wstaking): ignore inactive unused fixed-deposit con | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #844 | fix(kyc): reject zero-level updates | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #839 | fix(rollapp): restrict skip-delay control to GlobalDao | x0tta6bl4-ai | human contributor | 2026-06-04 | 3 |
| #838 | fix(did): revoke KYC-derived state when disabling DID | x0tta6bl4-ai | human contributor | 2026-06-04 | 9 |
| #837 | fix(rollapp): require success ack before opening transf | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #836 | fix(gravity): remove BSC fallback from routed queries | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #835 | fix(kyc): reject protocol region writes | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #829 | fix(did): persist partial filter logger deletions | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #828 | fix(rollapp): enable transfers on successful recv ack | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #827 | fix(delayedack): cap FinalizeRollappPackets queue proce | sureshchouksey8 | human contributor | 2026-06-04 | 1 |
| #826 | fix(wgov): fix governance tally bypass issue (#792) | sureshchouksey8 | human contributor | 2026-06-04 | 1 |
| #825 | fix(delayedack): log hook error in finalizeRollappPacke | sureshchouksey8 | human contributor | 2026-06-04 | 1 |
| #824 | fix(wmint): skip zero mint after cap | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #823 | fix(eibc): prevent self-fulfillment of demand orders | sureshchouksey8 | human contributor | 2026-06-04 | 3 |
| #822 | fix(wstaking): reject nil region queries | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #821 | fix(wnft): paginate owner nft queries | x0tta6bl4-ai | human contributor | 2026-06-04 | 8 |
| #817 | fix(megroup): remove members after kyc downgrade | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #816 | fix(wstaking): stop rewards after mint cap | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #815 | fix(delayedack): record timeout-on-close proof height | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #814 | fix(kyc): require approve pubkey to match address | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #811 | fix: reject kyc remove vc through did | x0tta6bl4-ai | human contributor | 2026-06-04 | 3 |
| #810 | fix: remove empty rollapp finalization queues | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #809 | fix: pass through denom metadata ack for non-rollapps | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #808 | fix: enforce rollapp DRS violations | x0tta6bl4-ai | human contributor | 2026-06-04 | 6 |
| #807 | fix: preserve megroup state in genesis | x0tta6bl4-ai | human contributor | 2026-06-04 | 4 |
| #806 | fix: block relayer unbond while external set still trus | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #805 | fix: reject skip-delay post-genesis genesis_transfer pa | x0tta6bl4-ai | human contributor | 2026-06-04 | 2 |
| #803 | fix(wgov): make active tally query read-only | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #802 | fix(ante): reject malformed fee addresses without panic | x0tta6bl4-ai | human contributor | 2026-06-03 | 4 |
| #801 | fix(ante): include fee payer fees in fund check | x0tta6bl4-ai | human contributor | 2026-06-03 | 4 |
| #800 | fix(eibc): validate updated demand orders | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #798 | fix(gravity): guard offline relayer concentration | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #797 | fix(delayedack): log hook error instead of aborting pac | AnwarSup | human contributor | 2026-06-03 | 1 |
| #795 | fix(delayedack): cap FinalizeRollappPackets batch to 10 | AnwarSup | human contributor | 2026-06-03 | 1 |
| #791 | fix(kyc): avoid empty protocol regions | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #790 | fix(ibctesting): check rollapp packet commitments | x0tta6bl4-ai | human contributor | 2026-06-03 | 1 |
| #787 | fix: require active DID for KYC benefits | Aglcr7 | human contributor | 2026-06-03 | 7 |
| #786 | fix(sequencer): paginate rollapp sequencer queries | x0tta6bl4-ai | human contributor | 2026-06-03 | 6 |
| #785 | fix(wstaking): normalize validator region updates | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #784 | fix(rollapp): gate packet-forward transfers | x0tta6bl4-ai | human contributor | 2026-06-03 | 3 |
| #782 | fix(wstaking): check region interest treasury reserve | x0tta6bl4-ai | human contributor | 2026-06-03 | 5 |
| #781 | fix(wgov): tally proposals before execution | x0tta6bl4-ai | human contributor | 2026-06-03 | 1 |
| #780 | fix(wstaking): preserve signing info on key rotation | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #778 | fix(gravity): make CancelOutgoingTxBatch atomic to prev | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #776 | fix(gravity): validate relayer query addresses | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #775 | fix: [Bug]: [Bug]: [Low] [Explorer/UI] Testnet Block He | HMS091 | human contributor | 2026-06-03 | 9 |
| #772 | fix(wstaking): clamp DelegateInterest to zero instead o | sureshchouksey8 | human contributor | 2026-06-03 | 1 |
| #771 | fix(gravity): reject BSC USDT/USDC dust deposits | x0tta6bl4-ai | human contributor | 2026-06-03 | 4 |
| #770 | fix(eibc): validate and propagate expected fee in Fulfi | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #769 | fix(gravity): re-enable relayer slashing in EndBlocker  | sureshchouksey8 | human contributor | 2026-06-03 | 1 |
| #768 | fix(wdistri): remove Int64 overflow in AllocateBlockRew | sureshchouksey8 | human contributor | 2026-06-03 | 1 |
| #767 | fix(sequencer): replace fmt.Errorf with proper Cosmos S | sureshchouksey8 | human contributor | 2026-06-03 | 1 |
| #766 | fix(kyc): avoid phantom protocol regions | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #762 | fix(wgov): delete votes after tally iteration | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #761 | fix(gravity): preserve slash counter on delegate | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #757 | fix(gravity): count batch fees in pending supply | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #756 | fix(sequencer): auto-cancel stale ReplaceProposer reque | sureshchouksey8 | human contributor | 2026-06-03 | 1 |
| #755 | fix(wstaking): call SetInviterReward after sending invi | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #752 | fix: enqueue replaced proposer for unbonding and clean  | sureshchouksey8 | human contributor | 2026-06-03 | 3 |
| #751 | fix: emit kyc remove events and sync megroup membership | sureshchouksey8 | human contributor | 2026-06-03 | 3 |
| #750 | fix(kyc): preserve third-party issuers during DAO addre | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #749 | fix(did): delete old DID filters on VC update (#669) | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #748 | fix(did): reserve kyc credential route | x0tta6bl4-ai | human contributor | 2026-06-03 | 3 |
| #747 | fix(did): persist updated FilterLogger on partial filte | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #746 | fix(kyc): prevent empty region entries in GetAllRegions | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #745 | fix(eibc): implement rollappHasPacketCommitment verific | sureshchouksey8 | human contributor | 2026-06-03 | 1 |
| #744 | fix(gravity): mint delegate coins to module accounts du | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #743 | fix(eibc): validate demand order after fee/price update | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #742 | fix(bridgingfee): fail packet receipt on charge bridgin | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #741 | fix(gravity): enforce proposal relayer unbonding govern | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #740 | fix(gravity): correct outgoing batch timeout check in B | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #739 | fix(rollapp): automatically open SkipDelay RollApp brid | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #738 | fix(kyc): reject duplicate SBT creation | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #737 | fix(gravity): reject duplicate bridge tokens in genesis | sureshchouksey8 | human contributor | 2026-06-03 | 3 |
| #736 | fix(ante): enforce minimum base fee policy during Deliv | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #735 | fix(rollapp): trim rollapp-id in CreateRollapp to preve | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #734 | fix(did): ignore filter prefix collisions | x0tta6bl4-ai | human contributor | 2026-06-03 | 4 |
| #733 | fix(gravity): reject zero external address in MsgBonded | sureshchouksey8 | human contributor | 2026-06-03 | 2 |
| #730 | fix(did): replace VC filter indexes on update | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #729 | fix(wstaking): remove MeID NFT region index entries | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #728 | fix(wstaking): reject stake when region is missing | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #727 | fix(wstaking): execute transfer region message | x0tta6bl4-ai | human contributor | 2026-06-03 | 4 |
| #726 | fix(eibc): resolve errack timeout demand order denom | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #725 | fix(eibc): delete pending demand orders when packets ar | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #724 | fix(eibc): reject malformed fee amounts | x0tta6bl4-ai | human contributor | 2026-06-03 | 2 |
| #721 | fix(bridgingfee): reject transfers when fee charge fail | x0tta6bl4-ai | human contributor | 2026-06-03 | 1 |
| #720 | fix(gravity): add per-token supply cap for SendToMeClai | AnwarSup | human contributor | 2026-06-03 | 2 |
| #716 | fix(sequencer): enqueue replaced proposer for unbonding | AnwarSup | human contributor | 2026-06-03 | 1 |
| #658 | fix: reject duplicate sequencer dymint pubkeys | NikYak228 | human contributor | 2026-06-02 | 2 |
| #657 | fix: reject multiple rollapp genesis proposers | NikYak228 | human contributor | 2026-06-02 | 2 |
| #653 | fix: reject duplicated outgoing tx genesis | NikYak228 | human contributor | 2026-06-02 | 2 |
| #652 | fix: reject cross-rollapp replace proposer | NikYak228 | human contributor | 2026-06-02 | 2 |
| #651 | fix: preserve pending replace proposer genesis | NikYak228 | human contributor | 2026-06-02 | 7 |
| #648 | fix: preserve skip-delay rollapps in genesis | NikYak228 | human contributor | 2026-06-02 | 6 |
| #647 | fix: preserve fixed deposit runtime genesis | NikYak228 | human contributor | 2026-06-02 | 6 |
| #645 | fix: preserve kyc service status on genesis import | NikYak228 | human contributor | 2026-06-02 | 4 |
| #643 | fix: preserve dao free gas genesis accounts | NikYak228 | human contributor | 2026-06-02 | 9 |
| #642 | fix: restore gravity outgoing genesis counters | NikYak228 | human contributor | 2026-06-02 | 3 |
| #640 | fix: allow replace proposer at latest state height | NikYak228 | human contributor | 2026-06-02 | 2 |
| #638 | fix: reserve kyc wnft class id | NikYak228 | human contributor | 2026-06-02 | 5 |
| #637 | fix: reject invalid wnft owner filters | NikYak228 | human contributor | 2026-06-02 | 2 |
| #635 | fix: count authz messages in ante limit | NikYak228 | human contributor | 2026-06-02 | 4 |
| #634 | fix: reject non-canonical wnft token ids | NikYak228 | human contributor | 2026-06-02 | 2 |
| #633 | fix: bind DID credentials to original issuer | eliterdav09-creator | human contributor | 2026-06-02 | 4 |
| #632 | fix: restore remove region handling | NikYak228 | human contributor | 2026-06-02 | 4 |
| #631 | fix: validate global dao fee withdrawal amount | NikYak228 | human contributor | 2026-06-02 | 2 |
| #629 | fix: bound sequencer rollapp index queries | NikYak228 | human contributor | 2026-06-02 | 3 |
| #628 | fix: accept historical relayer set claims | NikYak228 | human contributor | 2026-06-02 | 2 |
| #627 | fix: reject oversized EIP155 rollapp ids | eliterdav09-creator | human contributor | 2026-06-02 | 2 |
| #626 | fix: make kyc genesis issuer import idempotent | NikYak228 | human contributor | 2026-06-02 | 2 |
| #625 | fix(megroup): prevent group admin from self-joining own | sureshchouksey8 | human contributor | 2026-06-02 | 1 |
| #623 | fix(wstaking): validate region DelegateInterest balance | sureshchouksey8 | human contributor | 2026-06-02 | 2 |
| #622 | fix: preserve wmint minted counters in genesis | NikYak228 | human contributor | 2026-06-02 | 3 |
| #621 | fix: restrict rollapp channel updates | NikYak228 | human contributor | 2026-06-02 | 6 |
| #618 | fix: restrict send-to-module targets | NikYak228 | human contributor | 2026-06-02 | 4 |
| #615 | fix: require sender ownership for bridge fee increases | NikYak228 | human contributor | 2026-06-02 | 2 |
| #614 | fix: enforce attestation quorum with ceiling division | NikYak228 | human contributor | 2026-06-01 | 2 |
| #613 | fix: reject overflowing fixed-deposit terms | NikYak228 | human contributor | 2026-06-01 | 4 |
| #612 | fix: remove stale validator power index during zero-hei | NikYak228 | human contributor | 2026-06-01 | 2 |
| #611 | fix: make outgoing batch cancellation atomic | NikYak228 | human contributor | 2026-06-01 | 2 |
| #606 | fix(wstaking): call SetInviterReward after sending invi | AnwarSup | human contributor | 2026-06-01 | 1 |
| #604 | fix: reject invalid delayedack packet status | NikYak228 | human contributor | 2026-06-01 | 2 |
| #603 | fix: gate denom metadata on downstream ack | NikYak228 | human contributor | 2026-06-01 | 3 |
| #602 | fix: clear stale replace proposer requests | NikYak228 | human contributor | 2026-06-01 | 2 |
| #590 | fix: wrap replace proposer errors | NikYak228 | human contributor | 2026-06-01 | 2 |
| #589 | fix: enforce KYC eligibility during group migration | Aglcr7 | human contributor | 2026-06-01 | 2 |
| #583 | fix(gravity): enable relayer slashing | jamilahmadzai | human contributor | 2026-06-01 | 2 |
| #582 | fix: restrict proposal relayer authority | jamilahmadzai | human contributor | 2026-06-01 | 3 |
| #579 | fix(rollapp): correct tendermint client state freeze he | sureshchouksey8 | human contributor | 2026-06-01 | 2 |
| #578 | fix: validate delegate interest before kyc rewards | NikYak228 | human contributor | 2026-06-01 | 5 |
| #577 | fix: validate tron address checksum | NikYak228 | human contributor | 2026-06-01 | 2 |
| #576 | fix: honor fixed deposit query pagination | NikYak228 | human contributor | 2026-06-01 | 2 |
| #575 | fix: reject invalid demand order denoms | NikYak228 | human contributor | 2026-06-01 | 2 |
| #574 | fix: reject invalid fee payer addresses | NikYak228 | human contributor | 2026-06-01 | 5 |
| #573 | fix: reject invalid fixed deposit rates | NikYak228 | human contributor | 2026-06-01 | 2 |
| #572 | fix: block module fee receivers | NikYak228 | human contributor | 2026-06-01 | 4 |
| #571 | fix: ignore expired outgoing tx batches | NikYak228 | human contributor | 2026-06-01 | 3 |
| #570 | fix: cap relayer delegate power increases | NikYak228 | human contributor | 2026-06-01 | 2 |
| #555 | [BUG BOUNTY] fix(wstaking): correct RegionShare account | soongyintong | human contributor | 2026-06-01 | 2 |
| #554 | fix(gravity): atomic CancelOutgoingTxBatch with cache c | MeowMeow1230 | human contributor | 2026-06-01 | 1 |
| #543 | fix: reject sequencer self replacement | jamilahmadzai | human contributor | 2026-05-31 | 4 |
| #542 | fix: require dao update addresses | jamilahmadzai | human contributor | 2026-05-31 | 3 |
| #541 | fix: check authz rollapp transfers | jamilahmadzai | human contributor | 2026-05-31 | 2 |
| #509 | fix: reject none level KYC approvals | Aglcr7 | human contributor | 2026-05-31 | 4 |
| #506 | fix: validate relayer addresses in gravity queries | Aglcr7 | human contributor | 2026-05-31 | 4 |
| #504 | fix: reject nil Gravity relayer set queries | Aglcr7 | human contributor | 2026-05-31 | 2 |
| #503 | fix: avoid phantom KYC protocol regions | Aglcr7 | human contributor | 2026-05-31 | 2 |
| #491 | fix(gravity): resolve BSC USDT/USDC decimal conversion  | sureshchouksey8 | human contributor | 2026-05-30 | 4 |
| #490 | fix(wgov): resolve delegator vote silent discarding and | sureshchouksey8 | human contributor | 2026-05-30 | 1 |
| #489 | fix(wstaking): resolve RegionShare accounting overwrite | sureshchouksey8 | human contributor | 2026-05-30 | 2 |
| #487 | fix(wstaking): resolve staker unbonding timelock bypass | sureshchouksey8 | human contributor | 2026-05-30 | 1 |
| #479 | fix: reject duplicate outgoing tx ids in genesis Valida | yoshiha-ji | human contributor | 2026-05-30 | 2 |
| #478 | fix: allow higher KYC levels to join MEGroup | Aglcr7 | human contributor | 2026-05-30 | 2 |
| #473 | Fix: Prevent BridgeToken denom/contract duplication in  | yoshiha-ji | human contributor | 2026-05-30 | 1 |
| #464 | fix: reject whitespace rollapp ids | Aglcr7 | human contributor | 2026-05-30 | 3 |
| #452 | fix(rollapp): normalize RollappId by trimming whitespac | zeroknowledge0x | human contributor | 2026-05-30 | 1 |
| #446 | fix: validate wasm multisend messages | jamilahmadzai | human contributor | 2026-05-29 | 2 |
| #442 | fix(gravity): validate duplicate BridgeTokens in genesi | aglichandrap | human contributor | 2026-05-29 | 3 |
| #440 | fix(ante): enforce minimumFee and base denom validation | aglichandrap | human contributor | 2026-05-29 | 1 |
| #431 | fix(ante): enforce umec base-denom fee during DeliverTx | zeroknowledge0x | human contributor | 2026-05-29 | 1 |
| #429 | Fix UnBondRegion operator clear issue | lqkhanh295 | human contributor | 2026-05-29 | 2 |
