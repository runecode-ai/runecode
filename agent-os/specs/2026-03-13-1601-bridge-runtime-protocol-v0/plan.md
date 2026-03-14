# Bridge Runtime Protocol v0 — Post-MVP

User-visible outcome: RuneCode can integrate user-installed local provider runtimes through a shared, auditable bridge contract that keeps those runtimes in explicit LLM-only mode and avoids ad-hoc RPC/auth shapes.

## Task 1: Save Spec Documentation

Create `agent-os/specs/2026-03-13-1601-bridge-runtime-protocol-v0/` with:
- `plan.md` (this file)
- `shape.md`
- `standards.md`
- `references.md`
- `visuals/` (empty)

Parallelization: docs-only; safe to do anytime.

## Task 2: Shared Bridge Object Families

- Define the shared bridge/runtime object families used by later provider specs:
  - `BridgeRuntimeIdentity`
  - `BridgeCompatibilityProbe`
  - `BridgeSessionPosture`
  - `BridgeTokenChallenge`
  - typed bridge RPC error-mapping objects
- Shared bridge objects must include explicit `schema_id` and `schema_version` and fail closed on unknown versions/fields.
- `BridgeRuntimeIdentity` must make runtime identity auditable, including provider/runtime identifiers, version/build identity, transport kind, and stable feature-set descriptors.
- `BridgeCompatibilityProbe` must record required features, tested-range posture, probe inputs, and pass/fail status rather than a freeform compatibility note.
- `BridgeSessionPosture` must make LLM-only posture explicit, including workspace-isolation posture, sandbox posture, and any approved persistence mode.
- `BridgeTokenChallenge` must model refresh-needed or re-auth-needed signals without exposing raw credential material.
- Provider-specific RPC methods and wire details stay in provider specs; this spec owns the shared object families, invariants, and versioning rules.

Parallelization: define these object families before provider-specific bridge implementations so OpenAI/Copilot-style integrations do not diverge.

## Task 3: Compatibility + Version Policy

- Define a shared compatibility policy for user-installed vendor runtimes:
  - maintain a tested version range
  - probe required methods/schema shapes/features before first use
  - allow newer untested versions only when the probe passes
  - fail closed with a clear remediation when the probe fails
- Untested-but-probe-passing runtimes require explicit user acknowledgment before use.
- Audit metadata must record runtime identity/version, tested-range posture, probe result, and whether the session required an explicit untested-version acknowledgment.

Parallelization: can be designed in parallel with provider runtime selection work; it depends on stable audit and TUI posture surfaces.

## Task 4: LLM-Only Runtime + Secret Delivery Invariants

- Bridge runtimes run only behind `model-gateway` and must remain in explicit LLM-only mode:
  - no workspace mounts or direct workspace filesystem access
  - empty/scratch `cwd`
  - isolated `HOME`/tool directories limited to an allowlisted sandbox path
  - deny command execution, patch application, and arbitrary read/write requests as policy violations
- Prefer spawned child-process IPC or stdio over listening ports.
- Secret/token delivery requirements:
  - no environment-variable secret injection
  - no raw secret logging
  - token delivery uses `secretsd` leases, file descriptors/stdin, or a secretsd-managed encrypted directory when runtime persistence is unavoidable
- Session posture rules:
  - default to ephemeral sessions
  - any persisted bridge state requires explicit manifest + policy opt-in and an auditable posture record

Parallelization: can be implemented in parallel with auth-gateway and provider work once the shared bridge posture contract is fixed.

## Task 5: Audit, Policy, and TUI Integration

- Treat bridge attempts to execute commands, apply patches, or access workspace files as policy violations with stable reason codes.
- Audit events for bridge sessions must record the active runtime identity, compatibility posture, session posture, and any token-refresh challenge events without leaking secret values.
- TUI/runtime surfaces must make these follow-on postures visible:
  - untested-but-probe-passing runtime use
  - LLM-only posture
  - persisted-vs-ephemeral bridge session posture

Parallelization: can be implemented in parallel across policy, audit, and TUI once the shared object families are fixed.

## Task 6: Fixtures + Provider Integration

- Add checked-in fixtures for:
  - runtime identity objects
  - compatibility probe pass/fail cases
  - session posture objects
  - token challenge/refresh-needed signals
  - bridge RPC error-mapping objects
- Provider specs consume these shared bridge object families and add only provider-specific RPC details.

Parallelization: fixtures can be developed in parallel with provider adapters so long as they validate against the same shared schemas.

## Acceptance Criteria

- Later provider bridge specs share the same typed bridge/runtime object families instead of inventing ad-hoc control messages.
- Untested vendor runtime use is explicit, auditable, and blocked unless the shared compatibility probe policy permits it.
- Bridge runtimes remain in explicit LLM-only posture and cannot gain workspace/file/patch capabilities through provider-specific behavior.
- Secret delivery avoids env vars and raw secret logging.
