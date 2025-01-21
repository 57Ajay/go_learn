package methodsInterfaces

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func change(v *Vertex) {
	v.X = 2 * v.X
	v.Y = 3 * v.X
}

// interfaces are implemented implicitly

func interfaces() {
	var s interface{} = "Hello"
	s = 121
	i, ok := s.(string)
	if ok {
		fmt.Println(i)
	} else {
		fmt.Println("Not a string")
	}
}

type Person struct {
	name string
	age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.name, p.age)
}

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.
// String() string will be called when fmt.Println is called on IPAddr
func (ip IPAddr) String() string {
	var val []string
	for _, v := range ip {

		val = append(val, strconv.Itoa(int(v)))

	}
	return strings.Join(val, ".")
}

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.When, e.What)
}

func run() error {
	return &MyError{
		time.Now(),
		"it didn't work",
	}
}

type ErrNegativeSqrt float64

// ErrNegativeSqrt is a type that implements the error interface
// Error() string is a method that is called when fmt.Println is called on ErrNegativeSqrt
func (e ErrNegativeSqrt) Error() string {
	if e < 0 {
		return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
	} else {
		return fmt.Sprintf("The square root of %v is %v", float64(e), math.Sqrt(float64(e)))
	}
}

func Sqrt(x float64) ErrNegativeSqrt {
	return ErrNegativeSqrt(x)
}

func readers() {
	r := strings.NewReader("Hello, Reader!")
	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}

type MyReader struct{}

func (r MyReader) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = 'A'
	}
	return len(b), nil
}

type rot13Reader struct {
	r io.Reader
}

func (rot13 rot13Reader) Read(b []byte) (int, error) {
	n, err := rot13.r.Read(b)
	if err != nil {
		return 0, err
	}
	for i := 0; i < n; i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] = 'A' + (b[i]-'A'+13)%26
		} else if b[i] >= 'a' && b[i] <= 'z' {
			// c := b[i]
			// fmt.Println("This is c: ", c)
			b[i] = 'a' + (b[i]-'a'+13)%26
		}
	}
	return n, nil
}

type Image struct {
	Width, Height int
	Data          [][]uint8
}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img Image) At(x, y int) color.Color {
	v := img.Data[y][x]
	return color.RGBA{v, v, 255, 255}
}

func MethodsInterfacesMain() {
	v := Vertex{3, 4}
	v.Scale(10)
	change(&v)
	fmt.Println("The Abs: ", v.Abs())
	interfaces()
	person := Person{"Arthur Dent", 42}
	fmt.Println(person)
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
	if err := run(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
	readers()

	reader := MyReader{}
	b := make([]byte, 10)
	n, err := reader.Read(b)
	fmt.Printf("Read %d bytes: %s, error: %v\n", n, b, err)
	b2 := make([]byte, 5)
	n2, err2 := reader.Read(b2)
	fmt.Printf("Read now %d bytes: %s, error: %v\n", n2, b2, err2)

	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(m.Bounds())
	defer fmt.Println(m.At(0, 0).RGBA())

	width := 1250
	height := 1250
	data := make([][]uint8, height)
	for y := range data {
		data[y] = make([]uint8, width)
		for x := range data[y] {
			data[y][x] = uint8(math.Sin(float64(x)/float64(width)*2*math.Pi) * 127.0)
		}
	}

	img := Image{Width: width, Height: height, Data: data}

	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, img)
	fmt.Println("Image saved to image.png")

}
