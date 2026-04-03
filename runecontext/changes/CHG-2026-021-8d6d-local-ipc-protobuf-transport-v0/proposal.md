## Summary
RuneCode can migrate local broker IPC to protobuf without changing the logical protocol or local-only trust posture.

## Problem
The legacy Local IPC Protobuf Transport v0 plan still lives under `agent-os/specs/2026-03-13-1602-local-ipc-protobuf-transport-v0/`, so the roadmap and related planning references do not yet point at a canonical RuneContext change record.

## Proposed Change
- Proto Mapping for the Existing Logical Model.
- Local IPC Transport Requirements.
- Optional Local-Only gRPC Profile.
- Migration and Compatibility Rules.

## Why Now
This work remains scheduled for v0.2, and Phase 5 needs a canonical RuneContext change so later delivery and verification no longer depend on legacy Agent OS folders.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runtime implementation of the feature during this migration step.
- Preserving the legacy `agent-os/specs/*` folder as an active planning source of truth.

## Impact
Keeps Local IPC Protobuf Transport v0 reviewable as a RuneContext-native change and removes the need for a second semantics rewrite later.
