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
	App.lookup.plans = LoadPlanAlpha()
	App.lookup.products = LoadProducts()
	App.lookup.prices = LoadPrices()
	App.lookup.filters = LoadFilters()
	App.lookup.planAddons = LoadPlanAddons()
}

type Categ_t struct {
	categId CategId_t
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
		if rows.HasError() { panic(rows.Message()) }
		out = out.Add(int(x.categId), x)
	}
	return out
}

type LevelId_t int
type Level_t struct {
	levelId LevelId_t
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
		if rows.HasError() { panic(rows.Message()) }
		out = out.Add(int(x.levelId), x)
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
		if rows.HasError() { panic(rows.Message()) }
		if isPast || !exists || x.year <= 0 { continue }
		out = out.Add(x.year, x)
	}
	return out
}

func LoadPrices() map[YAP_t]Price_t {
	prices := make(map[YAP_t]Price_t)

	rows := App.DB.Call(`quo_product_prices`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var yap YAP_t
	var pr Price_t
	for rows.Next() {
		rows.Scan(&yap.year, &yap.age, &yap.productId, &pr.base, &pr.surcharge)
		if rows.HasError() { panic(rows.Message()) }
		prices[yap] = pr
	}

	return prices
}

type ProductId_t int
type PlanId_t ProductId_t
type AddonId_t ProductId_t
type CategId_t int

func LoadProducts() map[ProductId_t]Product_t {
	products := make(map[ProductId_t]Product_t)
	rows := App.DB.Call(`quo_product_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var p Product_t
	for rows.Next() {
		rows.Scan(&p.productId, &p.providerId, &p.name, &p.categ, &p.level, &p.segs)
		if rows.HasError() { panic(rows.Message()) }
		products[p.productId] = p
	}
	return products
}

type Filters_t struct {
	plan PlanId_t
	segmask, priorcov, noexam, referral int
    vis_dec2 int
	vis_pct bool
	tempvisa bool
    hospital, dental int
	ad_value, ch_value EuroFlat_t
}

func LoadFilters() map[PlanId_t]Filters_t {
	filters := make(map[PlanId_t]Filters_t)
	rows := App.DB.Call(`quo_plan_filters_query`)
	if rows.HasError() { panic(rows.Message()) }
	defer rows.Close()
	var p Filters_t
	for rows.Next() {
		rows.Scan(
			&p.plan, &p.segmask, &p.priorcov, &p.noexam, &p.referral, &p.vis_pct, &p.vis_dec2,
			&p.tempvisa, &p.hospital, &p.dental, &p.ad_value, &p.ch_value,
		)
		if rows.HasError() { panic(rows.Message()) }
		filters[p.plan] = p
	}
	return filters
}

type PlanCateg_t struct {
	plan	PlanId_t
	categ	CategId_t
}

type CatChoice_t struct {
	addon	AddonId_t
	level	int
	isdef	bool
	label	string
}

func LoadPlanAddons() map[PlanCateg_t]CatChoice_t {
	choices := make(map[PlanCateg_t]CatChoice_t)

	rows := App.DB.Call(`quo_plan_categ_addons`)
	defer rows.Close()
	if rows.HasError() { panic(rows.Message()) }
	var k PlanCateg_t
	var v CatChoice_t 
	for rows.Next() {
		rows.Scan(&k.plan, &k.categ, &v.addon, &v.level, &v.isdef, &v.label)
		if rows.HasError() { panic(rows.Message()) }
		choices[k] = v
	}

	return choices
}

type PlanName_t struct { plan PlanId_t; exact_age bool; name string }
func LoadPlanAlpha() []PlanName_t {
	var plans []PlanName_t

	rows := App.DB.Call(`quo_plan_name_query`)
	defer rows.Close()
	if rows.HasError() { panic(rows.Message()) }
	var k PlanName_t
	for rows.Next() {
		rows.Scan(&k.plan, &k.exact_age, &k.name)
		if rows.HasError() { panic(rows.Message()) }
		plans = append(plans, k)
	}

	return plans
}

