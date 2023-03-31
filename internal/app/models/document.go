package models

import "time"

type Document struct {
	ID         int       `json:"ID" db:"id"`
	Car        string    `json:"car" db:"car"`
	CarID      string    `json:"car_id" db:"car_id"`
	Waybill    string    `json:"waybill" db:"waybill"`
	DriverName string    `json:"driver_name" db:"driver_name"`
	GasAmount  int       `json:"gas_amount" db:"gas_amount"`
	GasType    string    `json:"gas_type" db:"gas_type"`
	IssueDate  time.Time `json:"issue_date" db:"issue_date"`
}

type CreateDocInput struct {
	Car        string    `json:"car" db:"car"`
	CarID      string    `json:"car_id" db:"car_id"`
	Waybill    string    `json:"waybill" db:"waybill"`
	DriverName string    `json:"driver_name" db:"driver_name"`
	GasAmount  int       `json:"gas_amount" db:"gas_amount"`
	GasType    string    `json:"gas_type" db:"gas_type"`
	IssueDate  time.Time `json:"issue_date" db:"issue_date"`
}
