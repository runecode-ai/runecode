# RuneCode (Project Idea)

RuneCode is a security-first, Qubes-inspired (“Qubes-ish”) agentic automation platform for software engineering.

The core premise is simple: assume any single AI/agent component can be compromised or behave maliciously, then design the system so compromise has a tiny blast radius and is always auditable.

This document captures the project idea and architecture direction based on our discussion.

## Name and Positioning

Project name decision:

- Name: RuneCode
- Brand promise: isolation + cryptographic provenance are co-equal first-class concepts.
- Open-source first: all features remain fully open-source; a future managed/enterprise offering (hosting/support/convenience) may exist without feature gatekeeping.

Tagline guidance:

- The tagline should appear everywhere (README, website, `runecode --help`) and explicitly state both pillars.
- Suggested taglines:
  - "RuneCode — Isolated execution with cryptographic provenance"
  - "RuneCode — Sandboxed AI engineering, signed and auditable"
  - "RuneCode — Spec-first workflows, isolated and attested"

Product vocabulary and commands should reinforce the same message (the tool should teach the model and the user):

- Nouns: isolates, isolate manifests, gateways, artifacts, attestations, provenance, verify, audit.
- Commands (examples): `runecode audit`, `runecode attest`, `runecode verify`, `runecode isolate ls`.

## Primary Goals

- Security is the primary product feature, not an add-on.
- Deny-by-default capabilities everywhere.
- No single component (agent, orchestrator, worker) can simultaneously:
  - access the web,
  - read or write the repo/workspace,
  - run arbitrary shell commands,
  - and hold third-party API/auth keys.
- Strong isolation between “isolates” (isolated agent/task execution environments).
- Secure, structured, auditable communication between isolates.
- Layered workflows: each layer decomposes an objective into smaller tasks, either implementing directly or orchestrating sub-isolates.
- Cross-platform single security model from day one: Linux, macOS, Windows (native Windows preferred; WSL2 as fallback if needed).
- Same capability model and workflow semantics regardless of isolation backend (microVM preferred; hardened containers as fallback), while being explicit about the assurance level of the active backend.
- Scale performance by adding hardware: minimal hardware runs serially with minimal isolate slots; more CPU/RAM allows more concurrency.
- TUI-first local dashboard (separate client) for approvals, audit, and run control.

## Non-Goals (Initial)

- Self-hosted LLM inference is not required.
- RuneCode should not require users to install Qubes OS.
- “Prove the LLM is correct” is not realistic for arbitrary code generation; instead, enforce deterministic guardrails and verifiable provenance.

## Security Model (Qubes-ish)

RuneCode borrows the architectural philosophy of Qubes OS:

- Treat the host as a tiny, trusted launcher (“dom0-like”) that does not do risky work.
- Run risky work in strongly isolated “domains” (RuneCode isolates).
- Force explicit data movement between domains via controlled channels and artifacts.
- Policy is central: what can run, with what capabilities, is explicit and reviewable.

### Core Design Principles

- Default-deny capabilities:
  - network egress,
  - filesystem read,
  - filesystem write,
  - shell execution,
  - secrets access,
  - git push/PR creation,
  - third-party APIs.
- No lateral movement: isolates do not talk to each other directly.
- Star topology: isolates communicate only with a broker/orchestrator.
- LLM outputs are untrusted input:
  - cross-boundary communication is structured (schemas), validated, logged.
  - no freeform text directly triggers privileged actions.
- All actions are auditable:
  - model prompts/responses,
  - tool calls,
  - diffs,
  - test runs,
  - approvals and user decisions.

## Isolation Strategy (Single Mode, Pluggable Backend)

RuneCode should feel like a single product mode to the user.

- The capability model and workflow semantics are consistent across platforms.
- The isolation backend is not security-equivalent across backends; RuneCode surfaces the active backend and its assurance level prominently and records it in the audit log.
- The isolation backend is an implementation detail selected at runtime based on host support.
- The chosen backend is recorded in the audit log and surfaced in the UI.

### Assurance Levels

RuneCode keeps one user-facing workflow model, but the underlying isolation backend impacts the strength of the boundary.

- MicroVM backend: high assurance isolation boundary.
- Container backend: best-effort hardening with a weaker boundary (shared kernel). RuneCode must require explicit user acknowledgment when running with the container backend, and should default to stricter policies (shorter secret lease TTLs, more frequent user-presence approvals) for higher-risk actions (eg, enabling third-party models, pushing to remotes).

### Baseline Approach

Prefer hypervisor-backed microVMs for isolates.

- The security boundary is the VM boundary (not merely tool permissions).
- Containers may still be used inside an isolate for packaging/reproducibility.

Fallback: hardened containers when microVMs are unavailable.

- Fallback aims to preserve the same capability separation and artifact-based data movement.
- The security boundary is weaker than microVMs (shared kernel), but should still be “decently secure” via aggressive hardening.

### Cross-Platform Isolation Backends

MicroVM backend (preferred): QEMU microVMs with hardware acceleration when available.

- Linux: KVM
- macOS: HVF (option to adopt Virtualization.framework later for better UX/perf)
- Windows: WHPX/Hyper-V

Each isolate runs a minimal Linux guest image tailored to its role/capability profile.

Container backend (fallback): strongly hardened containers with deny-by-default defaults.

Hardening targets (illustrative):

- rootless where possible; user namespaces
- seccomp + capability dropping (no `SYS_ADMIN`, no privileged)
- AppArmor/SELinux profiles where available
- read-only root filesystem; ephemeral writable layers
- no Docker socket; no host PID/IPC namespaces
- explicit egress policy (provider domains only when allowed)

Windows note: WSL2 can be used as an operational fallback path if native hypervisor integrations are impractical on a given machine.

### No Host Filesystem Mounts

RuneCode must not mount the host filesystem into isolates.

Instead, RuneCode uses explicit artifact movement:

- Per-job encrypted virtual disk images for workspace state.
- Content-addressed artifacts (patches, logs, plans) for handoffs.
- Controlled RPC channels (vsock/virtio-serial/mtls-over-local) rather than shared host mounts.

This is intentionally “Qubes-ish”: explicit copy/move of data, not implicit shared access.

## Capability Model (No Escalation-in-Place)

RuneCode uses an immutable, signed “capability manifest” per run/stage.

- A run is created with an explicit manifest that defines:
  - which isolate roles may be instantiated,
  - which privileges each role gets (network on/off, mounts, command allowlists, broker APIs),
  - which secrets scopes may be issued (if any),
  - which git targets are allowed (repo identities, remotes, branch patterns, and operations),
  - which third-party model endpoints (if any) may be called.

