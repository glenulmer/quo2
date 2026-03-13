package main

import (
	"net/http"
	"strconv"
	"strings"

	. "klec/lib/dec2"
	. "klec/lib/htmlHelper"
	. "klec/lib/output"
	. "klec/pkg.Global"
)

const postFiltersState = `/post/filters/state`
const postCustomerState = `/post/customer/state`
const postStateReset = `/post/state/reset`

func Page2FiltersGet(w0 http.ResponseWriter, req *http.Request) {
	sessionID := App.EnsureSession(w0, req)
	epoch := App.SessionEpochGet(sessionID)
	customerState, customerLookups := LoadCustomerPageState(req)
	filterState, filterLookups := LoadFiltersPageState(req)
	FiltersPage(w0, customerState, customerLookups, filterState, filterLookups, epoch)
}

func Page2CustomerPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sessionID := App.EnsureSession(w, req)
	epoch := App.SessionEpochGet(sessionID)
	if ParseFormInt(req.FormValue(`epoch`)) != epoch {
		state, lookups := LoadCustomerPageState(req)
		Rewrites(w, RewriteRow(`customer-record`, CustomerRecord(state, lookups, epoch)))
		return
	}
	state, lookups := LoadCustomerPageState(req)

	if _, ok := req.Form[`name`]; ok { state.name = NormalizeCustomerName(req.FormValue(`name`)) }
	if _, ok := req.Form[`birth`]; ok { state.birth = NormalizeDateInput(req.FormValue(`birth`)) }
	if _, ok := req.Form[`buy`]; ok { state.buy = NormalizeDateInput(req.FormValue(`buy`)) }
	if _, ok := req.Form[`cover`]; ok {
		if cover, ok := NormalizeCoverInput(req.FormValue(`cover`), lookups.coverMax); ok {
			state.cover = cover
		}
	}
	if _, ok := req.Form[`segment`]; ok { state.segment = PickCodebook(ParseFormInt(req.FormValue(`segment`)), lookups.segmentAllowed, state.segment) }
	if _, ok := req.Form[`vision`]; ok { state.vision = ParseFormBool(req.FormValue(`vision`)) }
	if _, ok := req.Form[`temp-visa`]; ok { state.tempVisa = ParseFormBool(req.FormValue(`temp-visa`)) }
	if _, ok := req.Form[`no-pvn`]; ok { state.noPVN = ParseFormBool(req.FormValue(`no-pvn`)) }
	if _, ok := req.Form[`natural-med`]; ok { state.naturalMed = ParseFormBool(req.FormValue(`natural-med`)) }

	App.CustomerStateSet(sessionID, state)
	Rewrites(w, RewriteRow(`customer-record`, CustomerRecord(state, lookups, epoch)))
}

func Page2StateResetPost(w http.ResponseWriter, req *http.Request) {
	sessionID := App.EnsureSession(w, req)
	App.CustomerStateClear(sessionID)
	App.FilterStateClear(sessionID)
	App.SessionEpochBump(sessionID)
	w.WriteHeader(http.StatusOK)
}

