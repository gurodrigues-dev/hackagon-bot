package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gin/config"
	"gin/types"
	"log"
	"time"

	"github.com/google/uuid"
)

type Postgres struct {
	conn *sql.DB
}

func NewPostgres() (*Postgres, error) {

	conf := config.Get()

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name),
	)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &Postgres{
		conn: db,
	}

	return repo, nil
}

func (p *Postgres) CreateQuestion(ctx context.Context, question *types.Question) error {
	sql := `INSERT INTO questions (id, title, description, date, level) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.conn.Exec(sql, question.ID, question.Title, question.Description, question.Date, question.Level)
	return err
}

func (p *Postgres) ReadQuestion(ctx context.Context) (*types.Question, error) {
	sql := `SELECT id, title, description, date, level FROM questions WHERE date = $1 LIMIT 1`

	var question types.Question

	err := p.conn.QueryRow(sql, time.Now().Format("2006-01-02")).Scan(
		&question.ID,
		&question.Title,
		&question.Description,
		&question.Date,
		&question.Level,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &question, nil
}

func (p *Postgres) UpdateQuestion(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (p *Postgres) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.Exec("DELETE FROM questions WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateUser(ctx context.Context, user *types.User) error {
	return nil
}

func (p *Postgres) ReadUser(ctx context.Context, id uuid.UUID) (*types.User, error) {

	return nil, nil

}

func (p *Postgres) UpdateUser(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (p *Postgres) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return nil
}
