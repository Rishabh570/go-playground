package cmd

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/rishabh570/aescbc/internal/blockcipher"
	"github.com/spf13/cobra"
)

var (
	ciphertext string
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Decrypting ciphertext: %x\n", ciphertext)
		fmt.Printf("Using key: %s\n", key)

		ciphertextStr, _ := decodeIfHex(ciphertext)
		keyConverted, _ := decodeIfHex(key)
		plaintext, err := blockcipher.Decrypt(ciphertextStr, keyConverted)
		if err != nil {
			fmt.Printf("Error decrypting data: %v\n", err)
			return
		}
		fmt.Println("======================================================================================================")
		fmt.Printf("Recovered original text: %s\n", string(plaintext))
		fmt.Println("======================================================================================================")
	},
}

func init() {
	decryptCmd.Flags().StringVarP(&key, "key", "k", "", "Encryption key")
	decryptCmd.Flags().StringVarP(&ciphertext, "ciphertext", "c", "", "Ciphertext")
	rootCmd.AddCommand(decryptCmd)
}

// Improved isHex function that works with strings of any length
func isHex(str string) bool {
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", str)
	return match
}

// decodeIfHex takes an input string and returns bytes:
// - If the input is a valid hex string, it decodes the hex to bytes
// - Otherwise, it returns the input as bytes directly
func decodeIfHex(input string) ([]byte, error) {
	if isHex(input) {
		return hex.DecodeString(input)
	}

	return []byte(input), nil
}