func Page2FiltersPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sessionID := App.EnsureSession(w, req)
	epoch := App.SessionEpochGet(sessionID)
	if ParseFormInt(req.FormValue(`epoch`)) != epoch {
		state, lookups := LoadFiltersPageState(req)
		Rewrites(w, RewriteRow(`filters-record`, FiltersRecord(state, lookups, epoch)))
		return
	}
	state, lookups := LoadFiltersPageState(req)

	if _, ok := req.Form[`deduct-min`]; ok { state.deductMin = PickInt(ParseFormInt(req.FormValue(`deduct-min`)), lookups.deductValues, state.deductMin) }
	if _, ok := req.Form[`deduct-max`]; ok { state.deductMax = PickInt(ParseFormInt(req.FormValue(`deduct-max`)), lookups.deductValues, state.deductMax) }
	if _, ok := req.Form[`hospital-min`]; ok { state.hospitalMin = PickLevel(ParseFormInt(req.FormValue(`hospital-min`)), lookups.hospitalLevels, state.hospitalMin) }
	if _, ok := req.Form[`hospital-max`]; ok { state.hospitalMax = PickLevel(ParseFormInt(req.FormValue(`hospital-max`)), lookups.hospitalLevels, state.hospitalMax) }
	if _, ok := req.Form[`dental-min`]; ok { state.dentalMin = PickLevel(ParseFormInt(req.FormValue(`dental-min`)), lookups.dentalLevels, state.dentalMin) }
	if _, ok := req.Form[`dental-max`]; ok { state.dentalMax = PickLevel(ParseFormInt(req.FormValue(`dental-max`)), lookups.dentalLevels, state.dentalMax) }
	if _, ok := req.Form[`prior-cover`]; ok { state.priorCover = PickCodebook(ParseFormInt(req.FormValue(`prior-cover`)), lookups.priorCoverAllowed, state.priorCover) }
	if _, ok := req.Form[`exam`]; ok { state.exam = PickCodebook(ParseFormInt(req.FormValue(`exam`)), lookups.examAllowed, state.exam) }
	if _, ok := req.Form[`specialist`]; ok { state.specialist = PickCodebook(ParseFormInt(req.FormValue(`specialist`)), lookups.specialistAllowed, state.specialist) }

	if state.deductMin > state.deductMax { state.deductMin, state.deductMax = state.deductMax, state.deductMin }
	if state.hospitalMin > state.hospitalMax { state.hospitalMin, state.hospitalMax = state.hospitalMax, state.hospitalMin }
	if state.dentalMin > state.dentalMax { state.dentalMin, state.dentalMax = state.dentalMax, state.dentalMin }

	App.FilterStateSet(sessionID, state)
	Rewrites(w, RewriteRow(`filters-record`, FiltersRecord(state, lookups, epoch)))
}

func LoadCustomerPageState(req *http.Request) (state CustomerState_t, lookups CustomerLookups_t) {
	lookups = CopyCustomerLookups(App.CustomerLookups)
	state = DefaultCustomerState(lookups)
	if x, ok := App.CustomerStateGet(req); ok { state = NormalizeCustomerState(x, lookups) }
	return
}

func CopyCustomerLookups(x CustomerLookups_t) CustomerLookups_t {
	x.segments = append([]SegmentOption_t(nil), x.segments...)
	if x.segmentAllowed != nil {
		a := make(map[int]bool, len(x.segmentAllowed))
		for k, v := range x.segmentAllowed { a[k] = v }
		x.segmentAllowed = a
	}
	return x
}

func DefaultCustomerState(x CustomerLookups_t) CustomerState_t {
	return CustomerState_t{
		buy:     TodayISODate(),
		cover:   x.coverDefault,
		segment: x.segmentDefault,
	}
}

func NormalizeCustomerState(state CustomerState_t, x CustomerLookups_t) CustomerState_t {
	out := DefaultCustomerState(x)
	out.name = NormalizeCustomerName(state.name)
	out.birth = NormalizeDateInput(state.birth)
	out.buy = NormalizeDateInput(state.buy)
	if cover, ok := NormalizeCoverInput(state.cover, x.coverMax); ok {
		out.cover = cover
	}
	out.segment = PickCodebook(state.segment, x.segmentAllowed, out.segment)
	out.vision = state.vision
	out.tempVisa = state.tempVisa
	out.noPVN = state.noPVN
	out.naturalMed = state.naturalMed
	return out
}

func NormalizeCustomerName(raw string) string {
	out := Trim(raw)
	if len(out) > 100 { out = out[:100] }
	return out
}

func NormalizeDateInput(raw string) string {
	out := Trim(raw)
	if len(out) > 10 { out = out[:10] }
	return out
}

func NormalizeCoverInput(raw string, max int) (string, bool) {
	cover := OnlyDigits(Trim(raw))
	if cover < 0 { cover = 0 }
	if max > 0 && cover > max { return Str(cover), false }
	return Str(cover), true
}

func CoverDisplayEuro(raw string) string {
	return EuroFlat_t(OnlyDigits(Trim(raw))).OutEuro()
}

func TodayISODate() string {
	return QueryCurrentDateISO()
}

func LoadFiltersPageState(req *http.Request) (state FilterState_t, lookups FilterLookups_t) {
	lookups = CopyFilterLookups(App.FilterLookups)
	state = DefaultFilterState(lookups)
	if x, ok := App.FilterStateGet(req); ok { state = NormalizeFilterState(x, lookups) }
	return
}

