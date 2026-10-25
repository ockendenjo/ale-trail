package main

import (
	"errors"
	"net/mail"
)

type SetStasherRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *SetStasherRequest) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}
	if len(s.Name) > 50 {
		return errors.New("name is too long")
	}
	if s.Email == "" {
		return errors.New("email is required")
	}
	if len(s.Email) > 80 {
		return errors.New("email is too long")
	}
	if _, err := mail.ParseAddress(s.Email); err != nil {
		return errors.New("invalid email")
	}
	return nil
}
