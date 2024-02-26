package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gin/config"
	"gin/types"
	"log"
	"strings"
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

	params1 := strings.Join(question.Inputs.Test1.Params, ",")
	params2 := strings.Join(question.Inputs.Test2.Params, ",")
	params3 := strings.Join(question.Inputs.Test3.Params, ",")

	sql := `INSERT INTO questions (id, title, description, date, level, params1, response1, params2, response2, params3, response3) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := p.conn.Exec(sql, question.ID, question.Title, question.Description, question.Date, question.Level, params1, question.Inputs.Test1.Response, params2, question.Inputs.Test2.Response, params3, question.Inputs.Test3.Response)
	return err
}

func (p *Postgres) ReadQuestion(ctx context.Context) (*types.Question, error) {
	sql := `SELECT id, title, description, date, level, params1, response1, params2, response2, params3, response3 FROM questions WHERE date = $1 LIMIT 1`

	var (
		question   types.Question
		params1Str string
		params2Str string
		params3Str string
	)

	err := p.conn.QueryRow(sql, time.Now().Format("2006-01-02")).Scan(
		&question.ID,
		&question.Title,
		&question.Description,
		&question.Date,
		&question.Level,
		&params1Str,
		&question.Inputs.Test1.Response,
		&params2Str,
		&question.Inputs.Test2.Response,
		&params3Str,
		&question.Inputs.Test3.Response,
	)
	if err != nil {
		log.Fatal(err)
	}

	question.Inputs.Test1.Params = strings.Split(params1Str, ",")
	question.Inputs.Test2.Params = strings.Split(params2Str, ",")
	question.Inputs.Test3.Params = strings.Split(params3Str, ",")

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
	sql := `INSERT INTO users (nickname, email, password) VALUES ($1, $2, $3)`
	_, err := p.conn.Exec(sql, user.Nickname, user.Email, user.Password)
	return err
}

func (p *Postgres) ReadUser(ctx context.Context, id *int) (*types.User, error) {
	sql := `SELECT id, nickname, email FROM users WHERE id = $1 LIMIT 1`

	var user types.User

	err := p.conn.QueryRow(sql, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.Email,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &user, nil
}

func (p *Postgres) UpdateUser(ctx context.Context, id *int) error {
	return nil
}

func (p *Postgres) DeleteUser(ctx context.Context, nickname *string) error {
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

	_, err = tx.Exec("DELETE FROM users WHERE nickname = $1", nickname)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) VerifyLogin(ctx context.Context, user *types.User) error {
	var password string

	sql := `SELECT password FROM users WHERE nickname = $1`

	err := p.conn.QueryRow(sql, user.Nickname).Scan(&password)
	if err != nil {
		return err
	}

	match := password == user.Password

	if !match {
		return fmt.Errorf("nickname or password wrong")
	}

	return nil

}

func (p *Postgres) CreateAnswer(ctx context.Context, answer *types.Answer) error {
	sql := `INSERT INTO answers (id, nickname, questionid, status, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.conn.Exec(sql, answer.ID, answer.Nickname, answer.QuestionID, answer.Status, answer.CreatedAt)
	return err
}

func (p *Postgres) DeleteAnswer(ctx context.Context, id uuid.UUID) error {
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

	_, err = tx.Exec("DELETE FROM answers WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) VerifyAnswer(ctx context.Context, question *types.Question, nickname *string) (*types.Answer, error) {
	sqlQuery := `SELECT nickname, status, created_at FROM answers WHERE questionid = $1 AND nickname = $2 ORDER BY created_at DESC LIMIT 1`

	var answerResponse types.Answer

	err := p.conn.QueryRow(sqlQuery, question.ID, nickname).Scan(
		&answerResponse.Nickname,
		&answerResponse.Status,
		&answerResponse.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &answerResponse, nil
		}

		log.Fatal(err)
		return nil, err
	}

	return &answerResponse, nil

}
