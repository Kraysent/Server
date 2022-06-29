package entities

import "time"

type User struct {
	ID               int       `json:"id"`
	Login            string    `json:"login"`
	Salt             int       `json:"salt"`
	PasswordHash     string    `json:"password_hash"`
	Description      string    `json:"description"`
	CityID           int       `json:"city"`
	RegistrationDate time.Time `json:"registration_date"`
}

type Token struct {
	ID             int       `json:"id"`
	Value          string    `json:"value"`
	UserID         int       `json:"user_id"`
	StartDate      time.Time `json:"start_date"`
	ExpirationDate time.Time `json:"expiration_date"`
}
