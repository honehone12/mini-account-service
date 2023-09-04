package pwhasher

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
)

const (
	CurrentHashFuncVersion = 1
)

const (
	N            = 1 << 0xf
	R            = 16
	P            = 4
	HashLength32 = 32
)

type PasswordHasher struct {
	n          int
	r          int
	p          int
	hashLength int
	password   []byte
}

type HashedPassword struct {
	HashFuncVersion uint32
	Salt            []byte
	DK              []byte
}

func NewPasswordHasher(password string) *PasswordHasher {
	return &PasswordHasher{
		n:          N,
		r:          R,
		p:          P,
		hashLength: HashLength32,
		password:   []byte(password),
	}
}

func (h *PasswordHasher) Hash() (*HashedPassword, error) {
	salt := make([]byte, HashLength32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	dk, err := scrypt.Key(
		h.password,
		salt,
		h.n, h.r, h.p,
		h.hashLength,
	)
	if err != nil {
		return nil, err
	}

	return &HashedPassword{
		HashFuncVersion: CurrentHashFuncVersion,
		Salt:            salt,
		DK:              dk,
	}, nil
}

func (h *PasswordHasher) Verify(hash []byte, salt []byte) (bool, error) {
	dk, err := scrypt.Key(
		h.password,
		salt,
		h.n, h.r, h.p,
		h.hashLength,
	)
	if err != nil {
		return false, err
	}

	return bytes.Equal(hash, dk), nil
}
