package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	. "pm/lib/date"
	. "pm/lib/dec2"
	. "pm/lib/htmlHelper"
	. "pm/lib/output"
	. "pm/pkg.Global"
)

const postFiltersState = `/post/filters/state`
const postCustomerState = `/post/customer/state`
const postStateReset = `/post/state/reset`

func Page2FiltersGet(w0 http.ResponseWriter, req *http.Request) {
	sessionID := App.EnsureSession(w0, req)
	epoch := App.SessionEpochGet(sessionID)
	customerState := LoadCustomerPageState(req)
	filterState := LoadFiltersPageState(req)
	FiltersPage(w0, customerState, filterState, epoch)
}

func Page2CustomerPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sessionID := App.EnsureSession(w, req)
	epoch := App.SessionEpochGet(sessionID)
	if ParseFormInt(req.FormValue(`epoch`)) != epoch {
		state := LoadCustomerPageState(req)
		Rewrites(w, RewriteRow(`customer-record`, CustomerRecord(state, epoch)))
		return
	}
	state := LoadCustomerPageState(req)

	if _, ok := req.Form[`name`]; ok { state.name = NormalizeCustomerName(req.FormValue(`name`)) }
	if _, ok := req.Form[`birth`]; ok { state.birth = NormalizeDateInput(req.FormValue(`birth`)) }
	if _, ok := req.Form[`buy`]; ok { state.buy = NormalizeDateInput(req.FormValue(`buy`)) }
	if _, ok := req.Form[`cover`]; ok {
		if cover, ok := NormalizeCoverInput(req.FormValue(`cover`), CoverMaxValue()); ok {
			state.cover = cover
		}
	}
	if _, ok := req.Form[`segment`]; ok { state.segment = PickSegment(ParseFormInt(req.FormValue(`segment`)), state.segment) }
	if _, ok := req.Form[`vision`]; ok { state.vision = ParseFormBool(req.FormValue(`vision`)) }
	if _, ok := req.Form[`temp-visa`]; ok { state.tempVisa = ParseFormBool(req.FormValue(`temp-visa`)) }
	if _, ok := req.Form[`no-pvn`]; ok { state.noPVN = ParseFormBool(req.FormValue(`no-pvn`)) }
	if _, ok := req.Form[`natural-med`]; ok { state.naturalMed = ParseFormBool(req.FormValue(`natural-med`)) }

	App.CustomerStateSet(sessionID, state)
	Rewrites(w, RewriteRow(`customer-record`, CustomerRecord(state, epoch)))
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
	customerState := LoadCustomerPageState(req)
	if ParseFormInt(req.FormValue(`epoch`)) != epoch {
		state := LoadFiltersPageState(req)
		Rewrites(w, RewriteRow(`filters-record`, FiltersRecord(state, customerState, epoch)))
		return
	}
	state := LoadFiltersPageState(req)
	deductValues := CurrentDeductValues(customerState)

	if _, ok := req.Form[`deduct-min`]; ok { state.deductMin = PickEuroFlat(ParseFormInt(req.FormValue(`deduct-min`)), deductValues, state.deductMin) }
	if _, ok := req.Form[`deduct-max`]; ok { state.deductMax = PickEuroFlat(ParseFormInt(req.FormValue(`deduct-max`)), deductValues, state.deductMax) }
	if _, ok := req.Form[`hospital-min`]; ok { state.hospitalMin = PickOption(ParseFormInt(req.FormValue(`hospital-min`)), App.lookup.hospitalLevels, state.hospitalMin) }
	if _, ok := req.Form[`hospital-max`]; ok { state.hospitalMax = PickOption(ParseFormInt(req.FormValue(`hospital-max`)), App.lookup.hospitalLevels, state.hospitalMax) }
	if _, ok := req.Form[`dental-min`]; ok { state.dentalMin = PickOption(ParseFormInt(req.FormValue(`dental-min`)), App.lookup.dentalLevels, state.dentalMin) }
	if _, ok := req.Form[`dental-max`]; ok { state.dentalMax = PickOption(ParseFormInt(req.FormValue(`dental-max`)), App.lookup.dentalLevels, state.dentalMax) }
	if _, ok := req.Form[`prior-cover`]; ok { state.priorCover = PickOption(ParseFormInt(req.FormValue(`prior-cover`)), App.lookup.priorCoverOptions, state.priorCover) }
	if _, ok := req.Form[`exam`]; ok { state.exam = PickOption(ParseFormInt(req.FormValue(`exam`)), App.lookup.examOptions, state.exam) }
	if _, ok := req.Form[`specialist`]; ok { state.specialist = PickOption(ParseFormInt(req.FormValue(`specialist`)), App.lookup.specialistOptions, state.specialist) }

	if state.deductMin > state.deductMax { state.deductMin, state.deductMax = state.deductMax, state.deductMin }
	if state.hospitalMin > state.hospitalMax { state.hospitalMin, state.hospitalMax = state.hospitalMax, state.hospitalMin }
	if state.dentalMin > state.dentalMax { state.dentalMin, state.dentalMax = state.dentalMax, state.dentalMin }

	App.FilterStateSet(sessionID, state)
	Rewrites(w, RewriteRow(`filters-record`, FiltersRecord(state, customerState, epoch)))
}

