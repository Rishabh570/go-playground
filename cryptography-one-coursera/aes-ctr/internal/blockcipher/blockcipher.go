package blockcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// Encrypt encrypts the given plaintext using AES in CBC mode.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	// 1. Initialize cipher block
	blockcipher, err := InitCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}
	blockSize := blockcipher.BlockSize()
	fmt.Printf("Block size: %d\n", blockSize)

	// 2. Generate an Initialization Vector (IV) to use as a 'previous ciphertext' for the first block
	// IV lenght is exactly equal to block size (for AES, it is 16 bytes)
	// IV is NOT secret, but MUST BE unique to ensure randomness for the same plaintext input
	counterBlock := generateIncrementableNonce()

	fmt.Printf("counter block (hex): %x\n", counterBlock)

	// 4. Split plaintext into blocks
	blocks := SplitIntoBlocks(plaintext, blockSize)
	fmt.Printf("Len of blocks: %d\n", len(blocks))
	fmt.Printf("len of plaintext: %d\n", len(plaintext))

	// resulting ciphertext
	ciphertext := make([]byte, len(plaintext)+blockSize)
	// Copy the IV to the beginning of the ciphertext
	copy(ciphertext[:blockSize], counterBlock)

	fmt.Printf("initialized ciphertext len: %d\n", len(ciphertext))

	// define wait group
	// var wg sync.WaitGroup

	// cipherCh := make(chan []byte, len(blocks))

	for i, block := range blocks {
		currentBlock := block

		fmt.Printf("Block %d - Current block (hex): %x, counter: %x\n", i, currentBlock, counterBlock)

		// wg.Add(1)
		// go func(ind int, currentBlock []byte, counterBlock []byte) {
		// var result []byte
		// defer wg.Done()

		// fmt.Printf("Processing block %d, counter: %x\n", ind, counterBlock)

		keystream := make([]byte, blockSize)
		// need to increment the prevCiphertext for every encryption round
		blockcipher.Encrypt(keystream, counterBlock)

		// st := len(currentBlock)
		// fmt.Printf("st: %d, keystream len: %d\n", st, len(keystream))
		// keystreamModified := keystream[:st]
		// fmt.Printf("Keystream (hex): %x\n", len(keystreamModified))
		xoredBlock, err := XorBytes(currentBlock, keystream)
		if err != nil {
			fmt.Printf("XOR error: %v\n", err)
			// return
			// return nil, fmt.Errorf("encrypt: %w", err)
		}
		fmt.Printf("XORed block: %x\n", xoredBlock)

		// assign xoredBlock to the ciphertext
		startingIndex := blockSize + (i * blockSize)
		endingIndex := startingIndex + len(xoredBlock)
		fmt.Printf("Starting index: %d, Ending index: %d\n", startingIndex, endingIndex)
		copy(ciphertext[startingIndex:endingIndex], xoredBlock)
		// cipherCh <- xoredBlock

		// }(i, currentBlock, counterBlock)

		// increment counter
		counterBlock = incrementCounter(counterBlock)
	}

	// wg.Wait()

	// TODO: fetch from channel and concatenate to form ciphertext
	// read all values from buffered channel
	// close(cipherCh)
	// for i := 0; i < len(blocks); i++ {
	// 	select {
	// 	case xoredBlock, ok := <-cipherCh:
	// 		if !ok {
	// 			fmt.Println("Channel closed")
	// 			break
	// 		}
	// 		fmt.Printf("XORed block from channel: %x\n", xoredBlock)
	// 		// startingIndex := blockSize + (i * blockSize)
	// 		// endingIndex := startingIndex + blockSize
	// 		// copy(ciphertext[startingIndex:endingIndex], xoredBlock)
	// 		ciphertext = append(ciphertext, xoredBlock...)
	// 	default:
	// 		fmt.Println("No more data in channel")
	// 	}
	// }

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
