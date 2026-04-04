## Summary
RuneCode enforces deterministic quality and safety gates that produce verifiable evidence artifacts and fail closed by default.

## Problem
Gate behavior and evidence semantics were previously embedded in a broad combined change, limiting focused verification.

## Proposed Change
- Implement a deterministic gate framework.
- Produce hash-addressed gate evidence artifacts.
- Define failure, retry, and override semantics with audit coverage.

## Why Now
This split keeps gate correctness and evidence production independently reviewable while preserving end-to-end workflow traceability.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runner durable-state internals.
- Workspace role implementation details.

## Impact
Keeps gate determinism and evidence semantics as a dedicated feature boundary.
