## Summary
RuneCode deterministically allows or denies actions based on signed manifests and role rules, with explicit human approvals for elevated risk.

## Problem
The legacy Policy Engine v0 plan still lives under `agent-os/specs/2026-03-08-1039-policy-engine-v0/`, so the roadmap and related planning references do not yet point at a canonical RuneContext change record.

## Proposed Change
- Role + Run/Stage Policy Model.
- Invariants (Fail Closed).
- Approval Policy (MVP: Moderate).
- Backend Selection Rules.
- Decision Outputs.

## Why Now
This work remains scheduled for v0.1.0-alpha.3, and Phase 5 needs a canonical RuneContext change so later delivery and verification no longer depend on legacy Agent OS folders.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runtime implementation of the feature during this migration step.
- Preserving the legacy `agent-os/specs/*` folder as an active planning source of truth.

## Impact
Keeps Policy Engine v0 reviewable as a RuneContext-native change and removes the need for a second semantics rewrite later.
