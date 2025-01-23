package goexp

import (
	"fmt"
	"maps"
	"slices"
)

func someExp() {
	var s []string
	fmt.Println("uninit:", s, s == nil, len(s) == 0)
	s = make([]string, 3)
	s = append(s, "d")
	s = append(s, "e", "f")
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("emp:", s, "len:", len(s), "cap:", cap(s))
	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy:", c)
}

func slice() {
	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)
	t2 := []string{"g", "h", "i"}
	if slices.Equal(t, t2) {
		fmt.Println("t == t2")
	}
	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}

func mapsExp() {
	m := make(map[string]int)

	m["k1"] = 7
	m["k2"] = 13

	fmt.Println("map:", m)

	v1 := m["k1"]
	fmt.Println("v1:", v1)

	v3 := m["k3"]
	fmt.Println("v3:", v3)

	fmt.Println("len:", len(m))

	delete(m, "k2")
	fmt.Println("map:", m)

	clear(m)
	fmt.Println("map:", m)

	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)

	n2 := map[string]int{"foo": 1, "bar": 2}
	if maps.Equal(n, n2) {
		fmt.Println("n == n2")
	}
}

func Exp1Main() {
	fmt.Println("Hello World, from exp1.go")
	someExp()
	slice()
	mapsExp()
}
