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
