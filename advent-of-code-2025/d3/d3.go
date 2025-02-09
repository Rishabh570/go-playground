package d3

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunPart1() {
	file, fileCloseFn := getFile("/d3/input.txt")
	defer fileCloseFn()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	result := 0
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by whitespace to get the two numbers
		parts := strings.Fields(line)

		// mul(1,2), mul(123,456)
		for _, part := range parts {
			for i := 0; i < len(part); i++ {
				if i+8 < len(part) && (part[i] == 'm' && part[i+1] == 'u' && part[i+2] == 'l') {
					fmt.Printf("mul found at index %d\n", i)

					if part[i+3] != '(' {
						fmt.Printf("expected '(' char at index: %d\n", i+3)
						continue
					}

					num1, endIdx, err := extractNum(part, i+4, ',')
					if err != nil {
						fmt.Printf("Error: %+v\n", err)
						continue
					}
					fmt.Printf("num1: %d, endIdx: %d\n", num1, endIdx)
					if endIdx >= len(part) {
						fmt.Printf("endIdx out of bounds\n")
					}
					// if endIdx != ',' {
					// 	fmt.Printf("expected ',' at index %d\n", endIdx)
					// 	continue
					// }
					// endIdx++

					num2, endIdx2, err := extractNum(part, endIdx, ')')
					if err != nil {
						fmt.Printf("Error: %+v\n", err)
						continue
					}
					fmt.Printf("num2: %d, endIdx2: %d\n", num2, endIdx2)

					result += num1 * num2
					i = endIdx2 - 1
				}
			}
		}
	}

	fmt.Printf("Result: %d\n", result)
}

func RunPart2() {
	file, fileCloseFn := getFile("/d3/input.txt")
	defer fileCloseFn()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	result := 0
	var unifiedStrs []string
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by whitespace to get the two numbers
		parts := strings.Fields(line)

		// join parts into a single string
		unifiedStrs = append(unifiedStrs, strings.Join(parts, ""))
	}
	globalInputStr := strings.Join(unifiedStrs, "")

	// mul(1,2), mul(123,456)
	// for _, part := range parts {
	part := globalInputStr
	fmt.Printf("part: %s\n", part)
	i := 0
	value, endIdx := getResultsUntilNextDont(part, i)
	fmt.Printf("✅ value: %d, endIdx: %d\n", value, endIdx)
	i = endIdx
	result += value

	for i < len(part) {
		fmt.Printf("i: %d\n", i)
		// if "do()" is encountered, continue to look for mul(x,y) instructions
		// if "don't()" is seen, continue to next iteration of the loop until do() is seen
		if i+3 < len(part) && (part[i] == 'd' && part[i+1] == 'o' && part[i+2] == '(' && part[i+3] == ')') {
			fmt.Printf("do() found at index %d\n", i)
			value, endIdx := getResultsUntilNextDont(part, i+4)
			fmt.Printf("✅✅ value: %d, endIdx: %d\n", value, endIdx)
			i = endIdx - 1
			result += value
		}
		i++
	}
	// }

	fmt.Printf("Result: %d\n", result)
}

func checkDont(part string, i int) bool {
	if i+6 < len(part) && (part[i] == 'd' && part[i+1] == 'o' && part[i+2] == 'n' && part[i+3] == '\'' && part[i+4] == 't' && part[i+5] == '(' && part[i+6] == ')') {
		return true
	}
	return false
}

func getResultsUntilNextDont(part string, i int) (int, int) {
	if i >= len(part) {
		return 0, i
	}

	if checkDont(part, i) {
		fmt.Printf("don't found at index %d\n", i+7)
		fmt.Printf("found for part: %s\n", part)
		return 0, i + 7
	}

	if i+8 < len(part) && (part[i] == 'm' && part[i+1] == 'u' && part[i+2] == 'l') {
		fmt.Printf("mul found at index %d\n", i)

		if part[i+3] != '(' {
			fmt.Printf("expected '(' char at index: %d\n", i+3)
			return getResultsUntilNextDont(part, i+3)
		}

		num1, endIdx, err := extractNum(part, i+4, ',')
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return getResultsUntilNextDont(part, endIdx)
		}
		fmt.Printf("num1: %d, endIdx: %d\n", num1, endIdx)
		// if endIdx >= len(part) {
		// 	fmt.Printf("endIdx out of bounds\n")
		// }
		// if endIdx != ',' {
		// 	fmt.Printf("expected ',' at index %d\n", endIdx)
		// 	continue
		// }
		// endIdx++

		num2, endIdx2, err := extractNum(part, endIdx, ')')
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return getResultsUntilNextDont(part, endIdx2)
		}
		fmt.Printf("num2: %d, endIdx2: %d\n", num2, endIdx2)

		value, endIdx3 := getResultsUntilNextDont(part, endIdx2)
		return (num1 * num2) + value, endIdx3

		// i = endIdx2 - 1
	}

	return getResultsUntilNextDont(part, i+1)
}

func extractNum(str string, start int, terminalChar byte) (int, int, error) {
	num := 0
	if start >= len(str) {
		return 0, -1, fmt.Errorf("start index out of bounds")
	}

	// if str[start] != '(' {
	// 	return 0, -1, fmt.Errorf("expected '(' at index %d", start)
	// }
	// start++ // skip '('

	if !isDigit(str[start]) {
		return 0, start, fmt.Errorf("expected digit at index %d", start)
	}

	for start < len(str) && isDigit(str[start]) {
		num = num*10 + int(str[start]-'0')
		start++
	}

	if start >= len(str) || str[start] != terminalChar {
		return 0, start, fmt.Errorf("expected '%s' at index %d", terminalChar, start)
	}
	start++ // skip terminalChar
	return num, start, nil
}

func isDigit(num byte) bool {
	if num >= '0' && num <= '9' {
		return true
	}
	return false
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
