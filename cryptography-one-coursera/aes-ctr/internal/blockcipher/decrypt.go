package blockcipher

import "fmt"

// Decrypt decrypts the ciphertext using the AES CTR mode of operation.
func Decrypt(cipher, key []byte) ([]byte, error) {
	fmt.Printf("Decrypt received cipherStr: %x\n", cipher)

	// 1. Initialize cipher block
	blockcipher, err := initCipher(key)
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
	blocks := splitIntoBlocks(ciphertext, blockSize)
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

		xoredBlock, err := xorBytes(block, keystream)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}
		fmt.Printf("Block %d - Decrypted (hex): %x, XORed (hex): %x\n", i, keystream, xoredBlock)

		recoveredPlaintext = append(recoveredPlaintext, xoredBlock...)

		// increment counter
		counterBlock = incrementCounter(counterBlock)
	}

	fmt.Println("Decryption complete")
	return recoveredPlaintext, nil
}
