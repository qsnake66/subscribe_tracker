package domain

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

type Subscription struct {
	ID          string
	UserID      string
	ServiceName string
	BankName    string
	CardLast4   string
	Billing     string
	ChargeDate  time.Time
}
