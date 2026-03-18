package main 

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/date"
//	. "pm/lib/output"
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

func CustomerCard() Elem_t {
	today := CurrentDBDate()
	birth := DateFromYMD(today.Year()-32, 6, 15)
	buy := today.Days(40).ToWorkDay()
	// buyYear := buy.Year()

	body := Div().Class(`card-body`).Id(`Customer`).Wrap(
			Field(12).Wrap(ti().Name(`name`).Place(`Customer name`)),
			Field(`Segment`, 4).Wrap(Chooser(`quo_segments_chooser`).Name(`segment`)),
			Field(`Birth date`, 4).Wrap(di(`birth`, birth).Class(`right`).Value(birth.Hyphens())),
			Field(`Buy date`, 4).Wrap(di(`buy`, buy).Class(`right`).Value(buy.Hyphens())),
			Field(`Sick cover`, 4).Wrap(ni(`sickcover`, 75000, 0, 150000, 1000)),
			Field(`Prior cover`, 4).Wrap(Chooser(`quo_priorcov_chooser`).Name(`priorcov`)),
			Field(`Exam`,4).Wrap(Chooser(`quo_noexam_chooser`).Name(`exam`)),
		)
	title := Div(`Customer`).Class(`card-title`)
	return Div().Class(`card`).Wrap(title, body)
}
