package models

import "encoding/json"

type User struct {
	Id         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`

	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`

	Online bool `json:"online" db:"online"`
}

func (u *User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &u)
}
