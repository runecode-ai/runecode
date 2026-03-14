# Local IPC Protobuf Transport v0 — Shaping Notes

## Scope

Define the post-MVP local-transport migration from JSON-on-wire to protobuf while preserving local-only security and the logical message contract.

## Decisions

- The logical object model remains authoritative; protobuf is an alternate encoding, not a new protocol.
- Persisted-object hashing/signing semantics do not change.
- Local IPC safety requirements (binding, auth, framing, limits, deadlines, backpressure) remain explicit regardless of transport.
- gRPC is optional and local-only.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-broker-local-api-v0/`, `agent-os/specs/2026-03-08-1039-protocol-schemas-v0/`, `agent-os/specs/2026-03-08-1039-windows-microvm-runtime/`
- Product alignment: improves local IPC efficiency/ergonomics without expanding the trust boundary.

## Standards Applied

- `security/trust-boundary-interfaces` - local IPC encoding changes are part of the shared boundary contract.
- `global/deterministic-check-write-tools` - generated protobuf artifacts and fixtures must stay explicit and reviewable.
