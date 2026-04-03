---
schema_version: 1
id: DEC-2026-004-runecode-ux-ownership-and-runecontext-generic-boundary
title: RuneCode UX Ownership and RuneContext Generic Boundary
originating_changes:
  - CHG-2026-001-57d6-agent-os-to-runecontext-migration-umbrella
related_changes: []
---

# DEC-2026-004: RuneCode UX Ownership and RuneContext Generic Boundary

## Status
Accepted

## Date
2026-04-03

## Context
The migration requires clear ownership boundaries between RuneContext and RuneCode in future planning artifacts.
Without a durable decision, migrated feature docs may drift toward RuneContext-specific product semantics or direct end-user CLI dependency.

## Decision
- Future planning assumes RuneCode owns the normal user-facing command set and UX, while invoking RuneContext capabilities under the hood.
- RuneContext remains a generic, machine-friendly project-content layer rather than a RuneCode-only semantic surface.
- RuneContext may provide generic advisory consumer-compatibility warnings during upgrades, while hard compatibility enforcement for RuneCode-managed repos remains in RuneCode.

## Consequences
- Migrated artifacts should prefer generic RuneContext metadata/capability extensions over RuneCode-only RuneContext command semantics.
- Product integration responsibilities such as orchestration, compatibility gating, and user-facing flows stay in RuneCode planning.
- Out-of-band unsupported project upgrades are modeled as fail-closed normal-operation states in RuneCode, with only safe diagnostic/remediation flows available.
