package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

func Encrypt(path, file string) (string, []byte) {
	value := []byte(file)
	keyPhrase := generateKey()

	// Key phrase is hashed for increased security as assumed you are passing it in
	// As is randomly generated probably not required
	hashedKey := mdHashing(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedKey))
	if err != nil {
		fmt.Println(err)
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println(err)
	}

	// Encrypt
	cipheredText := gcmInstance.Seal(nonce, nonce, value, nil)

	// Store ciphered text

	// Return the key required to decrypt
	return hashedKey, cipheredText
}

// For hashing the key phrase
func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func generateKey() string {
	// TODO: make random
	return "placeholder"
}

func Decrypt(path string, hashedKey string, ciphered []byte) string {
	aesBlock, err := aes.NewCipher([]byte(hashedKey))
	if err != nil {
		log.Fatalln(err)
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}

	nonceSize := gcmInstance.NonceSize()
	// ciphered := []byte(getFileFromPath(path))

	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	fmt.Println(cipheredText)

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(originalText))

	return string(originalText)
}

func getFileFromPath(path string) string {
	// TODO: get file
	return "Lh\aʀ�D�@�[�d��\"~�ؖQڸ\u0019�啊k!�C|\ta/4�y��.f��/"
}
