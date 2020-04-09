package main

import (
	"crypto/rand"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		log.WithError(err).Fatal("GenerateRandomString error")
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes)
}

func (s *Service) dataToken(userID int) string {
	salt := s.getKey(userID, "salt")
	if salt == "" {
		salt = GenerateRandomString(16)
		s.setKey(userID, "salt", salt)
	}

	p := fmt.Sprintf("%d_%s", userID, salt)

	return p
}

// check token
func (s *Service) checkToken(token string, userID int) bool {
	return token == s.generateToken(userID)
}

// generate token
func (s *Service) generateToken(userID int) string {
	p := s.dataToken(userID)

	return s.verify.GenerateToken(p)
}

func (s *Service) regenerateToken(userID int) string {
	s.setKey(userID, "salt", GenerateRandomString(16))
	p := s.dataToken(userID)

	return s.verify.GenerateToken(p)
}
