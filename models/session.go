package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/jakeecolution/lenslocked/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID     int `db:"id"`
	UserID int `db:"user_id"`
	// Token is only set when creating a new session. When looking up a session
	// this will be left empty, as we only store the hash of a session token
	// in our database and we cannot reverse it into a raw token.
	Token     string
	TokenHash string `db:"token_hash"`
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	var token string
	var err error
	if ss.BytesPerToken <= MinBytesPerToken {
		token, err = rand.String(MinBytesPerToken)
	} else {
		token, err = rand.String(ss.BytesPerToken)
	}
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}
	row := ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1, $2) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) hash(token string) string {
	thash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(thash[:])
}

func (ss *SessionService) User(token string) (*User, error) {
	var user User
	row := ss.DB.QueryRow(`
		SELECT users.id, users.email, users.hashed_password, users.created_at, users.updated_at
		FROM users
		INNER JOIN sessions
		ON users.id = sessions.user_id
		WHERE sessions.token_hash = $1`, token)
	err := row.Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("User: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {

	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1`, token)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}