Manifests are part of RuneCode's cryptographic root of trust:

- Every step request, checkpoint, and audit event binds to the manifest hash.
- A run cannot be resumed against a different manifest.
- Any operation that changes the effective security posture (eg, enabling third-party models, enabling git push) requires a new stage with a new manifest.

### “Minting” Credentials (Safely)

RuneCode may mint short-lived, tightly scoped credentials **only as implied by the signed manifest**.

- No “upgrade” of a running isolate.
- If a new capability is needed:
  - hard stop and create a new stage/run with a new manifest.

This preserves deny-by-default while still avoiding long-lived broad tokens.

### Key Management and Root of Trust (Day One MVP)

RuneCode's signing model must be explicit and hardware-backed where possible.

Key hierarchy (recommended baseline):

- Offline root key (optional for MVP but supported): signs per-machine keys.
- Per-machine key: signs capability manifests and role manifest revisions.
- Per-isolate identity keys: isolates sign their own audit events (core must not be able to forge isolate-attributed events).
- Image signing key: signs microVM/container images and toolchain bundles (separate from manifest signing).

Storage and protection:

- Prefer TPM/Secure Enclave-backed keys where available.
- Support hardware security keys (FIDO2/YubiKey) for explicit user approvals ("tap to approve manifest").
- Manifest signing should require explicit user presence (hardware key tap or equivalent OS-secured confirmation). `runecode-launcher` must not be able to silently mint new signed manifests.
- If hardware-backed keys are unavailable:
  - use OS keychain/keyring where possible
  - record the degraded posture (key protection level) in the audit log and UI

Rotation and revocation:

- Keys have explicit rotation procedures and revocation lists.
- Revocation is treated as a security event and is logged and (when configured) externally anchored.

### Git Target Allowlist and Credential Scoping (Day One MVP)

RuneCode treats git operations as high-risk egress.

- Run/stage manifests include an explicit allowlist of git targets:
  - canonical repo identity (recommended): `{provider, owner, repo}`
  - allowed remotes (eg `origin` only)
  - allowed branch patterns (eg `bot/*` only)
  - allowed operations (eg create PR, push commits)
- URLs are not trusted policy inputs.
  - git URLs are canonicalized into a repo identity.
  - policy checks use the canonical identity to prevent URL tricks (https vs ssh variants).

Enforcement is not "git-gateway self-checks a signature".
Enforcement is defense-in-depth and happens outside the potentially compromised gateway:

- `runecode-secretsd` issues repo-scoped, operation-scoped credentials that only work for the allowlisted repo(s).
  - Prefer GitHub App installation tokens / fine-grained repo-scoped tokens over broad PATs.
  - Deny non-essential surfaces (eg gists) by default.
- `git-gateway` egress is allowlisted to the provider domains.
- `runecode-launcher` requires an outbound-verification artifact before permitting push:
  - the pushed diff/tree hash must match the signed patch artifact.
- All git target decisions and token leases are audited (without logging raw secrets).

### Emergency Stop and Revocation (Day One MVP)

RuneCode includes a defined emergency response mechanism:

- Emergency stop: a privileged local action that immediately terminates all isolates and blocks new isolate launches.
- Credential revocation: short-lived tokens plus an explicit revocation list checked by gateways ("deny even if TTL not expired").
- Workspace sealing: mark affected workspace disks/snapshots as quarantined and require explicit user action to re-open; optionally re-key encryption on next open.
- Audit: append a tamper-evident "emergency stop" event and anchor it.
- Recovery: resuming work requires a new stage manifest and re-validation of gates.

## Isolate Roles (Capability Profiles)

RuneCode treats “isolate” as a role with a narrowly-scoped ability set.

Example roles (illustrative; final set is configurable):

- `model-gateway`:
  - network egress allowed (strict allowlist)
  - holds provider API keys
  - no repo/workspace access
  - provides completions to other isolates over a constrained local channel
  - outputs structured responses only; logs and redacts at boundaries
- `web-research`:
  - network egress allowed (allowlist-only; domain sets/patterns are explicit per run/stage)
  - no repo/workspace access, no secrets
  - may only receive low-sensitivity artifacts (eg `spec_text`, `web_query`), never workspace/code-derived data classes
  - supports OpenCode-style crawling within the configured allowlist policy (MVP uses suffix-wildcard patterns; future supports curated bundles and domain-expansion via new stage manifests)
  - outputs citations/notes only (eg `web_citations`)
- `deps-fetch` (optional, for offline builds):
  - network egress allowed (strict allowlist to package registries)
  - no repo/workspace access; reads only lockfile artifacts
  - outputs a read-only dependency cache artifact for workspace isolates
- `workspace-read`:
  - access to workspace disk read-only
  - no network
  - outputs code maps/snippets/summaries
- `workspace-edit`:
  - access to workspace disk read/write
  - no network
  - runs OpenCode build/plan flows with strict tool allowlists
  - applies patch artifacts and produces diffs
- `workspace-test`:
  - runs on a disposable workspace snapshot
  - source tree is read-only; writes only to scratch
  - no network by default
  - uses purpose-built executors (no shell passthrough)
  - outputs signed logs and results
- `workspace-lsp`:
  - source tree is read-only; writes only to ephemeral index/scratch
  - no network
  - runs language servers / analyzers and emits diagnostics artifacts
- `git-gateway`:
  - network egress allowed
  - holds git/Jira/Slack creds
  - no direct access to `workspace-*` disk contents
  - consumes signed patch artifacts and produces PRs/updates
  - must not materialize full repo history by default; prefer sparse/partial operations limited to files touched by the patch, or push pre-constructed objects
  - must verify outbound changes match the signed patch artifact before pushing
  - operates only on git targets allowlisted by the run/stage manifest (repo identity + branch patterns), enforced via repo-scoped credentials
  - credentials should be fine-grained and tightly scoped (repo-scoped and branch-scoped where supported)
- `secrets-broker` (mandatory component; separate trust domain):
  - only component that stores long-lived secrets at rest
  - issues short-lived, scope-bound tokens/leases only as allowed by the manifest
  - every secret read/lease is an audit event
  - supports revocation and rotation; does not share key material with the launcher/PEP

Key rule: no single role combines public network egress + workspace disk access (especially RW) + long-lived secrets. Gateway roles that must interact with remotes are constrained to minimal, purpose-built operations and are heavily audited.

## “Many Roles, Few Running Isolates” (Minimal Hardware Mode)

RuneCode should support a minimal slot count (example: 1–2 concurrent running microVMs), while still enforcing many roles.

