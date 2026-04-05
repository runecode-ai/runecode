# Design

## Overview
Define the local-only protobuf transport as an alternate encoding for the existing logical local API without changing persisted hashing semantics.

## Key Decisions
- The logical object model remains authoritative; protobuf is an alternate encoding, not a new protocol.
- Persisted-object hashing and signing semantics do not change; protobuf transport continues to preserve the existing RFC 8785 JCS-based logical JSON hashing contract for persisted objects.
- Local IPC safety requirements (binding, auth, framing, limits, deadlines, backpressure) remain explicit regardless of transport.
- gRPC is optional and local-only.

## Main Workstreams
- Proto Mapping for the Existing Logical Model
- Local IPC Transport Requirements
- Optional Local-Only gRPC Profile
- Migration and Compatibility Rules

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
