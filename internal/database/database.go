package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"users/internal/user"
)

type Database struct {
	Client *sqlx.DB
}

func (d *Database) GetUser(ctx context.Context, s string) (user.User, error) {
	var userRow UserRow
	query := "SELECT id, email,password FROM users WHERE id = $1"
	err := d.Client.GetContext(ctx, &userRow, query, s)
	user := convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	var userRow UserRow
	query := "SELECT id, email,password FROM users WHERE email = $1"
	err := d.Client.GetContext(ctx, &userRow, query, email)
	user := convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (d *Database) GetUserAndSaltByEmail(ctx context.Context, email string) (user.User, string, error) {
	var userRow UserRow
	query := "SELECT id, email,password,salt FROM users WHERE email = $1"
	err := d.Client.GetContext(ctx, &userRow, query, email)
	user := convertUserRowToUser(userRow)
	if err != nil {
		return user, "", err
	}
	return user, userRow.Salt, nil
}

func (d *Database) PostUser(ctx context.Context, user user.User) (user.User, error) {
	var userRow UserRow
	salt, err := generateRandomSalt(saltLength)
	if err != nil {
		return user, err
	}
	saltedPassword := user.Password + salt
	encryptedPassword, err := encryptPassword(saltedPassword)
	if err != nil {
		return user, err
	}

	query := "INSERT INTO users (email, password, salt) VALUES ($1, $2, $3) RETURNING id, email, password"
	err = d.Client.GetContext(ctx, &userRow, query, user.Email, encryptedPassword, salt)

	user = convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d *Database) UpdateUser(ctx context.Context, user user.User) (user.User, error) {
	var userRow UserRow
	query := "UPDATE users SET email = $1, password = $2 WHERE id = $3 RETURNING id, email, password"
	err := d.Client.GetContext(ctx, &userRow, query, user.Email, user.Password, user.ID)
	user = convertUserRowToUser(userRow)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (d *Database) DeleteUser(ctx context.Context, s string) error {
	query := "DELETE FROM users WHERE id = $1"
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

const maxRetries = 5
const retryInterval = time.Second * 5

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

	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", connectionString)
		if err == nil {
			log.Info("Connected to database")
			return &Database{
				Client: db,
			}, nil
		}

		log.Errorf("Could not connect to database: %v", err)

		if i < maxRetries-1 {
			log.Infof("Retrying database connection in %s...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	return &Database{}, fmt.Errorf("failed to connect to database after %d retries", maxRetries)
}
