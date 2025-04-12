package blockcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

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
	if len(a) != len(b) {
		return nil, fmt.Errorf("xorBytes: length mismatch")
	}
	result := make([]byte, len(a))
	for i := range a {
		result[i] = a[i] ^ b[i]
	}
	return result, nil
}

// splitIntoBlocks splits the given plaintext into blocks of the specified size.
// It assumes that the plaintext is already padded to a multiple of blockSize.
func splitIntoBlocks(plaintext []byte, blockSize int) [][]byte {
	// plaintext is already padded to a multiple of blockSize at this stage, therefore we can simply divide
	numBlocks := len(plaintext) / blockSize

	blocks := make([][]byte, numBlocks)

	// Divide plaintext into blocks
	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := start + blockSize
		blocks[i] = plaintext[start:end]
	}

	return blocks
}

// generateIV generates a random initialization vector (IV) of the specified block size.
// The IV is used in block cipher modes of operation like CBC to ensure that the same plaintext
func generateIV(blockSize int) ([]byte, error) {
	iv := make([]byte, blockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("could not generate IV: %v", err)
	}
	return iv, nil
}
