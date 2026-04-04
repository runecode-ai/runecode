# Design

## Overview
Use a project-level change to coordinate secure model/provider access features and shared trust-boundary requirements.

## Key Decisions
- Shared security invariants apply to all child features.
- Secrets lifecycle, auth, bridge, and provider lanes remain separable feature boundaries.
- Verification remains feature-level, with this project change tracking sequencing and integration posture.

## Main Workstreams
- Shared foundation tracking (`secretsd` and model-gateway).
- Auth and bridge feature sequencing.
- Provider-specific feature sequencing and integration checks.
