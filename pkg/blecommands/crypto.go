package blecommands

import (
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

type Crypto interface {
	Encrypt(nonce, message []byte) ([]byte, error)
	Decrypt(nonce, ciphertext []byte) ([]byte, error)
}
type crypto struct {
	key []byte
}

func NewCrypto(sharedKey []byte) Crypto {
	return &crypto{
		key: sharedKey,
	}
}

func (c *crypto) Encrypt(nonce, message []byte) ([]byte, error) {
	var nonceArray [24]byte
	copy(nonceArray[:], nonce)

	var keyArray [32]byte
	copy(keyArray[:], c.key)

	encrypted := secretbox.Seal(nil, message, &nonceArray, &keyArray)
	return encrypted, nil
}

func (c *crypto) Decrypt(nonce, ciphertext []byte) ([]byte, error) {
	var nonceArray [24]byte
	copy(nonceArray[:], nonce)

	var keyArray [32]byte
	copy(keyArray[:], c.key)

	decrypted, ok := secretbox.Open(nil, ciphertext, &nonceArray, &keyArray)
	if !ok {
		return nil, fmt.Errorf("decryption failed")
	}

	return decrypted, nil
}

func CRC(val []byte) uint16 {
	crc := uint16(0xFFFF)
	poly := uint16(0x1021)

	for _, b := range val {
		for i := 0; i < 8; i++ {
			bit := (b>>(7-i))&1 == 1
			c15 := (crc>>15)&1 == 1
			crc <<= 1

			if c15 != bit {
				crc ^= poly
			}
		}
	}

	return crc & 0xFFFF
}
