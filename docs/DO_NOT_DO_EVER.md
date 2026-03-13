# DO NOT DO EVER

This is mandatory for `quo`. It captures the deliberate non-idiomatic Go style used in `klec`/`pmlib` and must be enforced on every change.

## 1) Standard Go Infrastructure Patterns Intentionally NOT Used

- Never use `html/template` or `text/template` for page rendering. Never use templates.
- Do not use JSON-first app flow (`encoding/json`) for core page behavior.
- Do not introduce REST/SPA API layers
- Do not add ORM/query-builder layers (`gorm`, `sqlx` query builders, etc.).
- Do not write ad-hoc SQL anywhere, think about what stored procedure might be needed instead.
- Do not force exported struct fields just for serialization/template compatibility.
- Do not add DTO/view-model duplication just to match common Go API conventions.
- Do not introduce `context.Context` plumbing across all function signatures unless explicitly required.
- Do not introduce middleware stacks for concerns that are currently deferred (auth enforcement, observability, rate limiting).
- Do not split logic into many packages only to satisfy “idiomatic” package layering.
- Do not introduce dependency-injection frameworks or interface-heavy abstraction around concrete app flows.
- Do not introduce code generation, reflection-based binding, or tag-based validation frameworks.
- Do not introduce alternate routers/frameworks beyond the existing minimal stack (`chi` + `net/http`).
- Do not normalize files to conventional formatting conventions if it rewrites the local house style.

## 2) Small Standard Surface We DO Use

- `net/http` handlers and request form parsing.
- `flag` for runtime config (`-port`, DB CLI overrides).
- `chi` for explicit route registration.
- direct DB calls through `App.DB.Call(...)` / `CallRow(...)`.
- plain server-rendered HTML output via `klec/lib/output` + `klec/lib/htmlHelper`.
- explicit static file serving under `/static/*`.

## 3) Local Recommendations To Follow

- Keep feature files in `package main` using numeric/page naming already present.
- Keep routes explicit in `0.main.go`, one per line.
- Keep app-global wiring in `0.bootstrap.go` + `0.app.types.go`.
- Keep business rules in `quote.rules.*.go`; handlers orchestrate only.
- Keep phone CSS centralized in `static/css/phone.quote.css`.
- Keep JS minimal and page-local (`static/js/1.choices.js` style).
- Keep IDs/selectors and form keys stable and explicit.
- Keep behavior source-of-truth split: infrastructure from `klec`, quote business rules from `redef`, canonical price rows from `plan_products_query`.
- For server-driven UI changes, use klec rewrite protocol first: `data-post` + `data-record` + rewrite response.

## 4) Absolute Never List (Enforcement)

- Never use gofmt (instead read CODING_STYLE.md before writing code)
- no go fmt
- Never introduce mutexes.
- Never run shallow tests when behavior depends on full request/session/UI flow.
- Never write `end;` in SQL routine blocks. Use `end` on its own line, then the delimiter line.
- Never rearchitect when a localized port patch is sufficient.
- Never redesign UI/flow beyond what is required for mobile parity.
- Never use templates. Never use templates.
- Never introduce JSON app-flow, ORM, or API-first layering.
- Never render page-sized UI with large raw HTML string blobs when `htmlHelper` builders can express the same structure clearly.
- Never export fields.
- Never move `/edit` maintenance scope into `quo`.
- Never bypass `plan_products_query` for Step 1 quote plan pricing.
- Never add new JS code without explicit user vetting/approval first.

