package models

import (
	"errors"
	"regexp"
)

var checkPhone = regexp.MustCompile(`^[\+]?[0-9]{3}\s[0-9]{2,3}\s[0-9]{3}-[0-9]{2}-[0-9]{2}$`)

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

type UpdateWorkerInput struct {
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	FathersName *string `json:"fathers_name"`
	Phone       *string `json:"phone"`
}

type WorkerAttributes struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	Name string `json:"name"`
}

func (w *CreateWorkerInput) Validate() error {
	if w.Name == "" || w.Surname == "" || w.FathersName == "" {
		return errors.New("name fields can't be empty")
	}

	if w.Phone == "" || !checkPhone.MatchString(w.Phone) {
		return errors.New("empty or invalid phone number")
	}

	if w.Role != "admin" && w.Role != "worker" {
		return errors.New("invalid role")
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

	if w.Phone == "" || !checkPhone.MatchString(w.Phone) {
		return errors.New("empty or invalid phone number")
	}

	if w.Password == "" {
		return errors.New("password field can't be empty")
	}

	return nil
}

func (w *UpdateWorkerInput) ToWorker(worker *Worker) {
	if w.Name != nil && worker.Name != *w.Name {
		worker.Name = *w.Name
	}

	if w.Surname != nil && worker.Surname != *w.Surname {
		worker.Surname = *w.Surname
	}

	if w.FathersName != nil && worker.FathersName != *w.FathersName {
		worker.FathersName = *w.FathersName
	}

	if w.Phone != nil && worker.Phone != *w.Phone && checkPhone.MatchString(*w.Phone) {
		worker.Phone = *w.Phone
	}
}
