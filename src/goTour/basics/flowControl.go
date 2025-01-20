package basics

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

// Works like a while loop
func forLoop() int {
	sum := 1

	for sum < 1000 {
		sum += sum
	}
	return sum
}

func pow(x, n, lim float64) float64 {
	if v := math.Pow(x, n); v < lim {
		return v
	}
	return lim
}

func sqrt(x float64) float64 {
	var z float64 = 1
	for i := 0; i < int(x/2); i++ {
		z -= (z*z - x) / (2 * z)
		// fmt.Println(z)
	}
	return z
}

func getOS() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "OS X."
	case "linux":
		return "Linux."
	case "windows":
		return "Windows."
	case "freebsd":
		return "FreeBSD."
	default:
		return "Unknown."
	}
}

func getWeekday() string {
	today := time.Now().Weekday()
	switch time.Saturday {
	case today + 0:
		return "Today."
	case today + 1:
		return "Tomorrow."
	case today + 2:
		return "In two days."
	default:
		return "Too far away."
	}

}

func FlowControlMain() {
	fmt.Println("for/while loop: ", forLoop())
	fmt.Println("power: ", pow(3, 5, 100))
	fmt.Println("sqrt: ", sqrt(9))
	fmt.Println("OS: ", getOS())
	println("When is Saturday? ", getWeekday())
}
