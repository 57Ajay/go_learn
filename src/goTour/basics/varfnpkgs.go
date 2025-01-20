package basics

import (
	"fmt"
	"math"
)

type shapefunc interface {
	Area() float32
	Perimeter() float32
}

const (
	circleType    = "Circle"
	rectangleType = "Rectangle"
)

type shape struct {
	shapeType string
	radius    float32
	length    float32
	width     float32
}

func (s shape) Area() float32 {
	switch s.shapeType {
	case circleType:
		{
			return math.Pi * s.radius * s.radius
		}
	case rectangleType:
		{
			return s.length * s.width
		}
	default:
		{
			return 0
		}
	}
}

func (s shape) Perimeter() float32 {
	switch s.shapeType {
	case circleType:
		{
			return 2 * math.Pi * s.radius
		}
	case rectangleType:
		{
			return 2 * (s.length + s.width)
		}
	default:
		{
			return 0
		}
	}
}

// function with naked return;
func split(sum int) (x, y int) {
	x = sum * (4 / 9)
	y = sum - x
	return
}

func VarfnpkgsMain() {
	circle := shape{shapeType: circleType, radius: 5}
	rectangle := shape{shapeType: rectangleType, length: 5, width: 10}
	circleArea := circle.Area()
	rectangleArea := rectangle.Area()
	circleCircum := circle.Perimeter()
	rectanglePerimeter := rectangle.Perimeter()
	fmt.Println("area and Perimeter of circle and rectangle are: rectangleArea, rectanglePerimeter, circleArea, circleCircum", rectangleArea, rectanglePerimeter, circleArea, circleCircum)
}
