# Workflow Extensibility v0 — Shaping Notes

## Scope

Define post-MVP workflow composition and shared-memory acceleration without changing RuneCode's trust boundaries or allowing user-defined capability expansion.

## Decisions

- `ProcessDefinition` is a typed, schema-validated composition surface, not a plugin system.
- Custom workflows compose an allowlist of existing RuneCode step types; they do not add new privileged operations.
- Selected process definitions are signed, hash-bound inputs to policy, approval, and audit flows.
- `ProcessDefinition` uses JSON as its runtime and canonical on-disk format.
- JSON Schema is the single validation source of truth for `ProcessDefinition` objects.
- Future authoring adapters must normalize to the same canonical JSON object before validation and hashing; direct runtime execution consumes JSON only.
- Shared memory is a rebuildable accelerator for derived artifacts only; authoritative state remains in the run DB, artifact store, and audit trail.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-workflow-workspace-roles-gates-v0/`, `agent-os/specs/2026-03-08-1039-policy-engine-v0/`, `agent-os/specs/2026-03-08-1039-protocol-schemas-v0/`
- Product alignment: enables workflow customization and reuse while keeping least-privilege and deterministic auditability intact.

## Standards Applied

- `security/trust-boundary-layered-enforcement` - user-authored workflow definitions must never bypass manifest, broker, or policy enforcement.
- `global/deterministic-check-write-tools` - schema and fixture tooling must remain explicit and deterministic.
