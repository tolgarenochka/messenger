package db_wizard

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"messenger/services/api/pkg/helpers/models"
)

type Store struct {
	config interface{}
	conn   *sqlx.DB
}

func NewConnect() (*Store, error) {
	conn, err := sqlx.ConnectContext(context.Background(), "pgx", "postgresql://localhost:5432/postgres")
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return &Store{conn: conn}, nil
}

func (s *Store) Quit() error {
	return s.conn.Close()
}

func Auth(mail string, pas string) (models.User, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return models.User{}, err
	}

	defer func() { log.Print(db.Quit()) }()

	user := models.User{}

	query := db.conn.Rebind(`SELECT * from "user" WHERE mail = ? and pas = ?;`)
	err = db.conn.QueryRowx(query, mail, pas).StructScan(&user)
	if err != nil {
		fmt.Printf(err.Error())
		return user, err
	}

	return user, nil
}

// UpdatePhoto TODO: photo format? is base64 ok for front?
func UpdatePhoto(photo string, id int) (int64, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return 0, err
	}

	defer func() { log.Print(db.Quit()) }()

	query := db.conn.Rebind(`UPDATE "user" SET photo = ? WHERE id = ?;`)
	res, err := db.conn.Exec(query, photo, id)
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Print("Failed connect to count result rows. Reason: ", err.Error())
		return 0, err
	}

	return count, nil
}

func GetUsersList() ([]models.User, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { log.Print(db.Quit()) }()

	users := make([]models.User, 0)

	query := db.conn.Rebind(`SELECT * from "user";`)
	rows, err := db.conn.Queryx(query)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	for rows.Next() {
		user := models.User{}
		err = rows.StructScan(&user)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetDialogsList(id int) ([]models.Dialog, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { log.Print(db.Quit()) }()

	dials := make([]models.DialogDB, 0)

	query := db.conn.Rebind(`SELECT dialog.id, user_1, user_2, text as last_mes_text, is_read, last_mes_sender FROM dialog
    JOIN "message" m on m.id = dialog.last_mes
WHERE user_1 = ? or user_2 = ?;`)

	rows, err := db.conn.Queryx(query, id, id)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	for rows.Next() {
		dial := models.DialogDB{}
		err = rows.StructScan(&dial)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		dials = append(dials, dial)
	}

	dialogs := make([]models.Dialog, 0)
	friendId := 0
	userFriend := models.User{}

	for _, d := range dials {
		dialog := models.Dialog{}
		dialog.Id = d.Id
		dialog.LastMes = d.LastMes

		if d.LastMesSender == id {
			dialog.AreYouLastMesSender = true
			dialog.IsRead = true
		} else {
			dialog.AreYouLastMesSender = false
			dialog.IsRead = d.IsRead
		}

		if d.UserOne == id {
			friendId = d.UserTwo
		} else {
			friendId = d.UserOne
		}

		dialog.FriendId = friendId

		userFriend, err = GetUserInfoById(friendId)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		dialog.FriendFullName = userFriend.SecondName + " " + userFriend.FirstName + " " + userFriend.ThirdName
		dialog.FriendPhoto = userFriend.Photo

		dialogs = append(dialogs, dialog)
	}

	return dialogs, nil
}

func GetUserInfoById(userId int) (models.User, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return models.User{}, err
	}

	defer func() { log.Print(db.Quit()) }()

	user := models.User{}

	query := db.conn.Rebind(`SELECT * from "user" where id = ?;`)
	err = db.conn.QueryRowx(query, userId).StructScan(&user)
	if err != nil {
		fmt.Printf(err.Error())
		return user, err
	}

	return user, nil

}

func GetMessagesList(dialogId int, UserId int) ([]models.Message, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { log.Print(db.Quit()) }()

	messes := make([]models.MessageDB, 0)

	query := db.conn.Rebind(`SELECT message.id, time, text, sender, recipient from message join dialog d on
message.dialog_id = d.id WHERE is_deleted = FALSE and d.id = ?
                             ORDER BY time;`)
	rows, err := db.conn.Queryx(query, dialogId)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	for rows.Next() {
		mes := models.MessageDB{}
		err = rows.StructScan(&mes)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		files, err := GetFilesList(mes.Id)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		mes.File = files

		messes = append(messes, mes)
	}

	mess := make([]models.Message, 0)
	for _, m := range messes {
		mes := models.Message{}
		mes.Id = m.Id
		mes.File = m.File
		mes.Time = m.Time
		mes.Text = m.Text

		if m.Sender == UserId {
			mes.AmISender = true
		} else {
			mes.AmISender = false
		}

		mess = append(mess, mes)
	}

	return mess, nil

}

func GetFilesList(mesId int64) ([]models.File, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { log.Print(db.Quit()) }()

	files := make([]models.File, 0)

	query := db.conn.Rebind(`SELECT path, name from file
JOIN message m on m.id = file.mes_id WHERE m.id = ?;`)
	rows, err := db.conn.Queryx(query, mesId)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	for rows.Next() {
		file := models.File{}
		err = rows.StructScan(&file)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

type Participants struct {
	UserOne int `db:"user_1"`
	UserTwo int `db:"user_2"`
}

func GetDialogParticipants(dialogId int) (int, int, error) {
	dialog := Participants{}

	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return 0, 0, err
	}

	defer func() { log.Print(db.Quit()) }()

	query := db.conn.Rebind(`select user_1, user_2 from dialog where id=?;`)
	err = db.conn.QueryRowx(query, dialogId).StructScan(&dialog)
	if err != nil {
		fmt.Printf(err.Error())
		return 0, 0, err
	}

	return dialog.UserOne, dialog.UserTwo, nil
}

func PostMessage(message models.MessageDB, dialogId int) (int, error) {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return 0, err
	}

	defer func() { log.Print(db.Quit()) }()

	var mesId int

	query := db.conn.Rebind(`INSERT INTO message (id, text, sender, recipient, is_deleted, is_read, dialog_id, time) 
VALUES (DEFAULT, ?, ?, ?, DEFAULT, DEFAULT, ?, ?) RETURNING id;`)
	err = db.conn.QueryRow(query, message.Text, message.Sender, message.Recipient, dialogId, message.Time).Scan(&mesId)
	if err != nil {
		fmt.Printf(err.Error())
		return 0, err
	}

	return mesId, nil
}

func UpdateLastMesInDialog(dialogId int, mesId int, senderId int) error {
	db, err := NewConnect()
	if err != nil {
		log.Print("Failed connect to db. Reason: ", err.Error())
		return err
	}

	defer func() { log.Print(db.Quit()) }()

	query := db.conn.Rebind(`UPDATE dialog SET last_mes = ?, last_mes_sender = ? WHERE id = ?;`)
	_, err = db.conn.Queryx(query, mesId, senderId, dialogId)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	return nil
}
