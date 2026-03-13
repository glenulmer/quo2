package main

import (
	"database/sql"
	"strconv"
	"strings"
)
import . "klec/lib/output"

const spSegmentsQuery = `quo_segments_query`
const spYearGet = `klec_year_get`
const spCurrentDateQuery = `quo_current_date_query`

func LoadCustomerLookups() CustomerLookups_t {
	out := CustomerLookups_t{}
	out.segments = QuerySegmentOptions()
	coverDefault := QueryDefaultCoverFlat()
	out.coverDefault = Str(coverDefault)
	out.coverMax = coverDefault * 2
	if out.coverMax < coverDefault { out.coverMax = coverDefault }
	out.segmentAllowed, out.segmentDefault = ResolveSegmentOptions(out.segments)
	if len(out.segments) == 0 { panic(Error(`empty static lookup from `, spSegmentsQuery)) }
	return out
}

func QuerySegmentOptions() (list []SegmentOption_t) {
	rows := App.DB.Call(spSegmentsQuery)
	if rows.HasError() { panic(Error(`call `, spSegmentsQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var x SegmentOption_t
		e := rows.Scan(&x.segment, &x.name, &x.code)
		if e != nil { panic(Error(`scan `, spSegmentsQuery, ` failed: `, e)) }
		x.name = strings.TrimSpace(x.name)
		x.code = strings.TrimSpace(x.code)
		if x.segment <= 0 || x.name == `` { continue }
		list = append(list, x)
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spSegmentsQuery, ` failed: `, e)) }
	return
}

func ResolveSegmentOptions(list []SegmentOption_t) (allowed map[int]bool, defaultSegment int) {
	allowed = make(map[int]bool, len(list))
	for k, x := range list {
		allowed[x.segment] = true
		if k == 0 { defaultSegment = x.segment }
	}
	return
}

func QueryDefaultCoverFlat() int {
	currYear := QueryCurrentYearDB()
	latestYear := QueryLatestExistingYear(currYear)
	if latestYear < currYear { latestYear = currYear }
	return QueryYearCoverFlat(latestYear)
}

func QueryCurrentDateISO() string {
	rows := App.DB.Call(spCurrentDateQuery)
	if rows.HasError() { panic(Error(`call `, spCurrentDateQuery, ` failed: `, rows.Message())) }
	defer rows.Close()
	if !rows.Next() { panic(Error(spCurrentDateQuery, ` returned no rows`)) }
	var out string
	e := rows.Scan(&out)
	if e != nil { panic(Error(`scan `, spCurrentDateQuery, ` failed: `, e)) }
	out = Trim(out)
	if len(out) >= 10 { out = out[:10] }
	if len(out) != 10 { panic(Error(`invalid db date from `, spCurrentDateQuery, `: `, out)) }
	return out
}

func QueryCurrentYearDB() int {
	day := QueryCurrentDateISO()
	y, e := strconv.Atoi(day[:4])
	if e != nil || y <= 0 { panic(Error(`invalid db year from `, spCurrentDateQuery, `: `, day)) }
	return y
}

func QueryLatestExistingYear(currYear int) int {
	out := 0
	rows := App.DB.Call(spYearGet, 0)
	if rows.HasError() { panic(Error(`call `, spYearGet, ` failed: `, rows.Message())) }
	defer rows.Close()
	for rows.Next() {
		var year, shared, coverCents, ltccap int
		var exists, isPast bool
		var created sql.NullString
		e := rows.Scan(&year, &shared, &coverCents, &ltccap, &exists, &isPast, &created)
		if e != nil { panic(Error(`scan `, spYearGet, ` failed: `, e)) }
		if exists && year > out { out = year }
	}
	if e := rows.Err(); e != nil { panic(Error(`rows `, spYearGet, ` failed: `, e)) }
	if out == 0 { return currYear }
	return out
}

func QueryYearCoverFlat(year int) int {
	var out struct {
		year, shared, coverCents, ltccap int
		exists, isPast                   bool
		created                          sql.NullString
	}
	row := App.DB.CallRow(spYearGet, year).Scan(
		&out.year, &out.shared, &out.coverCents, &out.ltccap,
		&out.exists, &out.isPast, &out.created,
	)
	if row.HasError() { panic(Error(`call `, spYearGet, ` failed: `, row.Message())) }
	if !out.exists { panic(Error(spYearGet, ` returned no active row for year `, year)) }
	if out.coverCents < 0 { panic(Error(spYearGet, ` returned negative cover for year `, year)) }
	return out.coverCents / 100
}
