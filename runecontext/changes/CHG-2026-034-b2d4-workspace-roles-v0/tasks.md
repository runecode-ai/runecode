# Tasks

## Role Set

- [ ] Implement `workspace-read`, `workspace-edit`, and `workspace-test` roles.
- [ ] Enforce explicit capability manifests for each role.

## Execution Boundaries

- [ ] Implement constrained executors with allowlisted operations.
- [ ] Block shell passthrough behavior.

## Offline Posture

- [ ] Enforce no direct network egress from workspace roles.
- [ ] Route required cross-boundary data movement through artifacts.

## Acceptance Criteria

- [ ] Role execution remains least-privilege and offline by default.
- [ ] Workspace roles cannot bypass runner, policy, or broker controls.
