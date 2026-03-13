# Coding Style Notes (Project-Specific)

This documents the style patterns used in this user's code so changes can match existing code, not generic Go conventions.

## Core Principle

Preserve local style and behavior first. Prefer consistency with surrounding code over standard Go idioms unless a functional, safety, or build issue requires deviation.

## File and Package Organization

- Use `package main` heavily for app layers.
- Keep many top-level files with numeric prefixes.
- Group features by page/workflow instead of strict package decomposition.

## Imports

- Dot imports are intentionally common (for example `. "klec/lib/output"`).
- Raw string literals are often used in imports and string constants (backticks).
- Import blocks may be mixed-tab/space aligned and must not be gofmt-normalized.

## Naming Conventions

- Mixed naming styles are used intentionally:
- CamelCase for exported-ish helpers and page handlers (`GetPrices`, `ProviderSet`).
- Snake-style segments in identifiers they are intended for limited use, not shared by all go source in the package.
- Type suffix patterns: `_t` (`ProviderSet_t`, `PriceInfo_t`) for prominent types.
- Prefix patterns: `tX` for structs intended for local use, `idX` for typed IDs.
- SQL routine names (procedures/functions) are snake_case and use application specific prefixes (eg klec, quote) for new definitions.
- For CRUD request/record keys, use canonical explicit names:
- `provider_id`, `family_id`, `plan_id`, `addon_id`, `account_id`
- `row_id`, `entity_type`, `verb`

## Formatting and Layout

- Single-line control flow is common:
- `if cond { return }`
- `if x == 0 { ...; return; }`
- `switch` with compact case bodies.
- Semicolons appear occasionally and should not be treated as mistakes if already present.
- Keep short dense expressions when they improve local readability.
- Never run gofmt. Never run gofmt.

## Code Style Patterns

- Prefer small, direct functions with straightforward control flow.
- Use utility aliases from `lib/output` (`Str`, `Atoi`, `OnlyDigits`, `Join`, etc.).
- Use builder-style HTML composition (`Div(...).Class(...).Post(...)`) instead of templates.
- Keep constants near usage (`const postPrices = ...`, `const tbodyFilters = ...`).
- Keep procedural page flow in handler functions (read form -> call DB -> rewrite HTML).

## Go Route Registration Style (`0.main.go`)

- Keep one route registration per line (`r.Get(...)` or `r.Post(...)`), without line wrapping unless arguments become unusually long.
- Preserve the current spacing/line-break style in `0.main.go`: compact route lines, with blank lines between logical route groups.
- Keep route path literals as backtick strings when writing inline paths (for example ``r.Get(`/providers`, ...)``).
- Keep auth wrapping explicit at each protected route (`App.Auth(...)`) rather than introducing helper indirection.
- For feature routes, keep GET route declarations before related POST handlers when practical.
- Keep top-level feature groups ordered to mirror the GET paths currently shown on the `KL Menu` (`/`) page.

## Error Handling Style

- Practical, compact checks are preferred over verbose wrapping.
- Panics are used in some data-load paths when state is invalid.
- In request handlers, response messages/notes are often preferred over deep error types.

## Data and DB Interaction

- DB procedures are the behavioral source of truth.
- App logic frequently maps directly to stored procedures.
- Preserve existing proc names, parameter ordering, and response handling patterns.
- In non-SQL code, use stored procedure calls instead of ad-hoc SQL queries.
- If a required DB operation does not exist yet, add/update a stored procedure in `sql/` first, then call it from app code.

## UI/Frontend Coupling Style

- Go handlers return HTML via helper DSL, paired with focused page JS files.
- IDs/constants in Go are expected to match JS selectors exactly.
- Incremental updates use response rewrite helpers (`RewriteHTML`, `SendResponse`).
- Prefer rewrite wrappers over raw method literals when possible:
- `RewriteTBody(...)`, `RewriteRow(...)`, `RemoveRow(...)`.
- For provider-scoped page links and query parsing, use `provider_id` (not mixed legacy aliases).

## Change Rules for Future Edits

- Match the style of the file you are editing first.
- Keep whitespace/quote/backtick conventions consistent with nearby code.
- Avoid introducing new naming schemes into existing files.
- Avoid broad refactors during functional changes.
- If style must be broken for correctness, keep the change minimal and local.

## SQL / Stored Procedure Style

- The `sql/` directory is an incremental migration stream for an already-existing database, not a full schema definition.
- SQL files must remain alphabetically ordered so shell-based loaders execute them in sequence.
- Filename numeric segments are zero-padded to two digits (for example `00`, `01`, `03.05`, `04.08`).
- Add new incremental SQL files immediately before `z.sql`.
- Use lowercase SQL keywords and identifiers (`create or replace procedure`, `select`, `update`, `where`).
- New stored procedure and function names must be prefixed with `klec_` (the repo/database namespace prefix for this project).
- Standard procedure suffixes:
- `_create`: user-initiated create allowed, update not allowed.
- `_update`: update allowed, user-initiated create not allowed.
- `_upsert`: default mutating pattern when create/update are both valid.
- `_softdel`: user-initiated soft delete (logical delete via update).
- `_delete`: user-initiated hard delete of a single row/entity.
- `_kill`: hierarchical or multi-table hard delete workflow.
- `_query`: select returning 0..N rows.
- `_get`: select returning 0..1 row.
- Do not introduce new `_ups` suffix names.
- Keep legacy routine names unchanged unless you are explicitly migrating them.
- Parameter naming style:
- Prefix parameters with `$` (`$prov`, `$year`, `$data`).
- Use local variables with `_` prefix (`_errs`, `_proc`, `_mess`, `_note`, `_exists`).
- Delimiter pattern:
- Wrap procedures with `delimiter ###` ... `delimiter ;`.
- Prefer `end` on its own line.
- Put closing delimiter `###` on its own next line (not `end###`).
- Error handling pattern:
- Use `declare exit handler for sqlexception` in mutating procedures.
- Treat workflow procedures as mutating CRUD-class procedures and apply the same handler pattern.
- In the handler: `get diagnostics ...` and return a single status row with `errs`, `nrows`, `id`, `note`.
- Do not add handlers to read/list procedures (`*_read`, lookup/list procedures).
- Existing exceptions without handlers (`klec_prices_update`, `klec_year_clone`, `klec_year_kill`, `klec_year_softdel`) are legacy patterns and should not define new style.
- Result contract pattern:
- Return one final `select` status row from mutating procedures (`select 0 errs, row_count() nrows, ... note;`).
- Favor compact `if`/`then` blocks and direct procedural flow over abstraction.
- Joins/clauses are vertically aligned with leading commas for selected columns in many queries.
- Query layout style (match existing files):
- Keep short `select` lists on one line; for wider lists, put one item per line with leading commas (`select a, b` then `, c`, `, d`).
- Put `from` on its own next line, aligned with the clause block (`select ...` then `from ...`).
- Put one `join` per line; keep short `join ... on ...` predicates on the same line.
- Start `where` on its own line; keep the first condition on that line and continue with `and` lines below it.
- Keep `group by` and `order by` on their own lines at the end of the query block.
- For nested queries/`union all`, indent each subquery block consistently and align closing `)` with the opening query level.
- Keep comments in SQL files lightweight and practical (`-- call ...`, short block examples).
