# Tasks

## `ProcessDefinition` Contract

- [ ] Define the `ProcessDefinition` object family for custom workflow composition.
- [ ] Limit it to approved existing step types and typed control-flow constructs; no new privileged operations.
- [ ] Keep selected process definitions as signed, hash-bound inputs to policy, approval, and audit flows.

## Validation + Canonicalization

- [ ] Keep JSON as the canonical on-disk and runtime format.
- [ ] Use JSON Schema as the single validation source of truth.
- [ ] Normalize any future authoring adapters to the same canonical JSON object before validation and hashing.

## Shared-Memory Accelerators

- [ ] Define rebuildable shared-memory accelerators for derived artifacts only.
- [ ] Keep authoritative state in the run DB, artifact store, and audit trail.

## Policy, Approval, and Audit Binding

- [ ] Bind selected process definitions into policy evaluation, approval requests, and audit evidence.
- [ ] Ensure custom workflows cannot bypass manifest, broker, or policy enforcement.

## Authoring + UX Surfaces

- [ ] Define authoring and review surfaces for process definitions.
- [ ] Keep machine validation deterministic and explicit.

## Acceptance Criteria

- [ ] Custom workflows remain schema-validated, hash-bound, and auditable.
- [ ] Workflow customization does not add new privileged operations or weaken existing trust boundaries.
