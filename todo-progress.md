# Todo progress

## Role of this file
- Persistent handoff file between sessions.
- Keep a short project context without reopening the whole history.
- Update it when a todo changes status or when a technical decision is locked.

## Current target
- Deliver the MVP "reservation de salles" from `enonce.md`.
- Aim for maximum grading points with minimum complexity.

## Stable decisions
- Monorepo Go with 2 binaries: `cmd/server` and `cmd/cli`
- TUI with `tview`
- REST API with `gin`
- PostgreSQL with `pgx/v5`
- Simple auth with login/password + JWT
- Local `config.json` and `state.json`
- Preloaded rooms and demo users
- Fixed 1h slots, no complex calendar
- Shared `internal/domain` package for API/CLI/TUI models
- Server-side validation for slot conflicts and ownership
- Async HTTP calls in TUI to avoid blocking the interface

## Source of truth
- Functional scope: `enonce.md`
- MVP implementation plan: `plan.md`
- This file: inter-session progress and handoff

## Status legend
- `pending`: not started
- `in_progress`: currently being worked on
- `done`: finished
- `blocked`: waiting for a decision or missing input

## Todo board

| Status | ID | Todo | Notes |
| --- | --- | --- | --- |
| done | setup-skeleton | Setting up project skeleton | Go module, commands, shared packages and Makefile created |
| done | setup-postgresql | Setting up PostgreSQL schema | Embedded migration and seed files added with demo rooms/users |
| done | build-auth-api | Building authentication API | `POST /auth/login`, `GET /auth/me` and JWT middleware implemented |
| done | build-booking-api | Building booking API | Rooms, availability, booking, my reservations and cancellation implemented |
| done | build-sync-api | Building sync API | `GET/PUT /me/config` and `GET/PUT /me/state` store JSONB per user |
| done | build-cli-core | Building CLI core | Local JSON stores and typed HTTP client implemented |
| done | build-tui-screens | Building TUI screens | Home, login, settings, rooms and reservations screens implemented |
| blocked | prepare-delivery | Preparing delivery assets | README, Makefile and GoReleaser are ready; actual GitHub release and free hosting need a real repo, credentials and target infra |
| done | copy-plan-local | Copying local plan file | `plan.md` copied to repo root |
| done | create-progress-tracker | Creating progress tracker | `todo-progress.md` created for handoff |

## Next recommended todo
- `prepare-delivery`

## Session log
- 2026-05-28: Initial MVP plan created from `enonce.md`.
- 2026-05-28: `plan.md` copied into the project root for local tracking.
- 2026-05-28: `todo-progress.md` created to keep inter-session context.
- 2026-05-28: Bootstrap started with shared domain models, env-driven server config and non-blocking TUI convention.
- 2026-05-28: Backend MVP completed with PostgreSQL schema, seeded data, JWT auth, room availability and reservation APIs.
- 2026-05-28: CLI MVP completed with local JSON storage, API client and TUI screens for login, settings, rooms and reservations.
- 2026-05-28: Delivery assets prepared (`README.md`, `Makefile`, `.goreleaser.yaml`), but external deployment/release remain blocked by missing remote infrastructure.
- 2026-05-28: `README.md` enriched with a step-by-step launch guide and a ready-to-send tester message.
- 2026-05-28: CLI startup deadlock fixed by replacing direct `QueueUpdateDraw` calls with a safe update helper that does not block before `tview` starts.

## Update rules
- When a todo starts, move it to `in_progress`.
- When a todo ends, move it to `done` and add a short result note.
- If work is blocked, switch to `blocked` and write the blocker clearly.
- Keep this file short: status, decision, next step, blocker.