func CopyFilterLookups(x FilterLookups_t) FilterLookups_t {
	x.deductValues = append([]int(nil), x.deductValues...)
	x.hospitalLevels = append([]LevelName_t(nil), x.hospitalLevels...)
	x.dentalLevels = append([]LevelName_t(nil), x.dentalLevels...)
	x.priorCoverOptions = append([]CodebookOption_t(nil), x.priorCoverOptions...)
	x.examOptions = append([]CodebookOption_t(nil), x.examOptions...)
	x.specialistOptions = append([]CodebookOption_t(nil), x.specialistOptions...)
	return x
}

func DefaultFilterState(x FilterLookups_t) FilterState_t {
	out := FilterState_t{
		priorCover: x.priorCoverDefault,
		exam: x.examDefault,
		specialist: x.specialistDefault,
	}
	if len(x.deductValues) > 0 {
		out.deductMin = x.deductValues[0]
		out.deductMax = x.deductValues[len(x.deductValues)-1]
	}
	if len(x.hospitalLevels) > 0 {
		out.hospitalMin = x.hospitalLevels[0].level
		out.hospitalMax = x.hospitalLevels[len(x.hospitalLevels)-1].level
	}
	if len(x.dentalLevels) > 0 {
		out.dentalMin = x.dentalLevels[0].level
		out.dentalMax = x.dentalLevels[len(x.dentalLevels)-1].level
	}
	return out
}

func NormalizeFilterState(state FilterState_t, x FilterLookups_t) FilterState_t {
	out := DefaultFilterState(x)
	out.deductMin = PickInt(state.deductMin, x.deductValues, out.deductMin)
	out.deductMax = PickInt(state.deductMax, x.deductValues, out.deductMax)
	out.hospitalMin = PickLevel(state.hospitalMin, x.hospitalLevels, out.hospitalMin)
	out.hospitalMax = PickLevel(state.hospitalMax, x.hospitalLevels, out.hospitalMax)
	out.dentalMin = PickLevel(state.dentalMin, x.dentalLevels, out.dentalMin)
	out.dentalMax = PickLevel(state.dentalMax, x.dentalLevels, out.dentalMax)
	out.priorCover = PickCodebook(state.priorCover, x.priorCoverAllowed, out.priorCover)
	out.exam = PickCodebook(state.exam, x.examAllowed, out.exam)
	out.specialist = PickCodebook(state.specialist, x.specialistAllowed, out.specialist)
	if out.deductMin > out.deductMax { out.deductMin, out.deductMax = out.deductMax, out.deductMin }
	if out.hospitalMin > out.hospitalMax { out.hospitalMin, out.hospitalMax = out.hospitalMax, out.hospitalMin }
	if out.dentalMin > out.dentalMax { out.dentalMin, out.dentalMax = out.dentalMax, out.dentalMin }
	return out
}

func PickInt(wanted int, values []int, fallback int) int {
	for _, x := range values { if x == wanted { return x } }
	return fallback
}

func PickLevel(wanted int, values []LevelName_t, fallback int) int {
	for _, x := range values { if x.level == wanted { return x.level } }
	return fallback
}

func PickCodebook(wanted int, allowed map[int]bool, fallback int) int {
	if allowed[wanted] { return wanted }
	return fallback
}

func ParseFormInt(raw string) int {
	raw = Trim(raw)
	if raw == `` { return 0 }
	out, e := strconv.Atoi(raw)
	if e != nil { return 0 }
	return out
}

func ParseFormBool(raw string) bool {
	switch strings.ToLower(Trim(raw)) {
	case `1`, `true`, `on`, `yes`, `y`:
		return true
	}
	return false
}

func Bool01(x bool) string {
	if x { return `1` }
	return `0`
}

