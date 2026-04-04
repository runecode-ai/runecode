# Verification

## Planned Checks
- `runectx validate --json`
- `runectx status --json`
- `just test`

## Verification Notes
- Confirm the split preserves model-gateway trust-boundary requirements from the prior combined change.
- Confirm provider-facing features reference this gateway feature for egress controls.

## Close Gate
Use the repository's standard verification flow before closing this change.
