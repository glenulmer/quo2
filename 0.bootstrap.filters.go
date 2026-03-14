package main

import "strings"
import . "pm/lib/output"

const spPriorCovQuery = `klec_priorcov_query`
const spReferralsQuery = `klec_referrals_query`
const spLevelChooser = `quo_level_chooser`
const spPlanDeductiblesDistinct = `plan_deductibles_distinct`
const specialistAnyCode = 2
const examNoExamCode = 1

func LoadFilterLookups() FilterLookups_t {
	out := FilterLookups_t{}

	out.priorCoverOptions = QueryPriorCoverOptions()
	out.priorCoverAllowed, out.priorCoverDefault = ResolveCodebookOptions(out.priorCoverOptions)

	out.examOptions = []CodebookOption_t{
		{id: 0, name: `Exam OK`},
		{id: examNoExamCode, name: `No exam`},
	}
	out.examAllowed, out.examDefault = ResolveCodebookOptions(out.examOptions)

	referralOptions := QueryReferralOptions()
	referralAllowed, _ := ResolveCodebookOptions(referralOptions)
	out.specialistOptions = []CodebookOption_t{
		{id: specialistAnyCode, name: `Not important`},
		{id: 1, name: `Always referral`},
		{id: 0, name: `No referral`},
	}
	out.specialistAllowed, out.specialistDefault = ResolveCodebookOptions(out.specialistOptions)
	for _, x := range out.specialistOptions {
		if x.id == specialistAnyCode { continue }
		if !referralAllowed[x.id] {
			panic(Error(`invalid codebook from `, spReferralsQuery, `: missing referral code `, x.id))
		}
	}

	hospitalCateg := CategIDByName(App.lookup.categs, `hospital`)
	dentalCateg := CategIDByName(App.lookup.categs, `dental`)
	if hospitalCateg <= 0 || dentalCateg <= 0 {
		panic(Error(`missing hospital/dental category IDs from quo_categs_query`))
	}
	out.hospitalLevels = QueryLevelChooser(hospitalCateg)
	out.dentalLevels = QueryLevelChooser(dentalCateg)

	ValidateFilterLookupsLoaded(out)

	return out
}

func ValidateFilterLookupsLoaded(x FilterLookups_t) {
	if len(x.priorCoverOptions) == 0 { panic(Error(`empty static lookup from `, spPriorCovQuery)) }
	if len(x.hospitalLevels) == 0 { panic(Error(`empty static lookup from `, spLevelChooser, ` for hospital category`)) }
	if len(x.dentalLevels) == 0 { panic(Error(`empty static lookup from `, spLevelChooser, ` for dental category`)) }
}

func CategIDByName(idMap IdMap_t[Categ_t], name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, id := range idMap.sort {
		x, ok := idMap.byId[id]
		if !ok { continue }
		if strings.ToLower(strings.TrimSpace(x.name)) == name { return x.categId }
	}
	return 0
}

func QueryPriorCoverOptions() (list []CodebookOption_t) {
	rows := App.DB.Call(spPriorCovQuery)
	if rows.HasError() { panic(Error(`call `, spPriorCovQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x CodebookOption_t
		e := rows.Scan(&x.id, &x.name)
		if e != nil { panic(Error(`scan `, spPriorCovQuery, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.name == `` { continue }
		list = append(list, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spPriorCovQuery, ` failed: `, e)) }
	return
}

func QueryReferralOptions() (list []CodebookOption_t) {
	rows := App.DB.Call(spReferralsQuery)
	if rows.HasError() { panic(Error(`call `, spReferralsQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x CodebookOption_t
		e := rows.Scan(&x.id, &x.name)
		if e != nil { panic(Error(`scan `, spReferralsQuery, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.name == `` { continue }
		list = append(list, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spReferralsQuery, ` failed: `, e)) }
	return
}

func ResolveCodebookOptions(list []CodebookOption_t) (allowed map[int]bool, defaultID int) {
	allowed = make(map[int]bool, len(list))
	for k, x := range list {
		allowed[x.id] = true
		if k == 0 { defaultID = x.id }
	}
	return
}

func QueryLevelChooser(categ int) (levels []LevelName_t) {
	rows := App.DB.Call(spLevelChooser, categ)
	if rows.HasError() { panic(Error(`call `, spLevelChooser, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x LevelName_t
		e := rows.Scan(&x.level, &x.name)
		if e != nil { panic(Error(`scan `, spLevelChooser, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		if x.level <= 0 || x.name == `` { continue }
		levels = append(levels, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spLevelChooser, ` failed: `, e)) }
	return
}
