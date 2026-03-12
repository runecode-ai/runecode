# RuneCode — Security-first AI coding: isolated execution, signed, auditable

[![CI](https://github.com/runecode-ai/runecode/actions/workflows/ci.yml/badge.svg)](https://github.com/runecode-ai/runecode/actions/workflows/ci.yml)
![Status: pre-alpha](https://img.shields.io/badge/status-pre--alpha-orange)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

RuneCode is a security-first agentic automation platform for software engineering.
It treats isolation and cryptographic provenance as co-equal pillars: work runs in tightly scoped isolates with deny-by-default capabilities, explicit artifact-based data movement, and a tamper-evident audit trail.

## Status

RuneCode is pre-alpha and not production-ready.
This repository contains scaffolding, specs, and early guardrails; most runtime enforcement (isolation, signing, auditing) is not yet implemented.

- Roadmap: `agent-os/product/roadmap.md`

## Why RuneCode

- **Isolation is the boundary:** risky work runs in tightly scoped isolates; the scheduler/runner is treated as untrusted.
- **Deny-by-default posture:** capabilities (egress, secrets, workspace writes) are explicit and intended to be policy-controlled.
- **Signed, auditable evidence:** the goal is a tamper-evident trail for actions, decisions, and artifacts (diffs/logs/results).
- **Explicit data movement:** handoffs are intended to happen via hash-addressed artifacts, not implicit shared state.

## Threat model (micro)

RuneCode is built around a pessimistic assumption: any single AI/agent component (including the workflow runner) can be compromised or behave maliciously.
The architecture aims to reduce blast radius and preserve forensics by:

- Separating trusted control-plane components from an untrusted scheduler/runner.
- Preventing any one component from having broad combined powers (network + workspace + long-lived secrets).
- Making cross-boundary interfaces schema-driven and auditable.

Design inspiration includes compartmentalization models (e.g., QubesOS), applied to agentic workflows.

## Security Model (High Level)

RuneCode is designed around two local trust domains:

- **Trusted domain:** Go control plane daemons + Go TUI client
- **Untrusted domain:** TS/Node workflow runner (scheduler)

Key invariants (design targets; enforcement is implemented incrementally):

- Deny-by-default capabilities; explicit opt-ins for higher-risk posture changes
- No single component combines public network egress + workspace access (especially RW) + long-lived secrets
- Cross-boundary communication is brokered and schema-validated (no ad-hoc JSON)

Details (diagram, allowed interfaces, prohibited bypasses, and CI guardrail): `docs/trust-boundaries.md`.

## Repository Layout

- `cmd/` — trusted Go binaries (launcher, broker, secretsd, auditd, TUI)
- `internal/` — trusted Go libraries
- `runner/` — untrusted TS/Node workflow runner package
- `protocol/` — cross-boundary schema + fixture roots (source of truth)
- `agent-os/` — product/spec/standards documents (git-native system of record)

## Development

Canonical local workflow uses Nix + `just` (Nix `>= 2.18`):

```sh
nix develop -c just ci
```

Common commands:

```sh
just fmt
just lint
just test
just ci
```

Optional: enable automatic dev-shell entry with `direnv` + `nix-direnv`:

```sh
direnv allow
```

Non-Nix fallback (e.g., Windows): install Go 1.25.x, Node `>=22.22.1 <25`, and `just`, then run:

```sh
just ci
```

## Components

The repo currently ships stub binaries that are scaffolded and intentionally do not start network listeners.
You can inspect their help output:

```sh
go run ./cmd/runecode-tui --help
go run ./cmd/runecode-launcher --help
go run ./cmd/runecode-broker --help
go run ./cmd/runecode-secretsd --help
go run ./cmd/runecode-auditd --help
```

## Docs

- Mission: `agent-os/product/mission.md`
- Roadmap: `agent-os/product/roadmap.md`
- Trust boundaries: `docs/trust-boundaries.md`

## Contributing

See `CONTRIBUTING.md`. DCO sign-off is required (`git commit -s`).

## Security

Please do not open public issues for security vulnerabilities. See `SECURITY.md`.

## License

Apache-2.0. See `LICENSE` and `NOTICE`.