func LoadCustomerPageState(req *http.Request) (state CustomerState_t) {
	state = DefaultCustomerState()
	if x, ok := App.CustomerStateGet(req); ok { state = NormalizeCustomerState(x) }
	return
}

func DefaultCustomerState() CustomerState_t {
	out := CustomerState_t{
		buy: CurrentDBDate(),
		cover: CoverDefaultValue(),
	}
	for _, segment := range App.lookup.segments.sort {
		if _, ok := App.lookup.segments.byId[segment]; ok {
			out.segment = segment
			break
		}
	}
	return out
}

func NormalizeCustomerState(state CustomerState_t) CustomerState_t {
	out := DefaultCustomerState()
	out.name = NormalizeCustomerName(state.name)
	if Valid(state.birth) { out.birth = state.birth }
	if Valid(state.buy) { out.buy = state.buy }
	if cover, ok := NormalizeCoverValue(state.cover, CoverMaxValue()); ok {
		out.cover = cover
	}
	out.segment = PickSegment(state.segment, out.segment)
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

func NormalizeDateInput(raw string) CalDate_t {
	out := Trim(raw)
	if len(out) > 10 { out = out[:10] }
	if out == `` { return 0 }
	return Parse(`yyyy-mm-dd`, out)
}

func DateInputValue(x CalDate_t) string {
	if !Valid(x) { return `` }
	return x.Format(`yyyy-mm-dd`)
}

func NormalizeCoverInput(raw string, max int) (EuroFlat_t, bool) {
	cover := OnlyDigits(Trim(raw))
	if cover < 0 { cover = 0 }
	value := EuroFlat_t(cover)
	if max > 0 && cover > max { return value, false }
	return value, true
}

func NormalizeCoverValue(raw EuroFlat_t, max int) (EuroFlat_t, bool) {
	value := raw
	if value < 0 { value = 0 }
	if max > 0 && int(value) > max { return value, false }
	return value, true
}

func CoverDisplayEuro(raw EuroFlat_t) string {
	return raw.OutEuro()
}

func CoverDefaultValue() EuroFlat_t {
	x, ok := App.lookup.years.byId[App.defaultYear]
	if !ok || x.cover <= 0 { return 0 }
	out := x.cover
	if out < 0 { return 0 }
	return out
}

func CoverMaxValue() int {
	x, ok := App.lookup.years.byId[App.defaultYear]
	if !ok { return 0 }
	out := int(x.maxCover())
	min := int(x.cover)
	if min < 0 { min = 0 }
	if out < min { return min }
	if out < 0 { return 0 }
	return out
}

func PickSegment(wanted, fallback int) int {
	if _, ok := App.lookup.segments.byId[wanted]; ok { return wanted }
	return fallback
}

func CurrentDeductValues(customer CustomerState_t) (out []EuroFlat_t) {
	ageCode := 0
	ageYears, ok := AgeYearsFromBirthDate(customer.birth, customer.buy)
	if !ok || IsAdultAgeYears(ageYears) { ageCode = 1 }
	if x, ok := App.lookup.deductibles.byId[ageCode]; ok {
		out = x
	}
	return
}

func AgeYearsFromBirthDate(birth, buy CalDate_t) (ageYears int, ok bool) {
	if !Valid(birth) || !Valid(buy) { return }
	birthDate := time.Date(birth.Year(), time.Month(birth.Month()), birth.Day(), 0, 0, 0, 0, time.UTC)
	buyDate := time.Date(buy.Year(), time.Month(buy.Month()), buy.Day(), 0, 0, 0, 0, time.UTC)
	if buyDate.Before(birthDate) { return }
	ageYears = buyDate.Year() - birthDate.Year()
	anniversary := time.Date(buyDate.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, buyDate.Location())
	if buyDate.Before(anniversary) { ageYears-- }
	if ageYears < 0 { return }
	ok = true
	return
}

func LoadFiltersPageState(req *http.Request) (state FilterState_t) {
	customerState := LoadCustomerPageState(req)
	state = DefaultFilterState(customerState)
	if x, ok := App.FilterStateGet(req); ok { state = NormalizeFilterState(x, customerState) }
	return
}

func DefaultFilterState(customer CustomerState_t) FilterState_t {
	out := FilterState_t{
		priorCover: FirstOptionID(App.lookup.priorCoverOptions),
		exam: FirstOptionID(App.lookup.examOptions),
		specialist: FirstOptionID(App.lookup.specialistOptions),
	}
	deductValues := CurrentDeductValues(customer)
	if len(deductValues) > 0 {
		out.deductMin = int(deductValues[0])
		out.deductMax = int(deductValues[len(deductValues)-1])
	}
	if len(App.lookup.hospitalLevels) > 0 {
		out.hospitalMin = App.lookup.hospitalLevels[0].id
		out.hospitalMax = App.lookup.hospitalLevels[len(App.lookup.hospitalLevels)-1].id
	}
	if len(App.lookup.dentalLevels) > 0 {
		out.dentalMin = App.lookup.dentalLevels[0].id
		out.dentalMax = App.lookup.dentalLevels[len(App.lookup.dentalLevels)-1].id
	}
	return out
}

func NormalizeFilterState(state FilterState_t, customer CustomerState_t) FilterState_t {
	out := DefaultFilterState(customer)
	deductValues := CurrentDeductValues(customer)
	out.deductMin = PickEuroFlat(state.deductMin, deductValues, out.deductMin)
	out.deductMax = PickEuroFlat(state.deductMax, deductValues, out.deductMax)
	out.hospitalMin = PickOption(state.hospitalMin, App.lookup.hospitalLevels, out.hospitalMin)
	out.hospitalMax = PickOption(state.hospitalMax, App.lookup.hospitalLevels, out.hospitalMax)
	out.dentalMin = PickOption(state.dentalMin, App.lookup.dentalLevels, out.dentalMin)
	out.dentalMax = PickOption(state.dentalMax, App.lookup.dentalLevels, out.dentalMax)
	out.priorCover = PickOption(state.priorCover, App.lookup.priorCoverOptions, out.priorCover)
	out.exam = PickOption(state.exam, App.lookup.examOptions, out.exam)
	out.specialist = PickOption(state.specialist, App.lookup.specialistOptions, out.specialist)
	if out.deductMin > out.deductMax { out.deductMin, out.deductMax = out.deductMax, out.deductMin }
	if out.hospitalMin > out.hospitalMax { out.hospitalMin, out.hospitalMax = out.hospitalMax, out.hospitalMin }
	if out.dentalMin > out.dentalMax { out.dentalMin, out.dentalMax = out.dentalMax, out.dentalMin }
	return out
}

func PickInt(wanted int, values []int, fallback int) int {
	for _, x := range values { if x == wanted { return x } }
	return fallback
}

func PickEuroFlat(wanted int, values []EuroFlat_t, fallback int) int {
	for _, x := range values {
		if int(x) == wanted { return wanted }
	}
	return fallback
}

func PickOption(wanted int, values []SelectOption_t, fallback int) int {
	for _, x := range values {
		if x.id == wanted { return wanted }
	}
	return fallback
}

func FirstOptionID(values []SelectOption_t) int {
	if len(values) == 0 { return 0 }
	return values[0].id
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

func FiltersPage(w0 http.ResponseWriter, customerState CustomerState_t, state FilterState_t, epoch int) {
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
					CustomerRecord(customerState, epoch),
				).Id(`customer-post`).Post(postCustomerState),
			Div(
				FiltersRecord(state, customerState, epoch),
			).Id(`filters-post`).Post(postFiltersState),
		),
		NL, head.Right(), NL,
	)
}

