# Design

## Overview
Establish Phase 1 durable governance artifacts before bulk migration by creating explicit cross-cutting decisions and linking this umbrella change to those decisions.

## Decisions Captured
- RuneContext replaces `agent-os/` as canonical repo-local planning and standards substrate.
- Migration runs in verified mode with assurance artifacts under `runecontext/assurance/`.
- Future planning assumes bundled-by-default RuneContext and verified-mode requirements for RuneCode-managed repos.
- Future planning assumes RuneCode owns user-facing UX while RuneContext remains generic and machine-friendly.
- Generic advisory consumer-compatibility warnings may exist in RuneContext, while hard compatibility enforcement remains in RuneCode.

## Traceability
- Decision records live under `runecontext/decisions/` and are referenced from this change status.
- This change acts as the umbrella tracker for migration sequencing and verification notes.

## Shape Rationale
- Large, ambiguous, or high-risk feature work should move to full mode early.
