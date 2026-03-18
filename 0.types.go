package main

import "net/http"
import . "pm/lib/wrapdb"

type Session_t struct {
	name     string
	path     string
	maxAge   int
	httpOnly bool
	secure   bool
	sameSite http.SameSite
}

type App_t struct {
	DB            *DB_t
	port          string
	staticVersion string
	session       Session_t
	lookup struct {
		years IdMap_t[YearVars_t]
		categs IdMap_t[Categ_t]
		levels IdMap_t[Level_t]
	}
}

var App App_t
