# Approval Profiles (Strict/Permissive) — Shaping Notes

## Scope

Define additional human-in-the-loop approval profiles beyond MVP `moderate`, and how they integrate with policy decisions and the TUI.

## Decisions

- Approval profiles affect *when* an allowed action requires explicit human approval; they must not change isolation boundaries or weaken invariants.
- Profiles must never convert `deny -> allow`.
- Profiles are signed inputs (part of the run/stage capability manifest) and are fully auditable.
- Adding or constraining profile values is a schema-versioned protocol change for every object family that carries the enum.
- MVP ships with `moderate` only; `strict` and `permissive` are post-MVP extensions.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-policy-engine-v0/`, `agent-os/specs/2026-03-08-1039-minimal-tui-v0/`
- Product alignment: Avoids approval fatigue while preserving security-first defaults.

## Standards Applied

- None yet.
