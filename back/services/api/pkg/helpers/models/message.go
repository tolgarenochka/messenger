package models

import "time"

type MessageDB struct {
	Id        int64     `json:"id" db:"id"`
	Time      time.Time `json:"time" db:"time"`
	Text      string    `json:"text" db:"text"`
	Sender    int       `json:"sender" db:"sender"`
	Recipient int       `json:"recipient" db:"recipient"`
	File      []File    `json:"files"`
}

type Message struct {
	Id        int64     `json:"id" db:"id"`
	Time      time.Time `json:"time" db:"time"`
	Text      string    `json:"text" db:"text"`
	AmISender bool      `json:"am_i_sender"`
	File      []File    `json:"files"`
}