func FiltersPage(w0 http.ResponseWriter, customerState CustomerState_t, customerLookups CustomerLookups_t, state FilterState_t, lookups FilterLookups_t, epoch int) {
	head := Head().
		CSS(Str(`/static/css/phone.quote.css?v=`, App.StaticVersion)).
		JSTail(Str(`/static/js/validate.js?v=`, App.StaticVersion)).
		JSTail(Str(`/static/js/2.filters.js?v=`, App.StaticVersion)).
		Title(`Filters - Quo2`).
		End()

	w := Writer(w0)
	w.Add(
		head.Left(), NL,
		Elem(`main`).Class(`ios-page`).Wrap(
			Div(
				CustomerRecord(customerState, customerLookups, epoch),
			).Id(`customer-post`).Post(postCustomerState),
			Div(
				FiltersRecord(state, lookups, epoch),
			).Id(`filters-post`).Post(postFiltersState),
		),
		NL, head.Right(), NL,
	)
}

func CustomerRecord(state CustomerState_t, lookups CustomerLookups_t, epoch int) Elem_t {
	return Div(
		Elem(`details`).Class(`ios-card`).KV(`open`).Id(`card-customer`).Wrap(
			Elem(`summary`).Class(`ios-title`).Wrap(
				Span(CustomerTitle(state.name)).Class(`ios-title-text`),
				Div(CustomerResetButton()).Class(`ios-title-right`),
			),
			Div(
				RenderCustomer(state, lookups),
			).Class(`ios-card-body`),
		),
	).Id(`customer-record`).Args(`state:1,epoch:`, epoch)
}

func CustomerResetButton() Elem_t {
	return Elem(`button`).KV(`type`, `button`).Name(`reset`).Id(`reset`).Class(`ios-reset`).Text(`Reset`)
}

func RenderCustomer(state CustomerState_t, lookups CustomerLookups_t) Elem_t {
	return Div(
		IOSFormField(`name`, `Name`,
			Elem(`input`).
				KV(`type`, `text`).
				Name(`name`).
				Id(`name`).
				Class(`ios-input`).
				KV(`maxlength`, `100`).
				KV(`data-orig`, state.name).
				Value(state.name),
		),
		Div(
			IOSFormField(`birth`, `Birth date`,
				Elem(`input`).
					KV(`type`, `date`).
					Name(`birth`).
					Id(`birth`).
					Class(`ios-input`).
					KV(`data-orig`, state.birth).
					Value(state.birth),
			),
			IOSFormField(`buy`, `Buy date`,
				Elem(`input`).
					KV(`type`, `date`).
					Name(`buy`).
					Id(`buy`).
					Class(`ios-input`).
					KV(`data-orig`, state.buy).
					Value(state.buy),
			),
		).Class(`ios-row2`, `ios-row-dates`),
		Div(
			IOSFormFieldWedge(`cover`, `Sick Cover`, `bar-left`,
				Elem(`input`).
					KV(`type`, `text`).
					Name(`cover`).
					Id(`cover`).
					Class(`ios-input`, `r`).
					KV(`maxlength`, `16`).
					KV(`data-cover-max`, Str(lookups.coverMax)).
					KV(`data-orig`, state.cover).
					Value(CoverDisplayEuro(state.cover)),
			),
			IOSFormField(`segment`, `Segment`,
				SegmentSelect(state.segment, lookups.segments),
			),
		).Class(`ios-row2`, `ios-row3`),
		Div(
			CustomerCheckCell(`vision`, `Vision`, state.vision),
			CustomerCheckCell(`temp-visa`, `Temp Visa`, state.tempVisa),
			CustomerCheckCell(`no-pvn`, `No PVN`, state.noPVN),
			CustomerCheckCell(`natural-med`, `Natural Med`, state.naturalMed),
		).Class(`ios-row4`, `ios-row-checks`, `customer-checks-row`),
	).Class(`ios-stack`)
}

func CustomerCheckCell(name, label string, checked bool) Elem_t {
	return Div(
		CustomerCheckBox(name, label, checked),
	).Class(`customer-check-cell`)
}

func CustomerCheckBox(name, label string, checked bool) Elem_t {
	in := Elem(`input`).
		KV(`type`, `checkbox`).
		Name(name).
		Id(name).
		Class(`ios-check-input`).
		KV(`data-orig`, Bool01(checked))
	if checked { in = in.KV(`checked`) }
	return Elem(`label`).Class(`ios-check`).KV(`for`, name).Wrap(
		in,
		Span(label).Class(`ios-check-label`),
	)
}

func CustomerTitle(name string) string {
	name = Trim(name)
	if name != `` { return name }
	return `Customer`
}

