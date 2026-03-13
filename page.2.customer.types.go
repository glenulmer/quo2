package main

type SegmentOption_t struct {
	segment int
	name    string
	code    string
}

type CustomerLookups_t struct {
	segments []SegmentOption_t
	coverDefault string
	coverMax int
	segmentDefault int
	segmentAllowed map[int]bool
}

type CustomerState_t struct {
	name       string
	birth      string
	buy        string
	cover      string
	segment    int
	vision     bool
	tempVisa   bool
	noPVN      bool
	naturalMed bool
}
