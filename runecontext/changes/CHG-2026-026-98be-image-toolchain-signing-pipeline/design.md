# Design

## Overview
Define the signing and verification pipeline for isolate images and toolchains with fail-closed enforcement.

## Key Decisions
- Image/toolchain signing keys are separate from manifest signing.
- Enforcement is fail-closed.

## Main Workstreams
- Signing Key Hierarchy
- Build + Publication Pipeline
- Launcher Enforcement
- Audit + Verification Integration

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
