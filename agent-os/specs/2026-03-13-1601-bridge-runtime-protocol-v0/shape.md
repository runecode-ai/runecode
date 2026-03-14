# Bridge Runtime Protocol v0 — Shaping Notes

## Scope

Define the shared protocol contract for subscription-backed and other user-installed local provider runtimes behind `model-gateway`.

## Decisions

- Shared bridge/runtime object families are defined once and reused by later provider specs.
- Compatibility is probe-driven and fail-closed; newer vendor versions are not trusted implicitly.
- Bridge runtimes remain LLM-only and never receive workspace or patch capabilities.
- Token delivery must avoid environment variables and raw secret logging.
- Audit and TUI surfaces must make untested-version and persisted-session posture visible.

## Context

- Visuals: None.
- References: `agent-os/specs/2026-03-08-1039-secretsd-model-gateway-v0/`, `agent-os/specs/2026-03-12-1030-auth-gateway-role-v0/`, `agent-os/specs/2026-03-11-1920-openai-chatgpt-subscription-provider-v0/`, `agent-os/specs/2026-03-11-1921-github-copilot-subscription-provider-v0/`
- Product alignment: keeps subscription-backed provider access within the same least-privilege, auditable trust model as direct HTTP providers.

## Standards Applied

- `security/trust-boundary-interfaces` - bridge contracts are part of the shared boundary surface and must stay explicit.
- `security/trust-boundary-layered-enforcement` - provider runtimes cannot weaken broker, policy, or gateway enforcement.
