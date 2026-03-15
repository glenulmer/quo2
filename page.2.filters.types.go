package main

import . "pm/lib/htmlHelper"

type SelectOption_t struct {
	id   int
	name string
}

func SelectFromOptions(name string, selected int, options []SelectOption_t) Elem_t {
	sel := Select().Name(name).Id(name).Class(`ios-select`)
	for _, x := range options {
		sel = sel.Wrap(Option().Value(x.id).Text(x.name))
	}
	return sel.Choose(selected)
}

type FilterState_t struct {
	deductMin int
	deductMax int
	hospitalMin int
	hospitalMax int
	dentalMin int
	dentalMax int
	priorCover int
	exam int
	specialist int
}
