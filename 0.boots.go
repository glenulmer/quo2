package main

import (
	"database/sql"
	"strings"
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
	. "pm/lib/dec2"
)

func LoadStaticData() {
	App.lookup.categs = LoadCategIdMap()
	App.lookup.segments = LoadSegmentIdMap()
	App.lookup.years = LoadYearVarsIdMap()
	App.lookup.deductibles = DeductiblesIdMap()
	LoadFilterControlLookups()
}

func LoadSelectElements() {
	App.selects.segment = SegmentSelectElem()
}

func SegmentSelectElem() Elem_t {
	optionSlice := make([]Elem_t, 0, len(App.lookup.segments.sort))
	for _, id := range App.lookup.segments.sort {
		x, ok := App.lookup.segments.byId[id]
		if !ok { continue }
		optionSlice = append(optionSlice, Option().KV(`value`, x.segment).Text(x.name))
	}
	return Select(optionSlice...).Name(`segment`).Id(`segment`).Class(`ios-select`)
}

func LoadFilterControlLookups() {
	App.lookup.priorCoverOptions = QueryPriorCoverOptions()

	App.lookup.examOptions = []SelectOption_t{
		{id: 0, name: `Exam OK`},
		{id: examNoExamCode, name: `No exam`},
	}

	referralOptions := QueryReferralOptions()
	App.lookup.specialistOptions = []SelectOption_t{
		{id: specialistAnyCode, name: `Not important`},
		{id: 1, name: `Always referral`},
		{id: 0, name: `No referral`},
	}
	for _, x := range App.lookup.specialistOptions {
		if x.id == specialistAnyCode { continue }
		if !OptionContainsID(referralOptions, x.id) {
			panic(Error(`invalid codebook from `, spReferralsQuery, `: missing referral code `, x.id))
		}
	}

	hospitalCateg := CategIDByName(App.lookup.categs, `hospital`)
	dentalCateg := CategIDByName(App.lookup.categs, `dental`)
	if hospitalCateg <= 0 || dentalCateg <= 0 {
		panic(Error(`missing hospital/dental category IDs from `, spCategsQuery))
	}
	App.lookup.hospitalLevels = QueryLevelChooser(hospitalCateg)
	App.lookup.dentalLevels = QueryLevelChooser(dentalCateg)

	if len(App.lookup.priorCoverOptions) == 0 { panic(Error(`empty static lookup from `, spPriorCovQuery)) }
	if len(App.lookup.hospitalLevels) == 0 { panic(Error(`empty static lookup from `, spLevelChooser, ` for hospital category`)) }
	if len(App.lookup.dentalLevels) == 0 { panic(Error(`empty static lookup from `, spLevelChooser, ` for dental category`)) }
}

func OptionContainsID(options []SelectOption_t, wanted int) bool {
	for _, x := range options {
		if x.id == wanted { return true }
	}
	return false
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
