package d4

import (
	"bufio"
	"fmt"
	"os"
)

func RunPart2() {
	file, fileCloseFn := getFile("/d4/input.txt")
	defer fileCloseFn()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	var linesArr [][]string
	// linesArr := make([][]string, 1000)
	ind := 0

	for scanner.Scan() {
		linesArr = append(linesArr, []string{})
		runes := []rune(scanner.Text())
		for _, r := range runes {
			linesArr[ind] = append(linesArr[ind], string(r))
		}

		ind++
	}

	n := len(linesArr)
	m := len(linesArr[0])
	answer := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if linesArr[i][j] == "A" {
				if i == 0 || j == 0 || i == n-1 || j == m-1 {
					continue // we cannot have answers on the boundary
				}
				if (linesArr[i-1][j-1] == "M" && linesArr[i+1][j+1] == "S") || (linesArr[i-1][j-1] == "S" && linesArr[i+1][j+1] == "M") {
					if (linesArr[i-1][j+1] == "M" && linesArr[i+1][j-1] == "S") || (linesArr[i-1][j+1] == "S" && linesArr[i+1][j-1] == "M") {
						answer++
					}
				}
			}
		}
	}

	fmt.Printf("Count of X-MAX: %d\n", answer)

}

func RunPart1() {
	file, fileCloseFn := getFile("/d4/input.txt")
	defer fileCloseFn()

	// Create a scanner to read through the file line by line
	scanner := bufio.NewScanner(file)

	var linesArr [][]string
	// linesArr := make([][]string, 1000)
	ind := 0

	for scanner.Scan() {
		linesArr = append(linesArr, []string{})
		runes := []rune(scanner.Text())
		for _, r := range runes {
			linesArr[ind] = append(linesArr[ind], string(r))
		}

		// fmt.Println(linesArr[ind])
		ind++
	}

	n := len(linesArr)
	m := len(linesArr[0])
	// fmt.Printf("n: %d, m: %d\n", n, m)
	dir := [][]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if linesArr[i][j] == "X" {
				// fmt.Printf("Found X at (%d, %d)\n", i, j)
				cnt += countOccurancesOfString(linesArr, i, j, "XMAS", 0, dir, "")
			}
		}
	}
	fmt.Printf("Count of XMAS: %d\n", cnt)
}

func countOccurancesOfString(arr [][]string, x, y int, str string, ind int, dir [][]int, path string) int {
	if ind >= len(str) {
		// fmt.Printf("path: %s\n", path)
		return 1 // Found the string
	}
	if x < 0 || y < 0 || x >= len(arr) || y >= len(arr[0]) {
		return 0
	}

	ans := 0
	if arr[x][y] == string(str[ind]) {
		// fmt.Printf("Found match at (%d, %d) for ind %d\n", x, y, ind)
		// path += fmt.Sprintf("(%d, %d) -> ", x, y)
		// fmt.Printf("Found %s at (%d, %d); checking neighbours for next char\n", string(str[ind]), x, y)
		for _, d := range dir {
			x1, y1 := x+d[0], y+d[1]
			// if x == 9 && y == 3 {
			// 	fmt.Printf("trying x1: %d, y1: %d", x1, y1)
			// 	if x1 >= 0 && y1 >= 0 && x1 < len(arr) && y1 < len(arr[0]) {
			// 		fmt.Printf("; val: %s\n", arr[x1][y1])
			// 	} else {
			// 		fmt.Printf("\n")
			// 	}
			// }
			ans += countOccurancesOfString(arr, x1, y1, str, ind+1, [][]int{d}, path)
		}
	}

	return ans
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
