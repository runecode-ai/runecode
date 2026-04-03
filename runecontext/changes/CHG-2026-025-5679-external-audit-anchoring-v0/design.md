# Design

## Overview
Define typed external anchoring targets, receipt verification, and explicit egress controls for non-local audit anchoring.

## Key Decisions
- Later non-MVP anchor targets use typed descriptors and receipt payloads.
- External anchoring is explicit opt-in with a clear egress model.
- Verification/reporting must distinguish valid external anchors from deferred, unavailable, or invalid states.

## Main Workstreams
- Later Anchor Target Model
- Egress + Trust Boundary Model
- Receipt, Audit, and Verification Integration
- Fixtures + Adapter Conformance

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
