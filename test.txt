package main

import (
	"fmt"
	"github.com/fatih/color"
	"go_learn/src/goTour/basics"
	"go_learn/src/goTour/concurrency"
	"go_learn/src/goTour/generics"
	"go_learn/src/goTour/methodsInterfaces"
	"go_learn/src/server"
)

func main() {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("\n%s\n", cyan("-----------------variables, functions, and packages-----------------"))
	basics.VarfnpkgsMain()
	fmt.Printf("\n%s\n", yellow("-----------------flow control-----------------"))
	basics.FlowControlMain()
	fmt.Printf("\n%s\n", green("-----------------more types: structs, slices, and maps-----------------"))
	basics.TypesMain()
	basics.IGmain()
	fmt.Printf("\n%s\n", magenta("-----------------methods and interfaces-----------------"))
	methodsInterfaces.MethodsInterfacesMain()
	fmt.Printf("\n%s\n", blue("-----------------generics-----------------"))
	generics.GenericsMain()
	fmt.Printf("\n%s\n", red("-----------------concurrency-----------------"))
	concurrency.ConcurrencyMain()
	concurrency.WebCrawlerMain()
	fmt.Printf("\n%s\n", color.HiGreenString("-----------------Simple Server-----------------"))
	defer server.StartServer()
}
