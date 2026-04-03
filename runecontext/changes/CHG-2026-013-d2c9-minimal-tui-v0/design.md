# Design

## Overview
Implement the local TUI for runs, approvals, artifacts, and audit posture over the broker local API.

## Key Decisions
- TUI is a separate least-privilege client; it does not embed privileged execution.
- Use Bubble Tea as the TUI framework.
- The assurance level (microVM vs container) must be prominent.
- The active approval profile is part of the user safety posture and should be visible and explained (MVP default: `moderate`).
- Approval requests must be explainable from structured data (reason codes + what changes if approved).

## Main Workstreams
- Bubble Tea App Skeleton
- Core Screens (MVP)
- Local API Integration
- Safety UX

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
