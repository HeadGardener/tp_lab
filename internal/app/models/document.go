package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type MyTime time.Time

var checkCarID = regexp.MustCompile(`^[0-9]{4}\s[A-Z]{2}-[1-7]$`)

type Document struct {
	ID         int    `json:"ID" db:"id"`
	Car        string `json:"car" db:"car"`
	CarID      string `json:"car_id" db:"car_id"`
	Waybill    int    `json:"waybill" db:"waybill"`
	DriverName string `json:"driver_name" db:"driver_name"`
	GasAmount  int    `json:"gas_amount" db:"gas_amount"`
	GasType    string `json:"gas_type" db:"gas_type"`
	IssueDate  MyTime `json:"issue_date" db:"issue_date"`
}

type CreateDocInput struct {
	Car        string `json:"car" db:"car"`
	CarID      string `json:"car_id" db:"car_id"`
	Waybill    int    `json:"waybill" db:"waybill"`
	DriverName string `json:"driver_name" db:"driver_name"`
	GasAmount  int    `json:"gas_amount" db:"gas_amount"`
	GasType    string `json:"gas_type" db:"gas_type"`
	IssueDate  MyTime `json:"issue_date" db:"issue_date"`
}

type UpdateDocInput struct {
	Car        *string `json:"car" db:"car"`
	CarID      *string `json:"car_id" db:"car_id"`
	Waybill    *int    `json:"waybill" db:"waybill"`
	DriverName *string `json:"driver_name" db:"driver_name"`
	GasAmount  *int    `json:"gas_amount" db:"gas_amount"`
	GasType    *string `json:"gas_type" db:"gas_type"`
	IssueDate  *MyTime `json:"issue_date" db:"issue_date"`
}

func (d *CreateDocInput) Validate() error {
	if d.Car == "" || d.DriverName == "" || d.GasType == "" {
		return errors.New("there can't be empty fields")
	}

	if !checkCarID.MatchString(d.CarID) {
		return errors.New("invalid car_id")
	}

	if d.Waybill < 1000 || d.Waybill > 9999 {
		return errors.New("invalid waybill value")
	}

	if d.GasAmount <= 0 {
		return errors.New("gas_amount can't be less than zero")
	}

	return nil
}

func (d *UpdateDocInput) ToDocument(doc *Document) {
	if d.Car != nil && doc.Car != *d.Car {
		doc.Car = *d.Car
	}

	if d.CarID != nil && doc.CarID != *d.CarID && checkCarID.MatchString(*d.CarID) {
		doc.CarID = *d.CarID
	}

	if d.Waybill != nil && doc.Waybill != *d.Waybill {
		doc.Waybill = *d.Waybill
	}

	if d.DriverName != nil && doc.DriverName != *d.DriverName {
		doc.DriverName = *d.DriverName
	}

	if d.GasAmount != nil && doc.GasAmount != *d.GasAmount {
		doc.GasAmount = *d.GasAmount
	}

	if d.GasType != nil && doc.GasType != *d.GasType {
		doc.GasType = *d.GasType
	}

	if d.IssueDate != nil && doc.IssueDate != *d.IssueDate {
		doc.IssueDate = *d.IssueDate
	}
}

func (mt *MyTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	*mt = MyTime(t)
	return nil
}

func (mt MyTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(mt).Format("2006-01-02") + `"`), nil
}

// Value implements the driver Valuer interface.
func (mt MyTime) Value() (driver.Value, error) {
	return driver.Value(time.Time(mt)), nil
}

// Scan implements the Scanner interface.
func (mt *MyTime) Scan(value interface{}) error {
	if value == nil {
		var m *MyTime
		mt = m
	} else {
		date := fmt.Sprintf("%s", value)
		parts := strings.Split(date, " ")
		t, err := time.Parse("2006-01-02", parts[0])
		if err != nil {
			return err
		}
		*mt = MyTime(t)
	}
	return nil
}
