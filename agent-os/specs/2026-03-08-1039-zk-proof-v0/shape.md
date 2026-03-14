# ZK Proof v0 (One Narrow Proof + Verify) — Shaping Notes

## Scope

Add a single ZK proof capability for narrow integrity attestations over deterministic data.

## Decisions

- ZK is used for integrity attestations of deterministic computations/records, not for proving arbitrary reasoning.
- Proof generation is an explicit workflow step; verification is always deterministic.
- The first ZK proof ships only if a concrete proving system can be selected with acceptable performance; otherwise release is deferred rather than weakening the contract.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-audit-log-verify-v0/`
- Product alignment: Enables sharing verifiable evidence with minimal disclosure.

## Standards Applied

- None yet.
