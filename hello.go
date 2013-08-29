package main

import (
        "fmt"
        "labix.org/v2/mgo"
        "labix.org/v2/mgo/bson"
)

type StatsChev struct {
	LastCote int `bson:"lastCote"`
	RefCote int `bson:"refCote"`
	Valeur int
}
type Course struct {
		Name string `bson:"name"`
        Location string `bson:"location"`
        NbPartants int `bson:"nbPartants"`
        Date string //`bson:"date"`
        Finish [5]int //`bson:"finish"`
        Gains map[string]float32 //`bson:"gains"`
        Pronos map[string][]int //`bson:"pronos"`
        StatsChev map[string]StatsChev `bson:"statsChev"`
        
}

func main() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
            panic(err)
    }
    defer session.Close()
    // Optional. Switch the session to a monotonic behavior.
    //session.SetMode(mgo.Monotonic, true)

    c := session.DB("pmu").C("courses")
    result := Course{}
    err = c.Find(bson.M{"date": "2013-08-28"}).One(&result)
    //nb, err := c.Find(bson.M{"date": "2013-08-27"}).Count()
    if err != nil {
            panic(err)
    }
    //fmt.Printf("%d\n", nb)
  	fmt.Printf("%#v\n", result)
  	//fmt.Printf("leturf: ", result.Pronos["le turf"])
    //fmt.Println("location: ", result.finish)
}	