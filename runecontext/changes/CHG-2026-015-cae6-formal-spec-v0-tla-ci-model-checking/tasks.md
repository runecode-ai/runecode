# Tasks

## Define Invariants to Specify (MVP Scope)

- [ ] Capability manifest semantics (what it means for an action/data flow to be permitted).
- [ ] Scheduler invariants (no escalation-in-place).
- [ ] Role isolation constraints (no forbidden capability combinations).
- [ ] Artifact routing invariants (only allowed data-class flows; only hash-addressed artifacts consumed).
- [ ] Audit invariants (hash chaining rules; required events before state advances).

Parallelization: can be done in parallel with implementation specs; keep it aligned by referencing the canonical schemas and invariants.

## Write TLA+ Specification

- [ ] Encode a bounded model of:
  - roles
  - actions
  - manifests
  - artifacts
  - audit events
- [ ] Define safety properties and invariants.

Parallelization: can be done in parallel with early implementation; update as interfaces stabilize.

## CI Model Checking

- [ ] Run the model checker in CI.
- [ ] Keep bounds small but meaningful (enough to cover multi-step and failure/timeout paths).

Parallelization: can be implemented in parallel with CI scaffolding; it touches CI workflows and should avoid conflicting edits.

## Traceability

- [ ] Add a simple mapping between spec concepts and runtime modules so failures are actionable.

Parallelization: can be done in parallel with code organization work; it depends on stable module naming.

## Acceptance Criteria

- [ ] CI fails on invariant violations.
- [ ] The spec covers the highest-risk invariants that enforce separation and audit integrity.
