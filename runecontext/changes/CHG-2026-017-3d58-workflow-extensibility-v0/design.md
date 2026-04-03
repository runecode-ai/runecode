# Design

## Overview
Define schema-validated workflow composition and rebuildable shared-memory accelerators without adding new privileged operations.

## Key Decisions
- `ProcessDefinition` is a typed, schema-validated composition surface, not a plugin system.
- Custom workflows compose an allowlist of existing RuneCode step types; they do not add new privileged operations.
- Selected process definitions are signed, hash-bound inputs to policy, approval, and audit flows.
- `ProcessDefinition` uses JSON as its runtime and canonical on-disk format.
- JSON Schema is the single validation source of truth for `ProcessDefinition` objects.
- Future authoring adapters must normalize to the same canonical JSON object before validation and hashing; direct runtime execution consumes JSON only.
- Shared memory is a rebuildable accelerator for derived artifacts only; authoritative state remains in the run DB, artifact store, and audit trail.

## Main Workstreams
- `ProcessDefinition` Contract
- Validation + Canonicalization
- Shared-Memory Accelerators
- Policy, Approval, and Audit Binding
- Authoring + UX Surfaces

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