- Roles remain distinct; execution is serialized.
- “Active but not running” means state is persisted, not kept resident:
  - Workflow state: persisted in the orchestrator DB (durable execution).
  - Workspace state: persisted as encrypted virtual disk image snapshots.
  - Outputs: persisted as signed artifacts (patches, logs, decisions) in an artifact store.
- VMs are terminated between steps rather than suspended, reducing RAM snapshot leakage and aligning with Qubes-ish hygiene.

Performance tradeoff:

- Increased overhead from microVM boot + artifact handoff.
- Acceptable in exchange for security; scaling is achieved by increasing `max_running_isolates` with more hardware.

## Secure Communication Between Isolates

RuneCode uses broker-mediated, structured communication.

- Isolates do not open inbound ports to each other.
- Communication patterns:
  - RPC via vsock/virtio-serial where possible (microVM local transport)
  - or mutual TLS with per-isolate identities (SPIFFE-like) on a constrained internal network
- Messages are schema-validated (JSON schema or protobuf).
- The broker enforces defensive limits:
  - per-isolate request/byte rate limits
  - maximum message sizes
  - maximum in-flight requests
  - backpressure so one isolate cannot DoS others
  - ability to unilaterally terminate a misbehaving isolate
- Artifact handoffs are hash-addressed:
  - isolate uploads artifact to store via one-time, write-only URL
  - isolate reports `{artifact_hash, metadata}` to broker
  - broker verifies hash and routes to next role

## Resource Quotas and Execution Limits

RuneCode enforces resource limits as part of the security model (DoS resistance and blast-radius control).

- Role manifests define quotas:
  - CPU shares/limits
  - memory limits (hard OOM kill)
  - disk quotas (workspace and artifact output)
  - network bandwidth limits (gateway roles)
  - maximum execution time per step
  - maximum artifact count and total artifact size per step
- Enforcement occurs at the isolation backend level (cgroups/seccomp for containers; QEMU limits and host controls for microVMs).
- `runecode-launcher` includes a watchdog that terminates isolates that exceed quotas, and records quota violations in the audit log.

## Image and Toolchain Supply Chain Security

RuneCode assumes the isolate image supply chain is a critical part of the threat model.

- Images and toolchain bundles are built via a reproducible build pipeline.
- Images are signed with a key separate from manifest signing keys.
- Expected image digests/signatures are pinned in role manifests.
- `runecode-launcher` refuses to start isolates if the image signature/digest does not match the pinned expectation.
- Images should be minimal and role-specific (reduce attack surface).
- Update strategy is explicit:
  - security updates can be applied to images without silently changing role capabilities
  - image updates are recorded in the audit log (old digest, new digest, signer identity)

## Encrypted Workspace Disks and Key Lifecycle

Workspace state lives in encrypted virtual disks (microVM) or encrypted storage volumes (container backend).

Baseline encryption design (day one):

- MicroVM disks: LUKS2/dm-crypt (or equivalent) with modern AEAD-capable cipher suites where supported.
- Per-workspace encryption keys are distinct and never reused across unrelated workspaces.
- Keys are derived/sealed using a KDF (eg, Argon2) and hardware-backed sealing (TPM/Secure Enclave) when available.
- `runecode-launcher` should not permanently store raw disk keys; keys are provisioned per run/stage and can be re-sealed or quarantined.

Lifecycle:

- Creation: new workspace disk gets a new key.
- Rotation: keys can be rotated when a workspace is sealed/quarantined or on schedule.
- Destruction: secure deletion is supported (destroy key material; wipe metadata) and recorded as an audit event.
- Retention: define retention policies for workspace snapshots and artifacts (default: least retention needed).

## Data Classes and Inter-Role Flow Policy

RuneCode's information-flow controls depend on explicit data classification.

Data classes are a first-class taxonomy that is:

- Explicit: every artifact has exactly one `data_class`.
- Immutable: `data_class` cannot be changed after creation.
- Validated: producing roles may only emit classes allowed by their role manifest; the broker validates the class assignment.
- Fail-closed: unknown/ambiguous artifacts are classified as the most restrictive class.

Compound artifacts:

- Prefer decomposition: emit separate artifacts per data class.
- If not possible, use the most restrictive class.

Inter-role flow policy:

- Run/stage manifests define an explicit flow matrix: which roles may consume which data classes.
- Example: workspace/code-derived classes must not flow into `web-research`.
- Third-party model egress is enforced at `model-gateway` and is driven by this flow matrix.

Minimum starter taxonomy (illustrative):

- `spec_text`: approved spec content
- `web_query`: low-sensitivity research queries
- `web_citations`: URLs + quoted excerpts
- `approved_file_excerpts`: explicitly approved code excerpts
- `diffs`: patch/diff artifacts
- `build_logs`: build/test outputs (may be sensitive)
- `audit_events`: structured audit records

## Web Research Egress Policy (OpenCode-Style Crawling)

RuneCode supports OpenCode-style web research that can follow discovered URLs, while still keeping network access deny-by-default.

RuneCode should support three crawling experiences (same underlying egress policy model, different UX/workflow):

1) Tight patterns (MVP; balanced)
   - The run/stage manifest provides an allowlist consisting of:
     - explicit domains (eg `docs.python.org`)
     - domain patterns (eg `*.readthedocs.io`, `*.github.io`)
   - Domain pattern syntax is intentionally constrained:
     - suffix wildcard only (eg `*.example.com`)
     - no regex, no CIDR, no arbitrary globbing
     - patterns are canonicalized to eTLD+1 boundaries where possible
   - `web-research` may crawl only within this allowlist.
   - Any discovered URL outside policy is blocked and reported.

2) Curated allowlist bundles (later)
   - The user selects one or more named bundles at run start (eg `language_docs`, `common_oss_hosts`, `cloud_provider_docs`).
   - Bundles expand into domains/patterns in the run/stage manifest (still explicit, auditable).

3) Domain-expansion workflow (later; strictest)
   - `web-research` performs a results-only pass first (collect URLs/titles/snippets).
   - RuneCode proposes a domain list to follow.
   - Following new domains requires a new stage manifest (no escalation-in-place).

Implementation requirement:

- The egress policy data model must support explicit domains, patterns, and bundle expansion, so MVP does not paint the system into a corner.
- All outbound requests are logged (role + manifest hash + destination + bytes) and are subject to quotas.

## Failure Handling and Recovery

RuneCode must fail closed and recover deterministically.

Day one behaviors:

