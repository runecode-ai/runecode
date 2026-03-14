# Initial Spec Suite (MVP + Post-MVP) — Shaping Notes

## Scope

Define a small set of initial specs for RuneCode, split into MVP vs post-MVP, and update the product roadmap accordingly.

## Decisions

- MVP uses Nix Flakes for the canonical local dev environment, with `direnv` auto-entry and `just` for common developer commands.
- MVP is single-user and single-machine (no multi-user daemon or remote control plane).
- MVP uses SQLite (WAL) for durable local state and indexing (pinned SQLite version when WAL is enabled).
- MVP includes formal verification (TLA+ model checking in CI).
- The first ZK proof ships only if a proving system can be selected with deterministic, fast verification; otherwise release is deferred rather than weakening the proof contract.
- MVP runtime targets Linux + KVM first; macOS is included in MVP only if it does not materially slow delivery.
- Windows support in MVP is enforced via CI workflows (lint/tests/integration where possible) to keep the codebase portable; later Windows microVM runtime work is tracked in `agent-os/specs/2026-03-08-1039-windows-microvm-runtime/`.
- MVP supports both microVM and container isolation backends.
- Container backend is explicit opt-in only and must never be an automatic fallback when microVMs fail.
- MVP starts with JSON messages validated by JSON Schema, while keeping the logical object model encoding-agnostic so later local transport work in `agent-os/specs/2026-03-13-1602-local-ipc-protobuf-transport-v0/` can adopt protobuf over local IPC.
- MVP uses vsock-first for isolate <-> host transport on Linux with a virtio-serial fallback; a message-level authenticated+encrypted session is always required.
- MVP UX is CLI + minimal TUI.
- MVP approval posture is "moderate" (checkpoint-style approvals); later profile expansion is tracked in `agent-os/specs/2026-03-10-1530-approval-profiles-v0/`.
- MVP model provider access is API-key based; later subscription-backed providers are tracked in `agent-os/specs/2026-03-11-1920-openai-chatgpt-subscription-provider-v0/`, `agent-os/specs/2026-03-11-1921-github-copilot-subscription-provider-v0/`, and `agent-os/specs/2026-03-13-1601-bridge-runtime-protocol-v0/`.
- Spec docs must not mention the source discovery doc filename/path; they should stand on their own.

## Context

- Visuals: None.
- References: No existing code references (repo is currently docs-only).
- Product alignment: Aligns tightly with the product mission (security-first) and tech direction (Go control plane + Go TUI + TS LangGraph runner treated as untrusted), with a clarified constraint that container isolation is opt-in only.

## Standards Applied

- product/roadmap-conventions — Applies to the `agent-os/product/roadmap.md` update.
