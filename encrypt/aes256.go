package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// DecodeAES decode data with aes-256/CBC/PKCS7
func AESDecode(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	iv := make([]byte, blockSize)
	blockMode := cipher.NewCBCDecrypter(block, iv)

	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)
	length := len(origData)
	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)], nil
}

// AESEncrypt str with AES256-CBC, padding with PKCS7
func AESEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	ivT := make([]byte, aes.BlockSize+len(plaintext))
	// initialization vector
	iv := ivT[:aes.BlockSize]

	// block
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	blockSize := block.BlockSize()

	// PKCS7 padding
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padtext...)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)

	return crypted, nil
}

// Use the AESEncrypt and output in Base64
func AESEncryptOutInBase64(plaintext []byte, key []byte) ([]byte, error) {
	content, err := AESEncrypt(plaintext, key)
	if err != nil {
		return []byte{}, err
	}

	return []byte(Base64Encode(content)), nil
}

func AESDecodeAfterBase64(data []byte, key []byte) ([]byte, error) {
	content, err := Base64Decode(string(data))
	if err != nil {
		return []byte{}, err
	}

	plainText, err := AESDecode([]byte(content), key)
	if err != nil {
		return []byte{}, err
	}

	return plainText, nil
}
