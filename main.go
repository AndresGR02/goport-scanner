package main

import (
	"fmt"

	"github.com/AndresGR02/goport-scanner/port"
)

func main() {
	results := port.InitialScan("localhost")
	fmt.Println(results)
}
