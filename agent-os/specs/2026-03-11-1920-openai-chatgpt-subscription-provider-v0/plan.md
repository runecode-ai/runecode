# OpenAI ChatGPT Subscription Provider (OAuth + Codex Bridge) — Post-MVP

User-visible outcome: RuneCode can access OpenAI GPT models using a user's ChatGPT subscription plan/rate limits via official OAuth, while preserving strict isolation (no workspace access in egress roles, `secretsd` as the only long-lived secret store, and complete auditability).

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-11-1920-openai-chatgpt-subscription-provider-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

## Task 2: Official OAuth Client Registration

- RuneCode maintains its own official OAuth client registration for "Sign in with ChatGPT".
- Use PKCE + `offline_access` to obtain refresh capability.
- Provide two official login paths:
  - browser-based flow with a localhost callback bound to loopback only
  - device-code fallback for headless/remote environments
- Fail closed on OAuth state/redirect mismatches.
- RuneCode must not rely on vendor-internal OAuth clients or "piggyback" registrations.

## Task 3: Auth-Gateway Role (Auth Egress Only)

- Introduce a dedicated `auth-gateway` role:
  - no workspace access
  - network egress allowlist limited to OpenAI auth endpoints only
  - emits typed login artifacts/events (no secrets in logs)
- Store refresh token material and rotation metadata only in `secretsd`.
- Issue short-lived, scope-bound leases for `idToken` and `accessToken` (or equivalent) to `model-gateway`.
- Disallow environment-variable secret injection.
  - tokens flow only via lease IPC.

## Task 4: Model-Gateway Bridge via Codex App-Server

- Policy constraint: RuneCode does not ship/bundle/redistribute vendor CLIs or proprietary runtimes.
  - Integrate only with an officially supported, user-installed Codex runtime.
- Run the official Codex app-server runtime under the `model-gateway` role as a local bridge (stdio JSON-RPC; no listening ports by default).
- Pin and audit the external runtime version and protocol surface.
  - log the detected version per request for forensic traceability
  - fail closed on unsupported versions
- Use Codex external token mode (`chatgptAuthTokens`):
  - `model-gateway` supplies `idToken` and `accessToken` at session start
  - Codex keeps tokens in memory only
  - on authorization failure, the bridge requests refreshed tokens and `model-gateway` satisfies the request by obtaining a fresh lease
- Enforce "LLM-only" capability scoping:
  - no workspace mounts; empty/scratch `cwd`
  - isolated `HOME`/tool dirs pointing at an allowlisted provider sandbox directory
  - deny command execution and patch-application approvals
  - deny-by-default tool/permission requests; treat any attempt to exec/write/read workspace as a policy violation
  - accept only assistant text + schema-validated structured outputs
- Default to ephemeral sessions.
  - do not persist conversation state unless explicitly enabled by signed manifest + policy
  - if enabled, persist state as RuneCode artifacts (not in the bridge runtime home directory)
- Prefer protocol-level contract tests over HTTP wire fixtures.
  - generate and pin app-server schema artifacts for the selected runtime version
  - add RPC request/response fixture tests and stable error taxonomy mapping

## Task 5: Policy + Audit Integration

- Default deny: enabling this provider is an explicit signed-manifest opt-in and must be surfaced as a high-risk approval.
- Audit requirements:
  - auth events: login start/completed/cancelled, token lease issuance/renewal/revocation
  - model events: provider/model identifiers, bytes, timing, and outcome (without logging secret values)
- Enforce model egress data-class policy at the RuneCode `LLMRequest` boundary.

## Acceptance Criteria

- GPT model access uses ChatGPT subscription quotas via official OAuth.
- No environment-variable secret injection is used.
- No second secrets store exists: only `secretsd` persists long-lived auth material.
- Workspace roles remain offline; all model egress remains behind `model-gateway`.
