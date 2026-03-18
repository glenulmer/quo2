package main

import "net/http"
import . "pm/lib/wrapdb"
import . "pm/lib/dec2"

type Session_t struct {
	name     string
	path     string
	maxAge   int
	httpOnly bool
	secure   bool
	sameSite http.SameSite
}

type YAP_t struct { year, age, productId int }
func YAP(y, a, p int) YAP_t { return YAP_t{ y, a, p } }

type Price_t struct { base, surcharge EuroCent_t }

type App_t struct {
	DB            *DB_t
	port          string
	staticVersion string
	session       Session_t
	lookup struct {
		years IdMap_t[YearVars_t]
		categs IdMap_t[Categ_t]
		levels IdMap_t[Level_t]
		prices map[YAP_t]Price_t
	}
}

var App App_t
