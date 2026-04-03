# Design

## Overview
Implement the secrets daemon and model egress gateway, including lease handling, typed model requests, egress hardening, and provider drift detection.

## Key Decisions
- Third-party model usage is explicit opt-in; deny by default.
- Model traffic goes only through model-gateway; workspace roles remain offline.
- Only `secretsd` stores long-lived secrets; other daemons/components must not persist secret values (leases only).
- Secret values are never accepted or delivered via environment variables; use stdin/file-descriptor onboarding and brokered lease IPC.
- Secrets storage fails closed by default if secure key storage is unavailable (no silent plaintext-on-disk fallback).
- Model gateway egress is hardened against SSRF/DNS rebinding and enforces TLS-only provider connections.
- Model-gateway uses a typed `LLMRequest`/`LLMResponse` boundary; inputs reference artifacts by hash (no freeform prompt blobs).
- Model-gateway fetches artifact bytes by hash (via broker-mediated CAS access) and fails closed on disallowed data classes.
- Model-gateway is implemented in Go for MVP to minimize TCB; provider request shaping stays inside the Go gateway.
- Official provider SDKs are used only for fixture generation and drift detection, not in the production egress path.
- Streaming and tool calling are supported only within the typed boundary; tool calls remain untrusted proposals.
- MVP default for model egress is `spec_text` only; allowing `diffs` or `approved_file_excerpts` is an explicit, auditable opt-in.

## Main Workstreams
- Secretsd MVP Interface
- Model-Gateway Role
- Provider Adapters + Drift Detection (MVP)
- Bridge Provider Follow-On Specs
- Data-Class Policy for Model Egress
- Audit + Quotas

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
