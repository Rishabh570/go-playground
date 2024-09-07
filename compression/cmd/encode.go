package cmd

import (
	"fmt"
	"os"

	"container/heap"

	"github.com/spf13/cobra"
)

var (
	filePath string
)

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode a file using Huffman coding",
	Run: func(cmd *cobra.Command, args []string) {
		// Read the file
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		// Create a map to store character frequencies
		charFreq := make(map[rune]int)

		// Count the frequency of each character
		for _, char := range string(content) {
			charFreq[char]++
		}

		// create a huffman tree based on the frequency table
		huffmanTree := createHuffmanTree(charFreq)

		// create a huffman code table
		huffmanCodeTable := createHuffmanCodeTable(huffmanTree)

		// encode the file using the huffman code table
		encodedFile := encodeFile(content, huffmanCodeTable)

		// Write the compressed file with header
		outputPath := "output.compressed"
		err = writeCompressedFile(outputPath, huffmanCodeTable, encodedFile)
		if err != nil {
			fmt.Printf("Error writing compressed file: %v\n", err)
			return
		}
		fmt.Printf("Compressed file written to: %s\n", outputPath)
	},
}

type huffmanNode struct {
	Char  rune
	Freq  int
	Left  *huffmanNode
	Right *huffmanNode
}

func createHuffmanTree(charFreq map[rune]int) *huffmanNode {
	// Create a priority queue based on character frequencies
	pq := make(priorityQueue, 0)
	for char, freq := range charFreq {
		pq = append(pq, &huffmanNode{Char: char, Freq: freq})
	}
	heap.Init(&pq)

	// Build the Huffman tree
	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*huffmanNode)
		right := heap.Pop(&pq).(*huffmanNode)
		parent := &huffmanNode{Freq: left.Freq + right.Freq, Left: left, Right: right}
		heap.Push(&pq, parent)
	}

	return heap.Pop(&pq).(*huffmanNode)
}

func createHuffmanCodeTable(root *huffmanNode) map[rune]string {
	huffmanCodeTable := make(map[rune]string)
	buildHuffmanCodeTable(root, "", huffmanCodeTable)
	return huffmanCodeTable
}

func buildHuffmanCodeTable(node *huffmanNode, code string, huffmanCodeTable map[rune]string) {
	if node == nil {
		return
	}
	if node.Char != 0 {
		huffmanCodeTable[node.Char] = code
	}
	buildHuffmanCodeTable(node.Left, code+"0", huffmanCodeTable)
	buildHuffmanCodeTable(node.Right, code+"1", huffmanCodeTable)
}

func encodeFile(content []byte, huffmanCodeTable map[rune]string) []byte {
	encodedFile := ""
	for _, char := range content {
		encodedFile += huffmanCodeTable[rune(char)]
	}
	return []byte(encodedFile)
}

// Implement priorityQueue methods (Len, Less, Swap, Push, Pop)
type priorityQueue []*huffmanNode

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].Freq < pq[j].Freq }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*huffmanNode))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func init() {
	rootCmd.AddCommand(encodeCmd)
	encodeCmd.Flags().StringVarP(&filePath, "file", "f", "data.txt", "Path to the file to encode")
}

func writeCompressedFile(outputPath string, huffmanCodeTable map[rune]string, encodedFile []byte) error {
	fmt.Printf("Attempting to create file: %s\n", outputPath)
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	// Write header (Huffman code table)
	for char, code := range huffmanCodeTable {
		_, err := fmt.Fprintf(file, "%d:%s\n", char, code)
		if err != nil {
			return err
		}
	}

	// Write separator between header and encoded content
	_, err = file.WriteString("---\n")
	if err != nil {
		fmt.Printf("Error writing separator: %v\n", err)
		return err
	}

	// Write encoded content
	_, err = file.Write(encodedFile)
	if err != nil {
		fmt.Printf("Error writing encoded content: %v\n", err)
		return err
	}

	return nil
}
