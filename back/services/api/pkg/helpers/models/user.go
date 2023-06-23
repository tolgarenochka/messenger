package models

type User struct {
	Id         int    `json:"id" db:"id"`
	FirstName  string `json:"-" db:"first_name"`
	SecondName string `json:"-" db:"second_name"`
	ThirdName  string `json:"-" db:"third_name"`
	Mail       string `json:"mail,omitempty" db:"mail"`
	Password   string `json:"pas,omitempty" db:"pas"`
	Photo      string `json:"photo" db:"photo"`
	FullName   string `json:"full_name"`
}
