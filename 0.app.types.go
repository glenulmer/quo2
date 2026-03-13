package main

import "net/http"
import . "klec/lib/wrapdb"

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
	CategOptions  []CategOption_t
	sessionState  map[string]SimpleState_t
	FilterLookups FilterLookups_t
	CustomerLookups CustomerLookups_t
	sessionFilters map[string]FilterState_t
	sessionCustomers map[string]CustomerState_t
	sessionEpoch map[string]int
}

var App App_t
