package main

import (
	"flag"
	"github.com/VirgileD/pmu-go/pmulibs"
	"time"
)

func main() {
	var date string
	const layout = "2006-01-02"
	var now = time.Now()
	flag.StringVar(&date, "date", now.Format(layout), "help message for date")
	flag.Parse()
	var pStats2 = pmulibs.PStats{}
	now = now.AddDate(0, 0, -1)
	now = now.AddDate(0, 0, -1)
	for {
		var pStats3 = pmulibs.GetStats(now.Format(layout), true)
		pStats2 = pmulibs.AddStats(pStats2, pStats3)
		pmulibs.ApplyStats(pStats2, pmulibs.GetCourse(now.Format(layout)))
		now = now.AddDate(0, 0, -1)
	}
	/*
		var pStats = getStats("2013-09-04", false)
		fmt.Println(pStats.Stats["paris-courses_com"])

		var pStats2 = getStats("2013-09-05", true)

		fmt.Println(pStats2.Stats["paris-courses_com"])

		var pStats3 = addStats(pStats, pStats2)
		fmt.Println(pStats3.Stats["paris-courses_com"])
		fmt.Println(pStats3.Nbr)
		fmt.Println(pStats3.Date)
		applyStats(pStats3, getCourse("2013-09-05"))
	*/
}
