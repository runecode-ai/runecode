# Tasks

## Auth Gateway Role Contract

- [ ] Define the dedicated auth-gateway role and its limited public egress surface.
- [ ] Keep auth egress isolated from model egress and workspace access.

## Provider-Agnostic Auth Objects

- [ ] Define shared typed auth object families and versioned control flow.
- [ ] Ensure provider-specific flows extend the shared objects instead of redefining core semantics.

## Secret Handling + Token Storage

- [ ] Keep `secretsd` as the only long-lived secrets store.
- [ ] Forbid environment-variable or CLI-arg secret injection.
- [ ] Keep login, refresh, and token delivery flows typed and fail closed.

## Audit + Policy Integration

- [ ] Record typed auth lifecycle events without leaking secret values.
- [ ] Bind auth flows into policy and approval surfaces where required.

## Acceptance Criteria

- [ ] Provider auth runs through an auth-only gateway role.
- [ ] Long-lived auth material remains isolated to `secretsd`.
- [ ] Auth failures remain typed, auditable, and fail closed.
