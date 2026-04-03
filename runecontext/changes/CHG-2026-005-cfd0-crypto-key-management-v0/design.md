# Design

## Overview
Establish the minimum cryptographic root-of-trust for signing manifests and audit events in MVP.

## Key Decisions
- Isolates sign their own audit events; the control plane must not be able to forge isolate-attributed events.
- Manifest signing requires explicit user presence.
- Isolate identity private keys are generated and stored inside the isolate boundary; the launcher/broker must never possess isolate private keys.
- If secure key storage is unavailable (hardware/OS keystore), the system must fail closed by default (no silent plaintext fallback).
- If passphrase-derived encryption is explicitly opted into (dev/portable mode), the KDF + passphrase policy is specified and audited (Argon2id; minimum strength requirements).
- Isolate key provisioning is TOFU for MVP; provisioning mode and handshake binding context are recorded and surfaced as a degraded posture.
- Signature envelopes include `{alg, key_id}` to keep algorithm agility feasible.

## Main Workstreams
- Define MVP Key Hierarchy
- Key Storage + Posture Recording
- User-Presence Approval Hook (MVP Baseline)
- Sign/Verify Primitives
- Rotation + Revocation (Minimal)

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
