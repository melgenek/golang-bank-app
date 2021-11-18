package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang_bank_demo/src/dto"
	"golang_bank_demo/src/errors"
	"golang_bank_demo/src/model"
	"golang_bank_demo/src/service"
	"net/http"
	"strconv"
)

type AccountApi struct {
	accountService service.AccountService
	auth           *AuthenticatedApi
}

func NewAccountApi(accountService service.AccountService, auth *AuthenticatedApi) *AccountApi {
	return &AccountApi{accountService: accountService, auth: auth}
}

func (api *AccountApi) Router() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/accounts", api.auth.Authenticated(api.createAccount)).Methods("POST")
	router.Handle("/accounts/{id:[1-9][0-9]*}", api.auth.Authenticated(api.getAccount)).Methods("GET")
	router.Handle("/top-up", api.auth.Authenticated(api.topUp)).Methods("POST")
	router.Handle("/transfer", api.auth.Authenticated(api.transfer)).Methods("POST")
	return router
}

func (api *AccountApi) createAccount(userId model.UserId) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if account, err := api.accountService.Create(userId); err == nil {
			writeResponse(w, dto.AccountFromModel(account), http.StatusCreated)
		} else {
			handleServiceError(w, err)
		}
	})
}

func (api *AccountApi) getAccount(userId model.UserId) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr, idFound := mux.Vars(r)["id"]
		if !idFound {
			writeResponse(w, &dto.ErrorResponse{Message: "The request is missing the account id"}, http.StatusBadRequest)
		} else if id, err := strconv.Atoi(idStr); err != nil {
			writeResponse(w, &dto.ErrorResponse{Message: "The account id must be a number"}, http.StatusBadRequest)
		} else if account, err := api.accountService.Get(model.AccountId(id), userId); err == nil {
			writeResponse(w, dto.AccountFromModel(account), http.StatusOK)
		} else {
			handleServiceError(w, err)
		}
	})
}

func (api *AccountApi) topUp(userId model.UserId) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request dto.TopUpRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeResponse(w, &dto.ErrorResponse{Message: "The request is not a valid json"}, http.StatusBadRequest)
		} else if err := api.accountService.TopUp(&request, userId); err == nil {
			writeResponse(w, "{}", http.StatusOK)
		} else {
			handleServiceError(w, err)
		}
	})
}

func (api *AccountApi) transfer(userId model.UserId) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request dto.TransferRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeResponse(w, &dto.ErrorResponse{Message: "The request is not a valid json"}, http.StatusBadRequest)
		} else if err := api.accountService.Transfer(&request, userId); err == nil {
			writeResponse(w, "{}", http.StatusOK)
		} else {
			handleServiceError(w, err)
		}
	})
}

func handleServiceError(w http.ResponseWriter, err error) {
	errResponse := &dto.ErrorResponse{Message: err.Error()}
	switch err.(type) {
	case *errors.AccountDoesNotExistError:
		writeResponse(w, errResponse, http.StatusNotFound)
	case *errors.BalanceTooLowError:
		writeResponse(w, errResponse, http.StatusBadRequest)
	case *errors.DuplicateAccountError:
		writeResponse(w, errResponse, http.StatusConflict)
	case *errors.ForbiddenAccountAccessError:
		writeResponse(w, errResponse, http.StatusForbidden)
	case *errors.InternalServerError:
		writeResponse(w, errResponse, http.StatusInternalServerError)
	case *errors.ValidationError:
		writeResponse(w, errResponse, http.StatusBadRequest)
	default:
		writeResponse(w, &dto.ErrorResponse{Message: "Unhandled error"}, http.StatusInternalServerError)
	}
}
