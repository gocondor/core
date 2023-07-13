// Copyright 2021 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hashing struct{}

func (h *Hashing) HashPassword(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (h *Hashing) CheckPasswordHash(hashedPassword string, originalPassowrd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(originalPassowrd))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		fmt.Print(err)
		loggr.Debug(err.Error())
		return false, errors.New("failed checking password hash")
	}
	return true, nil
}
