package blockcipher

import (
	"bytes"
	"fmt"
)

// use PKCS5 to create padding functions: padPKCS5 and unpadPKCS5
// PKCS5 padding is a method of adding padding to data to make it a multiple of the block size.
// PKCS5 padding is used in block ciphers like DES and AES.
// PKCS5 padding works by adding a number of bytes, each of which is the number of bytes added.
// For example, if the block size is 8 bytes and the data is 5 bytes long, 3 bytes of padding will be added.
// The padding will be 0x03 0x03 0x03.
// If the data is already a multiple of the block size, 8 bytes of padding will be added.
// The padding will be 0x08 0x08 0x08 0x08 0x08 0x08 0x08 0x08.
// This padding value is used during decryption to remove the padding.

func padPKCS5(data []byte, blockSize int) []byte {
	mod := len(data) % blockSize
	padding := blockSize - mod
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}

func unpadPKCS5(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	padding := data[len(data)-1]
	if int(padding) > len(data) || padding == 0 {
		return nil, fmt.Errorf("no padding or more than data length")
	}
	for i := len(data) - int(padding); i < len(data); i++ {
		if data[i] != padding {
			return nil, fmt.Errorf("not all padding bytes are the same")
		}
	}
	return data[:len(data)-int(padding)], nil
}
