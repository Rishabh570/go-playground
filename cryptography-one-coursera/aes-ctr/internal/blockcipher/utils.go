package blockcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// splitIntoBlocks splits the given plaintext into blocks of the specified size.
func splitIntoBlocks(plaintext []byte, blockSize int) [][]byte {
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

// incrementCounter increments the last 4 bytes of the counterBlock
// and returns the updated counterBlock.
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

// initCipher initializes a new AES cipher block with the provided key.
// The key must be either 16, 24, or 32 bytes long.
// If the key length is invalid, an error is returned.
// The function returns a cipher.Block interface that can be used for encryption and decryption.
func initCipher(key []byte) (cipher.Block, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("initCipher: invalid key length %d", len(key))
	}
	return aes.NewCipher(key)
}

// xorBytes performs a bitwise XOR operation on two byte slices of equal length.
func xorBytes(a, b []byte) ([]byte, error) {
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
