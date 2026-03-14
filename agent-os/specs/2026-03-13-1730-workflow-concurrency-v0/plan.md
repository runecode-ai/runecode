# Workflow Concurrency v0 — Post-MVP

User-visible outcome: RuneCode can safely support explicitly designed concurrent runs without letting multiple workflows silently race over the same workspace, approvals, or artifacts.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-13-1730-workflow-concurrency-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

Parallelization: docs-only; safe to do anytime.

## Task 2: Workspace Concurrency Model

- Define the supported concurrency modes explicitly:
  - multiple runs across distinct workspaces
  - multiple runs against one workspace only when an explicit concurrency mode is enabled
- Keep fail-closed defaults:
  - one active run per workspace remains the default posture
  - unsupported or ambiguous workspace-sharing modes are rejected rather than best-effort scheduled
- Define the lock/lease model for workspace ownership, including acquisition, renewal, expiry, and crash recovery semantics.

Parallelization: can be designed in parallel with runner persistence work once run-state and manifest identity rules are stable.

## Task 3: Conflict Detection + Isolation Rules

- Define deterministic conflict handling for concurrent runs that touch the same logical workspace state.
- Concurrency design must account for:
  - workspace writes and gate outputs
  - artifact publication and promotion
  - approval requests/decisions bound to specific runs and inputs
  - branch/checkpoint selection where workflows share a repo
- Require explicit policy/audit recording when a run uses any non-default concurrency posture.

Parallelization: can be designed in parallel with artifact-store and policy follow-on work once run/resource identity rules are stable.

## Task 4: Runner, Broker, and TUI Integration

- Define how the runner, broker, and local API expose concurrency posture, lock ownership, waits, and conflicts.
- TUI/CLI surfaces must make it obvious when a run is blocked by another active run or is sharing a workspace under an explicit concurrency mode.
- Audit events must record lock acquisition/release, contention, overrides, and conflict-triggered failures.

Parallelization: can be implemented in parallel across runner, broker, and TUI once the shared lock/event model is fixed.

## Task 5: Fixtures + Recovery Cases

- Add checked-in fixtures and test cases for:
  - default single-run-per-workspace behavior
  - lock contention
  - stale lock recovery after crash
  - explicit concurrent-run opt-in
  - conflict-triggered failure and retry paths

Parallelization: fixtures can be developed in parallel with implementation so long as they validate against the same shared lock/event model.

## Acceptance Criteria

- Default behavior remains one active run per workspace.
- Concurrent use of one workspace requires an explicit design and fail-closed posture.
- Locking, contention, and recovery are auditable and deterministic.
- Approval, artifact, and gate semantics stay bound to the correct run under concurrency.
