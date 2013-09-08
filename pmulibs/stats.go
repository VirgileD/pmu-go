package pmulibs

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"sort"
	"strconv"
)

type PStat struct {
	In4 float32
	In6 float32
	In8 float32
}

type PStats struct {
	Stats map[string]*PStat
	Date  string
	Nbr   int
}

type StatsChev struct {
	LastCote int `bson:"lastCote"`
	RefCote  int `bson:"refCote"`
	Valeur   int
}

func contains(haystack []int, needles []int) (nbr int) {
	for _, needle := range needles {
		for _, a := range haystack {
			if a == needle {
				nbr += 1
			}
		}
	}
	return nbr
}

func calcStats(date string) (pStats PStats) {
	fmt.Println("calcStats(" + date + ")")
	result := PStats{}
	result.Date = date
	result.Nbr = 1
	result.Stats = make(map[string]*PStat)
	course := GetCourse(date)
	if len(course.Finish) == 0 {
		panic("[calcStats] no finish in course of " + date)
	}
	for pronoName, prono := range course.Pronos {
		nbr4 := contains(course.Finish[0:4], prono[0:4])
		nbr6 := contains(course.Finish[0:4], prono[0:6])
		nbr8 := contains(course.Finish[0:4], prono[0:8])
		if _, ok := pStats.Stats[pronoName]; !ok {
			result.Stats[pronoName] = &PStat{In4: float32(nbr4), In6: float32(nbr6), In8: float32(nbr8)}
		} else {
			result.Stats[pronoName].In4 = (3*result.Stats[pronoName].In4 + float32(nbr4)) / 4
			result.Stats[pronoName].In6 = (3*result.Stats[pronoName].In6 + float32(nbr6)) / 4
			result.Stats[pronoName].In8 = (3*result.Stats[pronoName].In8 + float32(nbr8)) / 4
		}
	}
	return result
}

func GetStats(date string, force bool) (pStats PStats) {
	fmt.Println("getStats(" + date + "," + strconv.FormatBool(force) + ")")
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic("[getStats] error while creating session " + err.Error())
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)

	c := session.DB("pmu").C("pstats")
	result := PStats{}
	if force {
		result = calcStats(date)
		SetStats(result)
		return result
	}
	err = c.Find(bson.M{"date": date}).One(&result)
	//nb, err := c.Find(bson.M{"date": date}).Count()
	if err != nil {
		if err.Error() != "not found" {
			panic("[getStats] error retireving stats " + date + ":" + err.Error())
			panic(err)
		} else {
			if err.Error() == "not found" {
				result = calcStats(date)
				SetStats(result)
			}
		}
	}

	return result
}

func AddStats(pStats1 PStats, pStats2 PStats) (pStats3 PStats) {
	fmt.Println("addStats")
	pStats3.Stats = make(map[string]*PStat)
	pStats3.Date = pStats1.Date + "," + pStats2.Date
	pStats3.Nbr = pStats1.Nbr + pStats2.Nbr
	for name, stats := range pStats2.Stats {
		if stat1, ok := pStats1.Stats[name]; !ok {
			fmt.Println("[addStats] " + name + " only in " + pStats2.Date)
			pStats3.Stats[name] = &PStat{In4: stats.In4 + (float32(pStats1.Nbr) * stats.In4), In6: stats.In6 + (float32(pStats1.Nbr) * stats.In6),
				In8: stats.In8 + (float32(pStats1.Nbr) * stats.In8)}
		} else {
			pStats3.Stats[name] = &PStat{In4: stats.In4 + stat1.In4, In6: stats.In6 + stat1.In6, In8: stats.In8 + stat1.In8}
		}
	}
	for name, stats := range pStats1.Stats {
		if _, ok := pStats2.Stats[name]; !ok {
			fmt.Println("[addStats] " + name + " only in " + pStats1.Date)
			pStats3.Stats[name] = &PStat{In4: stats.In4 + (float32(pStats2.Nbr) * stats.In4), In6: stats.In6 + (float32(pStats2.Nbr) * stats.In6),
				In8: stats.In8 + (float32(pStats2.Nbr) * stats.In8)}
		} else {
			// nothing to do as it has alredy been done in the first loop
		}
	}
	return pStats3
}

func SetStats(pStats PStats) {
	fmt.Println("setStats")
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic("[setStats] error while creating session " + err.Error())
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("pmu").C("pstats")
	_, err = c.Upsert(bson.M{"date": pStats.Date}, pStats)
	//nb, err := c.Find(bson.M{"date": "2013-08-27"}).Count()
	if err != nil {
		panic(err)
	}
}

// A data structure to hold a key/value pair.
type Pair struct {
	Key   int
	Value float32
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[int]float32) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

func ApplyStats(pStats PStats, course Course) {
	final := make(map[int]float32)
	for pronoName, prono := range course.Pronos {
		fmt.Println(pronoName, prono)
		if stats, ok := pStats.Stats[pronoName]; !ok {
			fmt.Println(pronoName + " is unknown in stats!")
		} else {
			for index, chev := range prono {
				//fmt.Println(chev, index)
				if course.StatsChev[strconv.Itoa(chev)].RefCote > 20 && stats.In4 >= 3 {
					fmt.Println("setting 0 to ", chev)
					y, _ := strconv.Atoi((strconv.FormatFloat(float64(stats.In8), 'g', -1, 32)))
					index = -y
				}
				if valChev, ok := final[chev]; !ok {
					final[chev] = (stats.In8 / float32(pStats.Nbr)) * (float32(8) - float32(index))
				} else {
					final[chev] = float32(valChev) + (stats.In8/float32(pStats.Nbr))*(float32(8)-float32(index))
				}
			}
		}
	}
	test := sortMapByValue(final)
	i := 0
	for _, pair := range test {
		fmt.Println(strconv.Itoa(i) + " " + strconv.Itoa(pair.Key) + " (" + strconv.FormatFloat(float64(pair.Value), 'g', -1, 32) + ")")
		i++
	}
	fmt.Println(course.Finish, course.NbPartants, course.Gains["m7"], course.Gains["4d"])
}
