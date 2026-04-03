---
schema_version: 1
id: DEC-2026-001-runecontext-canonical-planning-system
title: RuneContext as Canonical Planning and Standards System
originating_changes:
  - CHG-2026-001-57d6-agent-os-to-runecontext-migration-umbrella
related_changes: []
---

# DEC-2026-001: RuneContext as Canonical Planning and Standards System

## Status
Accepted

## Date
2026-04-03

## Context
RuneCode currently stores repo-local planning, standards, and product-governance material under `agent-os/`.
The migration target is a RuneContext-native structure under `runecontext/` with `runecontext.yaml` at the repo root.

## Decision
- RuneContext replaces `agent-os/` as the canonical repo-local planning and standards system for this repository.
- Canonical project planning content for this repository lives under `runecontext/`.
- The migration is foundation-first and direct-to-final-state, not a long-lived dual-track model.

## Consequences
- Future canonical references for repo-local planning and standards must target `runecontext/` artifacts.
- `agent-os/` content is treated as legacy material during migration and is deleted only after validated replacement and assurance capture.
- Migrated feature content is rewritten directly to RuneContext-era meaning during migration rather than via a second semantic rewrite.
