---
schema_version: 1
id: DEC-2026-002-verified-assurance-and-assurance-path
title: Verified Assurance Adoption and Assurance Path
originating_changes:
  - CHG-2026-001-57d6-agent-os-to-runecontext-migration-umbrella
related_changes: []
---

# DEC-2026-002: Verified Assurance Adoption and Assurance Path

## Status
Accepted

## Date
2026-04-03

## Context
This repository already has a RuneContext root config and baseline assurance artifact.
The migration plan requires immediate verified-mode operation and a stable assurance location.

## Decision
- This repository adopts RuneContext in `verified` mode from the start of migration.
- Assurance artifacts for this repository remain under `runecontext/assurance/`.
- There is no temporary plain-mode migration phase for canonical migration work.

## Consequences
- Migration steps must keep `runectx validate` and `runectx status` clean while in verified mode.
- Assurance history for migration work is captured natively under the active assurance path.
- Legacy `agent-os/` material may be represented as imported history when needed; it must not be rewritten as if it were native verified provenance.
