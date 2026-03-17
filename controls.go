package main 

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/date"
	. "pm/lib/output"
)

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

func ni() Elem_t {
return Div().Class(`euro-wrap`).Wrap(
	Elem(`input`).Type(`number`))
}

func CustomerCard() Elem_t {
	today := CurrentDBDate()
	birth := DateFromYMD(today.Year()-32, 6, 15)
	buy := today.Days(40).ToWorkDay()
	// buyYear := buy.Year()

	body := Div().Class(`card-body`).Id(`Customer`).Wrap(
		Field(`Name`, 12).Wrap(ti().Name(`name`).Place(`Customer name`)),
		Field(`Birth date`, 6).Wrap(di(`birth`, birth).Value(birth.Hyphens())),
		Field(`Buy date`, 6).Wrap(di(`buy`, buy).Value(buy.Hyphens())),
		Field(`Sick Cover`, 4).Wrap(
			ni().Name(`sickcover`).Value(75000).KV(`min`, 0).KV(`max`, 150000).KV(`step`, `1000`),
			Span(`€`).Class(`euro-mark`),
		),
	)
	title := Div(`Customer`).Class(`card-title`)
	return Div().Class(`card`).Wrap(title, body)
}

func Card(title any, body ...any) Elem_t {
	return Div(
		Elem(`h2`).Class(`card-title`).Text(Str(title)),
		Div(body...).Class(`card-body`),
	).Class(`card`)
}
