package blockcipher

import "fmt"

// Encrypt encrypts the given plaintext using AES in CBC mode.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	// 1. Initialize cipher block
	blockcipher, err := initCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	blockSize := blockcipher.BlockSize()
	fmt.Printf("Block size: %d\n", blockSize)

	// 2. Generate an Initialization Vector (IV) to use as a 'previous ciphertext' for the first block
	// IV lenght is exactly equal to block size (for AES, it is 16 bytes)
	// IV is NOT secret, but MUST BE unique to ensure randomness for the same plaintext input
	ivec, err := generateIV(blockSize)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	prevCiphertext := ivec

	fmt.Printf("IV (len: %d): %x\n", len(prevCiphertext), prevCiphertext)

	// 3. Pad plaintext to multiple of block size
	paddedPlaintext := padPKCS5(plaintext, blockSize)
	fmt.Printf("Padded plaintext (len: %d): %x\n", len(paddedPlaintext), paddedPlaintext)

	// 4. Split plaintext into blocks
	blocks := splitIntoBlocks(paddedPlaintext, blockSize)
	fmt.Printf("blocks num for padded plaintext: %d\n", len(blocks))

	// resulting ciphertext
	ciphertext := make([]byte, len(paddedPlaintext)+blockSize)
	// Copy the IV to the beginning of the ciphertext
	copy(ciphertext[:blockSize], prevCiphertext)

	for i, block := range blocks {
		// XOR the current plaintext block with the previous ciphertext block
		// If it's the first block, it XORs with IV
		xoredBlock, err := xorBytes(block, prevCiphertext)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}
		fmt.Printf("XORed block: %x\n", xoredBlock)

		// 5. Encrypt the block
		startingIndex := blockSize + (i * blockSize)
		endingIndex := startingIndex + blockSize
		blockcipher.Encrypt(ciphertext[startingIndex:endingIndex], xoredBlock)

		// 6. Update the previous ciphertext block, for next iteration
		prevCiphertext = ciphertext[startingIndex:endingIndex]
	}

	fmt.Println("Encryption complete")
	return ciphertext, nil
}
