package main

func AppendMake() {
	n := 10000
	r := make([]int, 1, n+1)
	// fmt.Printf("%p\n", r)
	r[0] = 1
	ad := [10000]int{}
	r = append(r, ad[:]...)
	// fmt.Printf("%p\n", r)
	_ = r
}

func AppendSimple() {
	r := []int{1}
	// fmt.Printf("%p\n", r)
	ad := [10000]int{}
	r = append(r, ad[:]...)
	// fmt.Printf("%p\n", r)
	_ = r
}

func main() {
	AppendMake()
	AppendSimple()
}
