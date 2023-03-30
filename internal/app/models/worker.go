package models

import (
	"errors"
	"regexp"
)

type Worker struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Surname      string `json:"surname" db:"surname"`
	FathersName  string `json:"fathers_name" db:"fathers_name"`
	Phone        string `json:"phone" db:"phone"`
	Role         string `json:"role" db:"role"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}
type CreateWorkerInput struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	FathersName string `json:"fathers_name"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	Password    string `json:"password"`
}

type LogWorkerInput struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type WorkerAttributes struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

func (w *CreateWorkerInput) Validate() error {
	if w.Name == "" || w.Surname == "" || w.FathersName == "" {
		return errors.New("name fields can't be empty")
	}

	re, _ := regexp.Compile(`^[\+]?[0-9]{3}[\s][0-9]{2,3}[\s][0-9]{3}[-][0-9]{2}[-][0-9]{2}$`)
	if w.Phone == "" || !re.MatchString(w.Phone) {
		return errors.New("empty or invalid phone number")
	}

	if w.Role == "" {
		return errors.New("role field can't be empty")
	}

	if w.Password == "" {
		return errors.New("password field can't be empty")
	}

	return nil
}

func (w *LogWorkerInput) Validate() error {
	if w.Name == "" || w.Surname == "" {
		return errors.New("name fields can't be empty")
	}

	re, _ := regexp.Compile(`^[\+]?[0-9]{3}[\s][0-9]{2,3}[\s][0-9]{3}[-][0-9]{2}[-][0-9]{2}$`)
	if w.Phone == "" || !re.MatchString(w.Phone) {
		return errors.New("empty or invalid phone number")
	}

	if w.Password == "" {
		return errors.New("password field can't be empty")
	}

	return nil
}
