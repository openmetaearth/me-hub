# Meta Earth Phase I — Scope Map & Risk Notes

## Overview

This document defines the security-critical attack surface for Meta Earth Phase I. It maps each major component to its associated risks, trust boundaries, and recommended testing strategies. Use this as the primary reference for bug bounty hunting, security review, and penetration testing.

---

## Scope Map

| Component | Risk Level | Primary Attack Vector | Impact |
|-----------|-----------|----------------------|--------|
| Bridge | Critical | Fund theft via forged messages, replay, or proof manipulation | Loss of all bridged assets |
| Sequencer Logic | Critical | Fake state submissions, invalid rollup state finalization | Chain state corruption, asset theft |
| Governance | High | Unauthorized control, vote inflation, proposal manipulation | Protocol takeover, fund drain |
| Reward Distribution | High | Infinite minting, double-claim, accounting errors | Unauthorized token creation, economic collapse |
| DID/KYC Module | High | Identity bypass, sybil attacks, credential forgery | Compliance failure, unauthorized access |
| Validator Logic | Critical | Consensus manipulation, equivocation, liveness attacks | Chain halt, reorganization, double-spend |
| RollApp Settlement | High | Invalid state transitions, fraud proof bypass | Rollup compromise, fund loss |
| Token Minting | Medium | Unauthorized minting, supply inflation | Token devaluation, economic manipulation |
| Gas Logic | Medium | Gas exhaustion, underpricing, fee manipulation | Network DoS, spam attacks |

---

## Component Deep-Dive: Bridge

### Risk: Fund Theft

The bridge is the highest-value target. All cross-chain message passing, proof verification, and asset locking/unlocking logic must be hardened.

### Test Vectors

| Test | Description |
|------|-------------|
| Replay Attacks | Can a valid message be executed more than once across chains? |
| Forged Messages | Can an attacker craft a message that passes signature verification? |
| Fake Proofs | Can invalid Merkle proofs or light client proofs be accepted? |
| Nonce Reuse | Does the bridge properly track and enforce message nonces? |
| Chain ID Confusion | Can messages intended for one chain be replayed on another? |
| Partial Verification | Are all fields in a message verified, or can some be omitted? |

### Key Questions

- Can a message execute twice?
- Can a fake sequencer submit fraudulent state to the bridge?
- Can a malicious relayer forge finality confirmations?
- What happens if the bridge contract is paused mid-transaction?
- Are there race conditions between deposit and withdrawal flows?

---

## Component Deep-Dive: Sequencer Logic

### Risk: Fake State Submissions

Sequencers are trusted to submit valid rollup state roots. If this trust is misplaced or can be exploited, the entire chain is compromised.

### Test Vectors

| Test | Description |
|------|-------------|
| Unauthorized Submission | Can a non-sequencer submit state roots? |
| State Root Manipulation | Can an attacker submit a fraudulent state root? |
| Batch Ordering | Can transactions be reordered to extract value? |
| Liveness Attacks | Can a sequencer censor transactions? |
| Fee Manipulation | Can sequencer fees be stolen or inflated? |

### Key Questions

- How is sequencer identity verified?
- What prevents a sequencer from submitting conflicting state roots?
- Is there slashing for misbehavior?
- Can sequencers be rotated or removed?

---

## Component Deep-Dive: Governance

### Risk: Unauthorized Control

Governance controls protocol parameters, upgrades, and treasury. Compromise means total protocol control.

### Test Vectors

| Test | Description |
|------|-------------|
| Vote Inflation | Can an attacker create votes without holding tokens? |
| Snapshot Manipulation | Can voting power be calculated incorrectly? |
| Proposal Replay | Can a proposal be executed multiple times? |
| Timing Attacks | Can proposals pass with insufficient voting period? |
| Quorum Bypass | Can proposals pass without meeting minimum participation? |
| Upgrade Hijacking | Can a malicious upgrade be forced through? |

### Key Questions

- Who can create proposals?
- What is the minimum quorum?
- Can proposals be canceled or vetoed?
- Is there a timelock on execution?
- How are upgrade contracts verified?

---

## Component Deep-Dive: Reward Distribution

### Risk: Infinite Rewards

