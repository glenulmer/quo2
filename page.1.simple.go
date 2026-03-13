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

func LoadSimplePageState(req *http.Request) (state SimpleState_t, categs []CategOption_t) {
	state, _ = App.SimpleStateGet(req)
	categs = append([]CategOption_t(nil), App.CategOptions...)
	state.categ = PickCateg(state.categ, categs)
	return
}

func PickCateg(wanted int, list []CategOption_t) int {
	for _, x := range list {
		if x.id == wanted { return wanted }
	}
	if len(list) > 0 { return list[0].id }
	return 0
}

func ParseSimpleFormInt(raw string) int {
	raw = Trim(raw)
	if raw == `` { return 0 }
	out, e := strconv.Atoi(raw)
	if e != nil { return 0 }
	return out
}

func SimplePage(w0 http.ResponseWriter, nickname string, categ int, categs []CategOption_t) {
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

func SimpleRecord(nickname string, categ int, categs []CategOption_t) Elem_t {
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

func CategSelect(name string, selected int, values []CategOption_t) Elem_t {
	sel := Select().Name(name).Id(name).Class(`input`)
	for _, x := range values {
		sel = sel.Wrap(Option().Value(x.id).Text(x.name))
	}
	return sel.SelO(selected)
}

func CategName(categ int, list []CategOption_t) string {
	for _, x := range list {
		if x.id == categ { return x.name }
	}
	return ``
}

func NicknamePreview(nickname string, categ int, categs []CategOption_t) Elem_t {
	show := Trim(Str(CategName(categ, categs), ` `, nickname))
	if show == `` { show = `(empty)` }
	return Div(
		Span(`Saved value: `).Class(`preview-label`),
		Bold(show).Class(`preview-value`),
	).Id(`nickname-preview`).Class(`preview`)
}
