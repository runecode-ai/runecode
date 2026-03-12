# GitHub Copilot Subscription Provider (Official Runtime Bridge) — Post-MVP

User-visible outcome: RuneCode can access Copilot-backed models using a user's GitHub Copilot subscription via an officially supported local runtime, while preserving strict isolation (no workspace access in egress roles, `secretsd` as the only long-lived secret store, and complete auditability).

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-11-1921-github-copilot-subscription-provider-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

## Task 2: Official Runtime + Protocol Selection

- Policy constraint: RuneCode does not ship/bundle/redistribute vendor CLIs or proprietary runtimes.
  - Integrate only with an officially supported, user-installed Copilot runtime.
- Use an officially supported Copilot local runtime (installed/managed by the user/admin).
- Select the bridge protocol surface that supports strict permission control and least privilege:
  - ACP over stdio, or
  - official SDK/JSON-RPC mode (if it provides equivalent controls)
- Prefer stdio spawning over listening ports.
- Pin and audit the external runtime version and protocol surface.
  - log the detected version per request for forensic traceability
  - fail closed on unsupported versions

## Task 3: Auth Model (No Env Vars, No Second Store)

- Introduce a dedicated `auth-gateway` flow for GitHub auth when required.
- If OAuth/device-code is required for this provider, RuneCode maintains its own official OAuth client registration.
- Store long-lived auth material only in `secretsd`.
- Disallow environment-variable token injection.
  - Define a token delivery mechanism that does not use env vars (e.g., stdin/FD, or a runtime-supported config file in a secretsd-managed directory).
  - If the runtime requires persisted auth state, it must be stored only in a secretsd-managed encrypted directory and treated as secret material.

## Task 4: Model-Gateway Bridge (LLM-Only)

- Run the runtime under `model-gateway` with:
  - no workspace mounts; empty/scratch `cwd`
  - isolated `HOME`/tool dirs pointing at an allowlisted provider sandbox directory
  - strict deny-by-default tool/permission requests (LLM-only mode)
  - treat any attempt to exec/write/read workspace as a policy violation
  - schema-validated structured outputs only for machine-consumed actions
- Enforce model egress data-class policy at the RuneCode `LLMRequest` boundary.
- Default to ephemeral sessions.
  - do not persist conversation state unless explicitly enabled by signed manifest + policy
  - if the runtime requires local state, it must be stored only in a secretsd-managed encrypted directory
- Prefer protocol-level contract tests over HTTP wire fixtures.
  - add RPC request/response fixture tests and stable error taxonomy mapping

## Task 5: Policy + Audit Integration

- Default deny: enabling this provider is an explicit signed-manifest opt-in and must be surfaced as a high-risk approval.
- Audit requirements:
  - auth events: login start/completed/cancelled, token lease issuance/renewal/revocation
  - model events: provider/model identifiers, bytes, timing, and outcome (without logging secret values)

## Acceptance Criteria

- Copilot subscription model access is possible via official mechanisms.
- No environment-variable secret injection is used.
- No second secrets store exists: only `secretsd` persists long-lived auth material.
- Workspace roles remain offline; all model egress remains behind `model-gateway`.
