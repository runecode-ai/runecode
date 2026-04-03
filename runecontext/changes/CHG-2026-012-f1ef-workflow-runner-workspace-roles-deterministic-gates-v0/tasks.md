# Tasks

## Workflow Runner Contract (Untrusted Scheduler)

- [ ] Implement a TS/Node workflow runner using LangGraph.
- [ ] Treat LangGraph as an internal orchestration detail.
  - Stable interfaces are RuneCode schemas + the broker/local API (not LangGraph internals).
  - Document runner <-> broker contract as the only supported integration surface.
  - Keep the door open to replace LangGraph without changing security boundaries.
- [ ] Do not use LangChain agents / Deep Agents as the core runtime.
  - Keep orchestration typed and step-based so capability boundaries remain deterministic and auditable.
- [ ] Distribution: ship the runner as a Node SEA (single executable) built from a bundled CommonJS script.
  - SEA config must ignore `NODE_OPTIONS` (set `execArgvExtension: "none"`).
  - The runner must not rely on runtime `node_modules` resolution (bundle dependencies into the injected script).
  - Local development may still run the runner under `node`/`npm`; SEA is the release/runtime artifact.
- [ ] Ensure the runner has no direct secrets and no direct workspace access.
- [ ] Ensure the runner has no public network egress (runner is not a gateway role).
- [ ] All actions are requested through the broker/local API and independently validated by the launcher/policy engine.
- [ ] MVP approval posture: the runner is built to pause/resume on typed approvals, with the run/stage manifest carrying an approval profile (`moderate` for MVP).
- [ ] Pause/resume semantics (MVP):
  - Runner pauses only when policy returns `require_human_approval` (or a typed input request).
  - Resume occurs only when the broker records a typed decision/response artifact.
  - Approval requests have explicit TTL/expiry; expired/stale approvals must be re-requested (see `runecontext/changes/CHG-2026-007-2315-policy-engine-v0/`).
- [ ] Persist run state durably so pause/resume and crash recovery are real (MVP: SQLite):
  - run state machine (proposed/validated/authorized/executing/awaiting_approval/failed/succeeded)
  - step attempts, artifact references, and approval records
  - idempotency/replay rules for retrying after crashes:
    - every step execution has a stable `step_attempt_id`
    - state transitions are two-phase (`authorized -> executing -> attested`) so partial execution is detectable
    - retry rules are explicit per step type (some steps may be non-retriable and must fail closed)
  - define upgrade/migration rules for the durable DB:
    - explicit schema versioning and migrations
    - fail closed on unknown versions
    - record DB schema version and SQLite library version in audit metadata
- [ ] Decide the durable state authority to avoid dual sources of truth:
  - either keep LangGraph checkpointing minimal and treat the broker/run DB as authoritative, or
  - define strict rules for what LangGraph may checkpoint (control-plane IDs/hashes only).
- [ ] Define MVP concurrency rules:
  - default: one active run per workspace (explicit workspace lock)
  - concurrent runs require explicit design and are specified in `runecontext/changes/CHG-2026-027-71ed-workflow-concurrency-v0/`
  - multiple runs across distinct workspaces are permitted by default (locks are per-workspace, not global)

Node SEA feasibility note (MVP):
- [ ] Validate SEA packaging early with a spike (LangGraph + required deps).
- [ ] If SEA bundling is blocked by ecosystem constraints, ship an alpha-friendly fallback (pinned Node runtime + bundled JS) without changing the trust model; keep SEA as the release target.

Parallelization: runner implementation can proceed in parallel with the broker/policy engine as long as the runner<->broker schema contract is finalized early.

## Workflow Extensibility Follow-On Spec

- [ ] Post-MVP user-configurable workflows and shared-memory accelerators now live in `runecontext/changes/CHG-2026-017-3d58-workflow-extensibility-v0/`.
- [ ] This MVP runner spec keeps the baseline runner contract and durable-state rules those later workflow extensions build on.

## Runner Persistence Rules (MVP)

- [ ] Restrict LangGraph persistence to runner control-plane state only (MVP):
  - thread/run IDs, step IDs, artifact hashes, approval handles, and other non-sensitive bookkeeping.
  - explicitly forbid storing raw workspace/code, unredacted excerpts, or secrets in runner persistence.

Parallelization: can be implemented in parallel with artifact store work; it depends on stable artifact reference schemas.

## Workspace Roles (MVP Set)

- [ ] Define and implement the MVP workspace roles:
  - `workspace-read` (RO)
  - `workspace-edit` (RW, offline)
  - `workspace-test` (snapshot + discard)
- [ ] Ensure command execution is via purpose-built executors/allowlists (no shell passthrough).

Parallelization: roles can be implemented in parallel with the runner and broker once the executor interfaces and artifact movement model are stable.

## Propose -> Validate -> Authorize -> Execute -> Attest Loop

- [ ] Treat model output as untrusted proposals.
- [ ] Support streaming as a first-class UX/performance feature.
  - streaming events update UI/telemetry only; only finalized, schema-valid outputs may drive actions
- [ ] Tool calling is supported only as structured proposal objects.
  - no tool execution directly from model output
  - execution requires deterministic validation + policy authorization + the correct offline executor
- [ ] Require structured JSON outputs for machine-consumed results.
  - plain text is reserved for human-facing explanations/summaries
- [ ] Validate proposals structurally (schema, size, artifact references).
- [ ] Authorize deterministically via policy engine.
- [ ] Execute inside the correct role isolate.
- [ ] Attest by producing signed artifacts (diffs, logs, gate results) and audit events.
- [ ] Prefer checkpoint-style approvals (stage sign-off and explicit posture changes) over per-action approvals in the MVP workflow design.

Parallelization: can be implemented in parallel with policy engine and audit/artifact subsystems; it depends on shared schemas and reason codes.

## Deterministic Gates (MVP)

- [ ] Implement a gate framework with evidence artifacts for:
  - build/type checks
  - tests
  - lint/format
  - secret scanning
  - policy compliance checks
- [ ] Define gate failure semantics:
  - default: gate failure fails the step/run deterministically
  - retries are explicit and recorded
  - any override requires a recorded human approval and produces an audit event

Parallelization: gate framework can be implemented in parallel with workspace-test role and artifact store; avoid conflicts by agreeing on evidence artifact schemas first.

## Minimal End-to-End Demo Run

- [ ] Provide a single demo workflow that runs on Linux:
  - creates a small change in a demo workspace
  - runs gates
  - produces audit + artifacts
  - requires at least one explicit approval (e.g., manifest sign-off)
  - produces a verifier artifact (run audit verification; data class: `audit_verification_report`) that can be shown in the TUI

Parallelization: can be developed after the minimal slices of runner/broker/policy/audit/artifacts exist; treat it as an integration milestone rather than a dependency.

## Acceptance Criteria

- [ ] A run can be started, paused for approval, resumed, and completed.
- [ ] A run can be recovered after a crash/restart of the scheduler process (no in-memory state required to resume).
- [ ] Gates are deterministic and produce verifiable artifacts.
- [ ] The scheduler cannot exceed policy or bypass gates.
