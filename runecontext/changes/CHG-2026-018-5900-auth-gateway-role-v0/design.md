# Design

## Overview
Define the dedicated auth gateway role, provider-agnostic auth object families, and secret-safe OAuth-style login and refresh flows.

## Key Decisions
- Auth egress and model egress are separated.
- `secretsd` is the only long-lived secrets store; there is no second credential cache.
- Shared auth object families are provider-agnostic, typed, and versioned; provider specs extend them rather than redefining the control flow.
- No environment-variable or CLI-arg secret injection.
- Auth flows are typed, auditable, and fail closed on state/protocol mismatches.

## Main Workstreams
- Auth Gateway Role Contract
- Provider-Agnostic Auth Objects
- Secret Handling + Token Storage
- Audit + Policy Integration

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
