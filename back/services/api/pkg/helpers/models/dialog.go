package models

type Dialog struct {
	Id                  int    `json:"id" db:"id"`
	LastMes             string `json:"last_mes" db:"last_mes"`
	AreYouLastMesSender bool
	FirstName           string
	SecondName          string
	ThirdName           string
}

func (d *Dialog) UpdateLastMes(mes string) {
	d.LastMes = mes
}