func CustomerRecord(state CustomerState_t, epoch int) Elem_t {
	return Div(
		Elem(`details`).Class(`ios-card`).KV(`open`).Id(`card-customer`).Wrap(
			Elem(`summary`).Class(`ios-title`).Wrap(
				Span(CustomerTitle(state.name)).Class(`ios-title-text`),
				Div(CustomerResetButton()).Class(`ios-title-right`),
			),
				Div(
					RenderCustomer(state),
				).Class(`ios-card-body`),
		),
	).Id(`customer-record`).Args(`state:1,epoch:`, epoch)
}

func CustomerResetButton() Elem_t {
	return Elem(`button`).KV(`type`, `button`).Name(`reset`).Id(`reset`).Class(`ios-reset`).Text(`Reset`)
}

func RenderCustomer(state CustomerState_t) Elem_t {
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
					KV(`data-orig`, DateInputValue(state.birth)).
					Value(DateInputValue(state.birth)),
			),
			IOSFormField(`buy`, `Buy date`,
				Elem(`input`).
					KV(`type`, `date`).
					Name(`buy`).
					Id(`buy`).
					Class(`ios-input`).
					KV(`data-orig`, DateInputValue(state.buy)).
					Value(DateInputValue(state.buy)),
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
						KV(`data-cover-max`, Str(CoverMaxValue())).
						KV(`data-orig`, Str(int(state.cover))).
						Value(CoverDisplayEuro(state.cover)),
				),
				IOSFormField(`segment`, `Segment`,
					SegmentSelect(state.segment),
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

func FiltersRecord(state FilterState_t, customer CustomerState_t, epoch int) Elem_t {
	return Div(
		Elem(`details`).Class(`ios-card`).KV(`open`).Id(`card-filters`).Wrap(
			Elem(`summary`).Class(`ios-title`).Wrap(
				Span(`Filters`).Class(`ios-title-text`),
			),
			Div(
				RenderFilters(state, customer),
			).Class(`ios-card-body`),
		),
	).Id(`filters-record`).Args(`state:1,epoch:`, epoch)
}

func RenderFilters(state FilterState_t, customer CustomerState_t) Elem_t {
	deductValues := CurrentDeductValues(customer)
	return Div(
		Div(
			IOSFormField(`deduct-min`, `Deductible Min`,
				DeductSelect(`deduct-min`, state.deductMin, deductValues),
			),
			IOSFormField(`deduct-max`, `Deductible Max`,
				DeductSelect(`deduct-max`, state.deductMax, deductValues),
			),
		).Class(`ios-row2`),
		Div(
			IOSFormField(`hospital-min`, `Hospital Level Min`,
				SelectFromOptions(`hospital-min`, state.hospitalMin, App.lookup.hospitalLevels),
			),
			IOSFormField(`hospital-max`, `Hospital Level Max`,
				SelectFromOptions(`hospital-max`, state.hospitalMax, App.lookup.hospitalLevels),
			),
		).Class(`ios-row2`),
		Div(
			IOSFormField(`dental-min`, `Dental Level Min`,
				SelectFromOptions(`dental-min`, state.dentalMin, App.lookup.dentalLevels),
			),
			IOSFormField(`dental-max`, `Dental Level Max`,
				SelectFromOptions(`dental-max`, state.dentalMax, App.lookup.dentalLevels),
			),
		).Class(`ios-row2`),
		Div(
			IOSFormField(`prior-cover`, `Prior Cover`,
				SelectFromOptions(`prior-cover`, state.priorCover, App.lookup.priorCoverOptions).Class(`ios-select-compact`),
			),
			IOSFormField(`exam`, `Exam`,
				SelectFromOptions(`exam`, state.exam, App.lookup.examOptions).Class(`ios-select-compact`),
			),
			IOSFormField(`specialist`, `Specialist`,
				SelectFromOptions(`specialist`, state.specialist, App.lookup.specialistOptions).Class(`ios-select-compact`),
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

func DeductSelect(name string, selected int, values []EuroFlat_t) Elem_t {
	sel := Select().Name(name).Id(name).Class(`ios-select`, `r`)
	for _, x := range values { sel = sel.Wrap(Option().Value(int(x)).Text(x.OutEuro())) }
	return sel.Choose(selected)
}

func SegmentSelect(selected int) Elem_t { return App.selects.segment.Clone().Choose(selected) }
