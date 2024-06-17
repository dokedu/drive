package api

import (
	"crypto/rand"
	"encoding/json"
	"example/internal/database/db"
	"fmt"
	"math/big"
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

func (s *Config) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := r.FormValue("email")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	organisation := r.FormValue("organisation")

	// Generate otp token
	var minOtp = 100000
	var maxOtp = 999999
	seed, err := rand.Int(rand.Reader, big.NewInt(int64(maxOtp)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	otp := fmt.Sprintf("%06d", seed.Int64()%(int64(maxOtp)-int64(minOtp)+1)+int64(minOtp))
	fmt.Println(otp)

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
