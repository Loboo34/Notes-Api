package middleware

import (
	"net/http"
	"notes/utils"
	"strings"
)

func AuthMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing Token", "")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := utils.ValidateJWT(token)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token", "")
			return
		}

		next (w,r)
	}
}
