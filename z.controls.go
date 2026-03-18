package main 

import (
	. "pm/lib/htmlHelper"
	. "pm/lib/date"
//	. "pm/lib/output"
)

func CustomerCard() Elem_t {
	today := CurrentDBDate()
	birth := DateFromYMD(today.Year()-32, 6, 15)
	buy := today.Days(40).ToWorkDay()
	// buyYear := buy.Year()
	adult := true
	max := true

	body := Div().Class(`card-body`).Id(`Customer`).Wrap(
		Field(8).Wrap(ti().Name(`custName`).Place(`Customer name`)),
		Field(4).Wrap(Chooser(`quo_segments_chooser`).Name(`segment`)),

		Field(`Birth date`, 4).Wrap(di(`birth`, birth).Class(`right`).Value(birth.Hyphens())),
		Field(`Buy date`, 4).Wrap(di(`buy`, buy).Class(`right`).Value(buy.Hyphens())),
		Field(`Sick cover`, 4).Wrap(ni(`sickCover`, 75000, 0, 150000, 1000)),

		Field(`Prior cover`, 4).Wrap(Chooser(`quo_priorcov_chooser`).Name(`priorCov`)),
		Field(`Exam`,4).Wrap(Chooser(`quo_noexam_chooser`).Name(`exam`)),
		Field(`Specialist`,4).Wrap(Chooser(`quo_specialist_chooser`).Name(`specialist`)),

		Field(1).Wrap(`&nbsp;`),
		Field(10).Wrap(
			Div().Class(`check-grid`).Wrap(
				CheckCell(`vision`, `Vision`, false),
				CheckCell(`tempVisa`, `Temp Visa`, false),
				CheckCell(`noPVN`, `No PVN`, false),
				CheckCell(`naturalMed`, `Natural Med`, false),
			),
		),
		Field(1).Wrap(`&nbsp;`),

		Field(`Deductible`, 4).Wrap(Chooser(`quo_deductibles_chooser`, adult, !max).Name(`deductibleMin`).Class(`right`)),
		Field(`Hospital`, 4).Wrap(Chooser(`quo_level_chooser_max`, 3, !max).Name(`hospitalMin`)), // use categ
		Field(`Dental`, 4).Wrap(Chooser(`quo_level_chooser_max`, 4, !max).Name(`dentalMin`)), // use categ


		Field(4).Wrap(Chooser(`quo_deductibles_chooser`, adult, max).Name(`deductibleMax`).Class(`right`)),
		Field(4).Wrap(Chooser(`quo_level_chooser_max`, 3, max).Name(`hospitalMax`)), // use categ
		Field(4).Wrap(Chooser(`quo_level_chooser_max`, 4, max).Name(`dentalMax`)), // use categ

	)

	title := Div(`Quote Information`).Class(`card-title`)
	return Div().Class(`card`).Wrap(title, body)
}

/*
*/