# Verification

## Planned Checks
- `runectx validate --json`
- `runectx status --json`
- `just test`

## Verification Notes
- Confirm the split preserves runner and durable-state requirements from the prior combined change.
- Confirm workspace roles and deterministic gates remain tracked as separate child features.

## Close Gate
Use the repository's standard verification flow before closing this change.
