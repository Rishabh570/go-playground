package blockcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"runtime"
	"sync"
)

// Encrypt encrypts the given plaintext using AES in CBC mode.
func EncryptConcurrent(plaintext, key []byte) ([]byte, error) {
	fmt.Printf("Encrypt received plaintext (len: %d): %x\n", len(plaintext), plaintext)
	// 1. Initialize cipher block
	blockcipher, err := InitCipher(key)
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
	blocks := SplitIntoBlocks(plaintext, blockSize)
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

				xoredBlock, err := XorBytes(currBlock, keystream)
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
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	return ciphertext, nil
}

// Encrypt encrypts the given plaintext using AES in CBC mode.
func Decrypt(cipher, key []byte) ([]byte, error) {
	fmt.Printf("Decrypt received cipherStr: %x\n", cipher)

	// 1. Initialize cipher block
	blockcipher, err := InitCipher(key)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	blockSize := blockcipher.BlockSize()
	fmt.Printf("Block size: %d\n", blockSize)

	fmt.Printf("Going to decrypt ciphertext: %x\n", cipher)

	// extract IV of size blockSize from the cipher
	ivec := cipher[:blockSize]
	fmt.Printf("IV: %x\n", ivec)
	fmt.Printf("IV length: %d\n", len(ivec))

	ciphertext := cipher[blockSize:]
	fmt.Printf("Ciphertext length: %d\n", len(ciphertext))
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// 4. Split plaintext into blocks
	blocks := SplitIntoBlocks(ciphertext, blockSize)
	fmt.Printf("Len of blocks: %d\n", len(blocks))

	// resulting ciphertext
	recoveredPlaintext := make([]byte, 0, len(ciphertext))

	fmt.Printf("Initialized recoveredPlaintext: %d\n", len(recoveredPlaintext))

	counterBlock := make([]byte, blockSize)
	copy(counterBlock, ivec)

	for i, block := range blocks {
		// Create a temporary buffer for decrypted block
		keystream := make([]byte, blockSize)
		blockcipher.Encrypt(keystream, counterBlock)

		// startInd := i * blockSize
		// endInd := min(startInd+blockSize, len(ciphertext))
		xoredBlock, err := XorBytes(block, keystream)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}
		fmt.Printf("Block %d - Decrypted (hex): %x, XORed (hex): %x\n", i, keystream, xoredBlock)
		// fmt.Printf("XORed block: %s\n", string(xoredBlock))

		// add the decrypted block to the recovered plaintext
		// startingIndex := i * blockSize
		// endingIndex := startingIndex + blockSize
		recoveredPlaintext = append(recoveredPlaintext, xoredBlock...)

		// increment counter
		counterBlock = incrementCounter(counterBlock)
	}

	fmt.Println("Decryption complete")
	fmt.Printf("Restored recoveredPlaintext (hex): %x\n", recoveredPlaintext)
	fmt.Printf("Restored recoveredPlaintext: %s\n", string(recoveredPlaintext))

	return recoveredPlaintext, nil
}

func incrementCounter(counterBlock []byte) []byte {
	// Start from the least significant byte (end of array)
	// and only increment the last 4 bytes
	for i := len(counterBlock) - 1; i >= len(counterBlock)-4; i-- {
		// Increment the current byte
		counterBlock[i]++

		// If no overflow occurred (byte didn't wrap to 0), we're done
		if counterBlock[i] != 0 {
			break
		}
		// Otherwise continue to the next byte (carry the 1)
	}
	return counterBlock
}

// InitCipher initializes a new AES cipher block with the provided key.
// The key must be either 16, 24, or 32 bytes long.
// If the key length is invalid, an error is returned.
// The function returns a cipher.Block interface that can be used for encryption and decryption.
func InitCipher(key []byte) (cipher.Block, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("initCipher: invalid key length %d", len(key))
	}
	return aes.NewCipher(key)
}

// XorBytes performs a bitwise XOR operation on two byte slices of equal length.
func XorBytes(a, b []byte) ([]byte, error) {
	fmt.Printf("XORing %x with %x\n", a, b)
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	result := make([]byte, len(a))
	for i := 0; i < n; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}
