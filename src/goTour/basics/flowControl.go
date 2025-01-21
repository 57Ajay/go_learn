package basics

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
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

func deferFunc() {
	fmt.Println("counting")

	for i := 0; i < 10; i++ {
		// Will print in reverse order
		defer fmt.Println(i)
	}

	fmt.Println("done")
}

func CopyFile(dstName, srcName string) (written int64, err error) {

	// Some extra code to get the current directory
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fmt.Println("Executable Path:", executablePath)

	// Get the directory containing the executable:
	executableDir := filepath.Dir(executablePath)
	fmt.Println("Executable Directory:", executableDir)
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("Current Working Directory:", currentDir)

	// Main code for the function
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func f() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	fmt.Println("Calling g.")
	g(0)
	fmt.Println("Returned normally from g.")
}

func g(i int) {
	if i > 3 {
		fmt.Println("Panicking!")
		panic(fmt.Sprintf("%v", i))
	}
	defer fmt.Println("Defer in g", i)
	fmt.Println("Printing in g", i)
	g(i + 1)
}

func panicDefer() {
	f()
	fmt.Println("Returned normally from f.")
}

func FlowControlMain() {
	fmt.Println("for/while loop: ", forLoop())
	fmt.Println("power: ", pow(3, 5, 100))
	fmt.Println("sqrt: ", sqrt(9))
	fmt.Println("OS: ", getOS())
	defer fmt.Println("defered function")
	fmt.Println("When is Saturday? ", getWeekday())
	deferFunc()
	fmt.Println(CopyFile("test.txt", "main.go"))
	panicDefer()
}
