package main

const AgeChildMaxYears = 19
const DeductibleLookupAgeYears = 28

func IsChildAgeYears(ageYears int) bool {
	return ageYears <= AgeChildMaxYears
}

func IsAdultAgeYears(ageYears int) bool {
	return !IsChildAgeYears(ageYears)
}
