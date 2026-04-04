## Summary
Track workflow execution as a project-level change while runner, roles, and gates ship through scoped child features.

## Problem
The prior combined feature mixed multiple independently deliverable components, which limited sequencing and verification granularity.

## Proposed Change
- Keep this change as the workflow execution parent project.
- Track `CHG-2026-033-6e7b-workflow-runner-durable-state-v0` for runner and durable-state boundaries.
- Track `CHG-2026-034-b2d4-workspace-roles-v0` for role execution boundaries.
- Track `CHG-2026-035-c8e1-deterministic-gates-v0` for gate determinism and evidence.

## Why Now
This work remains scheduled for v0.1.0-alpha.3 as the first honest end-to-end slice built strictly on the secure substrate, with remaining hardening and scope completed in v0.1.0-alpha.4. Keeping it as a parent project preserves roadmap traceability while allowing finer-grained feature delivery and verification.

## Assumptions
- `runecontext/changes/*` is the canonical planning surface for this repository.
- RuneCode keeps the end-user command surface while using bundled RuneContext capabilities under the hood where project context or assurance is involved.
- Context-aware delivery for this feature is planned directly against verified-mode RuneContext rather than a later retrofit from legacy Agent OS semantics.

## Out of Scope
- Runtime implementation details that belong in child feature changes.
- Re-introducing legacy Agent OS planning paths as canonical references.

## Impact
Keeps workflow execution reviewable as a parent project with explicit feature-level ownership and sequencing.
