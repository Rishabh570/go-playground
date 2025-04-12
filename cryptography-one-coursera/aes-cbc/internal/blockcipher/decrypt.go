package blockcipher

import "fmt"

// Encrypt encrypts the given plaintext using AES in CBC mode.
func Decrypt(cipher, key []byte) ([]byte, error) {
	// 1. Initialize cipher block
	blockcipher, err := initCipher(key)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	blockSize := blockcipher.BlockSize()
	fmt.Printf("Block size: %d\n", blockSize)

	// extract IV of size blockSize from the cipher
	ivec := cipher[:blockSize]
	fmt.Printf("IV (len: %d): %x\n", len(ivec), ivec)
	prevCiphertext := ivec

	ciphertext := cipher[blockSize:]
	fmt.Printf("Ciphertext (len: %d): %x\n", len(ciphertext), ciphertext)

	// 4. Split plaintext into blocks
	blocks := splitIntoBlocks(ciphertext, blockSize)
	fmt.Printf("Number of blocks for ciphertext: %d\n", len(blocks))

	recoveredPlaintext := make([]byte, 0, len(ciphertext))

	for i, block := range blocks {
		// Create a temporary buffer for decrypted block
		decrypted := make([]byte, blockSize)
		blockcipher.Decrypt(decrypted, block)

		xoredBlock, err := xorBytes(decrypted, prevCiphertext)
		if err != nil {
			return nil, fmt.Errorf("encrypt: %w", err)
		}
		fmt.Printf("Block %d - Decrypted (hex): %x, XORed (hex): %x\n", i, decrypted, xoredBlock)

		// add the decrypted block to the recovered plaintext
		recoveredPlaintext = append(recoveredPlaintext, xoredBlock...)

		// update the previous ciphertext block to the current ciphertext block
		prevCiphertext = block
	}

	// Remove PKCS#5 un-padding
	unpaddedRecoveredPlaintext, err := unpadPKCS5(recoveredPlaintext)
	if err != nil {
		return nil, fmt.Errorf("invalid padding: %v", err)
	}

	fmt.Println("Decryption complete")
	return unpaddedRecoveredPlaintext, nil
}
