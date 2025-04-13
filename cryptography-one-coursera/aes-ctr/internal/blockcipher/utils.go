package blockcipher

import (
	"crypto/rand"
	"fmt"
)

func SplitIntoBlocks(plaintext []byte, blockSize int) [][]byte {
	// Calculate how many blocks we need
	// numBlocks := (len(plaintext) + blockSize - 1) / blockSize

	// plaintext is already padded to a multiple of blockSize at this stage, therefore we can simply divide
	// numBlocks := len(plaintext) / blockSize
	numBlocks := (len(plaintext) + blockSize - 1) / blockSize

	fmt.Printf("num blocks: %d\n", numBlocks)

	// Create slice to hold blocks
	blocks := make([][]byte, numBlocks)

	// Divide plaintext into blocks
	for i := 0; i < numBlocks; i++ {
		// Calculate start and end positions
		start := i * blockSize
		end := start + blockSize

		// Handle the last block (might be smaller than blockSize)
		if end > len(plaintext) {
			end = len(plaintext)
		}

		// Create the block
		fmt.Printf("creating a block from st: %d to end: %d\n", start, end)
		blocks[i] = plaintext[start:end]

		fmt.Printf("block %d: %x\n", i, blocks[i])
	}

	fmt.Printf("blocks: %x\n", blocks)
	return blocks
}

// Create a nonce with space for counter
func generateIncrementableNonce() []byte {
	// Create 16-byte slice (128 bits for AES block size)
	nonce := make([]byte, 16)

	// Fill first 12 bytes with random data
	if _, err := rand.Read(nonce[:12]); err != nil {
		panic(err)
	}

	// Last 4 bytes are left as zeros (initial counter value)
	return nonce
}
