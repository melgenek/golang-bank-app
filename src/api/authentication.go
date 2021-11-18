package api

import (
	"golang_bank_demo/src/dto"
	"golang_bank_demo/src/model"
	"golang_bank_demo/src/service"
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

type AuthenticatedApi struct {
	authenticationService service.AuthenticationService
}

func NewAuthenticatedApi(authenticationService service.AuthenticationService) *AuthenticatedApi {
	return &AuthenticatedApi{authenticationService: authenticationService}
}

type AuthenticatedHandler func(id model.UserId) http.Handler

func (api *AuthenticatedApi) Authenticated(next AuthenticatedHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, bearerPrefix) {
			writeResponse(w, &dto.ErrorResponse{Message: "Unauthorized"}, http.StatusUnauthorized)
		} else if user, err := api.authenticationService.GetUser(strings.TrimPrefix(auth, bearerPrefix)); err != nil {
			writeResponse(w, &dto.ErrorResponse{Message: err.Error()}, http.StatusForbidden)
		} else {
			next(*user).ServeHTTP(w, r)
		}
	})
}
