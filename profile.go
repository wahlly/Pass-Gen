package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)


var key = "askdjasjdbreonfsdfibsdhfgsdfhboo"

var ErrMalformedEncryption = errors.New("malformed encryption")

//password in small letters so it is not stored
type profile struct {
	Enc, Platform, password string
}

func (p *profile) encrypt() error {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	enc := gcm.Seal(nonce, nonce, []byte(p.password), nil)
	p.Enc = hex.EncodeToString(enc)

	return nil
}

func (p *profile) decrypt() error {
	block, err := aes.NewCipher([]byte(key))
	if err != nil{
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nsize := gcm.NonceSize()

	if len(p.Enc) < nsize {
		return ErrMalformedEncryption
	}

	enc, err := hex.DecodeString(p.Enc)
	if err != nil {
		return err
	}

	password, err := gcm.Open(nil, enc[:nsize], enc[nsize:], nil)
	if err != nil {
		return err
	}

	p.password = string(password)

	return nil
}