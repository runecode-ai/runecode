# Standards for Workflow Runner + Workspace Roles + Deterministic Gates v0

## Runner Distribution (Node SEA)

- The workflow runner is packaged as a Node SEA (single executable) for release/runtime distribution.
- SEA is packaging, not a security boundary; the runner remains untrusted at runtime.
- SEA config must ignore `NODE_OPTIONS` (set `execArgvExtension: "none"`) so environment variables cannot silently extend Node runtime flags.
- Bundle the runner into a single injected CommonJS script; do not depend on runtime `node_modules` resolution.
