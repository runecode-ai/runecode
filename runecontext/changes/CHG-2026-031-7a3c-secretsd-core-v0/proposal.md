## Summary
RuneCode has a dedicated secrets daemon that stores long-lived secrets safely and issues short-lived, scope-bound leases with complete auditing.

## Problem
The previous combined change mixed shared secret lifecycle foundations with model egress behavior, making sequencing and verification less clear.

## Proposed Change
- Secretsd storage and key posture requirements.
- Lease issuance, renewal, and revocation rules.
- Safe secret onboarding/import path.
- Local-only health signals and minimal operational metrics.

## Why Now
This feature isolates reusable secret-management foundations so downstream gateway and provider features can depend on one reviewed contract.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Model provider egress and payload-shaping behavior.
- Re-introducing legacy Agent OS planning paths as canonical references.

## Impact
Keeps secrets lifecycle behavior independently reviewable while preserving project-level traceability to the broader secure model/provider access plan.
