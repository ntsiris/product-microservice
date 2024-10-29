package auth

import (
	"net/http"
	"ntsiris/product-microservice/internal/config"
	"ntsiris/product-microservice/internal/types"
	"ntsiris/product-microservice/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) error {
	scope := r.URL.Query().Get("scope")
	if scope == "" {
		return &types.APIError{
			Code:          http.StatusBadRequest,
			Message:       "Scope is required",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: "",
		}
	}

	claims := jwt.MapClaims{
		"scope": scope,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.EnvAPIServerConfig.JWTSecret))
	if err != nil {
		return &types.APIError{
			Code:          http.StatusInternalServerError,
			Message:       "Token generation failed",
			Operation:     types.FormatOperation(r.Method, r.URL.Path),
			EmbeddedError: err.Error(),
		}
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}
