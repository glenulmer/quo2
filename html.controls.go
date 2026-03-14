package main

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
)

func Card(title any, body ...any) Elem_t {
	return Div(
		Elem(`h2`).Class(`card-title`).Text(Str(title)),
		Div(body...).Class(`card-body`),
	).Class(`card`)
}

func FormField(id, label string, control Elem_t, hint ...any) Elem_t {
	field := Elem(`label`).Class(`field`).KV(`for`, id).Wrap(
		Span(label).Class(`field-label`),
		Div(
			Wedge(),
			control,
		).Class(`field-control`),
	)
	if len(hint) > 0 {
		msg := Trim(Str(hint...))
		if msg != `` {
			field = field.Class(`has-error`).Data(`error`, msg).Wrap(Elem(`small`).Class(`field-error`).Text(msg))
		}
	}
	return field
}

func TextInput(name, value string) Elem_t {
	return TextIn().
		CutClass(`is-small`).
		CutClass(`kledit`).
		Class(`input`).
		Type(`text`).
		Name(name).
		Id(name).
		Value(value)
}

