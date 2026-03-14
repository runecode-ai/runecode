# Workflow Extensibility v0 — Post-MVP

User-visible outcome: users can opt into schema-validated custom workflows and rebuildable shared-memory accelerators without weakening RuneCode's manifest- and policy-driven safety model.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-13-1600-workflow-extensibility-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

Parallelization: docs-only; safe to do anytime.

## Task 2: ProcessDefinition Object Model

- Define a schema-validated `ProcessDefinition` as the user extension point for post-MVP workflow composition.
  - `ProcessDefinition` is a top-level object family with explicit `schema_id` and `schema_version`.
  - Nested step/block objects stay typed and discriminator-driven; the schema must reject unknown block kinds and unknown fields.
- `ProcessDefinition` may choose:
  - which allowlisted RuneCode step types exist
  - ordering and branching
  - sequential vs parallel blocks
  - per-step provider/model selection where the signed manifest already allows it
  - retry/backoff posture where the step type explicitly supports it
- `ProcessDefinition` must not:
  - introduce a new capability not already available through signed manifests + policy
  - widen data-class, egress, secret, or filesystem access
  - define ad-hoc executable code or dynamic step plugins
- Adding a new allowlisted step type is a capability expansion and requires a schema version bump plus security review.
- The selected process definition hash must be bound into the run's signed inputs so approvals, audit events, and policy decisions can distinguish different workflow shapes.

Parallelization: can be designed in parallel with runner and policy follow-on work; implementation should wait until the shared object model and manifest binding rules are fixed.

## Task 3: ProcessDefinition Validation + Governance

- Support JSON and YAML authoring, but normalize both to the same canonical logical object before hashing/validation.
- Fail closed on:
  - unknown step types or block kinds
  - cycles or unsupported recursion
  - unsupported concurrency/fanout shapes
  - omitted required capability references
  - implicit defaults that would widen behavior in a security-sensitive way
- Define deterministic structural limits for definition depth, total blocks, branch fanout, and per-step parameter sizes.
- Require explicit references to signed manifest capabilities, approval profile posture, and relevant allowlists rather than inferring them from step names.

Parallelization: can be implemented in parallel with fixture/tooling work once the schema profile is agreed.

## Task 4: Shared Memory Accelerator Model

- Define optional "shared memory" as a rebuildable accelerator keyed by `(repo, commitSHA)` and any additional versioned inputs needed to avoid stale reuse.
- Shared memory may cache only derived summaries, maps, selections, embeddings, or similar helper artifacts.
- Shared memory must not store:
  - raw workspace/code content beyond already-approved artifact references
  - secrets or unredacted sensitive values
  - authoritative run state, approvals, or policy decisions
- Invalidate or rebuild shared memory when any relevant bound input changes, including commit, schema version, process definition hash, manifest hash, policy inputs, or model/adapter version where applicable.
- Artifact store objects and the runner's durable run DB remain authoritative; shared memory is advisory only and safe to discard.

Parallelization: can be designed in parallel with artifact-store and runner persistence work; it depends on stable artifact references and manifest-hash binding rules.

## Task 5: Runner, Policy, and TUI Integration

- The runner may execute a user-selected `ProcessDefinition` only after broker-side schema validation and policy authorization.
- Policy outputs and approval payloads must be able to reference the active process definition hash and relevant block/step identifiers.
- TUI/runtime surfaces must show:
  - the active process definition identity
  - whether shared memory is enabled
  - when a cached accelerator was invalidated or bypassed
- Audit events should record which process definition and shared-memory posture were active for each run/session.

Parallelization: can be implemented in parallel across runner, policy, and TUI once the schema contract is fixed.

## Task 6: Fixtures + Validation

- Add checked-in fixtures for:
  - valid `ProcessDefinition` examples (linear, branching, parallel)
  - invalid definitions (unknown step kinds, capability escalation attempts, invalid recursion/fanout)
  - JSON/YAML equivalence for the same logical workflow
  - shared-memory metadata/invalidation examples
- Fixture updates must stay explicit and reviewable; CI verifies but does not regenerate them implicitly.

Parallelization: fixtures can be created in parallel with implementation so long as they validate against the same schemas.

## Acceptance Criteria

- Users can opt into schema-validated custom workflows without introducing new capabilities outside signed manifests + policy.
- Unknown or malformed process definitions fail closed.
- The active process definition is hash-bound and auditable.
- Shared memory is optional, rebuildable, and never authoritative for security-sensitive state.
- Shared-memory reuse is invalidated deterministically when bound inputs change.
