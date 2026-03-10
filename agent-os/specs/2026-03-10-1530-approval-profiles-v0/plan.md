# Approval Profiles (Strict/Permissive) — Post-MVP

User-visible outcome: RuneCode supports selectable human-in-the-loop approval profiles (beyond MVP `moderate`) without weakening core security invariants.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-10-1530-approval-profiles-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

## Task 2: Approval Profile Model (Post-MVP)

- Define additional approval profiles:
  - `strict`: maximize human involvement (frequent, granular approvals)
  - `permissive`: minimize interruptions while preserving the same security invariants
- Define how profiles are selected:
  - profile is an explicit field in the run/stage capability manifest (signed input)
  - the system fails closed on unknown profile values
- Profiles must never convert `deny -> allow`; they only affect whether an otherwise-allowed action requires explicit human approval.
- Define the non-negotiable invariant set that profiles cannot bypass.

## Task 3: Strict Profile Semantics

- Define which action categories require approval in `strict` mode (illustrative):
  - step start/resume
  - workspace writes
  - command execution (even via allowlisted executors)
  - artifact publication beyond the current step
  - all egress-related opt-ins (model, git, web)
- Define batching rules to prevent UX deadlocks (e.g., "approve N related writes" in one approval request).

## Task 4: Permissive Profile Semantics

- Define `permissive` mode as "approve at milestones" while keeping the same enforcement boundaries:
  - stage manifest sign-off remains required
  - posture-changing actions (e.g., container backend, new egress scopes) remain explicit approvals
  - gate overrides remain explicit approvals
  - when git-gateway exists: require an explicit final approval for git remote state changes (push/tag/PR creation)
  - post-MVP review: consider a dedicated `git-remote-ops` approval trigger category once git-gateway exists

## Task 5: Policy + Runner + TUI Integration

- Extend the policy engine to interpret `strict` and `permissive` profiles.
- Ensure the workflow runner pauses only on policy-returned `require_human_approval` decisions.
- Ensure the TUI can:
  - display the active profile
  - explain why an approval is required (reason codes + structured payload)
  - show what changes if approved

## Acceptance Criteria

- The system supports `strict` and `permissive` profiles with deterministic behavior and clear audit events.
- Profiles cannot weaken core invariants (deny-by-default, no escalation-in-place, no host mounts, no unsafe role capability combinations).
- Profiles must never convert `deny -> allow`.
- Attempting to use an unknown profile value fails closed.
