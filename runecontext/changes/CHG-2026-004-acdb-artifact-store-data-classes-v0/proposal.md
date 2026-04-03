## Summary
all cross-boundary handoffs are explicit, hash-addressed artifacts with immutable data classification; the system enforces allowed flows between roles.

## Problem
The legacy Artifact Store + Data Classes v0 plan still lives under `agent-os/specs/2026-03-08-1039-artifact-store-data-classes-v0/`, so the roadmap and related planning references do not yet point at a canonical RuneContext change record.

## Proposed Change
- Define MVP Data Classes.
- Content-Addressed Artifact Store (CAS).
- Flow Matrix Enforcement.
- Quotas + Limits (Minimal).
- Garbage Collection + Retention (Minimal).

## Why Now
This work remains scheduled for v0.1.0-alpha.2, and Phase 5 needs a canonical RuneContext change so later delivery and verification no longer depend on legacy Agent OS folders.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runtime implementation of the feature during this migration step.
- Preserving the legacy `agent-os/specs/*` folder as an active planning source of truth.

## Impact
Keeps Artifact Store + Data Classes v0 reviewable as a RuneContext-native change and removes the need for a second semantics rewrite later.
