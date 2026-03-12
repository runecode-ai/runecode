# Protocol & Schema Bundle v0

User-visible outcome: cross-component and cross-isolate communication is structured, schema-validated, and hash-addressable, enabling deterministic policy and audit.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-08-1039-protocol-schemas-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

## Task 2: Define Core Object Model

Define the minimal canonical objects needed for MVP:
- Role manifests and run/stage capability manifests (including explicit opt-ins and an approval profile).
- Artifact references (hash, size, content type, data class, origin).
- Audit events (hash-chained, signed, typed).
- Approval requests/decisions (typed, structured payloads for TUI display and deterministic enforcement).
- Policy decisions (allow/deny/require_human_approval) with reason codes.
- Model gateway protocol objects:
  - `LLMRequest` and `LLMResponse` (including streaming event shapes where applicable)
  - provider/model selection fields that do not allow arbitrary capability escalation
  - inputs must reference artifacts by hash (no raw prompt blobs crossing boundaries)
  - outputs are untrusted proposals and must be representable as typed artifacts
- Model output features (MVP):
  - streaming: supported; define incremental event types and completion semantics
  - tool calling: supported only as typed proposal objects (never direct execution)
    - `LLMRequest` carries an explicit tool allowlist per request.
    - Tool-call args are schema-validated; unknown/extra fields are rejected.
    - Add conservative limits (e.g., cap tool calls per response and cap total tool-call bytes).
  - structured JSON outputs: required for any machine-consumed output that can drive actions
- Audit events must be gateway-role aware:
  - include role identity and role kind (workspace vs gateway)
  - include egress category metadata for outbound network activity (model, auth, git, web, deps)
  - include allowlist identifiers and stable destination descriptors (without logging secret values)
- Reserved (post-MVP): `ProcessDefinition` (JSON/YAML) as the user-configurable process surface:
  - a schema-validated step graph model (sequential + branching + optional parallel blocks)
  - allowlisted RuneCode step types only (cannot introduce new capabilities)
    - Define the initial allowlist (illustrative):
      - `llm_request`
      - `workspace_read`
      - `workspace_edit`
      - `workspace_test`
      - `gate_run`
      - `approval_checkpoint`
      - `git_gateway_pr_create` (post-MVP; requires git-gateway)
      - `web_research` (post-MVP; requires web-research gateway)
  - per-step provider/model selection and step-level limits
- Reserved (post-MVP) protocol surface for `bridge` providers (local runtimes behind model-gateway):
  - a typed request/response envelope, runtime identity/version fields, and stable error taxonomy
  - Define an explicit "LLM-only" capability mode for bridge runtimes.
    - bridge requests to execute commands, read/write workspace files, or apply patches are denied and treated as policy violations
  - explicit streaming support and backpressure/queueing signals
  - contract-test fixtures for request/response envelopes and error mapping

## Task 3: Choose Schema + Validation Strategy

- Use JSON Schema as the single source of truth for MVP:
  - on-wire local RPC messages (broker <-> isolates <-> clients) use JSON (MVP)
  - on-disk manifests and policy documents use JSON
- Generate/derive validators for both Go and TS from the same schema bundle.
- To keep post-MVP protobuf migration feasible (with an optional gRPC facade), restrict schemas to an MVP profile that maps cleanly to protobuf messages:
  - avoid regex-heavy schemas and dynamic keys (`patternProperties` / arbitrary maps) in on-wire messages
  - model unions via an explicit discriminator field (no ambiguous `oneOf` without a tag)
  - keep numeric ranges within I-JSON expectations; represent high-precision numbers as strings
- Fail closed at trust boundaries:
  - reject unknown fields (no permissive parsing)
  - enforce message size limits and structural complexity limits (depth / array length)
- Canonicalization for hashing/signing (MVP requirement):
  - Use RFC 8785 (JSON Canonicalization Scheme, JCS) for canonical bytes.
  - Prohibit floats/NaN/Infinity in hashed/signed objects; use integers or strings.
  - Encode bytes as base64 strings; timestamps as RFC 3339 strings; durations as integer milliseconds.
  - Hash/sign inputs are the canonical JSON bytes produced by JCS.
- Add field-level data classification metadata in schemas (`public | sensitive | secret`) to support structural redaction/boundary enforcement.

## Task 6: On-Wire Encoding Migration Plan (Post-MVP)

- Keep the logical object model stable and documented independent of encoding.
- Prefer protobuf message encoding for on-wire local RPC post-MVP without requiring gRPC:
  - define `.proto` message definitions that map 1:1 to the logical model
  - keep golden fixtures and cross-language tests so JSON and protobuf encodings are behaviorally equivalent
  - continue using local IPC transports (UDS / named pipes / vsock / virtio-serial); do not introduce a network API by default
  - keep message framing, size limits, deadlines/timeouts, and backpressure as explicit requirements regardless of transport
- gRPC is optional (post-MVP) and must remain local-only:
  - prefer gRPC over Unix domain sockets (Unix) and OS-native local IPC (e.g., named pipes on Windows) where supported
  - do not use TCP by default
  - if TCP loopback is used for compatibility, require one of:
    - mTLS with pinned/trusted local certificates, or
    - a strong, short-lived local token mechanism (stored with strict filesystem permissions)
  - binding safety is a security requirement: never bind privileged APIs to non-loopback interfaces
- Do not change hashing/signing semantics for persisted/signed objects (canonicalization remains defined by this spec).

## Task 4: Versioning + Compatibility Rules

- Every top-level object includes explicit `schema_id` and `schema_version` fields.
- Manifest hashes bind to the specific schema version used for validation/canonicalization.
- Compatibility model (MVP):
  - no "loose" parsing at trust boundaries (unknown fields are rejected)
  - changes require a schema version bump
  - older schema versions remain verifiable (verifier keeps old schemas)
- If the verifier encounters an unsupported schema version, verification fails closed with a clear reason code.

Approval profile versioning note:
- MVP supports a single approval profile value (`moderate`). Adding new profiles (e.g., `strict`, `permissive`) is a schema version bump and is post-MVP.

Approval profile semantics note:
- Approval profiles must never convert `deny -> allow`.
- Approval profiles only affect whether an otherwise-allowed action returns `allow` vs `require_human_approval`.
- Unknown profile values fail closed.

## Task 5: Reference Fixtures

- Add small, checked-in example manifests and events that validate against schemas.
- Include both a “microVM stage” and a “container stage (explicit opt-in)” fixture.
- Include an MVP approval profile fixture (`moderate`) embedded in the run/stage manifest.
- Include a minimal `LLMRequest`/`LLMResponse` fixture that uses only `spec_text` inputs.
- Include fixtures for:
  - streaming event sequences (including interruption/cancellation)
  - tool-call proposal outputs (schema-valid)
  - structured JSON output validation (schema pass/fail cases)
  - bridge provider envelope + error taxonomy examples (post-MVP)
  - ProcessDefinition example (post-MVP; validates but cannot expand capabilities)
- Add canonicalization + hashing fixtures:
  - canonical JSON bytes (golden)
  - expected hash outputs
  - (where relevant) expected signature verification outcomes

## Acceptance Criteria

- Go and TS validate the same fixtures deterministically and reject the same invalid inputs.
- Canonical bytes and hash inputs are stable across platforms (golden fixtures pass in CI).
- Schema versions are explicit and bound to hashes; verification fails closed on unknown versions.
- All cross-boundary messages used in MVP are schema-defined and validated.
- The schema/profile avoids constructs that would make post-MVP protobuf migration impractical.
