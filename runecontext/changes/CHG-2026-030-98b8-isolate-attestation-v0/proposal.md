## Summary
RuneCode can upgrade MVP TOFU isolate binding to measured, attestable provisioning without changing the core audit contract.

## Problem
The legacy Isolate Attestation v0 plan still lives under `agent-os/specs/2026-03-13-1731-isolate-attestation-v0/`, so the roadmap and related planning references do not yet point at a canonical RuneContext change record.

## Proposed Change
- Attestation Evidence Model.
- Launch, Verification, and Policy Integration.
- TUI + Audit Posture.
- Fixtures + Cross-Platform Considerations.

## Why Now
This work remains scheduled for vNext, and Phase 5 needs a canonical RuneContext change so later delivery and verification no longer depend on legacy Agent OS folders.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runtime implementation of the feature during this migration step.
- Preserving the legacy `agent-os/specs/*` folder as an active planning source of truth.

## Impact
Keeps Isolate Attestation v0 reviewable as a RuneContext-native change and removes the need for a second semantics rewrite later.
