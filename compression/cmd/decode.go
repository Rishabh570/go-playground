package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	decodeFilePath string
)

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a compressed file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Decoding file: %s\n", decodeFilePath)

		// Read the compressed file
		file, err := os.Open(decodeFilePath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		// Read the Huffman code table from the header
		scanner := bufio.NewScanner(file)
		huffmanCodeTable := make(map[string]rune)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "---" {
				break
			}
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				fmt.Printf("Invalid header format: %s\n", line)
				return
			}
			charCode, err := strconv.Atoi(parts[0])
			if err != nil {
				fmt.Printf("Error parsing character code: %v\n", err)
				return
			}
			huffmanCodeTable[parts[1]] = rune(charCode)
		}

		// Read the encoded content
		var encodedContent strings.Builder
		for scanner.Scan() {
			encodedContent.WriteString(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		// Decode the content
		decodedContent := decodeContent(encodedContent.String(), huffmanCodeTable)

		// Write the decoded content to a new file
		outputPath := strings.TrimSuffix(decodeFilePath, ".compressed") + ".decoded"
		err = os.WriteFile(outputPath, []byte(decodedContent), 0644)
		if err != nil {
			fmt.Printf("Error writing decoded file: %v\n", err)
			return
		}

		fmt.Printf("Decoded file written to: %s\n", outputPath)
	},
}

func decodeContent(encodedContent string, huffmanCodeTable map[string]rune) string {
	var decodedContent strings.Builder
	currentCode := ""

	for _, bit := range encodedContent {
		currentCode += string(bit)
		if char, found := huffmanCodeTable[currentCode]; found {
			decodedContent.WriteRune(char)
			currentCode = ""
		}
	}

	return decodedContent.String()
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	decodeCmd.Flags().StringVarP(&decodeFilePath, "file", "f", "output.compressed", "Path to the compressed file to decode")
	decodeCmd.MarkFlagRequired("file")
}
