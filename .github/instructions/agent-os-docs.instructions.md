---
applyTo: "runecontext/project/**/*.md,agent-os/specs/**/*.md,agent-os/standards/**/*.md"
---

Use these references for planning and roadmap review comments:

- `/agent-os/standards/index.yml`
- `/agent-os/standards/product/roadmap-conventions.md`
- `/runecontext/project/roadmap.md`

When reviewing changes in this scope, focus on:

- Roadmap structure remains valid (`Upcoming Features`, `Unscheduled (Needs Specs)`, `Completed Features`).
- The roadmap remains a human-facing summary rather than the lifecycle source of truth.
- Upcoming and completed entries stay outcome-focused and do not reintroduce `agent-os/specs/*` as canonical roadmap links.
- Active lifecycle state stays in `runecontext/changes/*/status.yaml`, with durable completed outcomes in `runecontext/specs/*.md`.
- Standards index entries stay accurate and concise.

Prefer comments that preserve traceability from roadmap items to RuneContext changes or specs when those canonical artifacts exist.
