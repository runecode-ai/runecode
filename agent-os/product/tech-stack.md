# Tech Stack

## Frontend

- Go (TUI-first dashboard for runs, approvals, diffs, artifacts, and audit timeline)

## Backend

- Go (local control plane / security kernel: launcher/PEP, broker, secrets daemon, audit daemon)
- TypeScript + Node.js (LangGraph-based workflow runner/scheduler; treated as untrusted at runtime)

## Database

- To be defined (durable local state for runs/artifacts/audit indexing)

## Other

- Isolation: microVMs preferred (QEMU with KVM/HVF/WHPX); hardened containers fallback with explicit reduced-assurance UX
- Security/provenance: signed capability manifests, per-isolate identity keys, hash-chained audit events, content-addressed artifacts
- Storage: encrypted workspace disks/snapshots (LUKS2/dm-crypt baseline)
- Key management (where available): TPM/Secure Enclave; user-presence approvals via OS-confirmation and/or FIDO2/YubiKey
- CI: GitHub Actions (Linux/macOS/Windows)
- Workflow/docs: AgentOS (`agent-os/**`) as git-native source of truth; OpenCode runs inside offline workspace roles
