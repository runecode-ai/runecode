# Design

## Overview
Implement a content-addressed artifact store and a minimal data classification system with enforced role-to-role flow rules.

## Key Decisions
- Artifacts are hash-addressed and immutable.
- Unknown or ambiguous artifacts are classified as the most restrictive class (fail-closed).
- Artifact contents are stored on encrypted-at-rest storage by default (no silent plaintext mode).
- Artifact retention/GC is required to avoid unbounded growth.
- `approved_file_excerpts` are only created via explicit human approval; unapproved excerpts use a more restrictive class (`unapproved_file_excerpts`) and are not eligible for third-party egress.
- Promotions are hardened: approvals are explicit, reviewable, rate-limited, and revocable via policy (no history rewriting).
- Derived evidence is stored as explicit artifacts with their own data class (e.g., `audit_verification_report`).

## Main Workstreams
- Define MVP Data Classes
- Content-Addressed Artifact Store (CAS)
- Flow Matrix Enforcement
- Quotas + Limits (Minimal)
- Garbage Collection + Retention (Minimal)

## RuneContext Migration Notes
- Canonical references now point at `runecontext/project/`, `runecontext/specs/`, and `runecontext/changes/` paths.
- Future-facing planning assumptions are rewritten to use RuneContext as the canonical planning substrate for this repository.
- Where this feature touches project context, approvals, assurance, or typed contracts, the migrated plan assumes bundled verified-mode RuneContext integration from the feature surface rather than a later retrofit.
