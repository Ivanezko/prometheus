package main

func AppendMake() {
	n := 10000
	r := make([]int, 1, n+1)
	r[0] = 1
	ad := [100]int{}
	r = append(r, ad[:]...)
	_ = r
}

func AppendSimple() {
	r := []int{1}
	ad := [10000]int{}
	r = append(r, ad[:]...)
	_ = r
}
