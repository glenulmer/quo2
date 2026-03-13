# KLEC Web Build Style (for `quo`)

This is the target architecture and coding style for new UI work in `quo`, based on how `klec` is built.

## Non-negotiable

- Never use templates. Never use templates.
- Never use JSON-first API design as the primary app flow.
- Never introduce REST/SPA layering to do page behavior already covered by server-rendered flows.

## Core Architecture

- Use `chi` + `net/http` only.
- Keep route registration explicit in `0.main.go` using one `r.Get(...)` or `r.Post(...)` per line.
- Keep only `GET` and `POST` methods for feature behavior.
- Keep auth wrapping explicit at route registration (`App.Auth(...)`).

## Page Construction Style

- Build HTML directly in Go with `klec/lib/htmlHelper` (`Head`, `Table`, `TR`, `TD`, `Div`, etc.).
- Write pages with `Writer(w).Add(...)` and explicit composition order:
- `head.Left()`
- `Navbar(...)`
- page title/content/tables/forms
- `head.Right()`
- Keep IDs, form keys, and selectors explicit and stable.
- Keep feature code in `package main` page files (`page.*.go`) and avoid unnecessary package decomposition.

## HTML Readability Rules

- Prefer semantic element builders (`main`, `section`, `form`, `label`, `input`, `select`, `button`) over raw HTML string blocks.
- Prefer `Elem(...).Class(...).Wrap(...)` composition over long text fragments that look like templates.
- Keep one UI block per builder call chain so structure is easy to scan in code review.
- Keep content text in `.Text(...)` and attributes in `.KV(...)`/`.Id(...)`/`.Name(...)`/`.Value(...)`.
- Use CSS classes for appearance; avoid inline style attributes in Go unless there is a strict one-off need.
- Use explicit labels linked with `for`/`id` for controls.
- Keep `Head()` setup concise and route-local: title, icon, CSS/JS includes only for the page being rendered.
- If a page function starts to read like pasted HTML, refactor into helper builders before adding more logic.

## Control-First Build Pattern

- Before writing page markup, check `klec/lib/htmlHelper` for existing control primitives (`TextIn`, `Select`, `Option`, `CBox`, `Wedge`, `Crud*` helpers).
- For repeated UI, add small control helpers in app code (for example `Card`, `FormField`, `DatePicker`) instead of repeating `Elem(...)` chains in handlers.
- Build specialized controls by extending existing helpers first (for example, derive `DatePicker` from `TextIn` and set only the differing attributes/classes).
- Keep control names semantic to humans (`Card`, `FormField`, `DatePicker`) rather than framework/device-centric names.
- Keep post contracts explicit at the control level: stable `name`/`id` values that map directly to `req.FormValue(...)`.

## Naming By Purpose

- Name structs and fields by intended use in the page flow, not by raw query mechanics.
- Prefer names that explain purpose (`planAttribs`, `plansLoadErr`, `RenderPlans`) over generic transport names.
- If a structure's purpose is unclear at design time, stop and ask for clarification before finalizing names.

## Error Signaling Pattern (Wedge + Message)

- Follow the `klec` visual pattern: render a wedge beside each control and use error state to turn it red.
- Use a field wrapper (`FormField`) to centralize label, wedge, control, and error message rendering.
- Show helper/error text only when invalid; keep quiet UI when valid.
- Validate on both client and server:
- Client validates on `change`/`blur`, toggles `has-error`, updates error text, and blocks submit while invalid.
- Server re-validates all posted values and re-renders authoritative error state.
- Clearing an error must clear both signals: wedge color and message text.

## GET Handler Pattern

- `GET` handlers parse query/form values as needed (`req.ParseForm()`).
- Validate required IDs early (`RequireID(...)`) and return early on invalid input.
- Load page data through stored procedure queries.
- Render the full page HTML server-side.
- Include only page-local JS/CSS needed for that page.

## POST Handler Pattern

There are two valid `POST` patterns:

1. CRUD table updates (AJAX rewrite contract)
- Parse form payload (`ParseMultipartForm` / `FormValue`).
- Resolve IDs with canonical keys (`provider_id`, `family_id`, `plan_id`, `addon_id`, `account_id`, `row_id`, `entity_type`, `verb`).
- Call stored procedures via `App.DB.CallRow(...)`.
- Return rewrite/message instructions using `SendResponse(...)`, `RewriteTBody(...)`, `RewriteRow(...)`, `RemoveRow(...)`, `Note(...)`.

2. File-upload workflow pages (full page re-render)
- Parse multipart upload.
- Validate and import through server-side workflow functions.
- Re-render the same page (full HTML) with message blocks and current controls.

## Database Contract

- Stored procedures are the behavior source of truth.
- Prefer procedure calls over ad-hoc SQL in Go handlers.
- Keep procedure call order and parameter shape explicit and readable.
- Use transaction wrappers for batch imports (`WithTx`) where needed.

## Frontend Contract

- JS remains thin and page-local.
- For table CRUD pages, JS posts form data and consumes server rewrite instructions.
- If a UI change may require server HTML rewrites, klec protocol is first choice:
- `data-post` + `data-record` request contracts
- `FormData` posts (not custom header/JSON contracts when avoidable)
- rewrite queue responses (`RewriteRow`, `RewriteTBody`, `RemoveRow`) consumed by `serverResponse(...)`
- Use `data-post` for endpoint binding and `data-record` for per-row identity/payload.
- Keep client validators focused on control-level behavior; business truth stays server/DB side.
- Any decision to add new JS must be reviewed with the user before implementation.

## Practical Checklist Before Merging

- No templates added.
- Routes remain explicit `GET`/`POST`.
- Handler logic is orchestration, not business-rule sprawl.
- Stored procedures back all durable writes.
- Page HTML is constructed with helper DSL, not string-concatenated full documents (except legacy static pages).
- IDs/keys/selectors align between Go and JS.
