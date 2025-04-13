package blockcipher

import (
	"bytes"
	"fmt"
)

// use PKCS5 to create padding functions: PadPKCS5 and UnpadPKCS5
// PKCS5 padding is a method of adding padding to data to make it a multiple of the block size.
// PKCS5 padding is used in block ciphers like DES and AES.
// PKCS5 padding works by adding a number of bytes, each of which is the number of bytes added.
// For example, if the block size is 8 bytes and the data is 5 bytes long, 3 bytes of padding will be added.
// The padding will be 0x03 0x03 0x03.
// If the data is already a multiple of the block size, 8 bytes of padding will be added.
// The padding will be 0x08 0x08 0x08 0x08 0x08 0x08 0x08 0x08.

func PadPKCS5(data []byte, blockSize int) []byte {
	fmt.Printf("PadPKCS5: data: %x\n", data)
	fmt.Printf("len(data): %d\n", len(data))
	mod := len(data) % blockSize
	fmt.Printf("mod: %d\n", mod)
	padding := blockSize - mod
	fmt.Printf("padding calculated: %d\n", padding)
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	fmt.Printf("pad: %x\n", pad)
	fmt.Printf("len(pad): %d\n", len(pad))
	fmt.Printf("len(data): %d\n", len(data))
	fmt.Printf("len(data) + len(pad): %d\n", len(data)+len(pad))
	return append(data, pad...)
}

func UnpadPKCS5(data []byte) ([]byte, error) {
	fmt.Printf("UnpadPKCS5: data: %x\n", data)
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	padding := data[len(data)-1]
	fmt.Printf("int(padding): %d\n", int(padding))
	fmt.Printf("padding: %d\n", padding)
	fmt.Printf("len(data): %d\n", len(data))
	if int(padding) > len(data) || padding == 0 {
		fmt.Printf("padding: %d, len(data): %d\n", padding, len(data))
		return nil, fmt.Errorf("no padding or more than data length")
	}
	for i := len(data) - int(padding); i < len(data); i++ {
		if data[i] != padding {
			return nil, fmt.Errorf("not all padding bytes are the same")
		}
	}
	return data[:len(data)-int(padding)], nil
}
