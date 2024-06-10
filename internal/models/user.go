package models

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    BirthDate string `json:"birth_date"`
    Gender    string `json:"gender"`
    Interests string `json:"interests"`
    City      string `json:"city"`
    Username  string `json:"username"`
    Password  string `json:"password"`
}
