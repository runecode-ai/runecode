# Tasks

## Runner Contract

- [ ] Implement untrusted runner orchestration with stable broker-facing contracts.
- [ ] Keep LangGraph internal and non-canonical.

## Durable State

- [ ] Implement persisted run-state transitions and step-attempt tracking.
- [ ] Implement explicit crash recovery and idempotency rules.

## Execution Loop

- [ ] Enforce propose, validate, authorize, execute, and attest transitions.
- [ ] Keep approvals typed, bounded, and resumable.

## Acceptance Criteria

- [ ] Runs can pause/resume and recover safely after process failure.
- [ ] Runner cannot bypass policy or direct execution boundaries.
