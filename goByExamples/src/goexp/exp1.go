package goexp

import (
	"fmt"
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

func Exp1Main() {
	fmt.Println("Hello World, from exp1.go")
	someExp()
}
