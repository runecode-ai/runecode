## Summary
RuneCode routes all third-party model traffic through a hardened gateway with typed request and response contracts, data-class controls, and quota auditing.

## Problem
The previous combined change coupled secrets foundations and model egress details, reducing review focus for network-bound trust boundary controls.

## Proposed Change
- Dedicated model-gateway role with allowlisted egress.
- Typed `LLMRequest` and `LLMResponse` boundary.
- Data-class policy enforcement for model egress.
- Audit and quota enforcement for outbound model traffic.

## Why Now
This feature keeps egress-bound controls independently reviewable while still aligning under the secure model/provider access project.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Secret storage internals and key posture recording.
- Provider-specific runtime bridge contracts.

## Impact
Keeps model egress controls and trust-boundary hardening reviewable as a standalone feature.
