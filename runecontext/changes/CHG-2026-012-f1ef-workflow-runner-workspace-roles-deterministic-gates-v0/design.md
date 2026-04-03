# Design

## Overview
Implement the baseline workflow runner contract, offline workspace roles, deterministic gates, and durable pause/resume state for end-to-end runs.

## Key Decisions
- The scheduler is treated as untrusted; the launcher/policy is the enforcement point.
- LangGraph is an internal implementation detail of the untrusted runner.
- Avoid LangChain agents / Deep Agents as the core runtime.
- The workflow runner is distributed as a Node SEA (single executable) built from a bundled CommonJS script.
- The runner has no public network egress and no direct secrets or workspace access.
- Runner persistence stores control-plane state only and must never store raw workspace/code or secrets.
- Pause/resume is implemented via a persisted run state machine rather than in-memory orchestration.
- Approval requests and decisions are hash-bound and time-bounded.
- MVP uses a moderate approval profile with checkpoint-style approvals.
- The workflow produces verifiable evidence artifacts, not just human-readable logs.
- Concurrency is locked per workspace by default; later shared-workspace concurrency is tracked separately.
- SEA feasibility is validated early; a pinned-Node plus bundled-JS fallback is acceptable for early alpha without changing the trust boundary.

## Main Workstreams
- Workflow Runner Contract (Untrusted Scheduler)
- Workflow Extensibility Follow-On Spec
- Runner Persistence Rules (MVP)
- Workspace Roles (MVP Set)
- Propose -> Validate -> Authorize -> Execute -> Attest Loop
- Deterministic Gates (MVP)
- Minimal End-to-End Demo Run

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
