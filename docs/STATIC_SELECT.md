# STATIC_SELECT Procedure

## Purpose
Define the standard process for converting select controls to startup-built static controls stored on `App.selects`.

## Core Rules
1. Build static select controls once at startup, never per request.
2. `Bootstrap()` must call `LoadSelectElements()` immediately after `LoadStaticData()`.
3. Every static select builder lives in `0.boots.go` and uses `pm/lib/htmlHelper`.
4. Members of `App.selects` are immutable after startup:
no request/page/render function may reassign, clear, or rebuild any `App.selects.*` member.
5. Render code must read from `App.selects.<name>`, then clone before mutation:
`App.selects.<name>.Clone().Choose(selected)`.
6. Do not call select builder functions from render/request paths.
7. For controls that do not use validation, do not depend on `data-orig` and do not use `SelO(...)`.

## Implementation Steps
1. Add a field to `App_t.selects` in `0.app.types.go`, typed as `Elem_t`.
2. In `0.boots.go`, add:
`LoadSelectElements()` assignment for the new field.
`<Control>SelectElem()` builder function.
3. In the builder function:
iterate static lookup data (`App.lookup.*`), build `[]Elem_t` options, then return:
`Select(optionSlice...).Name("<form-name>").Id("<id>").Class("<class>")`
4. In `0.bootstrap.go`, ensure `LoadSelectElements()` is called once after `LoadStaticData()`.
5. Replace page-level/manual select construction with:
`App.selects.<name>.Clone().Choose(selected)`
6. Remove any render-path fallback that rebuilds the select.

## Validation Guidance
1. Unless explicitly specified otherwise, assume a select control has no validator.
2. If the select is not validated, do not register a validator for that control name.
3. Keep validator `data-orig` setup scoped to registered controls only.
4. Use `.Clone().Choose(...)` for non-validated selects.
5. Use `.SelO(...)` only when a control explicitly requires validator-origin tracking.

## Review Checklist
1. `rg "LoadSelectElements|App.selects|SelectElem" *.go`
2. `rg "SelO\\(|data-orig" page.*.go static/js/*.js`
3. Confirm no render/request function assigns `App.selects.*`.
4. Run `./scripts/check-all.sh`.
