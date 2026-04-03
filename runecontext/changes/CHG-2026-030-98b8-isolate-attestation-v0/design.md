# Design

## Overview
Define isolate attestation evidence and verification that upgrades MVP TOFU provisioning to an attestable posture when required.

## Key Decisions
- MVP TOFU session metadata remains the compatibility baseline.
- Attestation adds stronger evidence; it does not replace the need for explicit session and image binding.
- Verifier, policy, and TUI surfaces must expose provisioning posture explicitly.
- Invalid or replayed attestation evidence fails closed when an attested posture is required.

## Main Workstreams
- Attestation Evidence Model
- Launch, Verification, and Policy Integration
- TUI + Audit Posture
- Fixtures + Cross-Platform Considerations

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