- Isolate boot failure: retry with backoff; surface a clear error; optionally fall back to container backend with explicit user acknowledgment.
- Artifact corruption: verify hashes on every transfer; never accept corrupted artifacts; re-request or fail the step.
- Disk full/corruption: seal the affected audit/workspace segment, start a new segment, alert the user, and record the event.
- Step timeouts: terminate isolate, record the timeout event, and require an explicit retry.
- Core crash recovery:
  - persist a minimal WAL/journal for `runecode-launcher` and `runecode-auditd`
  - detect and terminate orphaned isolates on restart
  - resume only from checkpoints bound to manifest hash and verified gate artifacts

## Project Architecture (Dom0-like Core + Clients)

RuneCode is structured as a local “control plane” plus multiple local clients.

The control plane is intentionally decomposed for defense-in-depth. Rather than one “god daemon,” RuneCode runs a small set of cooperating processes with least privilege. These processes may be shipped together as `runecode-core`, but they must be separable and independently hardenable.

Control plane components (security kernel):

- Each component should run under a distinct OS identity (separate user/service account) with minimal filesystem/network permissions, so compromise of one does not trivially imply compromise of all.

- `runecode-launcher` (Go, local daemon; PEP)
  - starts/stops isolates (microVM/container backend)
  - attaches encrypted workspace disks and read-only artifacts
  - enforces role manifests + run/stage capability manifests + step-ordering invariants
  - has no long-lived provider/git credentials
  - has no public network egress by default
- `runecode-broker` (Go, local daemon)
  - mediates isolate RPC and artifact routing
  - enforces rate limits, message size caps, and backpressure
  - holds no long-lived secrets and no workspace disks
- `runecode-secretsd` (Go, local daemon or dedicated isolate)
  - stores long-lived secrets at rest (hardware-backed where possible)
  - can integrate with established secret stores (OS keychain, TPM-sealed blobs, Vault) rather than inventing a new secrets format
  - issues short-lived, scope-bound leases only as allowed by manifests
  - enforces git target allowlists by minting repo-scoped, operation-scoped credentials (tokens are bound to allowed repos/operations)
  - supports explicit revocation/rotation
  - every secret read/lease is an audit event
- `runecode-auditd` (Go, local daemon or dedicated isolate)
  - append-only audit log writer
  - writes to append-only storage (WORM/append-only where available) under a separate OS identity
  - verifies isolate signatures and enforces audit schema
  - anchors audit roots externally when configured
  - `runecode-launcher` cannot forge isolate events and cannot rewrite history

Workflow runner and clients:

- `runecode-graph` (TypeScript/Node): LangGraph workflow runner.
  - treated as an untrusted scheduler (supply chain risk assumed)
  - runs with no secrets and no workspace access (ideally inside its own isolate)
  - dependencies are pinned/locked and kept minimal; build artifacts are reproducible and signed where possible
  - requests steps; `runecode-launcher` independently validates that requested steps are legal, ordered, and gate-complete
- `runecode-tui` (Go): TUI-first dashboard client.
  - pages for Runs, Approvals, Diffs, Artifacts, Audit Timeline, Logs, Settings
  - approves manifests/opt-ins, resumes runs, and triggers actions via the local API
- `runecode-web` (optional later): local-only web client for richer diff/audit views.
  - same API, no new privileges

Local API authentication (day one baseline):

- Local IPC only (Unix socket / named pipe). Socket permissions are restrictive (single-user by default).
- Authenticate clients using OS peer credentials where available.
- High-risk approvals require explicit user presence (passphrase confirmation and/or hardware key tap where available).

## Tech Stack (Recommended)

- Go: `runecode-core` and `runecode-tui`.
  - keeps the trusted computing base (TCB) smaller
  - good cross-platform daemon + local IPC story
  - good fit for isolation orchestration, encrypted storage, and streaming logs
- TypeScript: `runecode-graph`.
  - best fit for LangGraph JS and programmatic integration with OpenCode SDK

This split keeps the “security kernel” small and verifiable, while keeping workflow logic flexible.

## GitHub Organization and Repository Strategy

Recommended approach: start with a monorepo and publish multiple packages/artifacts from it.

Rationale:

- Early development requires frequent cross-cutting changes across the security kernel, role manifests, schemas, and workflow runner.
- A monorepo makes it easier to evolve boundaries safely (and to update policy + invariants everywhere consistently).
- Trust boundaries are enforced by process isolation, manifests, and sandboxing, not by repo boundaries.

### Repos to Create Under `runecode-ai`

- `runecode` (already)
  - Monorepo containing `runecode-launcher`, `runecode-broker`, `runecode-secretsd`, `runecode-auditd`, `runecode-tui`, and `runecode-graph` as separate packages/commands.
  - Produces multiple release artifacts (binaries, images, role manifests, schema bundles) from one source of truth.
- `.github`
  - Organization-level community health files (issue templates, PR templates, security policy, default labels).

Optional (only if you want independent docs deploy cadence):

- `runecode-docs` or `runecode-site`

### Monorepo Packaging Guidance

- Build and release each component independently (separate binaries) while keeping a unified source tree.
- Keep the workflow scheduler (`runecode-graph`) isolated at runtime and treated as untrusted, even if it lives in the same repo.
- Keep dependency boundaries explicit:
  - Go security kernel packages should not import Node/JS code.
  - Node dependencies are pinned/locked and kept minimal.
- Publish container/microVM images to `ghcr.io/runecode-ai/*` and pin expected digests in role manifests.

### When to Split Into Separate Repos (Later)

Split only when a stable interface exists and the split reduces real risk or maintenance burden:

- `runecode-graph` (separate repo) if you want Node/LangGraph supply chain to be physically separate and governed by different review rules.
- `runecode-protocol` if third parties integrate and you need stable, versioned schemas/protocols.
- `runecode-images` if image build definitions and signing pipelines become large and need tighter change control.

## CI (GitHub Actions)

RuneCode must continuously test across Linux, macOS, and Windows.

- Use separate GitHub Actions workflows per OS for unit tests.
- Use separate GitHub Actions workflows per OS for integration tests.
- Keep unit tests runnable on standard GitHub-hosted runners.
- Integration tests may be split into:
  - backend-agnostic tests that can run on GitHub-hosted runners (eg, manifest validation, policy invariants, artifact routing, audit verification)
  - backend-specific tests that may require self-hosted runners (eg, KVM microVM integration on Linux)
- CI outputs (test reports, logs) should be captured as artifacts where helpful and referenced in the audit/provenance story.

## Auditability (Everything, With Redaction)

RuneCode must produce a complete, tamper-evident audit trail.

### What Is Logged

- Every LLM call:
  - provider/model ID
  - prompt inputs (with policy-based redaction)
  - responses
