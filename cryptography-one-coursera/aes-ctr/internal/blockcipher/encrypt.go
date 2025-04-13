package blockcipher

import (
	"fmt"
	"runtime"
	"sync"
)

// EncryptConcurrent encrypts the plaintext using AES in CTR mode.
func EncryptConcurrent(plaintext, key []byte) ([]byte, error) {
	fmt.Printf("Encrypt received plaintext (len: %d): %x\n", len(plaintext), plaintext)
	// 1. Initialize cipher block
	blockcipher, err := initCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	blockSize := blockcipher.BlockSize()
	fmt.Printf("Block size: %d\n", blockSize)

	// 2. Generate an Initialization Vector (IV)
	// IV lenght is exactly equal to block size (for AES, it is 16 bytes)
	// IV is NOT secret, but MUST BE unique to ensure randomness for the same plaintext input (PRF theory)
	counterBlock := generateIncrementableNonce()
	fmt.Printf("counter block (hex): %x\n", counterBlock)

	// 4. Split plaintext into blocks
	blocks := splitIntoBlocks(plaintext, blockSize)
	fmt.Printf("Len of blocks: %d\n", len(blocks))

	// resulting ciphertext
	ciphertext := make([]byte, len(plaintext)+blockSize)
	// Copy the IV to the beginning of the ciphertext, needed for decryption
	copy(ciphertext[:blockSize], counterBlock)

	fmt.Printf("initialized ciphertext len: %d\n", len(ciphertext))

	// define wait group
	var wg sync.WaitGroup

	// Calculate blocks per worker
	numWorkers := runtime.NumCPU()
	numBlocks := (len(plaintext) + blockSize - 1) / blockSize
	blocksPerWorker := (numBlocks + numWorkers - 1) / numWorkers

	fmt.Printf("Number of workers: %d\n", numWorkers)
	fmt.Printf("Blocks per worker: %d\n", blocksPerWorker)

	// We want to process blocksPerWorker blocks per goroutine
	// We process blocksPerWorker sequentially in each goroutine, since it's a CPU-bound task
	// We run numWorksers goroutines, each processing blocksPerWorker blocks sequentially
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			// Global start and end index for blocks
			blocksStartInd := i * blocksPerWorker
			blocksEndInd := min((i+1)*blocksPerWorker, numBlocks)

			if blocksStartInd >= numBlocks {
				return
			}

			// Create a local counter for this goroutine
			// makes it easy to manage the counter state for each goroutine
			localCounter := make([]byte, blockSize)
			copy(localCounter, counterBlock)
			// Advance counter to starting position for this goroutine's first block
			for j := 0; j < blocksStartInd; j++ {
				localCounter = incrementCounter(localCounter)
			}

			for j := blocksStartInd; j < blocksEndInd; j++ {
				currBlock := blocks[j]
				fmt.Printf("Block %d - Current block (hex): %x, counter: %x\n", j, currBlock, localCounter)

				keystream := make([]byte, blockSize)
				blockcipher.Encrypt(keystream, localCounter)

				xoredBlock, err := xorBytes(currBlock, keystream)
				if err != nil {
					fmt.Printf("XOR error: %v\n", err)
					return
				}
				fmt.Printf("XORed block: %x\n", xoredBlock)

				// assign xoredBlock to the ciphertext
				startingIndex := blockSize + (j * blockSize)
				endingIndex := startingIndex + len(xoredBlock)
				copy(ciphertext[startingIndex:endingIndex], xoredBlock)

				localCounter = incrementCounter(localCounter)
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("Encryption complete")
	return ciphertext, nil
}
