# External Audit Anchoring v0 — Post-MVP

User-visible outcome: RuneCode can optionally anchor audit segment roots to later non-MVP targets with explicit egress where needed, typed receipts, and verifier-visible posture rather than keeping those follow-on targets as notes in the MVP anchoring spec.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-13-1603-external-audit-anchoring-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

Parallelization: docs-only; safe to do anytime.

## Task 2: Later Anchor Target Model

- Define the later non-MVP anchor targets, including at least the planned `anchor_kind` families:
  - `tpm_pcr`
  - `rfc3161`
  - `witness_service`
  - `transparency_log`
- Each target kind needs a typed target descriptor and receipt payload contract rather than a freeform blob.
- Receipt objects must preserve the shared anchor-receipt envelope from `agent-os/specs/2026-03-08-1039-audit-anchoring/` while adding target-specific typed fields.

Parallelization: can be designed in parallel with verifier work once the shared receipt envelope is fixed.

## Task 3: Egress + Trust Boundary Model

- External anchoring requires explicit signed-manifest opt-in and must never silently enable network access.
- External anchor traffic must use an explicit allowlist and a non-workspace execution pathway.
- Define how policy and audit distinguish:
  - anchor target selection
  - anchor-attempt approvals where required
  - temporary target unavailability vs invalid receipts
- Secret material for target authentication, if any, must follow the same no-env-var/no-raw-log posture as other gateway-style integrations.

Parallelization: can be designed in parallel with gateway/policy follow-on work; implementation should wait until the explicit egress model is fixed.

## Task 4: Receipt, Audit, and Verification Integration

- Store external anchor receipts as artifacts and reference them from anchoring audit events.
- Verification output must distinguish:
  - locally anchored only
  - externally anchored and valid
  - attempted but unavailable/deferred external anchoring
  - invalid or unverifiable external receipts
- Verification remains fail closed on invalid receipts and never rewrites existing audit history.

Parallelization: can be implemented in parallel with verifier/reporting work once receipt schemas and audit event references are stable.

## Task 5: Fixtures + Adapter Conformance

- Add checked-in fixtures for representative external anchor receipts and invalid cases for each supported target kind.
- Keep fixture updates explicit and reviewable; CI verifies but does not regenerate them implicitly.

Parallelization: fixtures can be developed in parallel with target adapters so long as they validate against the same receipt schemas.

## Acceptance Criteria

- External anchoring targets are defined in a dedicated later spec rather than remaining as a note in the MVP anchoring spec.
- External anchoring never silently enables network access.
- Receipt verification is typed, auditable, and fail closed on invalid data.
