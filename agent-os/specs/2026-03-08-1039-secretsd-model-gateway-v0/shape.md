# Secretsd + Model-Gateway v0 — Shaping Notes

## Scope

Implement secrets storage/lease issuance and a dedicated model-gateway that centralizes third-party model egress.

## Decisions

- Third-party model usage is explicit opt-in; deny by default.
- Model traffic goes only through model-gateway; workspace roles remain offline.
- Only `secretsd` stores long-lived secrets; other daemons/components must not persist secret values (leases only).
- Secrets storage fails closed by default if secure key storage is unavailable (no silent plaintext-on-disk fallback).
- Model gateway egress is hardened against SSRF/DNS rebinding and enforces TLS-only provider connections.
- MVP default for model egress is `spec_text` only; allowing `diffs` or `approved_file_excerpts` is an explicit, auditable opt-in.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-policy-engine-v0/`
- Product alignment: Prevents any single component from combining workspace access + public egress + long-lived secrets.

## Standards Applied

- None yet.
