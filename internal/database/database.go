package database

import (
	"auth/internal/user"
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	ErrorUserExists = errors.New("user already exists")
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

func (d *Database) PostUser(ctx context.Context, u user.User) (user.User, error) {
	userFound, err := d.GetUserByEmail(ctx, u.Email)

	if userFound != (user.User{}) {
		return u, ErrorUserExists
	}
	var userRow UserRow
	salt, err := generateRandomSalt(saltLength)
	if err != nil {
		return u, err
	}
	saltedPassword := u.Password + salt
	encryptedPassword, err := encryptPassword(saltedPassword)
	if err != nil {
		return u, err
	}

	query := "INSERT INTO users (email, password, salt) VALUES ($1, $2, $3) RETURNING id, email, password"
	err = d.Client.GetContext(ctx, &userRow, query, u.Email, encryptedPassword, salt)

	u = convertUserRowToUser(userRow)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (d *Database) UpdateUser(ctx context.Context, usr user.User) (user.User, error) {
	var userRow UserRow
	query := "UPDATE users SET email = $1, password = $2 WHERE id = $3 RETURNING id, email, password"
	err := d.Client.GetContext(ctx, &userRow, query, usr.Email, usr.Password, usr.ID)
	usr = convertUserRowToUser(userRow)
	if err != nil {
		return user.User{}, errors.New("could not update user")
	}
	return usr, nil
}

func (d *Database) DeleteUser(ctx context.Context, s string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := d.Client.ExecContext(ctx, query, s)
	if err != nil {
		return errors.New("could not delete user")
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
		getOrDefault("DB_HOST", "localhost"),
		getOrDefault("DB_PORT", "5432"),
		getOrDefault("DB_USERNAME", "postgres"),
		getOrDefault("DB_TABLE", "postgres"),
		getOrDefault("DB_PASSWORD", "postgres"),
		getOrDefault("SSL_MODE", "disable"),
	)

	fmt.Println(connectionString)
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
func getOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
