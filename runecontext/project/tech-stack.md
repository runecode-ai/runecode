# Tech Stack

## Frontend

- Go (TUI-first dashboard for runs, approvals, diffs, artifacts, and audit timeline)

## Backend

- Go (local control plane / security kernel: launcher or policy-enforcement point, broker, secrets daemon, audit daemon)
- TypeScript + Node.js (workflow runner and scheduler; treated as untrusted at runtime)

## Database

- SQLite (WAL) for MVP durable local state and indexing (runs, approvals, artifact metadata, audit indexing)
- Append-only files for large immutable blobs (CAS objects, audit log segments), with SQLite holding indexes and pointers
- SQLite version is pinned when WAL is enabled, including known WAL integrity fixes such as the WAL-reset fix in SQLite >= 3.52.0 or an equivalent backport; record the SQLite library version in audit metadata
- SQLite operational guidance for MVP: avoid cross-process write contention
  - Prefer separate SQLite databases per subsystem, or strict transaction discipline, if contention becomes an issue; record the chosen layout and WAL mode in audit metadata

## Other

- Isolation: microVMs are preferred for high-assurance roles; hardened containers remain an explicit reduced-assurance fallback
- MVP runtime target: Linux + KVM on a single-user local machine; later macOS and Windows runtime work is planned as RuneContext-managed product changes during this migration
- Isolate transport (MVP): vsock-first on Linux, with virtio-serial fallback; message-level authenticated and encrypted sessions are required
- Local APIs and protocols: JSON + JSON Schema for MVP; later local transport evolution is planned as RuneContext-managed product changes during this migration
- Security and provenance: signed capability manifests, per-isolate identity keys, hash-chained audit events, and content-addressed artifacts
- Storage: encrypted workspace disks and snapshots (LUKS2/dm-crypt baseline)
- Key management (where available): TPM or Secure Enclave, with user-presence approvals via OS confirmation and or FIDO2/YubiKey
- CI: GitHub Actions (Linux, macOS, Windows)
- Project planning and standards: RuneContext under `runecontext/**` is the canonical repo-local planning, standards, and product-governance source of truth for this repository
- Product UX: RuneCode owns the normal user-facing command surface and workflows while invoking bundled RuneContext capabilities under the hood
