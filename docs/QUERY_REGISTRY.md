# Query Registry

Purpose: single source-of-truth for all stored procedures called from Go.

## Rules

1. Every `App.DB.Call(...)` and `App.DB.CallRow(...)` target must be listed here.
2. Do not use unregistered proc names/symbols in code.
3. `Kind` is one of: `static`, `dynamic`.
4. `Allowed Layer` is one of: `bootstrap`, `page`, `any`, or `bootstrap|page`.
5. If `Kind=static`, prefer bootstrap cache + App field reads in handlers.

## Registry

| Symbol | DB Proc | Kind | Allowed Layer | Cache Field | Notes |
|---|---|---|---|---|---|
| `spCategsQuery` | `quo_categs_query` | `static` | `bootstrap` | `App.CategOptions` | category lookup loaded in bootstrap |
| `spSegmentsQuery` | `quo_segments_query` | `static` | `bootstrap` | `App.CustomerLookups.segments` | customer segment options |
| `spYearGet` | `klec_year_get` | `static` | `bootstrap` | `App.CustomerLookups.coverDefault` | default cover from active year row |
| `spCurrentDateQuery` | `quo_current_date_query` | `dynamic` | `bootstrap|page` | - | DB current date source for year/default-date logic |
| `spPriorCovQuery` | `klec_priorcov_query` | `static` | `bootstrap` | `App.FilterLookups.priorCoverOptions` | prior-cover codebook |
| `spReferralsQuery` | `klec_referrals_query` | `static` | `bootstrap` | `App.FilterLookups.specialistOptions` | referral codebook used to validate specialist options |
| `spLevelChooser` | `quo_level_chooser` | `static` | `bootstrap` | `App.FilterLookups.hospitalLevels` + `App.FilterLookups.dentalLevels` | level options for hospital/dental filters |
| `spPlanDeductiblesDistinct` | `plan_deductibles_distinct` | `static` | `bootstrap` | `App.lookup.deductibles` (`0=child`,`1=adult`) | deductible options grouped by child/adult flag |
