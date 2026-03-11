# Workflow Runner + Workspace Roles + Deterministic Gates v0 — Shaping Notes

## Scope

Build the end-to-end workflow engine and offline workspace execution roles, with deterministic gates and evidence artifacts.

## Decisions

- The scheduler is treated as untrusted; the launcher/policy is the enforcement point.
- The workflow runner is distributed as a Node SEA (single executable) built from a bundled CommonJS script.
  - SEA is packaging (not a sandbox) and does not change the runner's trust level.
  - SEA config ignores `NODE_OPTIONS` (set `execArgvExtension: "none"`) to prevent environment-driven runtime option injection.
- Workspace roles are offline; model egress (if enabled) is only via model-gateway.
- Pause/resume is implemented via a persisted run state machine (durable state), not in-memory orchestration.
- Gate failure semantics are explicit (fail/abort, retry, and any override requires recorded approval).
- MVP uses a "moderate" approval profile: approvals are checkpoint-style (stage sign-off and explicit posture changes), not per-action.
- The workflow produces verifiable evidence artifacts (including `audit_verification_report`), not just human-readable logs.

## Context

- Visuals: None.
- References: `agent-os/product/tech-stack.md`
- Product alignment: Spec-first, least-privilege automation with auditable evidence.

## Standards Applied

- None yet.
