# Design

## Overview
Implement the brokered local API contract that mediates all isolate RPC, local client access, and broker-managed artifact routing.

## Key Decisions
- Star topology only; no direct isolate-to-isolate communication.
- Schema validation at boundaries; rate limiting and backpressure are mandatory.
- Broker enforces concrete default limits (message size/complexity/in-flight/streaming) with audited overrides.
- The local API is per-user IPC with strict filesystem permissions; authentication fails closed when OS peer credentials are unavailable.
- Errors use a shared typed envelope and stable reason codes (no ad-hoc error shapes).
- MVP uses JSON on-wire; later transport migration is specified separately so this spec stays focused on the MVP broker/API contract.

## Main Workstreams
- Broker Responsibilities (MVP)
- Local Client API
- Local Auth
- Artifact Routing Integration

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
