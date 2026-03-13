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

query_registry="docs/QUERY_REGISTRY.md"
note "checking: query registry exists"
if [[ ! -f "$query_registry" ]]; then
	bad "missing $query_registry"
fi

query_tsv="/tmp/quo2.guardrail.query_registry.tsv"
awk -F'|' '
	/^\|/ {
		sym=$2; proc=$3; kind=$4; layer=$5;
		gsub(/`/, "", sym); gsub(/`/, "", proc); gsub(/`/, "", kind); gsub(/`/, "", layer);
		gsub(/^[ \t]+|[ \t]+$/, "", sym);
		gsub(/^[ \t]+|[ \t]+$/, "", proc);
		gsub(/^[ \t]+|[ \t]+$/, "", kind);
		gsub(/^[ \t]+|[ \t]+$/, "", layer);
		if (sym=="" || sym=="Symbol" || sym ~ /^-+$/) next;
		if (proc=="" || proc=="DB Proc" || proc ~ /^-+$/) next;
		if (kind=="" || kind=="Kind" || kind ~ /^-+$/) next;
		if (layer=="" || layer=="Allowed Layer" || layer ~ /^-+$/) next;
		print sym "\t" proc "\t" kind "\t" layer;
	}
' "$query_registry" >"$query_tsv"

note "checking: all DB calls are registered + layer-allowed"
hit_query=0
: >/tmp/quo2.guardrail.querycalls
while IFS= read -r line; do
	file="${line%%:*}"
	rest="${line#*:}"
	lineno="${rest%%:*}"
	text="${rest#*:}"

	token="$(printf '%s\n' "$text" | sed -nE 's/.*App\.DB\.(Call|CallRow)\([[:space:]]*(`[^`]+`|"[^"]+"|[A-Za-z_][A-Za-z0-9_]*)[[:space:]]*(,|\)).*/\2/p')"
	if [[ -z "$token" ]]; then
		printf '%s:%s: unparseable App.DB.Call/CallRow first arg\n' "$file" "$lineno" >>/tmp/quo2.guardrail.querycalls
		hit_query=1
		continue
	fi

	entry=""
	if [[ "$token" == \`* || "$token" == \"* ]]; then
		proc="${token:1:${#token}-2}"
		if ! entry="$(awk -F'\t' -v proc="$proc" '$2==proc {print; found=1} END{if(!found) exit 1}' "$query_tsv" 2>/dev/null)"; then
			printf '%s:%s: unregistered proc literal: %s\n' "$file" "$lineno" "$proc" >>/tmp/quo2.guardrail.querycalls
			hit_query=1
			continue
		fi
	else
		sym="$token"
		if ! entry="$(awk -F'\t' -v sym="$sym" '$1==sym {print; found=1} END{if(!found) exit 1}' "$query_tsv" 2>/dev/null)"; then
			printf '%s:%s: unregistered proc symbol: %s\n' "$file" "$lineno" "$sym" >>/tmp/quo2.guardrail.querycalls
			hit_query=1
			continue
		fi
	fi

	allowed_layer="$(printf '%s\n' "$entry" | awk -F'\t' '{print $4}')"
	base="$(basename "$file")"
	layer="other"
	case "$base" in
		0.bootstrap*.go) layer="bootstrap" ;;
		page.*.go) layer="page" ;;
	esac

	if [[ "$layer" != "other" ]]; then
		case "$allowed_layer" in
			any) ;;
			bootstrap|page) ;;
			bootstrap|page|any) ;;
			*)
				if [[ "$allowed_layer" != "$layer" && "$allowed_layer" != *"|$layer" && "$allowed_layer" != "$layer|"* && "$allowed_layer" != *"|$layer|"* ]]; then
					printf '%s:%s: layer %s not allowed by registry (%s)\n' "$file" "$lineno" "$layer" "$allowed_layer" >>/tmp/quo2.guardrail.querycalls
					hit_query=1
				fi
				;;
		esac
	fi
done < <(rg -n 'App\.DB\.(Call|CallRow)\(' . -g'*.go' -g'!go.sum')
if [[ "$hit_query" -ne 0 ]]; then
	bad "query registry violations found"
	sed -n '1,220p' /tmp/quo2.guardrail.querycalls
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
