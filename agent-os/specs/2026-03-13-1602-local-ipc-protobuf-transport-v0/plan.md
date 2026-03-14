# Local IPC Protobuf Transport v0 — Post-MVP

User-visible outcome: RuneCode's local broker API can migrate to protobuf-based IPC without changing its logical object model, security posture, or cross-language contract tests.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-13-1602-local-ipc-protobuf-transport-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

Parallelization: docs-only; safe to do anytime.

## Task 2: Proto Mapping for the Existing Logical Model

- Define `.proto` messages that map 1:1 to the existing logical local-API object model.
- Proto message definitions must preserve:
  - the same logical meaning as the JSON schema families
  - exact version binding for communicated object families
  - fail-closed handling of unsupported versions and unknown message kinds
- On-wire protobuf adoption must not change hashing/signing semantics for persisted or signed objects; canonicalization for persisted objects remains defined by the logical JSON model.

Parallelization: can be designed in parallel with broker work once the logical object model is stable.

## Task 3: Local IPC Transport Requirements

- Keep the transport local-only by default:
  - Unix domain sockets on Unix-like platforms
  - named pipes on Windows
  - other local-only transports such as vsock/virtio-serial only when explicitly justified by an isolation backend
- Message framing, size limits, deadlines/timeouts, streaming backpressure, and max in-flight posture remain explicit protocol requirements independent of encoding.
- Migration must preserve deterministic broker enforcement for message size/complexity limits rather than pushing them into transport-specific best-effort behavior.

Parallelization: transport work can be implemented in parallel with fixture generation once the mapping rules are fixed.

## Task 4: Optional Local-Only gRPC Profile

- gRPC is optional and must remain local-only.
- Prefer gRPC over Unix domain sockets or OS-native local IPC where supported.
- Do not use TCP by default.
- If compatibility requires loopback TCP:
  - require mTLS with pinned/trusted local certificates or a strong short-lived local token stored with strict filesystem permissions
  - never bind privileged APIs to non-loopback interfaces

Parallelization: can be designed in parallel with the base protobuf transport work; implementation should wait until the local-only auth/binding posture is fixed.

## Task 5: Dual-Encoding Fixtures + Migration Plan

- Keep golden fixtures and cross-language tests that prove JSON and protobuf encodings are behaviorally equivalent for the same logical messages.
- Define migration posture for local sessions:
  - coordinated local restarts rather than mixed-version live negotiation
  - deterministic rollback/downgrade behavior
  - explicit audit/version metadata for participating components
- CI must validate both encoding lanes without implicitly regenerating fixtures.

Parallelization: fixture and migration-plan work can proceed in parallel with implementation if they use the same authoritative message set.

## Acceptance Criteria

- Protobuf message definitions preserve the existing logical object model and fail-closed versioning posture.
- Local IPC remains local-only by default.
- Optional gRPC, if adopted, stays local-only and never weakens binding/authentication safety.
- JSON and protobuf lanes have shared cross-language fixtures proving behavioral equivalence.
