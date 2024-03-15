package user

import (
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	Plaintext    *string `json:"password"`
	Confirmation *string `json:"password_confirmation"`
	Hash         *[]byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = &hash
	return nil
}

func (p *password) Match(hash *[]byte) (bool, error) {
	if p.Plaintext == nil || hash == nil {
		return false, nil
	}
	if err := bcrypt.CompareHashAndPassword(
		*hash,
		[]byte(*p.Plaintext),
	); err != nil {
		return false, err
	}
	return true, nil
}
