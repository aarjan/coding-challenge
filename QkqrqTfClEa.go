package main

import (
	"fmt"
)

type trie struct {
	value int
	nodes []*trie
}

func mai() {
	t012 := &trie{2, nil}
	t015 := &trie{5, nil}
	t018 := &trie{8, nil}

	t002 := &trie{2, nil}
	t003 := &trie{3, nil}
	t004 := &trie{4, nil}
	t006 := &trie{6, nil}
	t007 := &trie{7, nil}
	t008 := &trie{8, nil}

	t01 := &trie{value: 1, nodes: []*trie{t012, t015, t018}}
	t00 := &trie{0, []*trie{t002, t003, t004, t006, t007, t008}}

	t0 := &trie{0, []*trie{t01, t00}}

	t111 := &trie{1, nil}
	t115 := &trie{5, nil}
	t116 := &trie{6, nil}
	t119 := &trie{9, nil}

	t101 := &trie{1, nil}
	t104 := &trie{4, nil}
	t105 := &trie{5, nil}
	t107 := &trie{7, nil}

	t11 := &trie{1, []*trie{t111, t115, t116, t119}}
	t10 := &trie{1, []*trie{t101, t104, t105, t107}}

	t1 := &trie{1, []*trie{t11, t10}}

	t211 := &trie{1, nil}
	t214 := &trie{4, nil}
	t2110 := &trie{10, nil}

	t21 := &trie{1, []*trie{t211, t214, t2110}}
	t2 := &trie{2, []*trie{t21}}

	//take the count of each leaf element for each root node
	t0m := mapper([]int{2, 2, 3, 4, 5, 6, 7, 8, 8, 10})
	t1m := mapper([]int{1, 1, 3, 4, 5, 5, 6, 7, 9})
	t2m := mapper([]int{1, 4, 10})
	fmt.Println(t0m, t1m, t2m)

	// fmt.Println(t0)
	// fmt.Println(t1)
	// fmt.Println(t2)

	request := [][]int{{0, 1}, {1, 0}, {2, 1}, {3, 0}, {4, 1}}

	for _, r := range request {
		switch r[1] {
		case 1:
			a := redundant(t0m)
			b := redundant(t1m)
			c := redundant(t2m)

			// fmt.Println("red values:",a,b,c)

			if a < b && a < c || a == b || a == c {
				fmt.Printf("redundant, device:%d, vland_id:%d\n", 0, a)
				delete(t0m, a)
				t0.eat()
				// fmt.Println(t0)
			} else if b < a && b < c || b == c {
				fmt.Printf("redundant, device:%d, vland_id:%d\n", 1, b)
				delete(t1m, b)
				t1.eat()
				// fmt.Println(t1)
			} else {
				fmt.Printf("redundant, device:%d, vland_id:%d\n", 2, c)
				delete(t2m, c)
				t2.eat()
				// fmt.Println(t2)
			}

		case 0:
			a := nonRedundant(t0)
			b := nonRedundant(t1)
			c := nonRedundant(t2)

			// fmt.Println("non-values:",a,b,c)
			if (a < b || a == b) && (a == c || a < c) {
				fmt.Printf("non, device:%d, vland_id:%d\n", 0, a)
				t0.nodes[0].nodes = t0.nodes[0].nodes[1:]
				// fmt.Println(t0)
			} else if b < a && (b < c || b == c) {
				fmt.Printf("non, device:%d, vland_id:%d\n", 1, b)
				t1.nodes[0].nodes = t1.nodes[0].nodes[1:]
				// fmt.Println(t1)
			} else {
				fmt.Printf("non, device:%d, vland_id:%d\n", 2, c)
				t2.nodes[0].nodes = t2.nodes[0].nodes[1:]
				// fmt.Println(t2)
			}
		}
	}
}

func (t *trie) eat() {
	t.nodes[0].nodes = t.nodes[0].nodes[1:]
	t.nodes[1].nodes = t.nodes[1].nodes[1:]
}
func nonRedundant(t *trie) int {
	return t.nodes[0].nodes[0].value
}

func redundant(t map[int]int) int {
	min := 0
	for key, val := range t {
		if val == 2 {
			if min == 0 {
				min = key
			} else if min > key {
				min = key
			}
		}
	}
	// return some execessively high number
	// Todo: make it less whacky
	if min == 0 {
		return 100000
	}
	return min
}

func (t *trie) String() string {
	s := ""
	s += fmt.Sprint(t.value)
	if t.nodes == nil {
		return "(" + s + ")"
	}
	for _, n := range t.nodes {
		s += " " + n.String()
	}
	return "(" + s + ")"

}

// func mapper(arr []int) map[int]int {
// 	res := make(map[int]int)
// 	for _, v := range arr {
// 		res[v] = res[v] + 1
// 	}
// 	return res
// }
