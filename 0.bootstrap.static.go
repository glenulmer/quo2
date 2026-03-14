package main

import (
	"database/sql"
	"strings"
	. "pm/lib/output"
	. "pm/lib/dec2"
)

func LoadStaticData() {
	App.lookup.categs = LoadCategIdMap()
	App.lookup.segments = LoadSegmentIdMap()
	App.lookup.years = LoadYearVarsIdMap()
	App.lookup.deductibles = DeductiblesIdMap()

	App.lookup.filters = LoadFilterLookups()
}

const spCategsQuery = `quo_categs_query`

func LoadCategIdMap() IdMap_t[Categ_t] {
	out := IdMap[Categ_t]()

	rows := App.DB.Call(spCategsQuery)
	if rows.HasError() { panic(Error(`call `, spCategsQuery, ` failed: `, rows.Message())) }
	defer rows.Close()

	for rows.Next() {
		var x Categ_t
		e := rows.Scan(&x.categId, &x.name, &x.catsur, &x.required, &x.display, &x.created, &x.updated)
		if e != nil { panic(Error(`scan `, spCategsQuery, ` failed: `, e)) }

		x.name = strings.TrimSpace(x.name)
		if x.categId <= 0 || x.name == `` { continue }

		out = out.Add(x.categId, x)
	}
	if rows.HasError() { panic(Error(`rows `, spCategsQuery, ` failed: `, rows.Message())) }

	return out
}

func LoadSegmentIdMap() IdMap_t[Segment_t] {
	out := IdMap[Segment_t]()

	rows := App.DB.Call(spSegmentsQuery)
	if rows.HasError() { panic(Error(`call `, spSegmentsQuery, ` failed: `, rows.Message())) }
	defer rows.Close()

	for rows.Next() {
		var x Segment_t
		e := rows.Scan(&x.segment, &x.name, &x.code)
		if e != nil { panic(Error(`scan `, spSegmentsQuery, ` failed: `, e)) }

		x.name = strings.TrimSpace(x.name)
		x.code = strings.TrimSpace(x.code)
		if x.segment <= 0 || x.name == `` { continue }

		out = out.Add(x.segment, x)
	}
	if rows.HasError() { panic(Error(`rows `, spSegmentsQuery, ` failed: `, rows.Message())) }

	return out
}

func LoadYearVarsIdMap() IdMap_t[YearVars_t] {
	out := IdMap[YearVars_t]()

	rows := App.DB.Call(spYearGet, 0)
	if rows.HasError() { panic(Error(`call `, spYearGet, ` failed: `, rows.Message())) }
	defer rows.Close()

	for rows.Next() {
		var x YearVars_t
		var exists, isPast bool
		var coverCents EuroCent_t
		e := rows.Scan(&x.year, &x.maxshare, &coverCents, &x.ltccap, &exists, &isPast, new(sql.NullString))
		if e != nil { panic(Error(`scan `, spYearGet, ` failed: `, e)) }
		x.cover = EuroFlatFromCent(coverCents)
		if isPast || !exists || x.year <= 0 { continue }
		out = out.Add(x.year, x)
	}
	if rows.HasError() { panic(Error(`rows `, spYearGet, ` failed: `, rows.Message())) }
	if len(out.sort) == 0 { panic(Error(`empty static lookup from `, spYearGet)) }
	App.defaultYear = out.sort[len(out.sort)-1]

	return out
}

func DeductiblesIdMap() IdMap_t[[]EuroFlat_t] {
	out := IdMap[[]EuroFlat_t]()
	var lists [2][]EuroFlat_t

	rows := App.DB.Call(spPlanDeductiblesDistinct)
	if rows.HasError() { panic(Error(`call `, spPlanDeductiblesDistinct, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var isAdult int
		var cents EuroCent_t
		e := rows.Scan(&isAdult, &cents)
		if e != nil { panic(Error(`scan `, spPlanDeductiblesDistinct, ` failed: `, e)) }
		if isAdult != 0 && isAdult != 1 {
			panic(Error(`invalid adult flag from `, spPlanDeductiblesDistinct, `: `, isAdult))
		}
		value := EuroFlatFromCent(cents)
		if value < 0 { continue }
		lists[isAdult] = append(lists[isAdult], value)
	}
	if rows.HasError() { panic(Error(`rows `, spPlanDeductiblesDistinct, ` failed: `, rows.Message())) }
	if len(lists[0]) == 0 || len(lists[1]) == 0 {
		panic(Error(`empty static lookup from `, spPlanDeductiblesDistinct))
	}
	return out.Add(0, lists[0]).Add(1, lists[1])
}
