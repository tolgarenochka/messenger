package models

import "time"

type Message struct {
	Id        int64     `json:"id" db:"id"`
	Time      time.Time `json:"time" db:"time"`
	Text      string    `json:"text" db:"text"`
	Sender    string    `json:"sender" db:"sender"`
	Recipient string    `json:"recipient" db:"recipient"`
	File      []File    `json:"files"`
}
