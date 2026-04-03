# Design

## Overview
Define the explicit reduced-assurance container backend, including opt-in UX, hardened defaults, artifact movement, and policy integration.

## Key Decisions
- Containers are never a silent fallback; they require explicit opt-in and acknowledgment.
- The active backend and its assurance level are treated as first-class audit data.
- Container networking is isolated by default (no egress); any allowed egress is enforced via explicit network namespace + firewall/proxy rules, not convention.

## Main Workstreams
- Opt-In UX + Audit
- Hardened Container Baseline
- No Host Mounts + Artifact Movement
- Policy Integration

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
