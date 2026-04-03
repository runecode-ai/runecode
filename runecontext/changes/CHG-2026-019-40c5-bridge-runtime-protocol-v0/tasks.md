# Tasks

## Bridge Runtime Contract

- [ ] Define shared bridge/runtime object families once for later provider specs.
- [ ] Keep bridge runtimes in explicit LLM-only mode with no workspace or patch capabilities.

## Compatibility + Probe Model

- [ ] Define probe-driven compatibility checks.
- [ ] Fail closed on unsupported or untested runtime versions instead of trusting newer vendor versions implicitly.

## Token Delivery + Session Rules

- [ ] Keep token delivery away from environment variables and raw secret logging.
- [ ] Define persisted-session posture and lifecycle rules explicitly.

## Audit + UX Surfaces

- [ ] Surface untested-version and persisted-session posture in audit and TUI flows.
- [ ] Keep bridge runtime behavior auditable and reviewable.

## Acceptance Criteria

- [ ] Shared bridge contracts are reusable by provider-specific changes.
- [ ] Bridge runtimes remain LLM-only and fail closed on unsupported versions or unsafe token-delivery paths.
