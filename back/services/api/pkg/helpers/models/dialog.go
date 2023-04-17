package models

type Dialog struct {
	Id                  int    `json:"id"`
	LastMes             string `json:"last_mes"`
	AreYouLastMesSender bool   `json:"are_you_last_mes_sender"`
	FriendFullName      string `json:"full_name"`
	IsRead              bool   `json:"is_read" db:"is_read"`
	FriendPhoto         string `json:"friend_photo"`
}

type DialogDB struct {
	Id            int    `json:"id" db:"id"`
	UserOne       int    `json:"user_1" db:"user_1"`
	UserTwo       int    `json:"user_2" db:"user_2"`
	LastMesSender int    `json:"last_mes_sender" db:"last_mes_sender"`
	LastMes       string `json:"last_mes_text" db:"last_mes_text"`
	IsRead        bool   `json:"is_read" db:"is_read"`
}

func (d *Dialog) UpdateLastMes(mes string) {
	d.LastMes = mes
}
