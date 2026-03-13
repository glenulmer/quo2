package main

type CodebookOption_t struct {
	id   int
	name string
}

type LevelName_t struct {
	level int
	name  string
}

type FilterLookups_t struct {
	deductValues []int
	hospitalLevels []LevelName_t
	dentalLevels []LevelName_t
	priorCoverOptions []CodebookOption_t
	priorCoverDefault int
	priorCoverAllowed map[int]bool
	examOptions []CodebookOption_t
	examDefault int
	examAllowed map[int]bool
	specialistOptions []CodebookOption_t
	specialistDefault int
	specialistAllowed map[int]bool
}

type FilterState_t struct {
	deductMin int
	deductMax int
	hospitalMin int
	hospitalMax int
	dentalMin int
	dentalMax int
	priorCover int
	exam int
	specialist int
}

