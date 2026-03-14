# Static Lookup Registry

Purpose: define DB lookups that are static for request-time behavior and must be cached at bootstrap.

## Usage Rule

1. If you are editing or adding `App.DB.Call(...)` / `App.DB.CallRow(...)`, read this file first.
2. If you are not touching DB calls, do not read this file.
3. Entries here must not be queried in `page.*.go` handlers.
4. Add new static entries before using them in page flow.

## Registry

| Symbol | DB Proc | Cache Field | Bootstrap Loader | Notes |
|---|---|---|---|---|
| `spCategsQuery` | `quo_categs_query` | `App.CategOptions` | `LoadCategOptions()` | category list used in simple page |
| `spSegmentsQuery` | `quo_segments_query` | `App.CustomerLookups.segments` | `LoadCustomerLookups()` | customer segment options |
| `spYearGet` | `klec_year_get` | `App.CustomerLookups.coverDefault` | `LoadCustomerLookups()` | default customer sick-cover from active year |
| `spPriorCovQuery` | `klec_priorcov_query` | `App.FilterLookups.priorCoverOptions` | `LoadFilterLookups()` | prior-cover filter options |
| `spReferralsQuery` | `klec_referrals_query` | `App.FilterLookups.specialistOptions` | `LoadFilterLookups()` | specialist option validation against referral codes |
| `spLevelChooser` | `quo_level_chooser` | `App.FilterLookups.hospitalLevels` + `App.FilterLookups.dentalLevels` | `LoadFilterLookups()` | hospital/dental level options |
| `spPlanDeductiblesDistinct` | `plan_deductibles_distinct` | `App.lookup.deductibles` (`0=child`,`1=adult`) | `LoadStaticData()` | deductible options grouped by child/adult flag |
