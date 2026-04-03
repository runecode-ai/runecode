## Summary
RuneCode can access Copilot models via an official local runtime bridge in LLM-only mode.

## Problem
The legacy GitHub Copilot Subscription Provider (Official Runtime Bridge) plan still lives under `agent-os/specs/2026-03-11-1921-github-copilot-subscription-provider-v0/`, so the roadmap and related planning references do not yet point at a canonical RuneContext change record.

## Proposed Change
- Official Runtime Bridge Integration.
- Auth-Gateway + Token Delivery.
- Compatibility + Probe Policy.
- Policy + Audit Integration.

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
Keeps GitHub Copilot Subscription Provider (Official Runtime Bridge) reviewable as a RuneContext-native change and removes the need for a second semantics rewrite later.
