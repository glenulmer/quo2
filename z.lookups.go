package main

import (
	"database/sql"
	. "pm/lib/dec2"
//	. "pm/lib/output"
)

func (x YearVars_t)maxCover() EuroFlat_t { return x.cover * 2 }

func LoadStaticData() {
	App.lookup.years = LoadYearVarsIdMap()
	App.lookup.categs = LoadCategIdMap()
	App.lookup.levels = LoadLevelIdMap()
}

type Categ_t struct {
	categId int
	name string
	catsur int
	required int
}

func LoadCategIdMap() IdMap_t[Categ_t] {
	out := IdMap[Categ_t]()
	rows := App.DB.Call(`quo_categs_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x Categ_t
		rows.Scan(&x.categId, &x.name, &x.catsur, &x.required)
		out = out.Add(x.categId, x)
	}
	return out
}

type Level_t struct {
	levelId int
	label string
	categId int
	segments int // bitmask
	canStack bool
}

func LoadLevelIdMap() IdMap_t[Level_t] {
	out := IdMap[Level_t]()
	rows := App.DB.Call(`klec_levels_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x Level_t
		rows.Scan(&x.levelId, &x.label, &x.categId, &x.segments, &x.canStack)
		out = out.Add(x.levelId, x)
	}
	return out
}

type YearVars_t struct {
	year int
	maxshare EuroFlat_t
	cover EuroFlat_t
	ltccap EuroFlat_t
}

func LoadYearVarsIdMap() IdMap_t[YearVars_t] {
	out := IdMap[YearVars_t]()

	rows := App.DB.Call(`klec_year_get` , 0)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	for rows.Next() {
		var x YearVars_t
		var exists, isPast bool
		rows.Scan(&x.year, &x.maxshare, &x.cover, &x.ltccap, &exists, &isPast, new(sql.NullString))
		if isPast || !exists || x.year <= 0 { continue }
		out = out.Add(x.year, x)
	}
	return out
}
