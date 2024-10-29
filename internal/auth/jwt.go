package auth

import (
	"fmt"
	"net/http"
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/types"
	"ntsiris/product-microservice/internal/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(requiredScope string) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer" {
				w.Header().Set("WWW-Authenticate", "Bearer")
				apiError := types.APIError{
					Code:          http.StatusForbidden,
					Message:       "Permission required",
					Operation:     types.FormatOperation(r.Method, r.URL.Path),
					EmbeddedError: fmt.Errorf("error: authHeader:%s", authHeader).Error(),
				}

				utils.WriteJSON(w, http.StatusForbidden, apiError)
				return
			}

			tokenString := authHeader[7:]
			token, err := validateJWT(tokenString)
			if err != nil || !token.Valid {
				w.Header().Set("WWW-Authenticate", "Bearer")
				apiError := types.APIError{
					Code:          http.StatusForbidden,
					Message:       "Permission required",
					Operation:     types.FormatOperation(r.Method, r.URL.Path),
					EmbeddedError: err.Error(),
				}
				utils.WriteJSON(w, http.StatusForbidden, apiError)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				w.Header().Set("WWW-Authenticate", "Bearer")
				apiError := types.APIError{
					Code:          http.StatusForbidden,
					Message:       "No claims found",
					Operation:     types.FormatOperation(r.Method, r.URL.Path),
					EmbeddedError: "",
				}
				utils.WriteJSON(w, http.StatusForbidden, apiError)
				return
			}

			scopeClaim, ok := claims["scope"].(string)
			if !ok || !scopeContains(scopeClaim, requiredScope) {
				w.Header().Set("WWW-Authenticate", "Bearer")
				apiError := types.APIError{
					Code:          http.StatusForbidden,
					Message:       "Incorrect scope",
					Operation:     types.FormatOperation(r.Method, r.URL.Path),
					EmbeddedError: "",
				}
				utils.WriteJSON(w, http.StatusForbidden, apiError)
				return
			}

			handler(w, r)
		})
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.EnvAPIServerConfig.JWTSecret), nil
	})
}

func scopeContains(scopeClaim, requiredScope string) bool {
	scopes := strings.Split(scopeClaim, ",")
	for _, scope := range scopes {
		if strings.TrimSpace(scope) == requiredScope {
			return true
		}
	}

	return false
}
