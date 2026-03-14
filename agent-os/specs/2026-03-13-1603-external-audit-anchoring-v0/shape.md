# External Audit Anchoring v0 — Shaping Notes

## Scope

Define later non-MVP anchoring targets and their receipt/egress model on top of the MVP local anchoring baseline.

## Decisions

- Later non-MVP anchor targets use typed descriptors and receipt payloads.
- External anchoring is explicit opt-in with a clear egress model.
- Verification/reporting must distinguish valid external anchors from deferred, unavailable, or invalid states.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-audit-anchoring/`, `agent-os/specs/2026-03-08-1039-audit-log-verify-v0/`, `agent-os/specs/2026-03-08-1039-policy-engine-v0/`
- Product alignment: extends tamper-evidence without weakening offline-by-default or trust-boundary posture.

## Standards Applied

- `security/trust-boundary-layered-enforcement` - external anchoring must not bypass explicit gateway/policy controls.
- `global/deterministic-check-write-tools` - receipt fixtures and verifier outputs must remain explicit and reviewable.
