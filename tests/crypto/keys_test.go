package crypto_test

import (
	"fmt"
	"testing"

	"github.com/BottleHub/Smart-Chain/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), 64)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), 32)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("Parse")

	sign := privKey.Sign(msg)
	assert.True(t, sign.Verify(pubKey, msg))

	// Invalid msg
	assert.False(t, sign.Verify(pubKey, []byte("Nil")))

	// Invalid pubKey
	invalidPrivKey := crypto.GeneratePrivateKey()
	invalidPubKey := invalidPrivKey.Public()
	assert.False(t, sign.Verify(invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()
	assert.Equal(t, len(address.Bytes()), 20)
	fmt.Println(address)
}

func TestGenerateKeyFromString(t *testing.T) {
	seed, addressString := "8e0d9d1838a6ba4573ed927d8ba9423d5e07883637a64e29115536104ae1b6d5", "728c3860a5d4f2d2ee39b24a9dd23aa02a4b57d0"
	privKey := crypto.GeneratePrivateKeyFromString(seed)
	assert.Equal(t, len(privKey.Bytes()), 64)

	address := privKey.Public().Address()
	assert.Equal(t, addressString, address.String())
}