func FiltersRecord(state FilterState_t, lookups FilterLookups_t, epoch int) Elem_t {
	return Div(
		Elem(`details`).Class(`ios-card`).KV(`open`).Id(`card-filters`).Wrap(
			Elem(`summary`).Class(`ios-title`).Wrap(
				Span(`Filters`).Class(`ios-title-text`),
			),
			Div(
				RenderFilters(state, lookups),
			).Class(`ios-card-body`),
		),
	).Id(`filters-record`).Args(`state:1,epoch:`, epoch)
}

func RenderFilters(state FilterState_t, lookups FilterLookups_t) Elem_t {
	return Div(
		Div(
			IOSFormField(`deduct-min`, `Deductible Min`,
				DeductSelect(`deduct-min`, state.deductMin, lookups.deductValues),
			),
			IOSFormField(`deduct-max`, `Deductible Max`,
				DeductSelect(`deduct-max`, state.deductMax, lookups.deductValues),
			),
		).Class(`ios-row2`),
		Div(
			IOSFormField(`hospital-min`, `Hospital Level Min`,
				LevelSelect(`hospital-min`, state.hospitalMin, lookups.hospitalLevels),
			),
			IOSFormField(`hospital-max`, `Hospital Level Max`,
				LevelSelect(`hospital-max`, state.hospitalMax, lookups.hospitalLevels),
			),
		).Class(`ios-row2`),
		Div(
			IOSFormField(`dental-min`, `Dental Level Min`,
				LevelSelect(`dental-min`, state.dentalMin, lookups.dentalLevels),
			),
			IOSFormField(`dental-max`, `Dental Level Max`,
				LevelSelect(`dental-max`, state.dentalMax, lookups.dentalLevels),
			),
		).Class(`ios-row2`),
		Div(
			IOSFormField(`prior-cover`, `Prior Cover`,
				CodebookSelect(`prior-cover`, state.priorCover, lookups.priorCoverOptions).Class(`ios-select-compact`),
			),
			IOSFormField(`exam`, `Exam`,
				CodebookSelect(`exam`, state.exam, lookups.examOptions).Class(`ios-select-compact`),
			),
			IOSFormField(`specialist`, `Specialist`,
				CodebookSelect(`specialist`, state.specialist, lookups.specialistOptions).Class(`ios-select-compact`),
			),
		).Class(`ios-row3f`, `ios-row-reg-top`),
	).Class(`ios-stack`)
}

func IOSFormField(id, label string, control Elem_t) Elem_t {
	return Elem(`label`).Class(`ios-field`).KV(`for`, id).Wrap(
		Span(label),
		Div(control).Class(`ios-control`),
	)
}

func IOSFormFieldWedge(id, label, sideClass string, control Elem_t) Elem_t {
	return Elem(`label`).Class(`ios-field`).KV(`for`, id).Wrap(
		Span(label),
		Div(
			Div(
				Div().Class(`wedge`),
				control,
			).Class(`ios-control-wedge`, sideClass),
		).Class(`ios-control`),
	)
}

func DeductSelect(name string, selected int, values []int) Elem_t {
	sel := Select().Name(name).Id(name).Class(`ios-select`, `r`)
	for _, x := range values { sel = sel.Wrap(Option().Value(x).Text(EuroFlatFromCent(EuroCent_t(x)).OutEuro())) }
	return sel.SelO(selected)
}

func SegmentSelect(selected int, values []SegmentOption_t) Elem_t {
	sel := Select().Name(`segment`).Id(`segment`).Class(`ios-select`).KV(`data-orig`, Str(selected))
	for _, x := range values { sel = sel.Wrap(Option().Value(x.segment).Text(x.name)) }
	return sel.SelO(selected)
}

func LevelSelect(name string, selected int, values []LevelName_t) Elem_t {
	sel := Select().Name(name).Id(name).Class(`ios-select`)
	for _, x := range values { sel = sel.Wrap(Option().Value(x.level).Text(x.name)) }
	return sel.SelO(selected)
}

func CodebookSelect(name string, selected int, values []CodebookOption_t) Elem_t {
	sel := Select().Name(name).Id(name).Class(`ios-select`)
	for _, x := range values { sel = sel.Wrap(Option().Value(x.id).Text(x.name)) }
	return sel.SelO(selected)
}
