# Audit Anchoring — Shaping Notes

## Scope

Add anchoring receipts for audit segment roots and integrate them with verification.
MVP includes a local-only anchoring mode (no network egress). Later external-anchoring work is tracked separately.

## Decisions

- Anchoring is an explicit step and produces receipts.
- Failures are recorded; no history rewriting.
- MVP baseline anchoring is local-only and requires explicit user presence to mint receipts.
- Verification distinguishes `verified_unanchored` vs `verified_anchored`; missing anchors are not a verification failure by default.
- Invalid receipts fail closed.
- External anchoring lives in `agent-os/specs/2026-03-13-1603-external-audit-anchoring-v0/` and requires an explicit egress model.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-audit-log-verify-v0/`, `agent-os/specs/2026-03-13-1603-external-audit-anchoring-v0/`
- Product alignment: Strengthens tamper-evidence for sharing and forensics.

## Standards Applied

- None yet.
