package d1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Data struct {
	Left  []int
	Right []int
}

func (d *Data) UnmarshalJSON(data []byte) error {
	var temp any
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	fmt.Printf("temp: %v\n", temp)

	return nil
}

func RunPart2() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the full path to the file
	filePath := dir + "/d1/input.txt"

	// read input.txt file
	// data, err := os.ReadFile(filePath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	var data Data

	// Read each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by whitespace to get the two numbers
		parts := strings.Fields(line)

		// Ensure each line has exactly two parts
		if len(parts) != 2 {
			log.Fatalf("Unexpected line format: %s", line)
		}

		// Convert each part to an integer
		num1, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalf("Failed to convert first number on line %s: %s", line, err)
		}

		num2, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Failed to convert second number on line %s: %s", line, err)
		}

		// Output the numbers
		fmt.Printf("Line: %s, Number 1: %d, Number 2: %d\n", line, num1, num2)

		data.Left = append(data.Left, num1)
		data.Right = append(data.Right, num2)
	}

	// Check for potential errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	fmt.Printf("Data: %+v\n", data)

	rightListFreq := make(map[int]int)
	for _, v := range data.Right {
		rightListFreq[v]++
	}

	result := 0
	for _, v := range data.Left {
		if rightListFreq[v] > 0 {
			result += (v * rightListFreq[v])
		}
	}

	fmt.Printf("Result: %d\n", result)

	// sum := 0
	// for i := 0; i < len(data.Left); i++ {
	// 	fmt.Printf("Left: %d, Right: %d\n", data.Left[i], data.Right[i])
	// 	diff := data.Right[i] - data.Left[i]
	// 	if diff < 0 {
	// 		diff = diff * -1
	// 	}
	// 	fmt.Printf("Diff: %d\n", diff)
	// 	sum += diff
	// }

	// fmt.Printf("Sum: %d\n", sum)
}

func RunPart1() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the full path to the file
	filePath := dir + "/d1/input.txt"

	// read input.txt file
	// data, err := os.ReadFile(filePath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	var data Data

	// Read each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by whitespace to get the two numbers
		parts := strings.Fields(line)

		// Ensure each line has exactly two parts
		if len(parts) != 2 {
			log.Fatalf("Unexpected line format: %s", line)
		}

		// Convert each part to an integer
		num1, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalf("Failed to convert first number on line %s: %s", line, err)
		}

		num2, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Failed to convert second number on line %s: %s", line, err)
		}

		// Output the numbers
		fmt.Printf("Line: %s, Number 1: %d, Number 2: %d\n", line, num1, num2)

		data.Left = append(data.Left, num1)
		data.Right = append(data.Right, num2)
	}

	// Check for potential errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	fmt.Printf("Data: %+v\n", data)

	// Sort the Left and Right slices
	sort.Ints(data.Left)
	sort.Ints(data.Right)

	fmt.Printf("Sorted Data: %+v\n", data)

	sum := 0
	for i := 0; i < len(data.Left); i++ {
		fmt.Printf("Left: %d, Right: %d\n", data.Left[i], data.Right[i])
		diff := data.Right[i] - data.Left[i]
		if diff < 0 {
			diff = diff * -1
		}
		fmt.Printf("Diff: %d\n", diff)
		sum += diff
	}

	fmt.Printf("Sum: %d\n", sum)
}
