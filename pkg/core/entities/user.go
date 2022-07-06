package entities

import "time"

const (
	TableUsers = "users"

	UserFieldID               = "id"
	UserFieldLogin            = "login"
	UserFieldSalt             = "salt"
	UserFieldPasswordHash     = "password_hash"
	UserFieldDescription      = "description"
	UserFieldCityID           = "city_id"
	UserFieldRegistrationDate = "registration_date"
)

type User struct {
	ID               int       `json:"id"`
	Login            string    `json:"login"`
	Salt             int       `json:"salt"`
	PasswordHash     string    `json:"password_hash"`
	Description      string    `json:"description"`
	CityID           int       `json:"city"`
	RegistrationDate time.Time `json:"registration_date"`
}
