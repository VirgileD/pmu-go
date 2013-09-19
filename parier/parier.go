package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/VirgileD/pmu-go/pmulibs"
	"time"
)

var _ = pmulibs.CalcStats
var _ = fmt.Println

func main() {
	var date string
	const layout = "2006-01-02"
	var now = time.Now()
	flag.StringVar(&date, "date", now.Format(layout), "help message for date")
	flag.Parse()
	now = now.AddDate(0, 0, -1)
	now = now.AddDate(0, 0, -1)
	//b := []byte(`{"date":"` + now.Format(layout) + `"}`)
	b := []byte(`{"gains.m7": { "$gt": 30, "$lt": 100} }`)

	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(pmulibs.GetCourses(f))
	courses := pmulibs.GetCourses(f)
	stats := pmulibs.CalcStats(courses)
	pmulibs.ApplyStats(courses[0], stats)
	pmulibs.ApplyStats(courses[1], stats)
	//fmt.Println(pmulibs.GetMeansAndVariances())
	/*
		var pStats2 = pmulibs.PStats{}
		for {
			var pStats3 = pmulibs.GetStats(now.Format(layout), true)
			if pStats3.Nbr != 0 {
				pStats2 = pmulibs.AddStats(pStats2, pStats3)
				now = now.AddDate(0, 0, -1)
			} else {
				break
			}
		}
		now = time.Now()
		now = now.AddDate(0, 0, -1)
		//fmt.Println(pStats2, now.Format(layout))
		pmulibs.ApplyStats(pStats2, pmulibs.GetCourse(now.Format(layout)))*/

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
