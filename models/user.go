package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Email          string
	HashedPassword string     `db:"hashed_password"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

func (u *User) String() string {
	return fmt.Sprintf("User<%d %s %s>", u.ID, u.Email, u.CreatedAt)
}

func (u *User) ComparePass(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (nu NewUser) String() string {
	return fmt.Sprintf("NewUser<%s %s>", nu.Email, nu.Password)
}

func HashPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+4)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(bytes), nil
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(user NewUser) (*User, error) {
	hp, err := HashPass(user.Password)
	if err != nil {
		return nil, err
	}
	row := us.DB.QueryRow(`INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING id, created_at, updated_at`, strings.ToLower(user.Email), hp)
	u := &User{Email: user.Email, HashedPassword: hp}
	err = row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *UserService) Update(user *User) error {
	row := us.DB.QueryRow(`UPDATE users SET email = $1, hashed_password = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`, strings.ToLower(user.Email), user.HashedPassword, user.ID)
	err := row.Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (us UserService) Authenticate(email, password string) (*User, error) {
	u := User{
		Email: email,
	}
	row := us.DB.QueryRow(`SELECT id, hashed_password FROM users WHERE email = $1`, strings.ToLower(email))
	err := row.Scan(&u.ID, &u.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !u.ComparePass(password) {
		return nil, fmt.Errorf("invalid password")
	}
	return &u, nil
}
