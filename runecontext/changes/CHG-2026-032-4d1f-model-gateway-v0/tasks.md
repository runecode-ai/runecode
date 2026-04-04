# Tasks

## Gateway Boundary

- [ ] Implement an egress allowlisted model-gateway role.
- [ ] Ensure no workspace access and no long-lived secret storage.

## Typed Model Contracts

- [ ] Enforce typed request/response boundaries.
- [ ] Ensure tool calls remain untrusted proposals.

## Egress Hardening

- [ ] Enforce destination validation and TLS requirements.
- [ ] Apply strict timeout and response-size limits.

## Data Class + Policy

- [ ] Enforce allowlisted egress data classes.
- [ ] Block disallowed classes at the boundary.

## Audit + Quotas

- [ ] Audit outbound destination, bytes, timing, and outcome.
- [ ] Enforce basic quota controls.

## Acceptance Criteria

- [ ] Model egress occurs only through the gateway boundary.
- [ ] Gateway behavior is policy-controlled, auditable, and fail-closed.
