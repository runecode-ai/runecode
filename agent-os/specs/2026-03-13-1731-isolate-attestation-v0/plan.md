# Isolate Attestation v0 — Post-MVP

User-visible outcome: RuneCode can upgrade MVP TOFU isolate key binding to measured, attestable isolate identity without changing the core audit/event model.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-13-1731-isolate-attestation-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

Parallelization: docs-only; safe to do anytime.

## Task 2: Attestation Evidence Model

- Define the attestation evidence objects that upgrade an isolate/session from TOFU to an attested binding.
- Evidence must bind at least:
  - isolate signing key identity
  - measured boot or equivalent launch evidence
  - image/toolchain digest identity
  - session nonce or equivalent replay-resistant binding
  - verifier result and posture metadata
- Preserve compatibility with the MVP audit/event envelope so previously recorded TOFU fields do not need a format break.

Parallelization: can be designed in parallel with launcher/image-signing work once the shared session metadata shape is stable.

## Task 3: Launch, Verification, and Policy Integration

- Define how the launcher obtains and verifies attestation evidence before trusting the upgraded isolate binding.
- Policy and verifier flows must distinguish:
  - TOFU-only sessions
  - attested-and-valid sessions
  - attestation required but unavailable
  - invalid or replayed attestation evidence
- Attestation failures must fail closed when an attested posture is required.

Parallelization: can be implemented in parallel with audit verification and image/toolchain signing follow-on work once the evidence model is fixed.

## Task 4: TUI + Audit Posture

- Audit metadata and TUI surfaces must make the provisioning posture explicit.
- Replace the MVP degraded TOFU-only posture with an attested posture only when verification succeeds.
- Record why an attested posture was unavailable or rejected without leaking sensitive local details.

Parallelization: can be implemented in parallel with audit/TUI follow-on work once posture fields are fixed.

## Task 5: Fixtures + Cross-Platform Considerations

- Add checked-in fixtures for valid, invalid, replayed, and unavailable attestation evidence.
- Account for platform-specific attestation sources without making the shared verifier contract platform-specific.

Parallelization: fixtures can be created in parallel with platform-specific attestation adapters so long as they validate against the same shared evidence model.

## Acceptance Criteria

- RuneCode can represent and verify an attested isolate/session binding without changing the MVP audit/event contract.
- Verifiers and TUI distinguish TOFU-only, attested-valid, unavailable, and invalid attestation states.
- Replay and invalid-evidence cases fail closed when attestation is required.
