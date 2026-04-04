# Design

## Overview
Implement the network egress boundary for model traffic as a dedicated gateway feature.

## Key Decisions
- Model traffic is explicit opt-in and deny-by-default.
- Gateway has no workspace access and uses leases only.
- Typed boundaries are required for machine-consumed traffic.
- SSRF/DNS-rebinding protections and TLS enforcement are mandatory.

## Main Workstreams
- Gateway role boundary and egress controls.
- Typed model schemas and policy integration.
- Data-class filtering and redaction enforcement.
- Quotas and audit event coverage.
