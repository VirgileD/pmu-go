// couses
package pmulibs

import (
	"fmt"
	"labix.org/v2/mgo"
)

var _ = fmt.Println

type Course struct {
	Name       string               `bson:"name"`
	Location   string               `bson:"location"`
	NbPartants int                  `bson:"nbPartants"`
	Date       string               //`bson:"date"`
	Finish     [5]int               //`bson:"finish"`
	Gains      map[string]float32   //`bson:"gains"`
	Pronos     map[string][]int     //`bson:"pronos"`
	StatsChev  map[string]StatsChev `bson:"statsChev"`
}

func GetCourse(findObject interface{}) (course Course) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic("[getCourse] error while ceating session: " + err.Error())
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)

	c := session.DB("pmu").C("courses")
	result := Course{}
	err = c.Find(findObject).One(&result)
	//nb, err := c.Find(bson.M{"date": "2013-08-27"}).Count()
	if err != nil {
		result.NbPartants = 0
		return result
	}
	//fmt.Printf("%d\n", nb)
	//fmt.Printf("%#v\n", result)
	//fmt.Printf("leturf: ", result.Pronos["turf tv"])
	//fmt.Println("location: ", result.finish)
	return result
}

func GetCourses(findObject interface{}) (courses []Course) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic("[getCourse] error while ceating session: " + err.Error())
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)

	c := session.DB("pmu").C("courses")
	results := []Course{}
	err = c.Find(findObject).All(&results)
	//nb, err := c.Find(bson.M{"date": "2013-08-27"}).Count()
	if err != nil {
		return results
	}
	//fmt.Printf("%d\n", nb)
	//fmt.Printf("%#v\n", result)
	//fmt.Printf("leturf: ", result.Pronos["turf tv"])
	//fmt.Println("location: ", result.finish)
	return results
}
