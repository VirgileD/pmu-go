package pmulibs

import (
//	"fmt"
)

type Combination struct {
	Comb  []int
	Value float32
}

type Combinations struct {
	Set        []int
	Of         int
	On         int
	Nbr        int
	Combs      []Combination
	CombsIndex int
}

func InitCombination(of int, set []int) (combs *Combinations) {
	lenset := len(set)
	NbrOfComb := NbrOfComb(lenset, of)
	//fmt.Println(NbrOfComb)
	combs = &Combinations{Of: of, Set: set, On: lenset, Nbr: NbrOfComb}
	combs.Combs = make([]Combination, NbrOfComb, NbrOfComb)
	combs.CombsIndex = 0
	combs.Combs[combs.CombsIndex].Comb = make([]int, of, of)
	for i := 0; i < of; i++ {
		combs.Combs[combs.CombsIndex].Comb[i] = i
	}
	//fmt.Println("initial comb: ", combs.Combs[combs.CombsIndex].Comb)
	i := 0
	for combs.Combs[combs.CombsIndex].Comb[0] < combs.On-combs.Of {
		//fmt.Println("next comb ", combs.Combs[combs.CombsIndex])
		combs.CombsIndex++
		combs.Combs[combs.CombsIndex].Comb = make([]int, of, of)
		for i = 0; i < of; i++ {
			combs.Combs[combs.CombsIndex].Comb[i] = combs.Combs[combs.CombsIndex-1].Comb[i]
		}
		for i = of - 1; i > 0 && combs.Combs[combs.CombsIndex].Comb[i] == combs.On-of+i; {
			i--
		}
		combs.Combs[combs.CombsIndex].Comb[i]++
		for j := i; j < of-1; j++ {
			combs.Combs[combs.CombsIndex].Comb[j+1] = combs.Combs[combs.CombsIndex].Comb[j] + 1
		}
	}
	combs.CombsIndex--
	//fmt.Println("last comb: ", combs.Combs[combs.CombsIndex].Comb)
	for i = 0; i < NbrOfComb; i++ {
		for j := 0; j < of; j++ {
			combs.Combs[i].Comb[j] = set[combs.Combs[i].Comb[j]]
		}
		//fmt.Printf("%#v\n", combs.Combs[i])
	}
	return combs
}

func NbrOfComb(on, of int) (ans int) {
	var delta, iMax int
	if on < 0 || of < 0 {
		panic("[NbrOfComb] Invalid negative parameter")
	}
	if on < of {
		return 0
	}
	if on == of {
		return 1
	}
	if of < on-of {
		delta = on - of
		iMax = of
	} else {
		delta = of
		iMax = on - of
	}
	ans = delta + 1
	for i := 2; i <= iMax; i++ {
		ans = (ans * (delta + i)) / i
	}
	//fmt.Println("[NbrOfComb] ", on, of, ans)
	return ans
}

/*
public class Combination {
private long n = 0;
private long k = 0;
private long[] data = null;

public Combination(long n, long k) {
	if (n < 0 || k < 0) // normally n >= k throw new Exception("Negative parameter in constructor");
	this.n = n; this.k = k; this.data = new long[k];
	for (long i = 0; i < k; ++i) this.data[i] = i;
}

public override string ToString() {
	StringBuilder sb = new StringBuilder();
	sb.Append("{ ");
	for (long i = 0; i < this.k; ++i) sb.AppendFormat("{0} ", this.data[i]);
	sb.Append("}");
	return sb.ToString;
}

public string[] ApplyTo(string[] strarr) {
	if (strarr.Length != this.n) throw new Exception("Bad array size");
	string[] result = new string[this.k];
	for (long i = 0; i < result.Length; ++i) result[i] = strarr[this.data[i]];
	return result;
}
*/
