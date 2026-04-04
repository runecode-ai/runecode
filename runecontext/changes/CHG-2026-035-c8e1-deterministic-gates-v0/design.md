# Design

## Overview
Implement deterministic gate execution with explicit, auditable evidence outputs and fail-closed semantics.

## Key Decisions
- Gates are deterministic and produce typed evidence artifacts.
- Gate failures fail the run by default.
- Any override requires explicit approval and audit events.

## Main Workstreams
- Gate framework and execution order.
- Evidence artifact schema and retention linkage.
- Retry and override policy integration.
