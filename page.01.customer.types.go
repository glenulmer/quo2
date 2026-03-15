package main

import . "pm/lib/dec2"
import . "pm/lib/date"

type CustomerState_t struct {
	name       string
	birth      CalDate_t
	buy        CalDate_t
	cover      EuroFlat_t
	segment    int
	vision     bool
	tempVisa   bool
	noPVN      bool
	naturalMed bool
}
