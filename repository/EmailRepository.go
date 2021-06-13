package repository

import (
	"context"

	"gitlab.jit.team/smurfs/email-service/entity"

	"github.com/jackc/pgx/v4"
)

type EmailRepository struct {
	db *pgx.Conn
}

func NewEmailRepository(conn *pgx.Conn) *EmailRepository {
	emailRepository := new(EmailRepository)
	emailRepository.db = conn
	return emailRepository
}

func (repo *EmailRepository) FindAll() []entity.Email {
	rows, err := repo.db.Query(context.Background(), "SELECT id, title, receiver, sending_user, sending_system FROM email;")
	emails := make([]entity.Email, 0)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var tmp entity.Email
		rows.Scan(&tmp.Id, &tmp.Title, &tmp.Receiver, &tmp.SendingUser, &tmp.SendingSystem)
		emails = append(emails, tmp)
	}

	return emails
}

func (repo *EmailRepository) Save(email entity.Email) {
	sqlStatement := `INSERT INTO email (title, content, receiver, sending_user, sending_system)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	id := 0
	err := repo.db.QueryRow(
		context.Background(), sqlStatement,
		email.Title, email.Content, email.Receiver, email.SendingUser, email.SendingSystem).Scan(&id)

	if err != nil {
		panic(err)
	}
}
