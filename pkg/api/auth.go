package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secret       = []byte("szq2tepV0oeG3sFBWSM89lGLeEH2COYa")
	authPassword = os.Getenv("TODO_PASSWORD")
)

type authRequest struct {
	Password string `json:"password"`
}

func getHash(password string) string {
	var hash [32]byte

	hash = sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		req authRequest
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, `Invalid auth request (JSON with "password" field is expected)`, http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		writeError(w, `Received "password" is empty`, http.StatusBadRequest)
		return
	}

	if req.Password != authPassword {
		writeError(w, "Password incorrect", http.StatusUnauthorized)
		return
	}

	var (
		token       *jwt.Token
		signedToken string
		claims      jwt.MapClaims
	)

	claims = jwt.MapClaims{
		"service": "Go Final Project",
		"year":    2026,
		"hash":    getHash(authPassword),
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err = token.SignedString(secret); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, struct {
		Token string `json:"token"`
	}{
		Token: signedToken,
	})
}

func verifyToken(jwtString string) bool {
	var (
		err   error
		token *jwt.Token
	)

	// Parse token
	if token, err = jwt.Parse(jwtString, func(*jwt.Token) (any, error) {
		return secret, nil
	}); err != nil {
		return false
	}

	// Check if valid
	if !token.Valid {
		return false
	}

	// Get claims
	var (
		ok     bool
		claims jwt.MapClaims
	)
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return false
	}

	// Check if "hash" field is in
	var (
		rawHash any
		hash    string
	)
	if rawHash, ok = claims["hash"]; !ok {
		return false
	}

	// And is a string
	if hash, ok = rawHash.(string); !ok {
		return false
	}

	// Verify "hash"
	if hash != getHash(authPassword) {
		return false
	}

	return true
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(authPassword) > 0 {
			var (
				err        error
				cookie     *http.Cookie
				jwtString  string
				authorized bool
			)

			// Getting a token cookie
			if cookie, err = r.Cookie("token"); err == nil {
				jwtString = cookie.Value

				// Verify token and hash
				if verifyToken(jwtString) {
					authorized = true
				}
			}
			if !authorized {
				writeError(w, "Unauthorized request", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
