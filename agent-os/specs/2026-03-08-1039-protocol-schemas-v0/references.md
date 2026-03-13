# References for Protocol & Schema Bundle v0

## Product Context

- **Mission:** `agent-os/product/mission.md`
- **Roadmap:** `agent-os/product/roadmap.md`
- **Tech stack:** `agent-os/product/tech-stack.md`
- **Trust boundaries:** `docs/trust-boundaries.md`

## Primary Downstream Specs

- `agent-os/specs/2026-03-10-1530-approval-profiles-v0/`
  - depends on stable approval request/decision objects, approval profile versioning, and trigger-code taxonomy
- `agent-os/specs/2026-03-12-1030-auth-gateway-role-v0/`
  - depends on reserved auth object families, stable principal identity, shared errors, and typed lease/audit handoffs
- `agent-os/specs/2026-03-11-1920-openai-chatgpt-subscription-provider-v0/`
  - depends on reserved bridge/runtime identity objects, compatibility-probe shapes, token-handoff invariants, and model-gateway request/response contracts
- `agent-os/specs/2026-03-11-1921-github-copilot-subscription-provider-v0/`
  - depends on the same bridge/runtime identity and error-mapping foundations plus explicit LLM-only posture objects

## Core MVP Dependencies

- `agent-os/specs/2026-03-08-1039-broker-local-api-v0/`
- `agent-os/specs/2026-03-08-1039-policy-engine-v0/`
- `agent-os/specs/2026-03-08-1039-secretsd-model-gateway-v0/`
- `agent-os/specs/2026-03-08-1039-artifact-store-data-classes-v0/`
- `agent-os/specs/2026-03-08-1039-audit-log-verify-v0/`
- `agent-os/specs/2026-03-08-1039-crypto-key-mgmt-v0/`

## Similar Implementations

None in-repo yet; this spec establishes the shared contract layer the rest of the system will implement.

## External References

- RFC 8259: The JavaScript Object Notation (JSON) Data Interchange Format
- RFC 8785: JSON Canonicalization Scheme (JCS)
- RFC 7493: I-JSON Message Format
- JSON Schema draft 2020-12
- JCS reference vectors / implementations: `https://github.com/cyberphone/json-canonicalization`
