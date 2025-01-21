package basics

import (
	"fmt"
	"golang.org/x/tour/pic"
	"golang.org/x/tour/wc"
	"strings"
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

func println(T ...any) {
	fmt.Println(T...)
}

func understandingSlicesIn_depth() {
	aSlice := make([]int, 5)
	a2dSlice := make([][]int, 5)
	a3dSlice := make([][][]int, 5)
	for i := 0; i < 5; i++ {
		a2dSlice[i] = make([]int, 5)
		for j := 0; j < 5; j++ {
			a2dSlice[i][j] = j + 1
		}
	}
	for i := 0; i < 5; i++ {
		a3dSlice[i] = make([][]int, 5)
		for j := 0; j < 5; j++ {
			a3dSlice[i][j] = make([]int, 5)
			for k := 0; k < 5; k++ {
				a3dSlice[i][j][k] = k + 1
			}
		}
	}
	println("Aslice: ", len(aSlice), cap(aSlice), aSlice)
	println("A2dSlice: ", len(a2dSlice), cap(a2dSlice), a2dSlice)
	println("A3dSlice: ", len(a3dSlice), cap(a3dSlice), a3dSlice)
}

// this function creates an N-dimensional slice with the given dimensions.
func createNDSlice(dims []int) interface{} {
	if len(dims) == 0 {
		return 0
	}

	dim := dims[0]
	rest := dims[1:]

	// Create the outer slice
	slice := make([]interface{}, dim)
	// println("Slice: ", slice)
	for i := range slice {
		// Recursively create inner slices
		slice[i] = createNDSlice(rest)
		// println("Slice[i]: ", slice[i])
	}

	return slice
}

// This function recursively prints an N-dimensional slice.
func printNDSlice(slice interface{}, level int) {
	switch s := slice.(type) {
	case []interface{}:
		fmt.Printf("[")
		for i, v := range s {
			printNDSlice(v, level+1)
			if i < len(s)-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("]")
		if level == 0 {
			fmt.Println()
		}
	default:
		fmt.Printf("%v", s)
	}
}

var wordCount map[string]int

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	wordCount = make(map[string]int)
	for _, word := range words {
		wordCount[word]++
	}
	return wordCount
}

func TypesMain() {
	pointers()

	ptrstr := ptrStruct{1, 2}
	ptr := &ptrstr
	ptr.X = 3
	fmt.Println(ptrstr)
	array()
	pic.Show(Pic)
	fmt.Println(m)
	understandingSlicesIn_depth()
	dims := []int{3, 4, 2}
	my3DSlice := createNDSlice(dims)
	s := my3DSlice.([]interface{})
	s0 := s[0].([]interface{})
	s00 := s0[0].([]interface{})
	s00[0] = 100
	s00[1] = 101

	fmt.Println("3D Slice:")
	printNDSlice(my3DSlice, 0)

	// Example: 2x2 slice
	dims2 := []int{2, 2}
	my2DSlice := createNDSlice(dims2)

	fmt.Println("\n2D Slice:")
	printNDSlice(my2DSlice, 0)

	// Example: 1x5 slice (effectively a 1D slice)
	dims3 := []int{1, 5}
	my1DSlice := createNDSlice(dims3)

	fmt.Println("\n1D Slice:")
	printNDSlice(my1DSlice, 0)
	wc.Test(WordCount)
}
