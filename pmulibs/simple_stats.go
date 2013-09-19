package pmulibs

import (
	"fmt"
	"math"
	"sort"
)

// A data structure to hold a key/value pair.
type Pair struct {
	Key   int
	Value float64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[int]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

type Stats map[string][]int

type MeanAndVariance struct {
	Mean     float64
	Variance float64
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

func CalcStats(courses []Course) (stats Stats) {
	stats = make(map[string][]int)
	nbStats := len(courses)
	for _, course := range courses {
		for pronoName, prono := range course.Pronos {
			//fmt.Println(pronoName+" :", course.Finish, "/", prono)
			nbr8 := contains(course.Finish[0:4], prono[0:8])
			if _, ok := stats[pronoName]; !ok {
				stats[pronoName] = make([]int, 1, 1)
			}
			stats[pronoName] = append(stats[pronoName], nbr8)
		}
	}
	for pronoName, stat := range stats {
		stats[pronoName] = stat[1:]
		//fmt.Println(pronoName, stats[pronoName])
	}
	for pronoName, stat := range stats {
		if float32(nbStats*75/100) > float32(len(stat)) {
			//fmt.Println(pronoName, float32(nbStats)*75/100, float32(len(stat)), "deleted")
			delete(stats, pronoName)
		}
	}
	return stats
}

func GetMeansAndVariances(stats Stats) (meansAndVariances map[string]MeanAndVariance) {
	var mean = 0.0
	var vari = 0.0
	meansAndVariances = make(map[string]MeanAndVariance)
	for pronoName, stat := range stats {
		//fmt.Println(pronoName, stat)
		for _, nb := range stat {
			mean = mean + float64(nb)
		}
		mean = mean / float64(len(stat))
		for _, nb := range stat {
			vari = vari + (mean-float64(nb))*(mean-float64(nb))
		}
		vari = math.Sqrt(vari / float64(len(stat)-1))
		meansAndVariances[pronoName] = MeanAndVariance{Mean: mean, Variance: vari}
	}
	return meansAndVariances
}

func ApplyStats(course Course, stats Stats) {
	meansAndVariances := GetMeansAndVariances(stats)
	//for pronoName, meanAndVariance := range meansAndVariances {
	//	fmt.Println(pronoName, meanAndVariance)
	//}
	mapChev := make(map[int]float64)
	for pronoName, prono := range course.Pronos {
		for index, chev := range prono {
			if meansAndVariances[pronoName].Variance == 0 {
				//fmt.Println(pronoName, meansAndVariances[pronoName])
			} else {
				if val, ok := mapChev[chev]; !ok {
					mapChev[chev] = float64(8-index) * meansAndVariances[pronoName].Mean / meansAndVariances[pronoName].Variance
					//fmt.Println(index, chev, mapChev[chev], float64(8-index), meansAndVariances[pronoName].Mean, meansAndVariances[pronoName].Variance)
				} else {
					mapChev[chev] = val + (float64(8-index) * meansAndVariances[pronoName].Mean / meansAndVariances[pronoName].Variance)
					//fmt.Println(index, chev, val, mapChev[chev], float64(8-index), meansAndVariances[pronoName].Mean, meansAndVariances[pronoName].Variance)
				}
			}
		}
	}
	fmt.Println(course.Finish, course.Gains["m7"])
	fmt.Println(sortMapByValue(mapChev))
}
