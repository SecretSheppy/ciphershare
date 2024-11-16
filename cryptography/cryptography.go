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

func Encrypt(value []byte) (string, []byte) {
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

	// Return the key required to decrypt and the cipher text
	return hashedKey, cipheredText
}

// For hashing the key phrase
func mdHashing(byteInput []byte) string {
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func generateKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println(err)
	}

	return key
}

func Decrypt(hashedKey string, ciphered []byte) string {
	aesBlock, err := aes.NewCipher([]byte(hashedKey))
	if err != nil {
		log.Fatalln(err)
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}

	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]

	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return string(originalText)
}
