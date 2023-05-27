package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"users/internal/user"
)

type Database struct {
	Client *sqlx.DB
}

func (d *Database) GetUser(ctx context.Context, s string) (user.User, error) {
	var userRow UserRow
	query := "SELECT id, email,password FROM auth WHERE id = $1"
	err := d.Client.GetContext(ctx, &userRow, query, s)
	user := convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d *Database) PostUser(ctx context.Context, user user.User) (user.User, error) {
	var userRow UserRow
	query := "INSERT INTO auth (id, email, password) VALUES ($1, $2, $3) RETURNING id, email, password"
	err := d.Client.GetContext(ctx, &userRow, query, user.ID, user.Email, user.Password)
	user = convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil

}

func (d *Database) UpdateUser(ctx context.Context, user user.User) (user.User, error) {
	var userRow UserRow
	query := "UPDATE auth SET email = $1, password = $2 WHERE id = $3 RETURNING id, email, password"
	err := d.Client.GetContext(ctx, &userRow, query, user.Email, user.Password, user.ID)
	user = convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d *Database) DeleteUser(ctx context.Context, s string) error {
	query := "DELETE FROM auth WHERE id = $1"
	_, err := d.Client.ExecContext(ctx, query, s)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Ping(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// NewDatabase - returns a pointer to a database object
func NewDatabase() (*Database, error) {
	log.Info("Setting up new database connection")

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_TABLE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}
	log.Info("connected to database")
	return &Database{
		Client: db,
	}, nil
}
