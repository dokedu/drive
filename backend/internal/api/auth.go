package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Config) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := s.DB.UserFindByEmail(ctx, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
