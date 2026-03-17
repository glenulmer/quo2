# AGENTS - Quo2 Guardrails

## 1) Non-Negotiables

1. Keep `package main` + page-oriented files (`0.*`, `page.*`).
2. Keep `chi` + `net/http`; explicit route lines.
3. Keep server-rendered HTML via `klec/lib/htmlHelper`.
4. Never use templates, JSON-first app flow, REST/SPA layering, ORM/DI/DTO architecture.
5. Never run `gofmt`/`go fmt`/`goimports`/IDE auto-formatters, and never do style-only "idiomatic Go" rewrites.
6. Never put `panic(...)` in request/page handlers (`page.*.go`).
7. No detailed error plumbing 

## 2) Rewrite Contract (`data-post` / `data-record`)

1. Mount one `data-post` container per POST endpoint.
2. Inside it, each editable unit must have exactly one `data-record` root.
3. Control `name` keys must map directly to `req.FormValue(...)` keys.
4. Buttons use `name=create|update|delete` and post through `validate.js`.
5. Server responses for incremental UI updates must use rewrite messages (`RewriteRow`, `RewriteTBody`, `RemoveRow`).
6. Rewrite target IDs/selectors are stable API: do not rename casually.

## 3) Static Lookup Policy

1. Static DB lookups are loaded at bootstrap and cached on `App`.
2. Request handlers must read cached values, not query static lookups.
3. Static registry lives in `docs/STATIC_LOOKUPS.md`.
4. If editing/adding `App.DB.Call(...)` or `CallRow(...)`, read `docs/STATIC_LOOKUPS.md` first; otherwise do not load it.
5. Every proc call in Go must be declared in `docs/QUERY_REGISTRY.md` with `Kind` + `Allowed Layer`.
6. Unknown/unregistered proc calls are guardrail failures.

## 4) Interaction Rule

1. Keep feedback technical and concrete.
2. Do not label the user with moral/political identity terms.

## 5) Quality Check

- You must achieve code clarity, simplicity, maintainability, and no needless duplication.
- I will ask another AI to criticize your for all three factors.
- Anticipate all criticisms, do not rush to a solution.

## 6) Enforcement

- `./scripts/check-guardrails.sh`
- `./scripts/check-all.sh`
- `./scripts/run-dev.sh -port 7777`

Before any completion claim, run `./scripts/check-all.sh`.
`run-dev.sh` refuses to start if guardrails fail.
