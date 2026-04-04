# Design

## Overview
Implement the untrusted runner orchestration and durable state authority for secure, resumable runs.

## Key Decisions
- Runner is untrusted and never directly executes privileged operations.
- Runner persistence stores control-plane state only.
- Pause/resume and crash recovery rely on durable state transitions.
- All real execution remains brokered and policy-authorized.

## Main Workstreams
- Runner contract and packaging constraints.
- Durable state schema and migration rules.
- Propose-to-attest execution loop integration.
