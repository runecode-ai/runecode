# Design

## Overview
Create a tamper-evident audit log with signed, hash-chained events and a verifier.

## Key Decisions
- Audit events are append-only and hash-chained.
- Isolate-attributed events must be signed by isolate keys; writers must verify signatures.
- Audit log storage is encrypted at rest by default (no silent plaintext mode).
- Audit logs are segmented for retention/archival without breaking verifiability.
- Segment roots may be anchored via verifiable receipts (see `runecontext/changes/CHG-2026-006-84f0-audit-anchoring-v0/`).
- Audit events are gateway-role aware (role identity + egress category metadata for outbound network activity), without logging secret values.
- Verification produces a machine-readable artifact (`audit_verification_report`) that can be stored and reviewed later (not just printed).
- Ordering integrity is defined using per-signer monotonic sequence + hash chaining; wall-clock timestamps are advisory metadata.

## Main Workstreams
- Audit Event Model
- Append-Only Audit Writer
- Redaction Boundaries (Minimal)
- Verify Command
- Segmentation + Retention (Minimal)

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
