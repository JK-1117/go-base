package controller

import (
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/jk1117/go-base/internal/database"
	logging "github.com/jk1117/go-base/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/sqlc-dev/pqtype"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	FirstName       *string   `json:"first_name"`
	LastName        *string   `json:"last_name"`
	Email           string    `json:"email"`
	IsAdministrator bool      `json:"is_admin"`
	Active          bool      `json:"active"`
}

type LoginSession struct {
	SessionID string      `json:"session_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	UserID    uuid.UUID   `json:"user_id"`
	LastLogin time.Time   `json:"last_login"`
	IpAddr    pqtype.Inet `json:"ip_addr"`
	UserAgent *string     `json:"user_agent,omitempty"`
	ExpiredAt time.Time   `json:"expired_at"`
}

type ROLE string

const CLIENT ROLE = "CLIENT"
const RIDER ROLE = "RIDER"
const MERCHANT ROLE = "MERCHANT"
const ADMIN ROLE = "ADMIN"

const MIN_PASSWORD_ENTROPY = 60

var base32RawStdEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

type SignUpParams struct {
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (controller *Controller) SignUp(c echo.Context, params SignUpParams) (uuid.UUID, error) {
	logger, _ := logging.GetLogger()
	err := validateCreateAccount(params.Password, params.Email)
	if err != nil {
		return uuid.Nil, err
	}
	_, err = controller.q.GetAccountByEmail(c.Request().Context(), params.Email)
	if err == nil {
		return uuid.Nil, errors.New(fmt.Sprintf("User with email: %v already exists, please try login.", params.Email))
	}

	tx, err := controller.db.Begin()
	if err != nil {
		logger.App.Err(fmt.Sprintf("Error Creating Account, error: %v, payload: %v", err, params))
		return uuid.Nil, errors.New("Something went wrong, try again later.")
	}
	defer tx.Rollback()
	qtx := controller.q.WithTx(tx)

	passHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		return uuid.Nil, err
	}
	dbFirstName := sql.NullString{}
	if params.FirstName != "" {
		dbFirstName.String = params.FirstName
		dbFirstName.Valid = true
	}
	dbLastName := sql.NullString{}
	if params.LastName != "" {
		dbLastName.String = params.LastName
		dbLastName.Valid = true
	}
	account, err := qtx.CreateAccount(c.Request().Context(), database.CreateAccountParams{
		ID:              uuid.New(),
		Password:        string(passHash),
		FirstName:       dbFirstName,
		LastName:        dbLastName,
		Email:           params.Email,
		IsAdministrator: false,
	})
	if err != nil {
		logger.App.Err(fmt.Sprintf("Error Creating Account, error: %v, payload: %v", err, params))
		return uuid.Nil, errors.New("Something went wrong, try again later.")
	}

	_, err = qtx.CreateUserRole(c.Request().Context(), database.CreateUserRoleParams{
		UserID: account.ID,
		Role:   string(CLIENT),
	})
	if err != nil {
		logger.App.Err(fmt.Sprintf("Error Creating User Role, error: %v, payload: %v", err, params))
		return uuid.Nil, errors.New("Something went wrong, try again later.")
	}

	if err = tx.Commit(); err != nil {
		logger.App.Err(fmt.Sprintf("Error Commiting SignUp, error: %v, payload: %v", err, params))
		return uuid.Nil, errors.New("Something went wrong, try again later.")
	}
	return account.ID, nil
}

func validateCreateAccount(password string, email string) error {
	if len(email) == 0 {
		return errors.New("Email is missing.")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("Invalid Email Format: " + err.Error())
	}

	if len(password) == 0 {
		return errors.New("Password is missing.")
	}
	if len(password) > 72 {
		return errors.New("Password is too long.")
	}
	err = passwordvalidator.Validate(password, MIN_PASSWORD_ENTROPY)
	if err != nil {
		return err
	}
	return nil
}

type VerifyAccountParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (controller *Controller) VerifyAccount(c echo.Context, params VerifyAccountParams) (*Account, error) {
	if params.Password == "" || params.Email == "" {
		return nil, errors.New("Incorrect email or password.")
	}
	account, err := controller.q.GetAccountByEmail(c.Request().Context(), params.Email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(params.Password)); err != nil {
		return nil, err
	}

	return parseAccount(account), nil
}

func parseLoginSession(dbLoginSession database.LoginSession) *LoginSession {
	var userAgent *string
	if dbLoginSession.UserAgent.Valid {
		userAgent = &dbLoginSession.UserAgent.String
	}

	return &LoginSession{
		SessionID: dbLoginSession.SessionID,
		CreatedAt: dbLoginSession.CreatedAt,
		UpdatedAt: dbLoginSession.UpdatedAt,
		UserID:    dbLoginSession.UserID,
		LastLogin: dbLoginSession.LastLogin,
		IpAddr:    dbLoginSession.IpAddr,
		UserAgent: userAgent,
		ExpiredAt: dbLoginSession.ExpiredAt,
	}
}

func parseAccount(dbAccount database.Account) *Account {
	var firstName *string
	if dbAccount.FirstName.Valid {
		firstName = &dbAccount.FirstName.String
	}
	var lastName *string
	if dbAccount.LastName.Valid {
		lastName = &dbAccount.LastName.String
	}

	return &Account{
		ID:              dbAccount.ID,
		CreatedAt:       dbAccount.CreatedAt,
		UpdatedAt:       dbAccount.UpdatedAt,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           dbAccount.Email,
		IsAdministrator: dbAccount.IsAdministrator,
		Active:          dbAccount.Active,
	}
}
