# Design

## Overview
Define the shared bridge/runtime contracts for user-installed provider runtimes in explicit LLM-only mode.

## Key Decisions
- Shared bridge/runtime object families are defined once and reused by later provider specs.
- Compatibility is probe-driven and fail-closed; newer vendor versions are not trusted implicitly.
- Bridge runtimes remain LLM-only and never receive workspace or patch capabilities.
- Token delivery must avoid environment variables and raw secret logging.
- Audit and TUI surfaces must make untested-version and persisted-session posture visible.

## Main Workstreams
- Bridge Runtime Contract
- Compatibility + Probe Model
- Token Delivery + Session Rules
- Audit + UX Surfaces

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
