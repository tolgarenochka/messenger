package db_wizard

import (
	"context"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"messenger/services/api/pkg/helpers/models"

	. "messenger/services/api/pkg/helpers/logger"
)

type Store struct {
	config interface{}
	conn   *sqlx.DB
}

// функция открытия подключения к БД
func NewConnect() (*Store, error) {
	conn, err := sqlx.ConnectContext(context.Background(), "pgx", "postgresql://localhost:5432/postgres")
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	return &Store{conn: conn}, nil
}

// функция акрытия подключения к бд
func (s *Store) Quit() error {
	return s.conn.Close()
}

// запрос в БД на авторизацию
func Auth(mail string, pas string) (models.User, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return models.User{}, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	user := models.User{}

	query := db.conn.Rebind(`SELECT * from "user" WHERE mail = ? and pas = ?;`)
	err = db.conn.QueryRowx(query, mail, pas).StructScan(&user)
	if err != nil {
		Logger.Error(err.Error())
		return user, err
	}

	return user, nil
}

// запрос в БД на получения списка пользователей, с кем еще нет диалога
func GetUsersList(id int) ([]models.User, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	users := make([]models.User, 0)

	query := db.conn.Rebind(`SELECT id, first_name, second_name, third_name, photo from "user" as u
         where not
             (u.id in (select user_1 from dialog where user_2 = ?)
            OR u.id in (select user_2 from dialog where user_1 = ?))
            AND not u.id = ?
			ORDER BY first_name, second_name, third_name;`)
	rows, err := db.conn.Queryx(query, id, id, id)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		user := models.User{}
		err = rows.StructScan(&user)
		if err != nil {
			Logger.Error(err.Error())
			return nil, err
		}
		user.FullName = user.SecondName + " " + user.FirstName + " " + user.ThirdName
		users = append(users, user)
	}

	return users, nil
}

// запрос в БД на получение списка диалогов пользователя
func GetDialogsList(id int) ([]models.Dialog, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	dials := make([]models.DialogDB, 0)

	query := db.conn.Rebind(`SELECT dialog.id, user_1, user_2, text as last_mes_text, is_read, last_mes_sender, time FROM dialog
    JOIN "message" m on m.id = dialog.last_mes
WHERE user_1 = ? or user_2 = ? ORDER BY time DESC;`)

	rows, err := db.conn.Queryx(query, id, id)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		dial := models.DialogDB{}
		err = rows.StructScan(&dial)
		if err != nil {
			Logger.Error(err.Error())
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
		dialog.Time = d.Time

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
			Logger.Error(err.Error())
			return nil, err
		}
		dialog.FriendFullName = userFriend.SecondName + " " + userFriend.FirstName + " " + userFriend.ThirdName
		dialog.FriendPhoto = userFriend.Photo

		dialogs = append(dialogs, dialog)
	}

	return dialogs, nil
}

// запрос в БД на получение информации о пользователе по его id
func GetUserInfoById(userId int) (models.User, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return models.User{}, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	user := models.User{}

	query := db.conn.Rebind(`SELECT * from "user" where id = ?;`)
	err = db.conn.QueryRowx(query, userId).StructScan(&user)
	if err != nil {
		Logger.Error(err.Error())
		return user, err
	}
	user.FullName = user.SecondName + " " + user.FirstName + " " + user.ThirdName

	return user, nil

}

// запрос в БД на получение всех сообщений в диалоге
func GetMessagesList(dialogId int, UserId int) ([]models.Message, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	messes := make([]models.MessageDB, 0)

	query := db.conn.Rebind(`SELECT message.id, time, text, sender, recipient from message join dialog d on
message.dialog_id = d.id WHERE is_deleted = FALSE and d.id = ?
                             ORDER BY time;`)
	rows, err := db.conn.Queryx(query, dialogId)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		mes := models.MessageDB{}
		err = rows.StructScan(&mes)
		if err != nil {
			Logger.Error(err.Error())
			return nil, err
		}
		files, err := GetFilesList(mes.Id)
		if err != nil {
			Logger.Error(err.Error())
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

// запрос в БД на получение всех файлов, относящихся к выбранному сообщению
func GetFilesList(mesId int64) ([]models.File, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return nil, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	files := make([]models.File, 0)

	query := db.conn.Rebind(`SELECT path, name from file
JOIN message m on m.id = file.mes_id WHERE m.id = ?;`)
	rows, err := db.conn.Queryx(query, mesId)
	if err != nil {
		Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		file := models.File{}
		err = rows.StructScan(&file)
		if err != nil {
			Logger.Error(err.Error())
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

// запрос в БД на получение списка участников диалога по заданному id диалога
func GetDialogParticipants(dialogId int) (int, int, error) {
	dialog := Participants{}

	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return 0, 0, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	query := db.conn.Rebind(`select user_1, user_2 from dialog where id=?;`)
	err = db.conn.QueryRowx(query, dialogId).StructScan(&dialog)
	if err != nil {
		Logger.Error(err.Error())
		return 0, 0, err
	}

	return dialog.UserOne, dialog.UserTwo, nil
}

// запрос в БД на запись сообщения в таблицу сообщений
func PostMessage(message models.MessageDB, dialogId int) (int, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return 0, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	var mesId int

	query := db.conn.Rebind(`INSERT INTO message (id, text, sender, recipient, is_deleted, is_read, dialog_id, time) 
VALUES (DEFAULT, ?, ?, ?, DEFAULT, DEFAULT, ?, ?) RETURNING id;`)
	err = db.conn.QueryRow(query, message.Text, message.Sender, message.Recipient, dialogId, message.Time).Scan(&mesId)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}

	return mesId, nil
}

// запрос в БД на обновление последнего сообщения в диалоге
func UpdateLastMesInDialog(dialogId int, mesId int, senderId int) error {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	query := db.conn.Rebind(`UPDATE dialog SET last_mes = ?, last_mes_sender = ? WHERE id = ?;`)
	_, err = db.conn.Queryx(query, mesId, senderId, dialogId)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	return nil
}

// запрос в БД на создание нового диалога
func CreateDialog(myId int, friendId int) (int, error) {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return 0, err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	query := db.conn.Rebind(`INSERT INTO dialog (user_1, user_2) VALUES (?,?) RETURNING id;`)

	var dId int
	err = db.conn.QueryRow(query, myId, friendId).Scan(&dId)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}

	return dId, nil
}

// запрос в БД на изменение статуса диалога (прочитано)
func ReadDialog(dialogId int) error {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	query := db.conn.Rebind(`UPDATE message SET is_read = true WHERE id = (SELECT m.id FROM dialog
    JOIN "message" m on m.id = dialog.last_mes
WHERE dialog.id=?);`)
	_, err = db.conn.Queryx(query, dialogId)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	return nil
}

// запрос в БД на запись файла, принадлежащему выбранному сообщению
func SaveFile(fileName string, userId int) error {
	db, err := NewConnect()
	if err != nil {
		Logger.Error("Failed connect to db. Reason: ", err.Error())
		return err
	}

	defer func() { Logger.Debug(db.Quit()) }()

	query := db.conn.Rebind(`insert into file (mes_id, path, name) values 
        ((select id from message where sender = ? order by time desc limit 1), ?, ?)`)
	_, err = db.conn.Queryx(query, userId, "files/"+fileName, fileName)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}

	return nil
}