- Every tool/action request and result:
  - command allowlist decisions
  - file patch application results
  - test/lint outputs
  - PR creation metadata
- Every user decision:
  - approvals/rejections
  - manifest sign-off
  - explicit opt-ins (egress, third-party models)

### Tamper Evidence

- Append-only event log per run, written by `runecode-auditd` (independent writer).
- Hash chain: each event includes the previous event hash.
- Isolates sign their own events with per-isolate identity keys before submission.
  - `runecode-launcher` must not be able to forge isolate-attributed events.
- Clients (TUI/web) independently verify the hash chain on read; they do not trust the writer's assertion.
- Periodically anchor audit roots externally when configured (eg, TPM PCR extension, RFC 3161 timestamping, or a lightweight witness service).

### Secret Scrubbing / Redaction

- Redact at boundaries (prevent secrets from leaving a trust zone), not only after.
- Maintain two views:
  - Encrypted forensic view (restricted access)
  - Redacted operational view (default)
- Explicitly block known secret sources from entering third-party egress paths.

## Third-Party Models (Explicit Opt-In)

Default: deny.

RuneCode should treat all third-party model usage as an explicit opt-in capability in the manifest.

- If opted-in, all provider traffic occurs only via `model-gateway` (or a pool of gateways).
- No repo/workspace access is granted to `model-gateway`.
- Worker isolates receive a permission like `llm_access`, defined as “may call the local model-gateway”, not “may reach the public internet”.
- The user explicitly approves which data classes may be sent (eg, "spec text only", "diffs only", "limited file excerpts").

## Workflow Direction (LangGraph + AgentOS + OpenCode)

RuneCode can use LangGraph, AgentOS, and OpenCode together without breaking the security model.

- LangGraph: durable orchestration (pause/resume, checkpoints, retries, branching/time-travel).
- AgentOS: specs/standards as the git-native system of record (`agent-os/`).
- OpenCode: controlled execution engine running inside role isolates.

Security positioning:

- OpenCode’s internal permission system is a guardrail, not the boundary.
- The boundary is the isolate capability profile (microVM preferred; hardened containers fallback), plus explicit artifact movement.
- `runecode-launcher` (part of `runecode-core`) is the policy enforcement point (PEP); LangGraph schedules but does not gain powers.

Checkpoint and replay safety (day one baseline):

- All checkpoints are bound to the manifest hash and role manifest versions.
- A checkpoint cannot be resumed if the manifest hash differs.
- Gate results are non-replayable: once a gate fails, the workflow cannot be rolled back past that point without re-running and re-passing the gate.
- Branching/time-travel/replay operations are privileged actions that require the same approval level as manifest approval.
- All checkpoint operations are audit events.

### AgentOS: Canonical Truth

- All long-lived product context, standards, and specs live as markdown files committed to git.
- Any proposed change to this truth is always a patch artifact and a PR (never a silent mutation).

### LangGraph Shared Memory: Ephemeral Cache + Working Set

Use LangGraph “shared memory” as a rebuildable, ephemeral accelerator:

- Parse/index AgentOS docs into structured objects.
- Cache “relevant standards” selection results.
- Cache code maps/summaries keyed by `(repo, commitSHA)`.
- Store decision packets, approvals, and run metadata.

Syncing back to AgentOS is explicit and auditable:

- If shared memory implies an update to `agent-os/**` (new standard, clarified spec text, new references), RuneCode generates a proposed patch against `agent-os/**`.
- The patch is applied inside a workspace isolate and opened as a PR via `git-gateway`.
- After merge, caches are invalidated and rebuilt from the new commit SHA.

### OpenCode Placement (Qubes-ish)

- OpenCode runs inside offline workspace-capable isolates (`workspace-read`/`workspace-edit`/`workspace-test`/`workspace-lsp`).
- Those isolates do not talk to the public internet.
- When an OpenCode session needs LLM completions, it uses `llm_access` to call `model-gateway` over a constrained local channel.
- `model-gateway` owns provider traffic and provider keys; it enforces data-class policy and boundary redaction.
- For performance, RuneCode can run a pool of `model-gateway` isolates and load-balance requests.

### Spec-First Workflow

- Generate spec artifacts in a dedicated branch/workspace.
- Create a spec PR for review.
- Implementation begins only from an approved spec (commit hash and spec folder path are frozen inputs).

### Layered Agent/Task Execution

- Higher-level orchestrator decomposes objectives into tasks.
- Each task executes in a role-specific isolate.
- Review and security checks run as separate, more constrained roles.

## How Common Actions Are Performed (Under This Security Model)

All actions are executed inside role isolates with explicit artifact handoffs; the host and `runecode-graph` do not directly operate on user workspaces.

- Web research
  - Runs in `web-research` (egress allowlist only; domains/patterns are explicit per run/stage; no repo/workspace; no secrets by default).
  - Outputs: citations/notes artifacts; optional structured summaries via `llm_access` to `model-gateway`.
  - Inputs are restricted by data class (eg `spec_text`, `web_query` only).
  - Crawling behavior is policy-driven:
    - MVP (balanced): tight domain patterns (eg `*.readthedocs.io`, `*.github.io`) plus explicit domains.
    - Later: curated allowlist bundles (named sets of domains/patterns selected at run start).
    - Later: domain-expansion workflow (results-only pass, then explicit user approval of new domains via a new stage manifest).
- File reads (repo)
  - Runs in `workspace-read` (workspace disk read-only; no network).
  - Outputs: excerpt/symbol-map artifacts gated by “data class” policy; always hashed/audited.
- File writes / code creates / code edits
  - Runs in `workspace-edit` (workspace disk read/write; no network).
  - Inputs: patch artifacts and/or structured edit instructions.
  - Outputs: diffs + updated workspace snapshots; all file changes are auditable as patch hashes and resulting tree hashes.
- Code execution (builds/tests)
  - Runs in `workspace-test` (no network by default; runs on a disposable snapshot; source tree read-only; writes only to scratch).
  - Uses purpose-built executors (no shell passthrough) and records tree hashes before/after.
  - Outputs: logs, junit/coverage artifacts; deterministic pass/fail gates.
- Dependency fetching (when offline execution needs it)
  - Prefer: prebuilt toolchains and caches baked into base images.
  - Optional: `deps-fetch` role fetches dependencies from allowlisted registries using only lockfile artifacts and emits a read-only cache artifact.
  - Workspace roles consume the cache read-only without gaining internet access.
- LSP and diagnostics (typos/errors)
  - Runs in `workspace-lsp` (no network; source tree read-only; writes only to ephemeral index/scratch).
  - Outputs: machine-readable diagnostics artifacts; can be displayed in the TUI.
