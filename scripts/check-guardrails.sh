#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

fail=0

note() { printf '[guardrails] %s\n' "$1"; }
bad()  { printf '[guardrails] FAIL: %s\n' "$1"; fail=1; }

note "checking: no panic in request/page handlers"
if rg -n 'panic\(' page.*.go >/tmp/quo2.guardrail.panic 2>/dev/null; then
	bad "panic(...) found in page.*.go"
	sed -n '1,120p' /tmp/quo2.guardrail.panic
fi

note "checking: no templates"
if rg -n 'html/template|text/template' . -g'*.go' -g'!go.sum' >/tmp/quo2.guardrail.templates 2>/dev/null; then
	bad "template package usage detected"
	sed -n '1,120p' /tmp/quo2.guardrail.templates
fi

note "checking: no JSON-first/app-layer framework imports"
if rg -n '"encoding/json"|"github.com/gin-gonic/gin"|"github.com/labstack/echo"|"gorm.io/gorm"|"github.com/jmoiron/sqlx"' . -g'*.go' -g'!go.sum' >/tmp/quo2.guardrail.imports 2>/dev/null; then
	bad "forbidden framework/import usage detected"
	sed -n '1,120p' /tmp/quo2.guardrail.imports
fi

note "checking: no context plumbing in handler signatures"
if rg -n 'context\.Context' page.*.go 0.main.go >/tmp/quo2.guardrail.context 2>/dev/null; then
	bad "context.Context found in page/main flow"
	sed -n '1,120p' /tmp/quo2.guardrail.context
fi

note "checking: no API/DTO-style file naming"
if rg -n '' . -g'*dto*.go' -g'*api*.go' -g'*controller*.go' -g'*service*.go' -g'*repository*.go' >/tmp/quo2.guardrail.names 2>/dev/null; then
	bad "forbidden layering-style filenames detected"
	sed -n '1,120p' /tmp/quo2.guardrail.names
fi

if [[ "$fail" -ne 0 ]]; then
	note "guardrail check failed"
	exit 1
fi

note "guardrail check passed"
