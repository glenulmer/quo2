package main

import . "klec/lib/output"

const spSegmentsQuery = `quo_segments_query`
const spYearGet = `klec_year_get`
const spCurrentDateQuery = `quo_current_date_query`

func QueryCurrentDateISO() string {
	rows := App.DB.Call(spCurrentDateQuery)
	if rows.HasError() { panic(Error(`call `, spCurrentDateQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	if !rows.Next() { panic(Error(spCurrentDateQuery, ` returned no rows`)) }
	var out string
	e := rows.Scan(&out)
	if e != nil { panic(Error(`scan `, spCurrentDateQuery, ` failed: `, e)) }
	out = Trim(out)
	if len(out) >= 10 { out = out[:10] }
	if len(out) != 10 { panic(Error(`invalid db date from `, spCurrentDateQuery, `: `, out)) }
	return out
}