- Git commits
  - Runs in `git-gateway` (egress allowed; holds git credentials; no workspace disk access).
  - Consumes a signed patch artifact (and metadata), applies it to a shallow/sparse/partial checkout by default, then commits.
  - Verifies the outbound diff/tree hash matches the signed patch artifact before `git push`.
  - Credentials are issued by `runecode-secretsd` and are scoped to the allowlisted repo identity and operations.
- PR creation
  - Runs in `git-gateway` using provider APIs; produces PR URL + metadata artifacts.
  - PR body can be generated from structured run artifacts (spec, diffs, test results) and optionally summarized via `model-gateway`.
  - PR targets are constrained by the same allowlisted repo identity and token scopes.
- OpenCode skills
  - Skills are treated as signed, versioned, read-only bundles (baked into images or attached as RO artifacts).
  - Manifests control which roles may load which skills; every skill invocation is logged.

## Role Manifests (Sketch)

RuneCode uses two levels of policy documents:

1) Role manifests (capability profiles): define what a role is allowed to do.
2) Run/stage manifests (instances): define which roles are instantiated for a specific run, what data classes are allowed, and what explicit opt-ins are enabled.

The role manifests are the same on all machines. Hardware only changes how many instances can run concurrently.

### Minimal Role Manifests (Capability Profiles)

The following sketches show the intent. Exact schemas are TBD.

```yaml
# role: model-gateway
id: role/model-gateway
isolation: { backend: microvm_or_container }
image:
  signature_required: true
  expected_digest: sha256:TBD
network:
  egress:
    mode: allowlist
    domains:
      - api.openai.com
      - api.anthropic.com
      - generativelanguage.googleapis.com
resources:
  cpu_max_vcpu: 2
  memory_max_mb: 2048
  disk_max_mb: 1024
  network_max_mbps: 10
  step_timeout_sec: 300
  artifacts_max_total_mb: 256
secrets:
  allow_scopes:
    - llm_provider_api_key
filesystem:
  workspace: none
  host_mounts: none
tools:
  provide:
    - llm_proxy_rpc
policy:
  redact_at_boundary: true
  accepted_data_classes:
    - spec_text
    - diffs
    - approved_file_excerpts

---
# role: workspace-edit
id: role/workspace-edit
isolation: { backend: microvm_or_container }
image:
  signature_required: true
  expected_digest: sha256:TBD
network:
  egress: denied
resources:
  cpu_max_vcpu: 4
  memory_max_mb: 4096
  disk_max_mb: 8192
  step_timeout_sec: 1800
  artifacts_max_total_mb: 1024
secrets:
  allow_scopes: []
filesystem:
  workspace: rw
  host_mounts: none
tools:
  allow:
    - opencode
    - apply_patch
    - read_files_under_workspace
    - write_files_under_workspace
  llm_access:
    via: model-gateway
executors:
  - id: git_status
  - id: git_diff
  - id: git_add
executor_policy:
  construct_argv_internally: true
  prohibit_shell_metacharacters: true
  env_mode: allowlist

---
# role: workspace-test
id: role/workspace-test
isolation: { backend: microvm_or_container }
image:
  signature_required: true
  expected_digest: sha256:TBD
network:
  egress: denied
filesystem:
  workspace:
    source_tree: ro
    scratch: rw_ephemeral
  host_mounts: none
resources:
  cpu_max_vcpu: 4
  memory_max_mb: 4096
  disk_max_mb: 8192
  step_timeout_sec: 3600
  artifacts_max_total_mb: 2048
executors:
  - id: run_tests
executor_policy:
  construct_argv_internally: true
  prohibit_shell_metacharacters: true
  env_mode: allowlist
policy:
  execution_mode: snapshot_and_discard
  verify_tree_hash_unchanged: true

---
# role: workspace-lsp
id: role/workspace-lsp
isolation: { backend: microvm_or_container }
image:
  signature_required: true
  expected_digest: sha256:TBD
network:
  egress: denied
filesystem:
  workspace:
    source_tree: ro
    scratch: rw_ephemeral
  host_mounts: none
resources:
  cpu_max_vcpu: 2
  memory_max_mb: 2048
  disk_max_mb: 4096
  step_timeout_sec: 1800
tools:
  allow:
    - language_server
    - diagnostics_export

---
# role: web-research
id: role/web-research
isolation: { backend: microvm_or_container }
image:
  signature_required: true
  expected_digest: sha256:TBD
network:
  egress:
    mode: allowlist
    # MVP (balanced): allowlist includes explicit domains and tight patterns.
    # Later: allowlist can be assembled from named bundles, or expanded via a new stage manifest.
    domains: []
    # Pattern syntax: suffix wildcard only (eg "*.example.com").
    domain_pattern_syntax: suffix_wildcard_only
    domain_patterns: []
    bundles: []
resources:
  cpu_max_vcpu: 1
  memory_max_mb: 1024
  disk_max_mb: 1024
  network_max_mbps: 5
  step_timeout_sec: 600
filesystem:
  workspace: none
  host_mounts: none
secrets:
  allow_scopes: []
tools:
  llm_access:
    via: model-gateway
policy:
  accepted_data_classes:
    - spec_text
    - web_query
  output_data_classes:
    - web_citations

---
# role: git-gateway
id: role/git-gateway
isolation: { backend: microvm_or_container }
image:
  signature_required: true
  expected_digest: sha256:TBD
network:
  egress:
    mode: allowlist
    domains:
      - github.com
      - api.github.com
      - gitlab.com
      - api.atlassian.com
secrets:
  allow_scopes:
    - git_provider_token
    - jira_token
    - slack_token
filesystem:
  workspace: none
  host_mounts: none
resources:
  cpu_max_vcpu: 2
  memory_max_mb: 2048
  disk_max_mb: 8192
  network_max_mbps: 10
  step_timeout_sec: 1800
tools:
  allow:
    - fetch_repo_patch_context  # shallow/partial/sparse by default
    - apply_patch_artifact
    - verify_outbound_matches_patch
    - git_commit
    - git_push
    - create_pr
  llm_access:
    via: model-gateway
policy:
  repo_identity:
    # Policy uses canonical repo identity (not raw URLs) to prevent URL trick bypasses.
    # Concrete values come from the run/stage manifest.
    allowed: []  # [{provider, owner, repo}]
    canonicalize_urls: true
    enforce_via_secretsd: true
  repo_access:
    clone_mode: shallow_sparse_partial
    max_history_depth: 1
    materialize: patch_touched_files_only
  push_restrictions:
    allowed_remotes: [origin]
    allowed_branches: ["bot/*"]
```

