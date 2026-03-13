package main

import "strings"
import . "klec/lib/output"

const spCategsQuery = `quo_categs_query`

type CategOption_t struct {
	id   int
	name string
}

func LoadCategOptions() (list []CategOption_t) {
	rows := App.DB.Call(spCategsQuery)
	if rows.HasError() { panic(Error(`call `, spCategsQuery, ` failed: `, rows.Message())) }
	defer rows.Close()

	for rows.Next() {
		var x CategOption_t
		var catsur, required, display int
		var created, updated string
		e := rows.Scan(&x.id, &x.name, &catsur, &required, &display, &created, &updated)
		if e != nil { panic(Error(`scan `, spCategsQuery, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.id <= 0 || x.name == `` { continue }
		list = append(list, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spCategsQuery, ` failed: `, e)) }
	return
}

