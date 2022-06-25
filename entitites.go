package main

type User struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	Salt         int    `json:"salt"`
	PasswordHash string `json:"password_hash"`
	Description  string `json:"description"`
}
