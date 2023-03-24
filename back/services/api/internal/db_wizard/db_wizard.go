package db_wizard

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
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

func (s *Store) GetDialogListByUserId(id int) ([]*models.Dialog, error) {
	dialogs := make([]*models.Dialog, 0)

	query := s.conn.Rebind(`SELECT text as last_mes_text FROM dialog
    JOIN "message" m on m.id = dialog.last_mes
WHERE user_1 = ? or user_2 = ?;`)

	rows, err := s.conn.Query(query, id, id)
	if err != nil {
		fmt.Printf(err.Error())
	}

	for rows.Next() {
		var lastMessage string
		err = rows.Scan(&lastMessage)
		if err != nil {
			fmt.Printf(err.Error())
			return nil, err
		}

		singleDialog := models.NewDialog()
		singleDialog.UpdateLastMes(lastMessage)
		dialogs = append(dialogs, singleDialog)
	}

	return dialogs, nil
}

func (s *Store) Auth(mail string, pas string) (models.User, error) {
	//user := make([]*models.User, 0)
	user := models.User{}

	query := s.conn.Rebind(`SELECT * from "user" WHERE mail = ? and pas = ?;`)

	err := s.conn.QueryRowx(query, mail, pas).StructScan(&user)
	if err != nil {
		fmt.Printf(err.Error())
		return user, err
	}

	return user, nil
}
