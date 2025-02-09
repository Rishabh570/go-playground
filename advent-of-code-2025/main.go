package main

import (
	"aoc25/d4"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Running AoC 2025 code...\n")

	start := time.Now()

	// Run the code
	d4.RunPart2()

	end := time.Now()

	duration := end.Sub(start)

	fmt.Printf("Time taken by myFunction: %v\n", duration)
}
