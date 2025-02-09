package d2

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Rows [][]int
}

func RunPart2() {
	file, fileCloseFn := getFile("/d2/input.txt")
	defer fileCloseFn()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	data := Data{
		Rows: [][]int{},
	}

	// Read each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by whitespace to get the two numbers
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		// Create a new row
		data.Rows = append(data.Rows, []int{})

		// store the parts as slice of integers
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				log.Fatalf("Failed to convert %s to int: %s", part, err)
			}
			data.Rows[len(data.Rows)-1] = append(data.Rows[len(data.Rows)-1], num)
		}
	}

	// // Output the numbers
	// for _, row := range data.Rows {
	// 	for _, num := range row {
	// 		fmt.Printf("%d ", num)
	// 	}
	// 	fmt.Printf("\n")
	// }

	safeRows := 0
	for _, row := range data.Rows {
		// Check if the row is already safe
		if isRowSafe(row) {
			safeRows++
			continue
		}

		// If not, try removing each element and check if the row becomes safe
		for j := 0; j < len(row); j++ {
			// fmt.Printf("removing %d th index\n", j)
			localRow := make([]int, len(row))
			copy(localRow, row)

			newRow := append(localRow[:j], localRow[j+1:]...)
			if isRowSafe(newRow) {
				// fmt.Printf("row is safe after removing %d th index\n", j)
				safeRows++
				break
			}
		}
	}

	fmt.Printf("safeRows: %d\n", safeRows)
}

func isRowSafe(row []int) bool {
	positives := 0
	negatives := 0

	for i := 1; i < len(row); i++ {
		diff := row[i] - row[i-1]
		if diff == 0 {
			return false
		}
		if diff > 0 && negatives > 0 {
			return false
		}
		if diff < 0 && positives > 0 {
			return false
		}

		if diff > 0 {
			if diff > 3 {
				return false
			}
			positives++
		} else {
			if diff < -3 {
				return false
			}
			negatives++
		}
	}

	return min(positives, negatives) == 0
}

func getFile(path string) (*os.File, func()) {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		// log.Fatal(err)
		return nil, nil
	}

	// Construct the full path to the file
	filePath := dir + path

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		// log.Fatalf("Failed to open file: %s", err)
		return nil, nil
	}

	return file, func() {
		file.Close()
	}
}

func RunPart1() {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the full path to the file
	filePath := dir + "/d2/input.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	data := Data{
		Rows: [][]int{},
	}

	// Read each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by whitespace to get the two numbers
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		// Create a new row
		data.Rows = append(data.Rows, []int{})

		// store the parts as slice of integers
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				log.Fatalf("Failed to convert %s to int: %s", part, err)
			}
			data.Rows[len(data.Rows)-1] = append(data.Rows[len(data.Rows)-1], num)
		}
	}

	// // Output the numbers
	// for _, row := range data.Rows {
	// 	for _, num := range row {
	// 		fmt.Printf("%d ", num)
	// 	}
	// 	fmt.Printf("\n")
	// }

	safeRows := 0
	// two checks:
	// 1. col's values must be either strictly increasing or strictly decreasing
	// 2. col's values must only differ by 1, 2, or 3 (it cannot be the same too)
	for _, row := range data.Rows {
		inc := true
		unsafelevels := 0

		if len(row) < 2 {
			continue
		}
		if row[0] > row[1] {
			inc = false

			if row[0]-row[1] > 3 || row[0] == row[1] {
				unsafelevels++
			}
		} else {
			if row[1]-row[0] > 3 || row[0] == row[1] {
				unsafelevels++
			}
		}

		// check the rest of the levels for the current row
		// entire row is unsafe if any of the levels is not safe
		for i := 1; i < len(row)-1; i++ {
			if inc {
				if row[i] > row[i+1] {
					unsafelevels++
					// log.Fatalf("Invalid input: %v", row)
				} else if row[i+1]-row[i] > 3 {
					unsafelevels++
					// log.Fatalf("Invalid input: %v", row)
				} else if row[i] == row[i+1] {
					unsafelevels++
					// log.Fatalf("Equal level values not allowed at row index %d and col index %d", rowInd, i)
				}
			} else {
				if row[i] < row[i+1] {
					unsafelevels++
					// log.Fatalf("Invalid input: %v", row)
				} else if row[i]-row[i+1] > 3 {
					unsafelevels++
					// log.Fatalf("Invalid input: %v", row)
				} else if row[i] == row[i+1] {
					unsafelevels++
					// log.Fatalf("Equal level values not allowed at row index %d and col index %d", rowInd, i)
				}
			}
		}

		if unsafelevels <= 1 {
			safeRows++
		}
	}

	fmt.Printf("safeRows: %d\n", safeRows)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
