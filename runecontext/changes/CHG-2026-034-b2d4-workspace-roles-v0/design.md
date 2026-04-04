# Design

## Overview
Implement the MVP workspace role set with strict capability boundaries and offline operation posture.

## Key Decisions
- Role capabilities remain explicit and least-privilege.
- Workspace roles are offline and non-gateway.
- Command execution uses constrained executors, not raw shell passthrough.

## Main Workstreams
- Role definitions and capability manifests.
- Executor contract and allowlist rules.
- Artifact handoff and output contracts.
