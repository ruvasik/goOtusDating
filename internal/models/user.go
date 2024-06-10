package models

type User struct {
	ID           int
	FirstName    string
	LastName     string
	BirthDate    string
	Gender       string
	Interests    string
	City         string
	Username     string
	PasswordHash string
}