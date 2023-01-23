package v1

import (
	"avito-tech/app/internal/controller/http/dto"
	"avito-tech/app/internal/domain/entity"
	"avito-tech/app/internal/domain/service/user"
	custom_error "avito-tech/app/internal/errors"
	"context"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

// TODO: Docker compose
// TODO: Unit Tests

type Service interface {
	CreateUser(ctx context.Context) (int, error)
	FindUserById(ctx context.Context, id int) (*entity.User, error)
	ChargeBalance(ctx context.Context, userChargeMoneyDTO dto.UserChargeMoneyDTO) error
	RemoveBalance(ctx context.Context, userRemoveMoneyDTO dto.UserRemoveMoneyDTO) error
	TransferMoney(ctx context.Context, data dto.TransferMoneyDTO) error
}

var (
	usr      = "/users/:id"
	users    = "/users/"
	add      = "/api/v1/charge"
	remove   = "/api/v1/remove"
	transfer = "/api/v1/transfer"
)

type UserHandler struct {
	service *user.UserService
}

func NewUserHandler(service *user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u UserHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, users, custom_error.Middleware(u.createUser))
	router.HandlerFunc(http.MethodGet, usr, custom_error.Middleware(u.getUser))
	router.HandlerFunc(http.MethodPost, add, custom_error.Middleware(u.chargeUserBalance))
	router.HandlerFunc(http.MethodPost, remove, custom_error.Middleware(u.removeUserBalance))
	router.HandlerFunc(http.MethodPost, transfer, custom_error.Middleware(u.transferMoney))
}

// Create user
// @Summary      Create user
// @Description  Create user with default values
// @Produce      json
// @Success      201
// @Failure      400
// @Failure      404
// @Failure      418
// @Failure      500
// @Router       /users [post]
func (u UserHandler) createUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Accept", "application/json")

	userId, err := u.service.CreateUser(r.Context())
	if err != nil {
		return err
	}

	resp := map[string]int{
		"user_id": userId,
	}

	marshalJson, err := json.Marshal(resp)
	if err != nil {
		return errors.New("failed to marshal")
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(marshalJson)
	return nil
}

// Get user balance
// @Summary      Get user balance
// @Description  Get user balance by id
// @Produce      json
// @Produce      json
// @Param		id	path int true "User id"
// @Success      201
// @Failure      500
// @Failure      404
// @Failure      400
// @Failure      418
// @Router       /users/{id} [get]
func (u UserHandler) getUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Accept", "application/json")
	params := httprouter.ParamsFromContext(r.Context())

	userId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		return errors.New("failed to convert")
	}

	usr, err := u.service.FindUserById(r.Context(), userId)
	if err != nil {
		return err
	}

	marshaledUser, err := usr.Marshal()
	if err != nil {
		return errors.New("failed to marshal")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUser)
	return nil
}

// Charge user balance
// @Summary      Charge balance
// @Description  add certain amount of money to user
// @Accept       json
// @Produce      json
// @Param		message body user.UserChargeMoneyDTO true "Charge money"
// @Success      201
// @Failure      400
// @Failure      404
// @Failure      500
// @Failure      418
// @Router       /api/v1/charge [post]
func (u UserHandler) chargeUserBalance(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read body")
	}

	userChargeMoneyDTO := new(dto.UserChargeMoneyDTO)

	if err := json.Unmarshal(body, userChargeMoneyDTO); err != nil {
		return errors.New("failed to unmarshal")
	}

	if err := u.service.ChargeBalance(r.Context(), user.UserChargeMoneyDTO(*userChargeMoneyDTO)); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// Decrease user balance
// @Summary      Decrease user balance
// @Description  remove certain amount of user`s money
// @Accept       json
// @Produce      json
// @Param		message body user.UserRemoveMoneyDTO true "Remove money"
// @Success      201
// @Failure      400
// @Failure      404
// @Failure      418
// @Failure      500
// @Router       /api/v1/remove [post]
func (u UserHandler) removeUserBalance(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read body")
	}

	userRemoveMoneyDTO := new(dto.UserRemoveMoneyDTO)

	if err := json.Unmarshal(body, userRemoveMoneyDTO); err != nil {
		return errors.New("failed to unmarshal")
	}

	if err := u.service.RemoveBalance(r.Context(), user.UserRemoveMoneyDTO(*userRemoveMoneyDTO)); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

// Transfer Money
// @Summary      Transfer Money
// @Description  transfer money from one user to another
// @Accept       json
// @Produce      json
// @Param		message body user.TransferMoneyDTO true "Remove money"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      418
// @Failure      500
// @Router       /api/v1/transfer [post]
func (u UserHandler) transferMoney(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read body")
	}

	transferMoneyDTO := new(dto.TransferMoneyDTO)

	if err := json.Unmarshal(body, transferMoneyDTO); err != nil {
		return errors.New("failed to unmarshal")
	}

	if err := u.service.TransferMoney(r.Context(), user.TransferMoneyDTO(*transferMoneyDTO)); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
