---
schema_version: 1
id: security/trust-boundary-change-checklist
title: Trust Boundary Change Checklist
status: active
suggested_context_bundles:
    - runner-boundary
---

# Trust Boundary Change Checklist

If you change boundary surfaces (runner access rules, protocol paths, broker API):

- Update `docs/trust-boundaries.md`
- Update guardrails to match:
  - `runner/scripts/boundary-check.js`
  - `runner/scripts/boundary-check.test.js`
- Update `protocol/schemas/` and `protocol/fixtures/` as needed (no ad-hoc formats)
- If the change adds trusted local persistence for cross-boundary evidence or artifacts, document storage sensitivity, backup/restore handling, and recovery behavior in `runecontext/standards/`
- Ensure CI parity (`just ci`) still exercises the boundary checks
- Treat as security-sensitive: keep `CODEOWNERS` coverage + required review
