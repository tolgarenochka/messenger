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
