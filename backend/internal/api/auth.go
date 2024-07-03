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

type OneTimeLoginResponse struct {
	Message string `json:"message"`
}

func (s *Config) OneTimeLogin(ctx context.Context, r *http.Request) ([]byte, error) {
	email := r.FormValue("email")

	user, err := s.DB.UserFindByEmail(ctx, email)
	if err != nil {
		return nil, ErrNotFound
	}

	// New random token
	token := gonanoid.Must(32)

	updateUserConfirmationTokenParams := db.UpdateUserConfirmationTokenParams{
		RecoveryToken:  pgtype.Text{String: token, Valid: true},
		RecoverySentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:             user.ID,
	}

	// Save the token in db
	_, err = s.DB.UpdateUserConfirmationToken(ctx, updateUserConfirmationTokenParams)
	if err != nil {
		return nil, ErrInternal
	}

	// Send email with token
	err = s.Mailer.SendToken(user.Email, user.FirstName, token)
	if err != nil {
		return nil, ErrInternal
	}

	return json.Marshal(OneTimeLoginResponse{
		Message: "Token sent",
	})
}

type SignInResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (s *Config) SignIn(ctx context.Context, r *http.Request) ([]byte, error) {
	token := r.FormValue("token")

	user, err := s.DB.UserFindByToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		return nil, ErrNotFound
	}

	// If the token is older than 5 minutes, delete it
	if time.Since(user.RecoverySentAt.Time) > 5*time.Minute {
		_, _ = s.DB.ResetUserConfirmationToken(ctx, user.ID)
		return nil, ErrUnauthorized
	}

	// Invalidate token
	_, err = s.DB.ResetUserConfirmationToken(ctx, user.ID)
	if err != nil {
		return nil, ErrInternal
	}

	sessionParams := db.CreateSessionParams{
		UserID: user.ID,
		Token:  gonanoid.Must(32),
	}

	// Create session
	session, err := s.DB.CreateSession(ctx, sessionParams)
	if err != nil {
		return nil, ErrInternal
	}

	response := SignInResponse{
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

	return json.Marshal(response)
}

type SignUpResponse struct {
	Message string `json:"message"`
}

func (s *Config) SignUp(ctx context.Context, r *http.Request) ([]byte, error) {
	email := r.FormValue("email")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	organisation := r.FormValue("organisation")

	switch {
	case email == "":
		return nil, ErrBadRequest
	case firstName == "":
		return nil, ErrBadRequest
	case lastName == "":
		return nil, ErrBadRequest
	case organisation == "":
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

	createUserParams := db.CreateUserParams{
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		OrganisationID: org.ID,
		Role:           db.UserRoleOwner,
	}

	// Create user
	user, err := s.DB.CreateUser(ctx, createUserParams)
	if err != nil {
		return nil, ErrInternal
	}

	// Generate token
	token := gonanoid.Must(32)

	confirmationTokenParams := db.UpdateUserConfirmationTokenParams{
		RecoveryToken:  pgtype.Text{String: token, Valid: true},
		RecoverySentAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:             user.ID,
	}

	// Save the token in db
	_, err = s.DB.UpdateUserConfirmationToken(ctx, confirmationTokenParams)
	if err != nil {
		return nil, ErrInternal
	}

	// Send email with token
	err = s.Mailer.SendToken(user.Email, user.FirstName, token)
	if err != nil {
		return nil, ErrInternal
	}

	return json.Marshal(SignUpResponse{
		Message: "Token sent",
	})
}

type LogOutResponse struct {
	Message string `json:"message"`
}

func (s *Config) LogOut(ctx context.Context, r *http.Request) ([]byte, error) {
	user, ok := middleware.GetUser(ctx, s.DB)
	if !ok {
		return nil, ErrUnauthorized
	}

	// Get the authorization token
	token := r.Header.Get("Authorization")

	removeSessionParams := db.RemoveSessionParams{
		Token:  token,
		UserID: user.ID,
	}

	_, err := s.DB.RemoveSession(ctx, removeSessionParams)
	if err != nil {
		return nil, ErrInternal
	}

	return json.Marshal(LogOutResponse{
		Message: "Logged out",
	})
}
