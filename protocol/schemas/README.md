# Protocol Schemas

- `protocol/schemas/manifest.json` is the authoritative bundle manifest for protocol object families and shared registries.
- `protocol/schemas/meta/manifest.schema.json` validates `protocol/schemas/manifest.json`.
- `protocol/schemas/meta/registry.schema.json` validates `protocol/schemas/registries/*.registry.json`.

## Status Semantics

- `mvp` means the object family is in MVP bundle scope. Some `mvp` families are intentionally minimal anchors until their owning spec task lands; those entries include a manifest `note` describing the pending task. In the current bundle, `PolicyDecision` and `Error` remain the main placeholder MVP families.
- `reserved` means the family is reserved for post-MVP extension work and must not expand capabilities without a later schema/task update.

## Current Lifecycle Coverage

- `PrincipalIdentity`, `RoleManifest`, and `CapabilityManifest` now carry the shared identity and lifecycle fields needed to bind requests, approvals, and audit records to active manifest context.
- `ApprovalRequest` and `ApprovalDecision` now bind immutable hash inputs, enforce explicit expiry semantics, and constrain MVP approval profiles to `moderate`.
- Approval trigger codes remain registry-owned values with fail-closed runtime validation; object schemas intentionally avoid hardcoding the full registry so new codes can land without a schema family bump.
- Timestamp ordering such as `requested_at < expires_at` remains a runtime validation rule even though both timestamps are required in the serialized protocol object.

## Schema Document IDs

- Object-schema `$id` values under `https://runecode.dev/protocol/schemas/...` are canonical schema identifiers for tooling and reference resolution.
- These `$id` values are not a network fetch contract. Validation and CI use the checked-in schema bundle as the source of truth.
