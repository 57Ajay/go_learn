package basics

import (
	"fmt"
	// "strings"
	"golang.org/x/tour/pic"
)

type ptrStruct struct {
	X int
	Y int
}

func pointers() {
	i, j := 42, 2701

	p := &i         // point to i
	fmt.Println(p)  // print the address of i
	fmt.Println(*p) // read i through the pointer -> i =  42
	*p = 21         // set i through the pointer -> i = 21 now
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j -> now p points to j
	*p = *p / 37   // divide j through the pointer -> 2701 / 37 = 73
	fmt.Println(j) // see the new value of j

}

func array() {
	var a_ [2]string
	a_[0] = "Hello"
	a_[1] = "World"
	fmt.Println(a_[0], a_[1])

	primes := [6]int{2, 3, 5, 7, 11, 13}

	// Slicing an array
	// Changing the elements of a slice modifies the corresponding
	// elements of its underlying array.
	// Other slices that share the same underlying array will see those changes.
	primesSlice := primes[1:4]
	fmt.Println(primesSlice)

	// Slices are like references to arrays
	// Struct and array
	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
		{7, true},
		{11, false},
		{13, true},
	}
	fmt.Println(s, len(s), cap(s))

	a := make([]int, 5)
	printSlice("a", a)

	b := make([]int, 0, 5)
	printSlice("b", b)

	c := b[:2]
	printSlice("c", c)

	d := c[2:5]
	printSlice("d", d)

	var pow = []int{1, 2, 4, 8, 16, 32, 64, 128}

	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}

}

func printSlice(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}

func Pic(dx, dy int) [][]uint8 {
	image := make([][]uint8, dy)
	for y := range image {
		image[y] = make([]uint8, dx)
		for x := range image[y] {
			grayValue := uint8((x + y) / 2)
			image[y][x] = grayValue
		}
	}
	return image
}

type Vertex struct {
	Lat, Long float64
}

var m = map[string]Vertex{
	"Bell Labs": {
		40.68433, -74.39967,
	},
	"Google": {
		37.42202, -122.08408,
	},
}

func TypesMain() {
	fmt.Println("-----------------types-----------------")
	pointers()

	ptrstr := ptrStruct{1, 2}
	ptr := &ptrstr
	ptr.X = 3
	fmt.Println(ptrstr)
	array()
	pic.Show(Pic)
	fmt.Println(m)
}
