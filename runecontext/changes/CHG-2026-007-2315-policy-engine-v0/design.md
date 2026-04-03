# Design

## Overview
Implement the core policy evaluator that enforces manifests, role invariants, and explicit approvals.

## Key Decisions
- Deny-by-default everywhere; allow only via signed manifest.
- No automatic fallback to containers; container mode is explicit opt-in.
- MVP policy language is declarative and schema-validated (no general-purpose code execution during evaluation).
- Core security invariants are non-negotiable; any approval policy or UX setting may only tighten policy, never loosen it.
- Network egress is a hard boundary: workspace roles are offline; public egress is only via explicit gateway roles (model inference via `model-gateway`), and non-gateway network egress is not approvable.
- MVP uses checkpoint-style approvals (stage sign-off and explicit posture changes) instead of per-action nags.
- MVP supports a single approval profile (`moderate`); later profile expansion lives in `runecontext/changes/CHG-2026-014-0c5d-approval-profiles-strict-permissive/`.
- Approvals are typed, hash-bound to immutable inputs, and time-bounded (TTL/expiry); stale approvals are invalid.
- Policy decisions and failures use a shared protocol error envelope and stable reason codes.

## Main Workstreams
- Role + Run/Stage Policy Model
- Invariants (Fail Closed)
- Approval Policy (MVP: Moderate)
- Backend Selection Rules
- Decision Outputs

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
