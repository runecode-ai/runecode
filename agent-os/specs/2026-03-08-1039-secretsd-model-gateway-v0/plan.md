# Secretsd + Model-Gateway v0

User-visible outcome: third-party model access is possible only via an explicitly allowed gateway role, using short-lived scoped secrets leases, with boundary redaction and complete auditing.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-08-1039-secretsd-model-gateway-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

## Task 2: Secretsd MVP Interface

- Provide a minimal secrets daemon that:
  - stores long-lived secrets at rest (prefer hardware-backed / OS key storage where available)
  - fails closed by default if secure key storage is unavailable (no silent plaintext fallback)
  - allows an explicit, audited opt-in to passphrase-derived encryption for developer/portable setups
  - issues short-lived, scope-bound leases only as allowed by the signed manifest
  - defines lease TTL bounds, renewal rules, and revocation semantics
  - records every lease as an audit event (without logging raw secrets)
- Define a safe secret onboarding/import flow (MVP):
  - secrets are provided via stdin or a file descriptor (never CLI args or environment variables)
  - only secret metadata/IDs are logged/audited (never secret values)

## Task 3: Model-Gateway Role

- Implement a dedicated gateway role with:
  - network egress allowlist (model provider domains only)
  - no workspace access
  - provider keys obtained only via secrets leases
  - schema-validated request/response boundary
  - RuneCode-native typed model requests/responses (no freeform prompt blobs cross the boundary; inputs reference artifacts by hash)
  - Support streaming responses within the typed boundary.
  - Support tool calling only as typed proposal objects; never execute tools from the gateway.
  - Require structured JSON outputs for any machine-consumed output.
- Implementation constraint (MVP): keep `model-gateway` implemented in Go to minimize the trusted computing base (TCB).
  - Avoid introducing npm supply-chain dependencies into a high-risk egress boundary.
  - Do not add a runtime Node "request builder" isolate for provider payload shaping.
- Model-gateway must fetch artifact bytes by hash (via broker-mediated artifact store APIs) and assemble provider requests only from allowlisted artifact data classes.
- Harden egress controls against SSRF and DNS rebinding:
  - resolve and validate destinations (block RFC1918/link-local/reserved ranges)
  - restrict redirects (or disable by default)
  - require TLS with certificate validation and SNI matching
  - apply strict timeouts and response size limits
  - Define streaming-specific limits (chunk sizes, total streamed bytes, idle timeouts) so streaming cannot bypass size/timeout controls.

## Task 3b: Provider Adapters + Drift Detection (MVP)

- Translate RuneCode-native `LLMRequest` into provider-specific HTTP payloads inside the Go model-gateway.
- Do not depend on LangChain provider packages for production egress payload shaping.
- Keep official provider SDK packages out of the production egress path.
  - Use them only for test/fixture generation to detect upstream request-shape drift.
- Commit non-sensitive "golden" fixtures and fail CI on drift.
  - Canonicalize away volatile fields (e.g., auth headers, timestamps, content-length) while remaining strict about semantically meaningful fields.
  - Add Go adapter conformance tests that load the same fixtures and assert `LLMRequest -> provider HTTP` matches after canonicalization.
  - Provide a Node `fixturegen` tool (non-production) that regenerates fixtures using official SDKs.
  - Run `fixturegen` locally via the Nix dev shell to update fixtures; commit the results.
  - Pin SDK versions in lockfiles; dependency upgrades require explicit fixture regeneration + review.
  - Use automated dependency updates so fixture drift is detected at upgrade time and requires explicit approval.
  - Add CI coverage (GitHub Actions) so upgrades fail closed when fixtures drift.

## Task 3c: Bridge Providers (Post-MVP)

- Support a second provider integration mode for subscription-backed and local runtimes:
  - policy constraint: RuneCode does not ship/bundle/redistribute vendor CLIs or proprietary runtimes; integrate with user-installed official runtimes
  - `http` providers (MVP): model-gateway translates `LLMRequest -> provider HTTP` directly.
  - `bridge` providers (post-MVP): model-gateway translates `LLMRequest -> local RPC` and the local runtime performs upstream network calls.
    - prefer spawned child processes over stdio (no listening ports)
    - require runtime version pinning + per-request version logging
    - enforce an explicit "LLM-only" mode (deny tool execution and file operations)
    - run with isolated `HOME`/tool dirs pointing at an allowlisted provider sandbox directory
    - default to ephemeral sessions (no persisted conversation state unless enabled by manifest+policy)
    - prefer protocol-level contract tests (RPC fixtures) over HTTP wire fixtures

## Task 4: Data-Class Policy for Model Egress

- Default deny for third-party model usage.
- When explicitly opted in, allow only specific data classes (MVP baseline: `spec_text` only).
- Expanding allowed egress classes beyond `spec_text` (e.g., `diffs`, `approved_file_excerpts`) requires an explicit signed manifest opt-in and must be surfaced as a high-risk approval in the `moderate` profile.
- Unapproved excerpts (`unapproved_file_excerpts`) are never eligible for model egress.
- Enforce redaction at the boundary structurally:
  - use schema field classification metadata (`secret` fields are rejected/stripped)
  - prefer allowlists of permitted fields/classes over heuristic redaction

## Task 5: Audit + Quotas

- Log outbound requests (destination, bytes, timing) as audit events.
- Enforce basic quotas (requests/bytes/time) for the gateway role.

## Acceptance Criteria

- No other role can directly reach the public internet for model traffic.
- Workspace roles have zero direct network egress.
- Network egress is limited to explicit gateway roles; in MVP the only egress-capable gateway role is `model-gateway`.
- Secrets are never persisted in the launcher/broker/scheduler; only leases are used.
- Model-gateway blocks SSRF/DNS rebinding classes of attacks (private IPs, unsafe redirects) by default.
- Opt-in model egress is explicit, enforceable, and auditable.
