# Verification

## Planned Checks
- `runectx validate --json`
- `runectx status --json`
- `just test`

## Verification Notes
- Confirm the migrated change preserves the legacy task breakdown and acceptance criteria in `tasks.md`.
- Confirm canonical references remain on RuneContext project, spec, and change paths, with no active workflow depending on legacy planning paths.
- Confirm the migrated text assumes RuneContext is canonical, RuneCode owns the user-facing UX, and verified-mode project state remains the expected operating posture.
- Confirm the change still matches its v0.1.0-alpha.2 roadmap bucket and title after migration.
- Confirm trusted artifact hashing, backup-manifest signing, and protocol fixture parity all use RFC 8785 JCS canonicalization semantics.
- Confirm the broker default artifact-store root resolves outside the repository when `RUNE_BROKER_STORE_ROOT` is unset.

## Close Gate
Use the repository's standard verification flow before closing this change.
