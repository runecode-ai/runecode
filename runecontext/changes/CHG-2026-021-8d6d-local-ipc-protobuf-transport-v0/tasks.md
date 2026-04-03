# Tasks

## Proto Mapping for the Existing Logical Model

- [ ] Define `.proto` messages that map 1:1 to the existing logical local-API object model.
- [ ] Preserve existing error envelope, hashes, and schema-versioning rules.
- [ ] Keep persisted-object hashing and signing semantics defined by the logical JSON model.

## Local IPC Transport Requirements

- [ ] Keep the transport local-only by default.
- [ ] Keep framing, limits, deadlines, streaming backpressure, and max in-flight posture explicit regardless of encoding.
- [ ] Preserve deterministic broker enforcement for size and complexity limits.

## Optional Local-Only gRPC Profile

- [ ] Define any optional local-only gRPC profile without widening the trust boundary.

## Migration and Compatibility Rules

- [ ] Keep migration from JSON encoding explicit and reviewable.
- [ ] Preserve compatibility rules for existing logical contracts and persisted objects.

## Acceptance Criteria

- [ ] Protobuf stays an alternate local transport encoding rather than a new protocol.
- [ ] Local IPC trust-boundary rules and persisted canonicalization semantics remain unchanged.
