package main

import "net/http"
import . "pm/lib/dec2"
import . "pm/lib/htmlHelper"
import . "pm/lib/wrapdb"

type Session_t struct {
	Name     string
	Path     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

type SimpleState_t struct {
	nickname string
	categ    int
}

type App_t struct {
	DB            *DB_t
	Port          string
	StaticVersion string
	Session       Session_t
	sessionState  map[string]SimpleState_t
	sessionFilters map[string]FilterState_t
	sessionCustomers map[string]CustomerState_t
	sessionEpoch map[string]int
	defaultYear int

	lookup struct {
		categs IdMap_t[Categ_t]
		segments IdMap_t[Segment_t]
		deductibles IdMap_t[[]EuroFlat_t]
		years IdMap_t[YearVars_t]

		hospitalLevels []SelectOption_t
		dentalLevels []SelectOption_t
		priorCoverOptions []SelectOption_t
		examOptions []SelectOption_t
		specialistOptions []SelectOption_t
	}

	selects struct {
		segment Elem_t
	}
}

var App App_t
