# Design

## Overview
Use this change as the project-level tracker for secrets and model-gateway work while feature-level implementation lands in child changes.

## Key Decisions
- Child features own runtime implementation detail.
- Parent project owns sequencing, boundaries, and integration posture.
- Security invariants remain deny-by-default, lease-only secret use, and typed/auditable egress.

## Main Workstreams
- `CHG-2026-031-7a3c-secretsd-core-v0`
- `CHG-2026-032-4d1f-model-gateway-v0`
- Cross-lane integration with auth/bridge/provider features

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
