package main

import (
	. "pm/lib/date"
	. "pm/lib/htmlHelper"
)

func CurrentDBDate() CalDate_t {
	var ymd int
	App.DB.CallRow(`quo_today_get`).Scan(&ymd)
	// if pack.HasError() { panic(pack.Message()) }
	return CalDate(ymd)
}

func Chooser(sp string, args ...any) Elem_t {
	var options []Elem_t
	var x struct { id int; label string }
	rows := App.DB.Call(sp, args...)
	// if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&x.id, &x.label)
		options = append(options, Option().Text(x.label).KV(`value`, x.id))
	}
	return Select(options)
}

func di(name string, date CalDate_t) Elem_t {
	return DateInput().
		Name(name).
		Req().
		KV(`inputmode`, `none`).
		KV(`onclick`, `if(this.showPicker){this.showPicker();}`).
		KV(`onfocus`, `if(this.showPicker){this.showPicker();}`).
		KV(`onkeydown`, `return false`).
		KV(`onpaste`, `return false`).
		KV(`ondrop`, `return false`).
		Value(date.Format(`yyyy-mm-dd`))
}

func ti() Elem_t { return Elem(`input`).KV(`autocomplete`,`off`).KV(`type`,`text`) }

func ni(name string, value, min, max, step int) Elem_t {
	return Div().Class(`euro-wrap`).Wrap(
		Elem(`input`).
			Type(`number`).
			Name(name).
			Value(value).
			KV(`min`, min).
			KV(`max`, max).
			KV(`step`, step).
			Class(`right`),
			Span(`€`).Class(`euro-mark`),
	)
}

func CheckCell(name, text string, varBool ...bool) Elem_t {
	var checked bool
	if len(varBool) > 0 { checked = varBool[0] }
	return Div().Class(`check-cell`, `center`).Wrap(
		CBox(name, checked),
		Span(text).Class(`check-text`),
	)
}
