package models

type Dialog struct {
	lastMes string
	//user    user.User
	//mail       string
}

func NewDialog() *Dialog {
	return &Dialog{
		lastMes: "",
	}
}

func (d *Dialog) UpdateLastMes(mes string) {
	d.lastMes = mes
}
