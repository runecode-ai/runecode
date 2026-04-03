# Tasks

## Role + Run/Stage Policy Model

- [ ] Define how role manifests and run/stage capability manifests combine into an effective policy.
- [ ] Bind all decisions to the manifest hash (no “implicit” policy).
- [ ] Define the MVP policy language and evaluation semantics:
  - declarative, schema-validated policy documents (no embedded scripting)
  - explicit precedence rules (e.g., explicit deny > require_human_approval > allow)
  - stable reason codes and structured decision details
  - use the shared protocol error taxonomy/envelope for denials and failures (see `runecontext/specs/protocol-schema-bundle-v0.md`)
- [ ] Define gateway-role concepts in the policy model (even if MVP ships only `model-gateway`):
  - gateway role kinds (model, auth, git, web, deps)
  - per-gateway destination allowlists as signed inputs (no URL-based ad-hoc decisions)
  - required hardening invariants for any egress-enabled gateway (SSRF/DNS rebinding protections, redirect rules, TLS requirements)
- [ ] Define policy loading/tamper resistance:
  - policy inputs are content-addressed and bound to signed manifests
  - policy evaluation rejects inputs that do not validate against the schema version

Parallelization: can be implemented in parallel with the broker once the policy decision request/response schema is fixed.

## Invariants (Fail Closed)

- [ ] Enforce MVP invariants:
  - no escalation-in-place
  - deny-by-default for network/filesystem/shell/secrets
  - network egress is a hard boundary:
    - workspace roles have zero direct network egress
    - public network egress is only allowed for explicit gateway roles (declared in signed manifests and enforced by policy)
      - model inference egress is only possible via the dedicated `model-gateway` role
      - other egress categories (git remote ops, web research, dependency fetch, provider auth) must be isolated behind dedicated gateway roles with strict allowlists and no workspace access
    - any non-gateway network egress attempt is denied and is not approvable
  - no single role combines workspace RW + public egress + long-lived secrets

Parallelization: can be implemented in parallel with gateway specs; invariants are the shared non-negotiable core.

## Approval Policy (MVP: Moderate)

- [ ] Define an approval policy model that controls when an otherwise-allowed action requires explicit human approval.
- [ ] Approval policy/profile is a signed input (part of the run/stage capability manifest).
- [ ] Approval profiles affect only *when* explicit human approval is required for actions already allowed by invariants + the signed manifest:
  - profiles must never convert `deny -> allow`
  - profiles may add additional `require_human_approval` decisions and/or additional denies, but cannot expand capabilities

Approval lifecycle semantics (MVP):
- [ ] Every approval request and decision is typed and binds to immutable inputs:
  - approval request includes `{manifest_hash, action_request_hash, relevant_artifact_hashes}`
  - approval decision references the approval request hash (no free-floating “approve by label”)
- [ ] Expiry/timeout:
  - approval requests have an explicit TTL (default: 30 minutes) and expire deterministically
  - expired approvals cannot be used; the runner must re-request approval
- [ ] Staleness:
  - if any bound hash changes while awaiting approval, the pending approval is invalidated and must be re-issued
- [ ] MVP implements a single approval profile: `moderate`:
  - approvals are checkpoint-style, not per-micro-action:
    - stage capability manifest sign-off (always; includes a structured summary of requested high-risk capability categories)
    - reduced-assurance opt-ins (e.g., container backend)
    - gate overrides
    - enabling gateway egress and/or expanding egress scope (MVP focus: enabling third-party model egress via `model-gateway` and expanding allowed model egress data classes beyond the baseline `spec_text` only)
  - moderate trigger categories (must be explicit in the stage manifest and surfaced at sign-off):
    - file writes outside the workspace volume/root (outside the declared workspace path allowlist)
    - secret access (issuing leases from `secretsd`)
    - dependency/package installation
    - system-modifying command execution (beyond ordinary workspace edit/test executors)
    - gateway egress scope changes (enable a gateway role; change allowlists; expand allowed egress data classes)
    - actions wholly inside the workspace sandbox and within the signed manifest execute without intermediate approvals
- [ ] Later approval-profile expansion lives in `runecontext/changes/CHG-2026-014-0c5d-approval-profiles-strict-permissive/`, and the dedicated `git-remote-ops` trigger category lives in `runecontext/changes/CHG-2026-002-33c5-git-gateway-commit-push-pr/`.

Parallelization: can be implemented in parallel with the workflow runner pause/resume work as long as the approval request/decision schemas are shared.

## Backend Selection Rules

- [ ] MicroVM is the default backend when available.
- [ ] Container backend is only allowed with an explicit opt-in recorded as an approval + audit event.
- [ ] The system must not automatically fall back from microVM to containers.

Parallelization: can be implemented in parallel with launcher backends; it depends only on emitting explicit posture + approval events.

## Decision Outputs

- [ ] Standardize policy decisions:
  - `allow | deny | require_human_approval`
  - stable reason codes
  - structured “required approvals” payloads
- [ ] Decision artifacts must include hashes of all evaluated inputs (manifest hash, request hash, relevant artifact hashes).

Parallelization: can be implemented in parallel with protocol schemas; it depends on a stable decision artifact schema and shared error taxonomy.

## Acceptance Criteria

- [ ] Every action request is evaluated deterministically and produces a policy decision artifact.
- [ ] Policy evaluation does not execute arbitrary code and is deterministic for identical inputs.
- [ ] Container usage is blocked unless explicitly opted in and recorded.
- [ ] Violations are auditable and do not partially execute.
- [ ] With the `moderate` approval profile, a typical offline edit+gate step can execute without intermediate approvals once the stage manifest is signed.
