# Workflow Concurrency v0 — Shaping Notes

## Scope

Define later workflow-run concurrency support without weakening RuneCode's deterministic run, approval, artifact, or audit model.

## Decisions

- One active run per workspace remains the default fail-closed posture.
- Shared-workspace concurrency requires an explicit model, not opportunistic scheduling.
- Concurrency state must be visible to the runner, broker, policy layer, and TUI.
- Approval and artifact bindings remain run-specific even when runs execute concurrently.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-workflow-workspace-roles-gates-v0/`, `agent-os/specs/2026-03-08-1039-policy-engine-v0/`, `agent-os/specs/2026-03-08-1039-artifact-store-data-classes-v0/`
- Product alignment: preserves predictable evidence and least-privilege defaults while allowing later throughput improvements.

## Standards Applied

- `security/trust-boundary-layered-enforcement` - concurrency must not weaken authorization or artifact-flow enforcement.
- `global/deterministic-check-write-tools` - lock-state fixtures and tests must remain explicit and deterministic.