### 2 Running Slots Machine (Minimal Practical Profile)

Goal: preserve the full role separation, but time-multiplex roles across two running isolate slots.

- Slot A (workspace slot): runs one workspace role at a time (`workspace-read`/`workspace-edit`/`workspace-test`/`workspace-lsp`) against the same encrypted workspace disk.
- Slot B (gateway slot): runs one networked role at a time (`model-gateway` or `git-gateway` or `web-research`).
- State is persisted between steps (graph checkpoints + encrypted workspace disk + artifacts). Isolates are terminated between steps.

Example pool sketch:

```yaml
hardware_profile: two_slots
max_running_isolates: 2
pools:
  workspace:
    warm: 1
    roles:
      - role/workspace-edit
      - role/workspace-test
      - role/workspace-lsp
  gateway:
    warm: 1
    roles:
      - role/model-gateway
      - role/git-gateway
      - role/web-research
notes:
  - gateway roles are mutually exclusive in the single gateway slot
  - git-gateway is started only for commit/push/PR phases
```

### 10 Running Slots Machine (Throughput Profile)

Goal: increase completion speed by parallelizing workspace slices and keeping gateways available.

- Multiple workspace roles run concurrently on separate workspace snapshots/worktrees.
- A small pool of `model-gateway` isolates handles concurrent completion traffic.
- A dedicated `git-gateway` isolate handles PR/commit operations.

Example pool sketch:

```yaml
hardware_profile: ten_slots
max_running_isolates: 10
pools:
  workspace_edit:
    warm: 6
    role: role/workspace-edit
  workspace_test:
    warm: 1
    role: role/workspace-test
  workspace_lsp:
    warm: 0  # on-demand, or reuse a workspace slot
    role: role/workspace-lsp
  model_gateway:
    warm: 2
    role: role/model-gateway
  git_gateway:
    warm: 1
    role: role/git-gateway
  web_research:
    warm: 0  # on-demand; can temporarily reduce workspace concurrency
    role: role/web-research
notes:
  - increase workspace_edit warm count as CPU/RAM allows
  - increase model_gateway pool if completions become the bottleneck
```

## “Prove Accuracy” Methods (Day One MVP Requirements)

RuneCode MVP includes all four methods from day one.

Important scoping note:

- These methods are applied where they are technically meaningful and provide real security/assurance value.
- In practice, this means neuro-symbolic controls and deterministic gates apply broadly, while formal methods and ZK proofs target the small “security kernel” and selected, high-value attestations.
- These methods are not a substitute for isolation. They harden and verify behavior inside the isolation model.

### 1) Neuro-Symbolic Guardrails (Day One MVP)

Neuro-symbolic guardrails mean: the LLM is only ever a proposer; deterministic systems are the only things that can authorize, execute, and advance workflow state.

#### The Control Loop: Propose -> Validate -> Authorize -> Execute -> Attest

Every action with side effects (network, filesystem writes, command execution, git operations, PRs) follows a single deterministic pipeline:

1) Propose (neural)
   - The LLM produces a structured proposal object.
   - The proposal must be schema-valid and must reference inputs only by artifact hashes.
2) Validate (symbolic: structural)
   - Schema validation (JSON schema/protobuf), size limits, quotas.
   - Reference integrity: only hash-addressed artifacts; no implicit hidden inputs.
3) Authorize (symbolic: policy)
   - Evaluate role manifest + run/stage capability manifest + current workflow state.
   - Deterministic decision: `allow | deny | require_human_approval`.
   - No escalation-in-place: if new capability is needed, the stage ends and a new manifest is required.
4) Execute (symbolic: deterministic runner)
   - A role-specific executor performs the action (patch applier, command runner, git runner, fetcher).
   - Execution happens inside the correct isolate with enforced sandbox constraints.
5) Attest + record (symbolic: evidence)
   - Produce artifacts (diffs, logs, diagnostics, test reports).
   - Hash and sign events; append to tamper-evident audit log.
   - Only advance the LangGraph state when postconditions are met.

#### The Core Objects

RuneCode uses explicit objects so that “guardrails” are enforceable and testable:

- Capability manifest (per run/stage; user-approved and signed)
  - defines which roles can run and which capabilities are enabled (including explicit opt-ins)
  - defines allowed data classes and allowed flows (eg, which artifacts may be sent to `model-gateway`)
  - defines command allowlists and risk gates
- Role manifest (capability profile)
  - defines what a role is allowed to do, independent of a specific run
- Artifact references (hash-addressed, typed)
  - all cross-role handoffs are artifacts: `{hash, data_class, origin, metadata}`
- Proposals (schema-validated)
  - the LLM can only influence the system by emitting proposals of known types (eg `PatchProposal`, `CommandProposal`, `LLMRequest`)

#### Data Flow Controls ("Taint" Without Magic)

- All artifacts have an explicit `data_class` (eg `spec_text`, `diffs`, `approved_file_excerpts`, `web_citations`, `test_logs`).
- Data flows are controlled by manifest rules, not by prompts.
- Third-party model egress is default-deny; if explicitly opted in, only approved data classes may cross the boundary.
- Workspace roles remain offline. `llm_access` means “call `model-gateway` over a constrained local channel”, not “reach the public internet”.

#### Enforcement Points

- `runecode-launcher` is the policy enforcement point (PEP).
  - it instantiates isolates, attaches disks/artifacts, and refuses any action outside policy.
- Role runners are deterministic executors.
  - they do not interpret freeform text into privileged actions.
- LangGraph schedules and checkpoints.
  - it cannot bypass `runecode-launcher`.
- OpenCode runs inside constrained workspace roles.
  - it can propose changes, but cannot exceed sandbox and policy.

#### Examples (Intended Behavior)

- A proposal to run `curl ... | bash` is denied by the command allowlist and/or network denial in workspace roles.
- A proposal to send `src/auth.ts` to a hosted model is denied unless (a) model opt-in is enabled and (b) the artifact is classified as an allowed data class.
- A proposal to `git push` is denied in workspace roles and must be routed through `git-gateway`.

### 2) Deterministic Assertion Creation and Gating (Day One MVP)

Deterministic assertion gating means RuneCode only accepts outcomes that are validated by deterministic checkers.

Key rules:

- Gates cannot be bypassed by LLM output.
- Gates run in role isolates (`workspace-test`, `workspace-lsp`, and policy checkers) with pinned toolchains.
- Gate results are artifacts and are part of the audit trail.

#### What Counts As a Gate

Day one gate set should include (repo-specific subsets are allowed, but the framework is universal):

