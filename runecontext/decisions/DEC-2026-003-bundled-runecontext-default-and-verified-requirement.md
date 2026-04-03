---
schema_version: 1
id: DEC-2026-003-bundled-runecontext-default-and-verified-requirement
title: Bundled RuneContext by Default and Verified Requirement
originating_changes:
  - CHG-2026-001-57d6-agent-os-to-runecontext-migration-umbrella
related_changes: []
---

# DEC-2026-003: Bundled RuneContext by Default and Verified Requirement

## Status
Accepted

## Date
2026-04-03

## Context
Future RuneCode planning in this repository needs consistent assumptions for project-context integration.
The migration plan defines bundled RuneContext as the normal integration model and verified mode as required for RuneCode-managed repos.

## Decision
- Future RuneCode planning in this repository assumes RuneCode uses a bundled-by-default RuneContext companion.
- External RuneContext support may exist later as an advanced option, but it is not the default planning assumption.
- Future RuneCode planning in this repository assumes RuneCode-managed repos require RuneContext `verified` mode for normal operation.

## Consequences
- Migrated artifacts should not require end users to drive raw `runectx` commands during normal RuneCode workflows.
- Planning for context-aware features assumes RuneContext integration is complete for that feature surface when delivered.
- Incompatible or non-verified project states are treated as blocked normal-operation states in RuneCode planning.
