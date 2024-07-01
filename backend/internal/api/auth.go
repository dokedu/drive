package api

import (
	"context"
	"encoding/json"
	"example/internal/database/db"
	"example/internal/middleware"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type UserPayload struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID             string      `json:"id"`
	Email          string      `json:"email"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	OrganisationID string      `json:"organisation"`
	Role           db.UserRole `json:"role"`
}

func (s *Config) HandleOneTimeLogin(ctx context.Context, r *http.Request) ([]byte, error) {
	email := r.FormValue("email")

	user, err := s.DB.UserFindByEmail(ctx, email)
	if err != nil {
		return nil, ErrNotFound
	}

	// Generate
	token := gonanoid.Must(32)

	// Save the token in db
	_, err = s.DB.UpdateUserConfirmationToken(ctx, db.UpdateUserConfirmationTokenParams{
		RecoveryToken:  pgtype.Text{String: token, Valid: true},
		RecoverySentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:             user.ID,
	})

	if err != nil {
		return nil, ErrInternal
	}

	// Send email with token
	err = s.Mailer.SendToken(user.Email, user.FirstName, token)
	if err != nil {
		return nil, ErrInternal
	}

	body := []byte(`{"message": "Token sent"}`)
	return body, nil
}

func (s *Config) HandleLogin(ctx context.Context, r *http.Request) ([]byte, error) {
	token := r.FormValue("token")

	user, err := s.DB.UserFindByToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		return nil, ErrNotFound
	}

	// Check if token isn't older than 5 minutes
	if time.Since(user.RecoverySentAt.Time) > 5*time.Minute {
		// Delete token
		_, _ = s.DB.ResetUserConfirmationToken(ctx, user.ID)
		return nil, ErrUnauthorized
	}

	// Invalidate token
	_, _ = s.DB.ResetUserConfirmationToken(ctx, user.ID)

	// Create session
	session, err := s.DB.CreateSession(ctx, db.CreateSessionParams{
		UserID: user.ID,
		Token:  gonanoid.Must(32),
	})

	if err != nil {
		return nil, ErrInternal
	}

	userPayload := UserPayload{
		Token: session.Token,
		User: User{
			ID:             user.ID,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			OrganisationID: user.OrganisationID,
			Role:           user.Role,
		},
	}

	return json.Marshal(userPayload)
}

func (s *Config) HandleSignUp(ctx context.Context, r *http.Request) ([]byte, error) {
	email := r.FormValue("email")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	organisation := r.FormValue("organisation")

	if email == "" || firstName == "" || lastName == "" || organisation == "" {
		return nil, ErrBadRequest
	}

	// Check if user doesn't exist
	_, err := s.DB.UserFindByEmail(ctx, email)
	if err == nil {
		return nil, ErrBadRequest
	}

	// Check if org doesn't already exist
	_, err = s.DB.OrganisationFindByName(ctx, organisation)
	if err == nil {
		return nil, ErrBadRequest
	}

	// Create organisation
	org, err := s.DB.CreateOrganisation(ctx, organisation)
	if err != nil {
		return nil, ErrInternal
	}

	// Create user
	user, err := s.DB.CreateUser(ctx, db.CreateUserParams{
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		OrganisationID: org.ID,
		Role:           "owner",
	})
	if err != nil {
		return nil, ErrInternal
	}

	// Generate token
	token := gonanoid.Must(32)

	// Save the token in db
	_, err = s.DB.UpdateUserConfirmationToken(ctx, db.UpdateUserConfirmationTokenParams{
		RecoveryToken:  pgtype.Text{String: token, Valid: true},
		RecoverySentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:             user.ID,
	})

	if err != nil {
		return nil, ErrInternal
	}

	// Send email with token
	err = s.Mailer.SendToken(user.Email, user.FirstName, token)
	if err != nil {
		return nil, ErrInternal
	}

	body := []byte(`{"message": "Token sent"}`)
	return body, nil
}

func (s *Config) HandleLogOut(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	// Get the authorization token
	token := r.Header.Get("Authorization")

	_, err := s.DB.RemoveSession(ctx, db.RemoveSessionParams{
		Token:  token,
		UserID: user.ID,
	})
	if err != nil {
		return nil, ErrInternal
	}

	body := []byte(`{"message": "Logged out"}`)
	return body, nil
}