- Build/type gates: compilation, type checking, static analysis.
- Test gates: unit/integration tests; optional smoke tests.
- Lint/format gates: lint rules and formatting checks.
- Security gates:
  - secret scanning (prevent committing/egressing secrets)
  - dependency scanning (known CVEs, license policy if desired)
  - SAST rules for common injection/auth issues
- Policy gates:
  - manifest compliance (no out-of-policy commands, no forbidden file touches)
  - artifact flow compliance (no forbidden data class egress)

#### Deterministic Assertion Creation

The LLM may propose new assertions (tests, invariants, static rules, migrations), but they are untrusted until proven by deterministic execution:

- LLM proposes: “add a regression test for bug X” -> results are accepted only if the test runs and passes.
- LLM proposes: “refactor module Y” -> accepted only if gates pass and diffs match policy.

#### Evidence and Reproducibility

- Toolchain versions are pinned per role image.
- Commands are executed from strict allowlists; environment variables are controlled.
- Outputs (test logs, coverage, diagnostics) are stored as signed, hash-addressed artifacts.

### 3) Formal Verification (Scoped, Day One MVP)

Formal verification is applied to the small security kernel: the pieces that enforce separation, capability rules, and audit integrity.

#### What Must Be Formally Specified

Day one target scope:

- Capability manifest semantics
  - what it means for an action/data flow to be permitted
- Scheduler invariants
  - no escalation-in-place
  - role isolation constraints (eg, no role gets both workspace RW and public egress + secrets)
- Artifact routing invariants
  - only allowed data class flows across roles
  - only hash-addressed artifacts may be consumed
- Audit invariants
  - event hash chaining rules
  - which events must exist before advancing state
- Broker/RPC invariants
  - request/response schemas and replay protection

#### Approach

- Check in a formal specification (eg, TLA+) covering the above invariants.
- CI runs a model checker against bounded scenarios.
- Align runtime code with spec via:
  - explicit invariant checks (fail closed)
  - property-based tests that generate action sequences and validate invariants

Day one deliverables:

- A checked-in formal spec for manifest + scheduler (“no escalation-in-place”) + artifact routing.
- Automated checking in CI.
- A traceability mapping between spec concepts and code modules (so violations are actionable).

### 4) Zero-Knowledge Proofs (Narrow, Day One MVP)

ZK is not intended to prove arbitrary LLM reasoning. RuneCode uses ZK to produce narrowly scoped attestations about deterministic computations and integrity claims, without revealing private inputs.

#### What ZK Is Used For (MVP Scope)

MVP ZK goals are intentionally narrow:

- Prove integrity claims about RuneCode’s deterministic records without revealing sensitive contents.
- Support sharing evidence externally (or across trust zones) while minimizing disclosure.

Candidate proof statements:

- Audit integrity: “I know a sequence of audit events (private witness) whose hashes form this published audit root under policy/version X.”
- Policy decision integrity: “Given manifest hash M and request hash R, the policy program P(version V) outputs decision D.”

#### How ZK Fits Into the System

- Proof generation is an explicit workflow step (not required to enforce policy at runtime).
- Proofs are stored as artifacts and verified deterministically.
- Proof artifacts are logged (hash + verifier result) in the audit chain.

Day one deliverable expectation:

- RuneCode can generate and verify at least one ZK proof type tied to RuneCode artifacts/logs (narrow scope, explicit command).
- Proof statements are based on hash-addressed artifacts and public policy/version identifiers.
- Verification is fast and deterministic; failure does not change history (it is recorded and the run is flagged).

Design requirement for feasibility:

- Artifacts and logs must be hash-addressed from day one.

## Hardware + Scaling Expectations

- Minimal hardware can run RuneCode with a single running isolate slot (fully serialized).
- Two running isolate slots provide a practical baseline (one workspace slot + one gateway slot).
- More CPU/RAM increases `max_running_isolates` and pool sizes; the trust model and manifests do not change.
- Throughput bottlenecks are typically in (a) tests/builds and (b) model completions; scale workspace slots and `model-gateway` pool independently.

Practical minimums discussed:

- Supported baseline: x86_64/ARM64 machine with hardware virtualization enabled, 16GB RAM, and fast storage (NVMe strongly preferred).
- Raspberry Pi: experimental/aspirational target. Useful for verification/audit viewing and low-concurrency runs; not expected to be pleasant for heavy builds, many isolates, or ZK proving.
- Consumer laptop/mini PC: 16GB workable; 32GB+ for smoother multi-isolate concurrency; NVMe preferred.

ZK practicality note:

- ZK proof verification should be lightweight and run on low-end machines.
- ZK proof generation may be resource-intensive depending on the proving system and statement; RuneCode should allow proving to run on more capable hardware while keeping verification and audit integrity universally available.

## Open Decisions / Next Design Milestones

- Implement how `runecode-core` selects/enforces the isolation backend (microVM preferred; hardened containers fallback) and how the assurance level + acknowledgments are surfaced.
- Implement microVM guest images, disk encryption plumbing (LUKS2 baseline), and vsock/transport wiring.
- Implement container hardening baseline and how egress allowlists are enforced cross-platform.
- Implement key management integration (TPM/Secure Enclave where available; hardware-key approvals; rotation + revocation).
- Implement emergency stop and credential revocation list enforcement across gateways.
- Implement git target allowlists (canonical repo identity) and map them to repo-scoped token issuance in `runecode-secretsd`.
- Finalize the data class taxonomy and inter-role flow matrix; enforce fail-closed classification.
- Implement the image/toolchain build pipeline, signing, and boot-time verification.
- Implement per-role resource quotas and watchdog enforcement.
- Decide/implement an audit anchoring mechanism (TPM PCR, RFC 3161 timestamps, or witness service).
- Decide the exact minimal role set for “runs at all” on low hardware.
- Define the manifest schema (capabilities, data classes, opt-ins).
- Choose transport details (vsock-first vs mtls-first) and artifact store design.
- Define the redaction policy framework (what is always redacted, what can be explicitly allowed).
- Define which “data classes” are allowed to be sent to third-party models when opted-in.
- Define how dependency fetching/caching works for offline workspace roles (image-baked toolchains vs `deps-fetch`).
- Implement web research egress policy model (domains + patterns + bundle expansion) with MVP support for tight patterns.
- Design the future web research UX for curated bundles and domain-expansion via new stage manifests.
- Select the ZK proving approach for MVP (eg, zkVM vs circuit) and define the first proof statement + public inputs.
- Implement GitHub Actions CI workflows (unit + integration) for Linux, macOS, and Windows.
