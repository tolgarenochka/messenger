package models

import (
	"encoding/json"
)

type User struct {
	Id         int    `json:"id" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	SecondName string `json:"second_name" db:"second_name"`
	ThirdName  string `json:"third_name" db:"third_name"`
	Mail       string `json:"mail,omitempty" db:"mail"`
	Password   string `json:"pas,omitempty" db:"pas"`
	Photo      string `json:"photo,omitempty" db:"photo"`
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &u)
}
