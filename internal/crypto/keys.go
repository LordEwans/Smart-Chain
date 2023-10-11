package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	seedLen = 32

	privKeyLen = 64
	pubKeyLen  = 32

	addressLen = 20
)

// Private Key
type PrivateKey struct {
	key ed25519.PrivateKey
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, pubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, seedLen)

	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func GeneratePrivateKeyFromString(s string) *PrivateKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedLen {
		panic("invalid seed, length must be 32")
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

// Public Key
type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

func (p *PublicKey) Address() *Address {
	return &Address{
		value: p.key[len(p.key)-addressLen:],
	}
}

// Signature
type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(p *PublicKey, msg []byte) bool {
	return ed25519.Verify(p.key, msg, s.value)
}

// Address
type Address struct {
	value []byte
}

func (a *Address) Bytes() []byte {
	return a.value
}

func (a *Address) String() string {
	return hex.EncodeToString(a.value)
}
