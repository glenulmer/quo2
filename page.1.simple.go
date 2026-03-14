package main

import (
	"net/http"
	"strconv"

	. "klec/lib/htmlHelper"
	. "klec/lib/output"
	. "klec/pkg.Global"
)

const postSimpleState = `/post/simple/state`

func Page1SimpleGet(w0 http.ResponseWriter, req *http.Request) {
	App.EnsureSession(w0, req)
	state, categs := LoadSimplePageState(req)
	SimplePage(w0, state.nickname, state.categ, categs)
}

func Page1SimplePost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	sessionID := App.EnsureSession(w, req)
	state, categs := LoadSimplePageState(req)

	if _, ok := req.Form[`nickname`]; ok {
		state.nickname = Trim(req.FormValue(`nickname`))
	}
	if _, ok := req.Form[`categ`]; ok {
		state.categ = PickCateg(ParseSimpleFormInt(req.FormValue(`categ`)), categs)
	}
	App.SimpleStateSet(sessionID, state)

	Rewrites(w, RewriteRow(`simple-record`, SimpleRecord(state.nickname, state.categ, categs)))
}

func LoadSimplePageState(req *http.Request) (state SimpleState_t, categs IdMap_t[Categ_t]) {
	state, _ = App.SimpleStateGet(req)
	categs = App.lookup.categs
	state.categ = PickCateg(state.categ, categs)
	return
}

func PickCateg(wanted int, idMap IdMap_t[Categ_t]) int {
	if _, ok := idMap.byId[wanted]; ok { return wanted }
	if len(idMap.sort) > 0 {
		return idMap.sort[0]
	}
	return 0
}

func ParseSimpleFormInt(raw string) int {
	raw = Trim(raw)
	if raw == `` { return 0 }
	out, e := strconv.Atoi(raw)
	if e != nil { return 0 }
	return out
}

func SimplePage(w0 http.ResponseWriter, nickname string, categ int, categs IdMap_t[Categ_t]) {
	head := Head().
		CSS(Str(`/static/css/phone.quote.css?v=`, App.StaticVersion)).
		JSTail(Str(`/static/js/validate.js?v=`, App.StaticVersion)).
		JSTail(Str(`/static/js/1.simple.js?v=`, App.StaticVersion)).
		Title(`Quo2 - Simple Session`).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		Elem(`main`).Class(`simple-page`).Wrap(
				Div(
					Card(`Simple Control`,
						P(`Type a nickname and press Save.`).Class(`muted`),
						SimpleRecord(nickname, categ, categs),
					),
				).Id(`simple-post`).Post(postSimpleState),
			),
		NL, head.Right(), NL,
	)
}

func SimpleRecord(nickname string, categ int, categs IdMap_t[Categ_t]) Elem_t {
	return Div(
		FormField(`nickname`, `Nickname`,
			TextInput(`nickname`, nickname).KV(`maxlength`, `40`),
		),
		FormField(`categ`, `Category`,
			CategSelect(`categ`, categ, categs),
		),
		NicknamePreview(nickname, categ, categs),
		Div(
			Elem(`button`).KV(`type`, `button`).Name(`create`).Class(`save-btn`).Text(`Save`),
		).Class(`actions`),
	).Id(`simple-record`).Args(`state:1`)
}

func CategSelect(name string, selected int, idMap IdMap_t[Categ_t]) Elem_t {
	sel := Select().Name(name).Id(name).Class(`input`)
	for _, id := range idMap.sort {
		x, ok := idMap.byId[id]
		if !ok { continue }
		sel = sel.Wrap(Option().Value(x.categId).Text(x.name))
	}
	return sel.SelO(selected)
}

func CategName(categ int, idMap IdMap_t[Categ_t]) string {
	x, ok := idMap.byId[categ]
	if ok { return x.name }
	return ``
}

func NicknamePreview(nickname string, categ int, categs IdMap_t[Categ_t]) Elem_t {
	show := Trim(Str(CategName(categ, categs), ` `, nickname))
	if show == `` { show = `(empty)` }
	return Div(
		Span(`Saved value: `).Class(`preview-label`),
		Bold(show).Class(`preview-value`),
	).Id(`nickname-preview`).Class(`preview`)
}
