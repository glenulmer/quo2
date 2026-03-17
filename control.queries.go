package main

import . "pm/lib/date"

func CurrentDBDate() CalDate_t {
	var ymd int
	pack := App.DB.CallRow(`quo_today_get`).Scan(&ymd)
	if pack.HasError() { panic(pack.Message()) }
	return CalDate(ymd)
}

