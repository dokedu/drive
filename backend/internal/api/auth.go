package api

import (
	"encoding/json"
	"example/internal/database/db"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"net/http"
	"time"
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

func (s *Config) HandleOneTimeLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.FormValue("email")

	user, err := s.DB.UserFindByEmail(ctx, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send email with token
	err = s.Mailer.SendToken(user.Email, user.FirstName, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body := []byte(`{"message": "Token sent"}`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Config) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.FormValue("token")

	user, err := s.DB.UserFindByToken(ctx, pgtype.Text{String: token, Valid: true})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusInternalServerError)
		return
	}

	// Check if token isn't older than 5 minutes
	if time.Since(user.RecoverySentAt.Time) > 5*time.Minute {
		http.Error(w, "Token expired", http.StatusBadRequest)

		// Delete token
		_, _ = s.DB.ResetUserConfirmationToken(ctx, user.ID)
		return
	}

	// Invalidate token
	_, _ = s.DB.ResetUserConfirmationToken(ctx, user.ID)

	// Create session
	session, err := s.DB.CreateSession(ctx, db.CreateSessionParams{
		UserID: user.ID,
		Token:  gonanoid.Must(32),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	body, err := json.Marshal(userPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Config) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.FormValue("email")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	organisation := r.FormValue("organisation")

	// generate token

	// Create org if it doesn't exist
	org, err := s.DB.OrganisationFindByName(ctx, organisation)
	if err != nil {
		org, err = s.DB.CreateOrganisation(ctx, organisation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Org exists, error
		http.Error(w, "Organisation already exists", http.StatusBadRequest)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(user)

	// Send email with token

}
