package models

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    BirthDate string `json:"birth_date"`
    City      string `json:"city"`
}