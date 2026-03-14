package main

import . "klec/lib/dec2"

type Segment_t struct {
	segment int
	name string
	code string
}

type Categ_t struct {
	categId int
	name string
	catsur int
	required int
	display int
	created string
	updated string
}

type YearVars_t struct {
	year int
	maxshare EuroFlat_t
	ltccap EuroFlat_t
	cover EuroCent_t
}

func (x YearVars_t)maxCover() EuroCent_t { return x.cover * 2 }

