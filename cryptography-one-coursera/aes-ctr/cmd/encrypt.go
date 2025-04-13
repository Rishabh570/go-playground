package cmd

import (
	"fmt"

	"github.com/rishabh570/aesctr/internal/blockcipher"
	"github.com/spf13/cobra"
)

var (
	plaintext string
	key       string
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Encrypting data: %s\n", plaintext)
		fmt.Printf("Using key: %s\n", key)
		ciphertext, err := blockcipher.Encrypt([]byte(plaintext), []byte(key))
		if err != nil {
			fmt.Printf("Error encrypting data: %v\n", err)
			return
		}
		fmt.Printf("Ciphertext received: %x\n", string(ciphertext))
	},
}

func init() {
	encryptCmd.Flags().StringVarP(&key, "key", "k", "", "Encryption key")
	encryptCmd.Flags().StringVarP(&plaintext, "plaintext", "t", "", "Plaintext data to encrypt")
	rootCmd.AddCommand(encryptCmd)
}
