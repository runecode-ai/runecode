## Summary
Key security invariants are formally specified and continuously model-checked, reducing the chance of subtle privilege-escalation or routing bugs.

## Problem
The legacy Formal Spec v0 (TLA+ + CI Model Checking) plan still lives under `agent-os/specs/2026-03-08-1039-formal-spec-tla-v0/`, so the roadmap and related planning references do not yet point at a canonical RuneContext change record.

## Proposed Change
- Define Invariants to Specify (MVP Scope).
- Write TLA+ Specification.
- CI Model Checking.
- Traceability.

## Why Now
This work remains scheduled for v0.1.0-beta.1, and Phase 5 needs a canonical RuneContext change so later delivery and verification no longer depend on legacy Agent OS folders.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runtime implementation of the feature during this migration step.
- Preserving the legacy `agent-os/specs/*` folder as an active planning source of truth.

## Impact
Keeps Formal Spec v0 (TLA+ + CI Model Checking) reviewable as a RuneContext-native change and removes the need for a second semantics rewrite later.
