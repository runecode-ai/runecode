# Isolate Attestation v0 — Shaping Notes

## Scope

Define later measured-boot/attestation support that upgrades MVP TOFU isolate identity to a stronger, verifiable posture.

## Decisions

- MVP TOFU session metadata remains the compatibility baseline.
- Attestation adds stronger evidence; it does not replace the need for explicit session and image binding.
- Verifier, policy, and TUI surfaces must expose provisioning posture explicitly.
- Invalid or replayed attestation evidence fails closed when an attested posture is required.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-crypto-key-mgmt-v0/`, `agent-os/specs/2026-03-08-1039-launcher-microvm-backend-v0/`, `agent-os/specs/2026-03-08-1039-audit-log-verify-v0/`, `agent-os/specs/2026-03-08-1039-image-toolchain-signing/`
- Product alignment: strengthens isolate provenance and trust without weakening fail-closed defaults.

## Standards Applied

- `security/trust-boundary-layered-enforcement` - attestation must strengthen, not bypass, existing launch/auth/audit controls.
- `global/deterministic-check-write-tools` - evidence fixtures and verifier expectations must stay explicit and deterministic.
