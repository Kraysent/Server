package entities

import "time"

const (
	TableTokens = "tokens"

	TokenFieldID             = "id"
	TokenFieldUserID         = "user_id"
	TokenFieldValue          = "value"
	TokenFieldStartDate      = "start_date"
	TokenFieldExpirationDate = "expiration_date"
)

type Token struct {
	ID             int       `json:"id"`
	Value          string    `json:"value"`
	UserID         int       `json:"user_id"`
	StartDate      time.Time `json:"start_date"`
	ExpirationDate time.Time `json:"expiration_date"`
}
