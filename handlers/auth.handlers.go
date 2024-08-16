package handlers

import (
	// "fmt"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/auth_views"
	"github.com/gonext-tech/internal/views/components"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	auth_sessions_key string = "authenticate-sessions"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	email_key         string = "email"
	tzone_key         string = "time_zone"
)

/********** Handlers for Auth Views **********/

type AuthService interface {
	GetALL(limit, page int, orderBy, sortBy, searchTerm, statuss string) ([]models.User, models.Meta, error)
	CreateUser(u models.User, passwordConfirm string) error
	CheckEmail(email string) (models.User, error)
	// GetUserById(id int) (services.User, error)
}

func NewAuthHandler(us AuthService) *AuthHandler {

	return &AuthHandler{
		UserServices: us,
	}
}

type AuthHandler struct {
	UserServices AuthService
}

func (ah *AuthHandler) HomeHandler(c echo.Context) error {
	homeView := auth_views.Home(fromProtected)
	isError = false
	log.Println("frommProrected", fromProtected)
	return renderView(c, auth_views.HomeIndex(
		"| Home",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		homeView,
	))
}

func (ah *AuthHandler) RegisterHandler(c echo.Context) error {
	registerView := auth_views.Register(fromProtected)
	isError = false

	if c.Request().Method == "POST" {
		user := models.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
			//Username: c.FormValue("username"),
		}
		passwordConfirm := c.FormValue("password_confirm")

		err := ah.UserServices.CreateUser(user, passwordConfirm)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				err = errors.New("the email is already in use")
				setFlashmessages(c, "error", fmt.Sprintf(
					"something went wrong: %s",
					err,
				))

				return c.Redirect(http.StatusSeeOther, "/register")
			}

			setFlashmessages(c, "error", fmt.Sprintf(
				"something went wrong: %s",
				err,
			))

			return c.Redirect(http.StatusSeeOther, "/register")
		}

		setFlashmessages(c, "success", "You have successfully registered!!")

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return renderView(c, auth_views.RegisterIndex(
		"| Register",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		registerView,
	))
}

func (ah *AuthHandler) SearchUser(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 5
	}
	orderBy := c.QueryParam("orderBy")
	if orderBy == "" {
		orderBy = "desc"
	}
	sortBy := c.QueryParam("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	status := c.QueryParam("status")
	searchTerm := c.QueryParam("searchTerm")
	users, _, err := ah.UserServices.GetALL(limit, page, orderBy, sortBy, searchTerm, status)
	if err != nil {
		isError = false
		errorMsg = "can't fetch projects"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	return renderView(c, components.UserResult(users))
}

func (ah *AuthHandler) LoginHandler(c echo.Context) error {
	loginView := auth_views.Login(fromProtected)
	isError = false

	if c.Request().Method == "POST" {
		// obtaining the time zone from the POST request of the login form
		tzone := ""
		if len(c.Request().Header["X-Timezone"]) != 0 {
			tzone = c.Request().Header["X-Timezone"][0]
		}
		log.Println("authhhh")
		// Authentication goes here
		user, err := ah.UserServices.CheckEmail(c.FormValue("email"))
		log.Println("authhhh-user", user)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				setFlashmessages(c, "error", "There is no user with that email")

				return c.Redirect(http.StatusSeeOther, "/login")
			}

			log.Println("did we enter here?", user)
			setFlashmessages(c, "error", fmt.Sprintf(
				"something went wrong: %s",
				err,
			))
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)

		log.Println("errrr", err)
		if err != nil {
			// In production you have to give the user a generic message
			setFlashmessages(c, "error", "Incorrect password")

			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Get Session and setting Cookies
		sess, err := session.Get(auth_sessions_key, c)
		log.Println("sesss-errr", err)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   360000, // in seconds
			HttpOnly: true,
		}

		log.Println("sesss", sess)
		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:    true,
			user_id_key: user.ID,
			email_key:   user.Email,
			tzone_key:   tzone,
		}

		log.Println("sesss", sess)
		sess.Save(c.Request(), c.Response())

		setFlashmessages(c, "success", "You have successfully logged in!!")

		return c.Redirect(http.StatusSeeOther, "/project")
	}

	return renderView(c, auth_views.LoginIndex(
		"| Login",
		"",
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		loginView,
	))
}

func (ah *AuthHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			// fmt.Println(ok, auth)
			fromProtected = false

			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if email, ok := sess.Values[email_key].(string); ok && len(email) != 0 {
			c.Set(email_key, email) // set the username in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		fromProtected = true

		return next(c)
	}
}