Reward logic errors can lead to unlimited token minting, destroying tokenomics.

### Test Vectors

| Test | Description |
|------|-------------|
| Double Claim | Can a user claim rewards multiple times? |
| Overflow/Underflow | Can integer math errors create or destroy tokens? |
| Rounding Errors | Can rounding accumulate to significant value? |
| Rate Manipulation | Can reward rates be artificially inflated? |
| Eligibility Bypass | Can non-eligible users claim rewards? |

### Key Questions

- How are reward rates calculated?
- Is there a maximum supply cap?
- Can rewards be claimed before they vest?
- What happens to unclaimed rewards?

---

## Component Deep-Dive: DID/KYC Module

### Risk: Identity Bypass

DID/KYC is the gatekeeper for compliance and access control. Bypass means unauthorized participation.

### Test Vectors

| Test | Description |
|------|-------------|
| Credential Forgery | Can fake KYC documents be accepted? |
| Sybil Attacks | Can one entity create multiple identities? |
| Reuse Attacks | Can expired credentials be reused? |
| Verification Bypass | Can the KYC step be skipped entirely? |
| Data Tampering | Can stored identity data be modified? |

### Key Questions

- How are credentials verified?
- What prevents identity reuse across accounts?
- Is there a revocation mechanism?
- How is user privacy protected?

---

## Component Deep-Dive: Validator Logic

### Risk: Consensus Manipulation

Validators secure the network. Compromise means chain takeover.

### Test Vectors

| Test | Description |
|------|-------------|
| Equivocation | Can a validator sign conflicting blocks? |
| Liveness Attacks | Can validators halt block production? |
| Censorship | Can validators exclude transactions? |
| Key Compromise | What happens if a validator key is stolen? |
| Slashing Bypass | Can validators avoid penalties for misbehavior? |

### Key Questions

- How are validators selected?
- What is the minimum stake?
- How is misbehavior detected and punished?
- Can validators be forcibly removed?

---

## General Attack Patterns

### The Golden Rule

> **"The protocol trusts something it should verify."**

This single mindset finds:
- Bridge exploits
- Validator exploits
- Governance exploits
- Settlement exploits
- Mint exploits

### High-Value Patterns

| Pattern | Description |
|---------|-------------|
| Access Control | Missing or incorrect permission checks |
| Reward Accounting | Math errors in distribution logic |
| Replay Attacks | Missing nonce or chain ID enforcement |
| Missing Validation | Unchecked inputs or assumptions |
| Incorrect State Transitions | Wrong state machine logic |

---

## Recommended Testing Strategy

### For Beginners

1. Start with **access control** — who can call what?
2. Check **reward accounting** — can numbers be manipulated?
3. Test **replay attacks** — can messages be reused?
4. Look for **missing validation** — what isn't being checked?
5. Audit **state transitions** — can invalid states be reached?

### For Advanced Researchers

1. Focus on **bridge logic** — highest payout potential
2. Analyze **consensus mechanisms** — most critical to chain security
3. Review **governance systems** — protocol-level control
4. Examine **economic incentives** — game theory attacks
5. Study **cross-chain interactions** — composability risks

---

## Reference: Previous Critical Vulnerabilities

Study these patterns from public reports:

| Source | Common Findings |
|--------|-----------------|
| Immunefi | Bridge exploits, oracle manipulation, access control |
| Code4rena | Accounting errors, validation bypass, reentrancy |
| Sherlock | Reward logic flaws, governance attacks, signature issues |

---

## Risk Scoring Methodology

| Level | Description | Payout Range |
|-------|-------------|--------------|
| Critical | Direct fund loss, chain compromise | $50,000+ |
| High | Significant economic damage, protocol manipulation | $10,000 - $50,000 |
| Medium | Limited impact, requires specific conditions | $1,000 - $10,000 |
| Low | Minor issues, informational | Up to $1,000 |

---

## Immediate Action Items

1. **Audit bridge message verification** — highest risk component
2. **Review sequencer authorization** — critical trust boundary
3. **Test governance proposal execution** — protocol control point
4. **Validate reward distribution math** — economic integrity
5. **Check DID/KYC verification flow** — compliance gate

---

*This scope map is a living document. Update as new components are added or attack surfaces are discovered.*