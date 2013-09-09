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

type sortedMap struct {
	m map[string]float32
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]float32) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
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
	//fmt.Println("calcStats(" + date + ")")
	result := PStats{}
	result.Date = date
	result.Nbr = 1
	result.Stats = make(map[string]*PStat)
	course := GetCourse(date)
	//fmt.Println("[calcStats] ", course.Pronos, course.Finish)
	if course.NbPartants == 0 {
		result.Nbr = 0
		return result
	}
	for pronoName, prono := range course.Pronos {
		//fmt.Println("[calcStats] ", pronoName, prono)
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
	//fmt.Println("getStats(" + date + "," + strconv.FormatBool(force) + ")")
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
		if result.Nbr != 0 {
			SetStats(result)
		}
		return result
	}
	err = c.Find(bson.M{"date": date}).One(&result)
	//nb, err := c.Find(bson.M{"date": date}).Count()
	if err != nil {
		if err.Error() != "not found" {
			panic("[getStats] error retireving stats " + date + ":" + err.Error())
		} else {
			if err.Error() == "not found" {
				result = calcStats(date)
				if result.Nbr != 0 {
					SetStats(result)
				}
				SetStats(result)
			}
		}
	}

	return result
}

func AddStats(pStats1 PStats, pStats2 PStats) (pStats3 PStats) {
	//fmt.Println("addStats")
	if pStats1.Nbr == 0 {
		return pStats2
	}
	pStats3.Stats = make(map[string]*PStat)
	pStats3.Date = pStats1.Date + "," + pStats2.Date
	pStats3.Nbr = pStats1.Nbr + pStats2.Nbr
	for name, stats := range pStats2.Stats {
		if stat1, ok := pStats1.Stats[name]; !ok {
			//fmt.Println("[addStats] " + name + " only in " + pStats2.Date)
			pStats3.Stats[name] = &PStat{In4: stats.In4 + (float32(pStats1.Nbr) * stats.In4), In6: stats.In6 + (float32(pStats1.Nbr) * stats.In6),
				In8: stats.In8 + (float32(pStats1.Nbr) * stats.In8)}
		} else {
			pStats3.Stats[name] = &PStat{In4: stats.In4 + stat1.In4, In6: stats.In6 + stat1.In6, In8: stats.In8 + stat1.In8}
		}
	}
	for name, stats := range pStats1.Stats {
		if _, ok := pStats2.Stats[name]; !ok {
			//fmt.Println("[addStats] " + name + " only in " + pStats1.Date)
			pStats3.Stats[name] = &PStat{In4: stats.In4 + (float32(pStats2.Nbr) * stats.In4), In6: stats.In6 + (float32(pStats2.Nbr) * stats.In6),
				In8: stats.In8 + (float32(pStats2.Nbr) * stats.In8)}
		} else {
			// nothing to do as it has alredy been done in the first loop
		}
	}
	return pStats3
}

func SetStats(pStats PStats) {
	//fmt.Println("setStats")
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

func applyStat(comb *Combination, pronos map[string][]int, stats PStats) (value float32, stringComb string) {
	stringComb = ""
	comb.Value = 0
	for pronoName, prono := range pronos {
		if stat, ok := stats.Stats[pronoName]; ok {
			comb.Value += float32(contains(prono[0:4], comb.Comb)) - (float32(stat.In4) / float32(stats.Nbr))
			comb.Value += float32(contains(prono[0:6], comb.Comb)) - (float32(stat.In6) / float32(stats.Nbr))
			comb.Value += float32(contains(prono[0:8], comb.Comb)) - (float32(stat.In8) / float32(stats.Nbr))
		}
	}
	for index, chev := range comb.Comb {
		stringComb += strconv.Itoa(chev)
		if index < len(comb.Comb)-1 {
			stringComb += "-"
		}
	}
	//fmt.Println("---------------", comb)
	return comb.Value, stringComb
}

func UnifyComb(comb *Combination, mapComb map[string]float32) (stringComb string, value float32) {
	combsOf4 := InitCombination(4, comb.Comb)
	for _, combOf4 := range combsOf4.Combs {
		stringComb := ""
		for index, chev := range combOf4.Comb {
			stringComb += strconv.Itoa(chev)
			if index < 3 {
				stringComb += "-"
			}
		}
		comb.Value += mapComb[stringComb]
	}
	stringComb = ""
	for index, chev := range comb.Comb {
		stringComb += strconv.Itoa(chev)
		if index < 6 {
			stringComb += "-"
		}
	}
	return stringComb, comb.Value
}

func ApplyStats(pStats PStats, course Course) {
	// calculate the pondered wiehgt of each horse
	final := make(map[int]float32)
	for pronoName, prono := range course.Pronos {
		//fmt.Println(pronoName, prono)
		if stats, ok := pStats.Stats[pronoName]; !ok {
			fmt.Println(pronoName + " is unknown in stats!")
		} else {
			for index, chev := range prono {
				if valChev, ok := final[chev]; !ok {
					final[chev] = (stats.In8 / float32(pStats.Nbr)) * (float32(8) - float32(index))
				} else {
					final[chev] = float32(valChev) + (stats.In8/float32(pStats.Nbr))*(float32(8)-float32(index))
				}
			}
		}
	}
	// init the set of horse (horse that have not been inlcuded in any prono are removed)
	//fmt.Println(final)

	/*set := make([]int, len(final), len(final))
	i := 1
	j := 0
	for i < len(final)+1 {
		_, ok := final[i]
		if ok {
			//fmt.Println("cev ", i, "(", value, "):", j)
			set[j] = i
			j++
		}
		i++
	}
	fmt.Println(set)
	combs4 := InitCombination(4, set)
	mapComb := make(map[string]float32)
	for indexComb, comb := range combs4.Combs {
		//fmt.Println(indexComb, comb)
		value, stringComb := applyStat(&comb, course.Pronos, pStats)
		fmt.Println(indexComb, comb, stringComb)
		mapComb[stringComb] = value
	}
	combs7 := InitCombination(7, set)
	mapComb7 := make(map[string]float32)
	for _, comb := range combs7.Combs {
		stringComb, value := UnifyComb(&comb, mapComb)
		//fmt.Println(comb)
		mapComb7[stringComb] = value
	}
	sorted := sortedKeys(mapComb7)
	fmt.Println(len(sorted))
	//for _, key := range sorted {
	//	fmt.Println(mapComb7, mapComb7[key])
	//}
	*/
	test := sortMapByValue(final)
	for index, pair := range test {
		fmt.Println(strconv.Itoa(len(test)-index) + " " + strconv.Itoa(pair.Key) + " (" + strconv.FormatFloat(float64(pair.Value), 'g', -1, 32) + ")")
	}
	fmt.Println(course.Finish, course.NbPartants, course.Gains["m7"], course.Gains["4d"])
}
